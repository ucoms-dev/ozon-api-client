package ozon

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSupplyOrderMethodsUseSwaggerHTTPVerbs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		call func(context.Context, *Client) error
	}{
		{
			name: "status counter",
			path: "/v1/supply-order/status/counter",
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().GetSupplyOrdersByStatus(ctx)
				return err
			},
		},
		{
			name: "timeslot get",
			path: "/v1/supply-order/timeslot/get",
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().GetSupplyTimeslots(ctx, &GetSupplyTimeslotsParams{})
				return err
			},
		},
		{
			name: "timeslot update",
			path: "/v1/supply-order/timeslot/update",
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().UpdateSupplyTimeslot(ctx, &UpdateSupplyTimeslotParams{})
				return err
			},
		},
		{
			name: "timeslot status",
			path: "/v1/supply-order/timeslot/status",
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().GetSupplyTimeslotStatus(ctx, &GetSupplyTimeslotStatusParams{})
				return err
			},
		},
		{
			name: "pass create",
			path: "/v1/supply-order/pass/create",
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().CreatePass(ctx, &CreatePassParams{})
				return err
			},
		},
		{
			name: "pass status",
			path: "/v1/supply-order/pass/status",
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().GetPass(ctx, &GetPassParams{})
				return err
			},
		},
		{
			name: "bundle",
			path: "/v1/supply-order/bundle",
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().GetSupplyContent(ctx, &GetSupplyContentParams{})
				return err
			},
		},
		{
			name: "cancel",
			path: "/v1/supply-order/cancel",
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().CancelSuppyOrder(ctx, &CancelSuppyOrderParams{})
				return err
			},
		},
		{
			name: "cancel status",
			path: "/v1/supply-order/cancel/status",
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().StatusCancelledSupplyOrder(ctx, &StatusCancelledSupplyOrderParams{})
				return err
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			handler := requestContractHandler(t, http.MethodPost, test.path, "", "{}")
			server := httptest.NewServer(handler)
			defer server.Close()

			if err := test.call(context.Background(), NewClient(WithURI(server.URL))); err != nil {
				t.Fatalf("call %s: %v", test.path, err)
			}
		})
	}
}

func TestGetSupplyOrdersByStatusSendsNoRequestBody(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}
		if len(body) != 0 {
			t.Errorf("request body = %q, want empty", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[]}`))
	}))
	defer server.Close()

	if _, err := NewClient(WithURI(server.URL)).FBO().GetSupplyOrdersByStatus(context.Background()); err != nil {
		t.Fatalf("GetSupplyOrdersByStatus: %v", err)
	}
}

func TestGetSupplyContentUsesCurrentSwaggerContract(t *testing.T) {
	t.Parallel()

	handler := requestContractHandler(
		t,
		http.MethodPost,
		"/v1/supply-order/bundle",
		`{"bundle_ids":["bundle-1"],"is_asc":true,"last_id":"100","limit":50,"query":"offer","sort_field":"SKU","item_tags_calculation":{"dropoff_warehouse_id":"10","storage_warehouse_ids":["20","30"]}}`,
		`{"items":[{"offer_id":"offer-1","placement_zone":"ambient","sku":123,"tags":["oversize"]}],"total_count":1}`,
	)
	server := httptest.NewServer(handler)
	defer server.Close()

	response, err := NewClient(WithURI(server.URL)).FBO().GetSupplyContent(context.Background(), &GetSupplyContentParams{
		BundleIds: []string{"bundle-1"},
		IsAsc:     true,
		LastId:    "100",
		Limit:     50,
		Query:     "offer",
		SortField: "SKU",
		ItemTagsCalculation: &GetSupplyContentItemTagsCalculation{
			DropoffWarehouseId:  "10",
			StorageWarehouseIds: []string{"20", "30"},
		},
	})
	if err != nil {
		t.Fatalf("GetSupplyContent: %v", err)
	}
	if len(response.Items) != 1 || response.Items[0].OfferId != "offer-1" || response.Items[0].PlacementZone != "ambient" || len(response.Items[0].Tags) != 1 || response.Items[0].Tags[0] != "oversize" {
		t.Fatalf("current supply item fields were not decoded: %+v", response.Items)
	}
}
