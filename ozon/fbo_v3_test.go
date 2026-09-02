package ozon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetFBOShipmentsListV3UsesCursorContract(t *testing.T) {
	t.Parallel()

	since := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC)
	handler := requestContractHandler(
		t,
		http.MethodPost,
		"/v3/posting/fbo/list",
		`{"cursor":"next-1","filter":{"order_numbers":["42"],"posting_numbers":["100-1"],"since":"2026-08-01T00:00:00Z","statuses":["delivering"],"to":"2026-08-02T00:00:00Z"},"limit":100,"sort_dir":"ASC","translit":true,"with":{"analytics_data":true,"financial_data":true,"legal_info":true}}`,
		`{"cursor":"next-2","has_next":true,"postings":[{"order_id":42,"order_number":"42","posting_number":"100-1","status":"delivering","substatus":"posting_on_way_to_city","analytics_data":{"city":"Moscow","client_delivery_date_begin":"2026-08-02T10:00:00Z","client_delivery_date_end":"2026-08-02T18:00:00Z","warehouse_name":"WH"},"cancellation":{"cancel_reason":"buyer request"},"external_order":{"is_external":true,"platform_name":"Ozon"},"financial_data":{"products":[{"commission":{"amount":12.5,"currency":"RUB","percent":10},"price":125.0,"product_id":123,"quantity":2}]},"legal_info":{"company_name":"UCOMS"},"products":[{"offer_id":"offer-1","price":{"amount":"123.45","currency":"RUB"},"quantity":2,"sku":123}]}]}`,
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	client := NewClient(WithURI(server.URL))

	response, err := client.FBO().GetShipmentsListV3Current(context.Background(), &GetFBOShipmentsListV3Params{
		Cursor: "next-1",
		Filter: GetFBOShipmentsListV3Filter{
			OrderNumbers:   []string{"42"},
			PostingNumbers: []string{"100-1"},
			Since:          since,
			Statuses:       []string{"delivering"},
			To:             to,
		},
		Limit:    100,
		SortDir:  Order("ASC"),
		Translit: true,
		With: &GetFBOShipmentsListV3With{
			AnalyticsData: true,
			FinancialData: true,
			LegalInfo:     true,
		},
	})
	if err != nil {
		t.Fatalf("GetShipmentsListV3: %v", err)
	}
	if response.Cursor != "next-2" || !response.HasNext {
		t.Fatalf("pagination = (%q, %v), want (%q, true)", response.Cursor, response.HasNext, "next-2")
	}
	if len(response.Postings) != 1 {
		t.Fatalf("postings length = %d, want 1", len(response.Postings))
	}
	posting := response.Postings[0]
	if posting.Substatus != "posting_on_way_to_city" {
		t.Errorf("substatus = %q", posting.Substatus)
	}
	if posting.Cancellation.CancelReason != "buyer request" {
		t.Errorf("cancel reason = %q", posting.Cancellation.CancelReason)
	}
	if !posting.ExternalOrder.IsExternal || posting.ExternalOrder.PlatformName != "Ozon" {
		t.Errorf("external order = %+v", posting.ExternalOrder)
	}
	if posting.LegalInfo.CompanyName != "UCOMS" {
		t.Errorf("legal company name = %q", posting.LegalInfo.CompanyName)
	}
	if posting.AnalyticsData.City != "Moscow" || posting.AnalyticsData.ClientDeliveryDateBegin.IsZero() || posting.AnalyticsData.ClientDeliveryDateEnd.IsZero() {
		t.Errorf("FBO current analytics data was not decoded: %+v", posting.AnalyticsData)
	}
	if len(posting.Products) != 1 || posting.Products[0].Price.Amount != "123.45" || posting.Products[0].Price.Currency != "RUB" {
		t.Fatalf("FBO product money was not decoded: %+v", posting.Products)
	}
	if len(posting.FinancialData.Products) != 1 || posting.FinancialData.Products[0].Commission.Amount != 12.5 || posting.FinancialData.Products[0].Commission.Currency != "RUB" || posting.FinancialData.Products[0].Commission.Percent != 10 {
		t.Errorf("FBO nested commission was not decoded: %+v", posting.FinancialData.Products)
	}
}

func TestGetFBOShipmentsListV3PreservesV116ResponseContract(t *testing.T) {
	t.Parallel()

	handler := requestContractHandler(
		t,
		http.MethodPost,
		"/v3/posting/fbo/list",
		`{"filter":{"since":"0001-01-01T00:00:00Z","to":"0001-01-01T00:00:00Z"},"limit":1}`,
		`{"cursor":"next","has_next":true,"postings":[{"posting_number":"100-1","analytics_data":{"warehouse_name":"WH"},"financial_data":{"products":[{"commission":{"amount":12.5,"currency":"RUB","percent":10},"product_id":123}]},"products":[{"offer_id":"offer-1","price":{"amount":"123.45","currency":"RUB"},"sku":123}]}]}`,
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	response, err := NewClient(WithURI(server.URL)).FBO().GetShipmentsListV3(context.Background(), &GetFBOShipmentsListV3Params{Limit: 1})
	if err != nil {
		t.Fatalf("GetShipmentsListV3: %v", err)
	}
	var _ []GetFBOShipmentsListResult = response.Postings
	posting := response.Postings[0]
	if posting.Products[0].Price != "123.45" || posting.Products[0].CurrencyCode != "RUB" {
		t.Fatalf("legacy product money was not adapted: %+v", posting.Products[0])
	}
	if posting.FinancialData.Products[0].CommissionAmount != 12.5 || posting.FinancialData.Products[0].CommissionsCurrencyCode != "RUB" {
		t.Fatalf("legacy commission was not adapted: %+v", posting.FinancialData.Products[0])
	}
}
