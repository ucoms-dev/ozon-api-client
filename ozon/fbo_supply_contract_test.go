package ozon

import (
	"context"
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
