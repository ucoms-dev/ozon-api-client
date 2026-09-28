package ozon

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAvailablePromotionsAutoAddDates(t *testing.T) {
	const date = "2026-10-13T03:00:00.123456789+03:00"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/actions" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"result":[{"id":1,"auto_add_dates":["`+date+`"]},{"id":2}]}`)
	}))
	defer server.Close()
	response, err := NewClient(WithURI(server.URL)).Promotions().GetAvailablePromotions(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Result) != 2 || len(response.Result[0].AutoAddDates) != 1 || len(response.Result[1].AutoAddDates) != 0 {
		t.Fatalf("dates not preserved: %+v", response.Result)
	}
	if got := response.Result[0].AutoAddDates[0].Format(time.RFC3339Nano); got != date {
		t.Fatalf("date = %q, want %q", got, date)
	}
}

func TestPromotionLegacyCompatibility(t *testing.T) {
	for _, path := range []string{"/v1/actions/products/activate", "/v1/actions/products/deactivate", "/v1/actions/products", "/v1/actions/candidates", "/v1/actions/discounts-task/approve", "/v1/actions/discounts-task/decline"} {
		t.Run(path, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != http.MethodPost || r.URL.Path != path {
					t.Errorf("request = %s %s, want POST %s", r.Method, r.URL.Path, path)
				}
				var body map[string]json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				if path == "/v1/actions/products/activate" {
					var products []map[string]json.RawMessage
					if err := json.Unmarshal(body["products"], &products); err != nil || len(products) != 1 || string(products[0]["action_price"]) != "123.45" {
						t.Errorf("legacy price must remain a number: %s", body["products"])
					}
				}
				_, _ = io.WriteString(w, `{"result":{}}`)
			}))
			defer server.Close()
			p := NewClient(WithURI(server.URL)).Promotions()
			var err error
			switch path {
			case "/v1/actions/products/activate":
				_, err = p.AddToPromotion(context.Background(), &AddProductToPromotionParams{ActionId: 1, Products: []AddProductToPromotionProduct{{ProductId: 2, ActionPrice: 123.45, Stock: 3}}})
			case "/v1/actions/products/deactivate":
				_, err = p.RemoveProduct(context.Background(), &RemoveProductFromPromotionParams{ActionId: 1, ProductIds: []float64{2}})
			case "/v1/actions/products":
				_, err = p.ProductsInPromotion(context.Background(), &ProductsInPromotionParams{ActionId: 1, Limit: 100})
			case "/v1/actions/candidates":
				_, err = p.ProductsAvailableForPromotion(context.Background(), &ProductsAvailableForPromotionParams{ActionId: 1, Limit: 100})
			case "/v1/actions/discounts-task/approve":
				_, err = p.ApproveDiscountRequest(context.Background(), &DiscountRequestParams{})
			case "/v1/actions/discounts-task/decline":
				_, err = p.DeclineDiscountRequest(context.Background(), &DiscountRequestParams{})
			}
			if err != nil || calls != 1 {
				t.Fatalf("calls=%d err=%v", calls, err)
			}
		})
	}
}
