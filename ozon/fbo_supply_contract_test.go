package ozon

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSupplyOrderMethodsUseSwaggerHTTPVerbs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		path     string
		wantBody string
		call     func(context.Context, *Client) error
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
			name:     "timeslot get",
			path:     "/v1/supply-order/timeslot/get",
			wantBody: `{"supply_order_id":101}`,
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().GetSupplyTimeslots(ctx, &GetSupplyTimeslotsParams{SupplyOrderId: 101})
				return err
			},
		},
		{
			name:     "timeslot update",
			path:     "/v1/supply-order/timeslot/update",
			wantBody: `{"supply_order_id":101,"timeslot":{"from":"2026-09-03T10:00:00Z","to":"2026-09-03T11:00:00Z"}}`,
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().UpdateSupplyTimeslot(ctx, &UpdateSupplyTimeslotParams{
					SupplyOrderId: 101,
					Timeslot: SupplyTimeslotValueTimeslot{
						From: time.Date(2026, 9, 3, 10, 0, 0, 0, time.UTC),
						To:   time.Date(2026, 9, 3, 11, 0, 0, 0, time.UTC),
					},
				})
				return err
			},
		},
		{
			name:     "timeslot status",
			path:     "/v1/supply-order/timeslot/status",
			wantBody: `{"operation_id":"op-1"}`,
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().GetSupplyTimeslotStatus(ctx, &GetSupplyTimeslotStatusParams{OperationId: "op-1"})
				return err
			},
		},
		{
			name:     "pass create",
			path:     "/v1/supply-order/pass/create",
			wantBody: `{"supply_order_id":101,"vehicle":{"driver_name":"Driver","driver_phone":"+70000000000","vehicle_model":"Van","vehicle_number":"A001AA"}}`,
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().CreatePass(ctx, &CreatePassParams{
					SupplyOrderId: 101,
					Vehicle: GetSupplyRequestInfoVehicle{
						DriverName:    "Driver",
						DriverPhone:   "+70000000000",
						VehicleModel:  "Van",
						VehicleNumber: "A001AA",
					},
				})
				return err
			},
		},
		{
			name:     "pass status",
			path:     "/v1/supply-order/pass/status",
			wantBody: `{"operation_id":"op-1"}`,
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().GetPass(ctx, &GetPassParams{OperationId: "op-1"})
				return err
			},
		},
		{
			name:     "bundle",
			path:     "/v1/supply-order/bundle",
			wantBody: `{"bundle_ids":["bundle-1"],"is_asc":true,"last_id":"100","limit":50,"query":"offer","sort_field":"SKU"}`,
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().GetSupplyContent(ctx, &GetSupplyContentParams{BundleIds: []string{"bundle-1"}, IsAsc: true, LastId: "100", Limit: 50, Query: "offer", SortField: "SKU"})
				return err
			},
		},
		{
			name:     "cancel",
			path:     "/v1/supply-order/cancel",
			wantBody: `{"order_id":101}`,
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().CancelSuppyOrder(ctx, &CancelSuppyOrderParams{OrderId: 101})
				return err
			},
		},
		{
			name:     "cancel status",
			path:     "/v1/supply-order/cancel/status",
			wantBody: `{"operation_id":"op-1"}`,
			call: func(ctx context.Context, client *Client) error {
				_, err := client.FBO().StatusCancelledSupplyOrder(ctx, &StatusCancelledSupplyOrderParams{OperationId: "op-1"})
				return err
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			handler := requestContractHandler(t, http.MethodPost, test.path, test.wantBody, "{}")
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
		`{"items":[{"offer_id":"offer-1","placement_zone":"PRODUCTS","sku":123,"tags":["OVERSIZE"]}],"total_count":1}`,
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
	if len(response.Items) != 1 || response.Items[0].OfferId != "offer-1" || response.Items[0].PlacementZone != "PRODUCTS" || len(response.Items[0].Tags) != 1 || response.Items[0].Tags[0] != "OVERSIZE" {
		t.Fatalf("current supply item fields were not decoded: %+v", response.Items)
	}
}
