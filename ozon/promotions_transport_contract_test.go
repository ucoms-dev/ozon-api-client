package ozon

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func promotionTransportCases() map[string]func(context.Context, *Promotions) (interface{}, error) {
	date := time.Date(2026, 10, 13, 0, 0, 0, 0, time.UTC)
	return map[string]func(context.Context, *Promotions) (interface{}, error){
		"/v2/actions/products": func(ctx context.Context, p *Promotions) (interface{}, error) {
			return p.ProductsInPromotionV2(ctx, &ProductsInPromotionV2Params{ActionID: 1, Limit: 100})
		},
		"/v2/actions/candidates": func(ctx context.Context, p *Promotions) (interface{}, error) {
			return p.ProductsAvailableForPromotionV2(ctx, &ProductsAvailableForPromotionV2Params{ActionID: 1, Limit: 100})
		},
		"/v1/actions/products/update": func(ctx context.Context, p *Promotions) (interface{}, error) {
			return p.UpdateProducts(ctx, &UpdatePromotionProductsParams{ActionID: 1, Products: []UpdatePromotionProduct{{ProductID: 2, ActionPrice: PromotionMoney{Amount: "12.34", Currency: "RUB"}}}})
		},
		"/v2/actions/products/deactivate": func(ctx context.Context, p *Promotions) (interface{}, error) {
			return p.RemoveProductV2(ctx, &RemoveProductFromPromotionV2Params{ActionID: 1, ProductIDs: []PromotionProductID{2}})
		},
		"/v2/actions/auto-add/products/list": func(ctx context.Context, p *Promotions) (interface{}, error) {
			return p.ListAutoAddProductsV2(ctx, &ListAutoAddProductsV2Params{ActionID: 1, AutoAddDate: date, Limit: 100})
		},
		"/v2/actions/auto-add/products/candidates": func(ctx context.Context, p *Promotions) (interface{}, error) {
			return p.ListAutoAddCandidatesV2(ctx, &ListAutoAddCandidatesV2Params{ActionID: 1, AutoAddDate: date, Limit: 100})
		},
		"/v2/actions/auto-add/products/update": func(ctx context.Context, p *Promotions) (interface{}, error) {
			return p.UpdateAutoAddProductsV2(ctx, &UpdateAutoAddProductsV2Params{ActionID: 1, AutoAddDate: date, Products: []UpdateAutoAddProductV2{{ID: 2, ActionPrice: PromotionMoney{Amount: "12.34", Currency: "RUB"}}}})
		},
		"/v2/actions/auto-add/products/delete": func(ctx context.Context, p *Promotions) (interface{}, error) {
			return p.DeleteAutoAddProductsV2(ctx, &DeleteAutoAddProductsV2Params{ActionID: 1, AutoAddDate: date, ProductIDs: []PromotionProductID{2}})
		},
	}
}

func TestPromotionV2TransportContract(t *testing.T) {
	for path, invoke := range promotionTransportCases() {
		for _, status := range []int{400, 403, 404, 409, 429, 500} {
			t.Run(fmt.Sprintf("%s/%d", path, status), func(t *testing.T) {
				var calls atomic.Int64
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					if r.URL.Path != path || r.Method != http.MethodPost {
						t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
					}
					if r.Header.Get("Client-Id") != "client" || r.Header.Get("Api-Key") != "key" {
						t.Error("missing credential headers")
					}
					w.WriteHeader(status)
					_, _ = io.WriteString(w, `{"code":9,"message":"provider rejected request","details":[{"typeUrl":"test","value":"detail"}]}`)
				}))
				defer server.Close()
				response, err := invoke(context.Background(), NewClient(WithURI(server.URL), WithClientId("client"), WithAPIKey("key")).Promotions())
				if err != nil {
					t.Fatal(err)
				}
				value := reflect.ValueOf(response).Elem()
				if int(value.FieldByName("StatusCode").Int()) != status || value.FieldByName("Code").Int() != 9 || value.FieldByName("Message").String() != "provider rejected request" || value.FieldByName("Details").Len() != 1 {
					t.Fatalf("lost HTTP error metadata: %#v", response)
				}
				if calls.Load() != 1 {
					t.Fatalf("unexpected retries: %d", calls.Load())
				}
			})
		}
	}
}

func TestPromotionV2CancelledContext(t *testing.T) {
	for path, invoke := range promotionTransportCases() {
		t.Run(path, func(t *testing.T) {
			var calls atomic.Int64
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
			defer server.Close()
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err := invoke(ctx, NewClient(WithURI(server.URL)).Promotions())
			if !errors.Is(err, context.Canceled) || calls.Load() != 0 {
				t.Fatalf("err=%v calls=%d", err, calls.Load())
			}
		})
	}
}

func TestPromotionV2InvalidCall(t *testing.T) {
	for path, invoke := range promotionTransportCases() {
		t.Run(path, func(t *testing.T) {
			var calls atomic.Int64
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1) }))
			defer server.Close()
			client := NewClient(WithURI(server.URL)).Promotions()
			for _, input := range []struct {
				ctx        context.Context
				promotions *Promotions
			}{{nil, client}, {context.Background(), &Promotions{}}} {
				_, err := invoke(input.ctx, input.promotions)
				var inputErr *PromotionInputError
				if !errors.As(err, &inputErr) {
					t.Fatalf("invalid call: err=%v", err)
				}
			}
			if calls.Load() != 0 {
				t.Fatalf("invalid call reached server: %d", calls.Load())
			}
		})
	}
}

func TestPromotionV2DeadlineWithoutRetry(t *testing.T) {
	for path, invoke := range promotionTransportCases() {
		t.Run(path, func(t *testing.T) {
			var calls atomic.Int64
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				_, _ = io.Copy(io.Discard, r.Body)
				<-r.Context().Done()
			}))
			defer server.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()
			_, err := invoke(ctx, NewClient(WithURI(server.URL)).Promotions())
			if !errors.Is(err, context.DeadlineExceeded) || calls.Load() > 1 {
				t.Fatalf("err=%v calls=%d", err, calls.Load())
			}
		})
	}
}

func TestPromotionV2MalformedHTTPResponse(t *testing.T) {
	for path, invoke := range promotionTransportCases() {
		for _, body := range []string{"", "null", "{}", "{", `{"result":{}}`} {
			t.Run(path+"/"+body, func(t *testing.T) {
				var calls atomic.Int64
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); _, _ = io.WriteString(w, body) }))
				defer server.Close()
				_, err := invoke(context.Background(), NewClient(WithURI(server.URL)).Promotions())
				if err == nil || calls.Load() != 1 {
					t.Fatalf("invalid success body %q: err=%v calls=%d", body, err, calls.Load())
				}
			})
		}
	}
}
