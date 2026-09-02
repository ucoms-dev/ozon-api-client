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
		`{"cursor":"page-2","has_next":true,"postings":[{"order_id":42,"order_number":"42","posting_number":"100-1","status":"awaiting_deliver","analytics_data":{"city":"Moscow","client_delivery_date_begin":"2026-08-02T10:00:00Z","client_delivery_date_end":"2026-08-02T18:00:00Z"},"container_sort_type":"SORT","customer":{"customer_email":"buyer@example.test","customer_id":77},"delivery_schema":"FBS","destination_place_id":17,"destination_place_name":"SC","financial_data":{"products":[{"commission":{"amount":23.4,"currency":"RUB","percent":10},"customer_price":{"amount":"234.56","currency":"RUB"},"price":210.0,"product_id":456,"quantity":1}]},"integration_type_flow":"ozon","is_click_and_collect":true,"is_presortable":true,"optional":{"products_with_possible_mandatory_mark":["456"]},"products":[{"offer_id":"offer-1","price":{"amount":"234.56","currency":"RUB"},"quantity":1,"sku":456}],"requirements":{"products_requiring_change_country":["456"],"products_requiring_country":["456"],"products_requiring_gtd":["456"],"products_requiring_imei":["456"],"products_requiring_jw_uin":["456"],"products_requiring_mandatory_mark":["456"],"products_requiring_rnpt":["456"],"products_requiring_weight":["456"]},"require_blr_traceable_attrs":true,"shipment_date_without_delay":"2026-08-02T03:00:00Z","tariffication":{"current_tariff_charge":{"amount":"12.50","currency":"RUB"},"current_tariff_rate":2.5,"current_tariff_type":"discount"},"volume_weight":2.5}]}`,
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	client := NewClient(WithURI(server.URL))

	response, err := client.FBS().GetFBSShipmentsListV4Current(context.Background(), &GetFBSShipmentsListV4Params{
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
	if posting.AnalyticsData.City != "Moscow" || posting.AnalyticsData.ClientDeliveryDateBegin.IsZero() || posting.AnalyticsData.ClientDeliveryDateEnd.IsZero() {
		t.Errorf("FBS current analytics data was not decoded: %+v", posting.AnalyticsData)
	}
	if posting.Customer.CustomerEmail != "buyer@example.test" {
		t.Errorf("FBS current customer email was not decoded: %+v", posting.Customer)
	}
	if posting.Tariffication.CurrentTariffCharge.Amount != "12.50" || posting.Tariffication.CurrentTariffCharge.Currency != "RUB" {
		t.Errorf("current tariff charge = %+v", posting.Tariffication.CurrentTariffCharge)
	}
	if len(posting.Products) != 1 || posting.Products[0].Price.Amount != "234.56" || posting.Products[0].Price.Currency != "RUB" {
		t.Fatalf("FBS product money was not decoded: %+v", posting.Products)
	}
	if len(posting.Requirements.ProductsRequiringWeight) != 1 || posting.Requirements.ProductsRequiringWeight[0] != "456" {
		t.Errorf("FBS string SKU requirements were not decoded: %+v", posting.Requirements)
	}
	if len(posting.Optional.ProductsWithPossibleMandatoryMark) != 1 || posting.Optional.ProductsWithPossibleMandatoryMark[0] != "456" {
		t.Errorf("FBS optional string SKU was not decoded: %+v", posting.Optional)
	}
	if len(posting.FinancialData.Products) != 1 || posting.FinancialData.Products[0].CustomerPrice.Amount != "234.56" || posting.FinancialData.Products[0].CustomerPrice.Currency != "RUB" {
		t.Fatalf("FBS customer money was not decoded: %+v", posting.FinancialData.Products)
	}
	if posting.FinancialData.Products[0].Commission.Amount != 23.4 || posting.FinancialData.Products[0].Commission.Percent != 10 || posting.FinancialData.Products[0].Commission.Currency != "RUB" {
		t.Errorf("FBS nested commission was not decoded: %+v", posting.FinancialData.Products[0])
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

	response, err := client.FBS().ListUnprocessedShipmentsV4Current(context.Background(), &ListUnprocessedShipmentsV4Params{
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

func TestGetFBSShipmentsListV4PreservesV116ResponseContract(t *testing.T) {
	t.Parallel()

	handler := requestContractHandler(
		t,
		http.MethodPost,
		"/v4/posting/fbs/list",
		`{"filter":{"since":"0001-01-01T00:00:00Z","to":"0001-01-01T00:00:00Z"},"limit":1}`,
		`{"postings":[{"posting_number":"100-1","financial_data":{"products":[{"commission":{"amount":23.4,"currency":"RUB","percent":10},"customer_price":{"amount":"234.56","currency":"RUB"},"product_id":456}]},"optional":{"products_with_possible_mandatory_mark":["456"]},"products":[{"offer_id":"offer-1","price":{"amount":"234.56","currency":"RUB"},"sku":456}],"requirements":{"products_requiring_gtd":["456"]},"tariffication":{"current_tariff_charge":{"amount":"12.50","currency":"RUB"}}}]}`,
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	response, err := NewClient(WithURI(server.URL)).FBS().GetFBSShipmentsListV4(context.Background(), &GetFBSShipmentsListV4Params{Limit: 1})
	if err != nil {
		t.Fatalf("GetFBSShipmentsListV4: %v", err)
	}
	var _ []FBSPostingV4 = response.Postings
	posting := response.Postings[0]
	if posting.Products[0].Price != "234.56" || posting.Products[0].CurrencyCode != "RUB" {
		t.Fatalf("legacy product money was not adapted: %+v", posting.Products[0])
	}
	if posting.FinancialData.Products[0].ClientPrice != "234.56" || posting.FinancialData.Products[0].CommissionAmount != 23.4 {
		t.Fatalf("legacy financial data was not adapted: %+v", posting.FinancialData.Products[0])
	}
	if posting.Requirements.ProductsRequiringGTD[0] != 456 || posting.Optional.ProductsWithPossibleMandatoryMark[0] != 456 {
		t.Fatalf("legacy product IDs were not adapted: requirements=%+v optional=%+v", posting.Requirements, posting.Optional)
	}
	if posting.Tariffication.CurrentTariffCharge.Amount != "12.50" {
		t.Fatalf("legacy v4 tariffication was not preserved: %+v", posting.Tariffication)
	}
}

func TestFBSV4LegacyAdapterRejectsNonNumericProductIDs(t *testing.T) {
	t.Parallel()

	_, err := adaptFBSPostingV4(FBSPostingV4Current{
		Requirements: FBSPostingV4Requirements{ProductsRequiringGTD: []string{"not-a-number"}},
	})
	if err == nil {
		t.Fatal("adaptFBSPostingV4 accepted a non-numeric legacy product ID")
	}
}
