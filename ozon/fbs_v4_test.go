package ozon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetFBSShipmentsListV4UsesCursorContract(t *testing.T) {
	t.Parallel()

	since := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC)
	handler := requestContractHandler(
		t,
		http.MethodPost,
		"/v4/posting/fbs/list",
		`{"cursor":"page-1","filter":{"delivery_method_ids":["11"],"order_numbers":["42"],"provider_ids":["12"],"since":"2026-08-01T00:00:00Z","statuses":["awaiting_deliver"],"to":"2026-08-02T00:00:00Z","warehouse_ids":["13"]},"limit":100,"sort_dir":"ASC","translit":true,"with":{"analytics_data":true,"barcodes":true,"financial_data":true,"legal_info":true}}`,
		`{"cursor":"page-2","has_next":true,"postings":[{"order_id":42,"order_number":"42","posting_number":"100-1","status":"awaiting_deliver","container_sort_type":"SORT","delivery_schema":"FBS","destination_place_id":17,"destination_place_name":"SC","integration_type_flow":"ozon","is_click_and_collect":true,"is_presortable":true,"require_blr_traceable_attrs":true,"shipment_date_without_delay":"2026-08-02T03:00:00Z","tariffication":{"current_tariff_charge":{"amount":"12.50","currency":"RUB"},"current_tariff_rate":2.5,"current_tariff_type":"discount"},"volume_weight":2.5}]}`,
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	client := NewClient(WithURI(server.URL))

	response, err := client.FBS().GetFBSShipmentsListV4(context.Background(), &GetFBSShipmentsListV4Params{
		Cursor: "page-1",
		Filter: GetFBSShipmentsListV4Filter{
			DeliveryMethodIds: []string{"11"},
			OrderNumbers:      []string{"42"},
			ProviderIds:       []string{"12"},
			Since:             since,
			Statuses:          []string{"awaiting_deliver"},
			To:                to,
			WarehouseIds:      []string{"13"},
		},
		Limit:    100,
		SortDir:  Order("ASC"),
		Translit: true,
		With: &GetFBSShipmentsListV4With{
			AnalyticsData: true,
			Barcodes:      true,
			FinancialData: true,
			LegalInfo:     true,
		},
	})
	if err != nil {
		t.Fatalf("GetFBSShipmentsListV4: %v", err)
	}
	if response.Cursor != "page-2" || !response.HasNext {
		t.Fatalf("pagination = (%q, %v), want (%q, true)", response.Cursor, response.HasNext, "page-2")
	}
	if len(response.Postings) != 1 {
		t.Fatalf("postings length = %d, want 1", len(response.Postings))
	}
	posting := response.Postings[0]
	if posting.ContainerSortType != "SORT" || posting.DeliverySchema != "FBS" {
		t.Errorf("posting v4 fields were not decoded: %+v", posting)
	}
	if posting.VolumeWeight != 2.5 {
		t.Errorf("volume weight = %v, want 2.5", posting.VolumeWeight)
	}
	if posting.Tariffication.CurrentTariffCharge.Amount != "12.50" || posting.Tariffication.CurrentTariffCharge.Currency != "RUB" {
		t.Errorf("current tariff charge = %+v", posting.Tariffication.CurrentTariffCharge)
	}
}

func TestListUnprocessedShipmentsV4UsesCursorContract(t *testing.T) {
	t.Parallel()

	cutoffFrom := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	cutoffTo := time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC)
	handler := requestContractHandler(
		t,
		http.MethodPost,
		"/v4/posting/fbs/unfulfilled/list",
		`{"cursor":"page-1","filter":{"cutoff_from":"2026-08-01T00:00:00Z","cutoff_to":"2026-08-02T00:00:00Z","statuses":["awaiting_packaging"]},"limit":50,"sort_dir":"ASC","translit":true,"with":{"barcodes":true}}`,
		`{"count":1,"cursor":"page-2","has_next":false,"postings":[{"posting_number":"100-1","status":"awaiting_packaging"}]}`,
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	client := NewClient(WithURI(server.URL))

	response, err := client.FBS().ListUnprocessedShipmentsV4(context.Background(), &ListUnprocessedShipmentsV4Params{
		Cursor: "page-1",
		Filter: ListUnprocessedShipmentsV4Filter{
			CutoffFrom: &cutoffFrom,
			CutoffTo:   &cutoffTo,
			Statuses:   []string{"awaiting_packaging"},
		},
		Limit:    50,
		SortDir:  Order("ASC"),
		Translit: true,
		With: &ListUnprocessedShipmentsV4With{
			Barcodes: true,
		},
	})
	if err != nil {
		t.Fatalf("ListUnprocessedShipmentsV4: %v", err)
	}
	if response.Count != 1 || response.Cursor != "page-2" || response.HasNext {
		t.Fatalf("response pagination = (%d, %q, %v)", response.Count, response.Cursor, response.HasNext)
	}
	if len(response.Postings) != 1 || response.Postings[0].PostingNumber != "100-1" {
		t.Fatalf("postings = %+v", response.Postings)
	}
}
