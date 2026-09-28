package ozon

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync/atomic"
	"testing"
)

func TestPromotionV2ListWire(t *testing.T) {
	t.Parallel()

	products := promotionCoreFixture(t, "products.json")
	candidates := promotionCoreFixture(t, "candidates.json")
	var calls atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.Header.Get("Client-Id") != "client-1" || r.Header.Get("Api-Key") != "key-1" {
			t.Errorf("auth headers = %q/%q", r.Header.Get("Client-Id"), r.Header.Get("Api-Key"))
		}
		body := promotionRequestBody(t, r)
		var responseBody string
		switch r.URL.Path {
		case "/v2/actions/products":
			if body != `{"action_id":42,"limit":100,"last_id":"cursor-1"}` {
				t.Errorf("products request body = %s", body)
			}
			responseBody = products
		case "/v2/actions/candidates":
			if body != `{"action_id":42,"limit":100,"last_id":"cursor-1"}` {
				t.Errorf("candidates request body = %s", body)
			}
			responseBody = candidates
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(responseBody))
	}))
	defer server.Close()
	api := NewClient(WithURI(server.URL), WithClientId("client-1"), WithAPIKey("key-1")).Promotions()

	participants, err := api.ProductsInPromotionV2(context.Background(), &ProductsInPromotionV2Params{ActionID: 42, Limit: 100, LastID: "cursor-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(participants.Products) != 1 || participants.Total != 1 || participants.LastID != "next-1" {
		t.Fatalf("participants = %+v", participants)
	}
	participant := participants.Products[0]
	if participant.ID != 9007199254740993 || participant.Price == nil || participant.Price.Amount != "10.000000000000000001" || participant.AddMode != "SELLER" || participant.Stock == nil || *participant.Stock != 3 {
		t.Fatalf("participant = %+v", participant)
	}
	if participant.WebsitePrices == nil || participant.WebsitePrices.PricesBySchema["FBS"].BlackPrice == nil || participant.IsQuarantined == nil || !*participant.IsQuarantined {
		t.Fatalf("participant optional prices/quarantine = %+v", participant)
	}

	available, err := api.ProductsAvailableForPromotionV2(context.Background(), &ProductsAvailableForPromotionV2Params{ActionID: 42, Limit: 100, LastID: "cursor-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(available.Products) != 1 || available.Total != 1 || available.LastID != "candidate-next" {
		t.Fatalf("candidates = %+v", available)
	}
	candidate := available.Products[0]
	if candidate.ID != PromotionProductID(^uint64(0)) || candidate.Price == nil || candidate.Price.Amount != "0.00" || candidate.PriceMinElastic == nil || candidate.PriceMinElastic.Amount != "9.99" {
		t.Fatalf("candidate = %+v", candidate)
	}
	if calls.Load() != 2 {
		t.Fatalf("HTTP calls = %d, want one per list method", calls.Load())
	}
}

func TestPromotionV2ListInvalidShape(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
	}{
		{name: "empty body", body: ""},
		{name: "null", body: `null`},
		{name: "missing fields", body: `{}`},
		{name: "nested result", body: `{"result":{"products":[],"total":0}}`},
		{name: "null products", body: `{"products":null,"total":0}`},
		{name: "negative total", body: `{"products":[],"total":-1}`},
		{name: "null cursor", body: `{"products":[],"total":0,"last_id":null}`},
		{name: "invalid money", body: `{"products":[{"id":1,"price":{"amount":"1e3","currency":"RUB"}}],"total":1}`},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(test.body))
			}))
			defer server.Close()
			_, err := NewClient(WithURI(server.URL)).Promotions().ProductsInPromotionV2(context.Background(), &ProductsInPromotionV2Params{ActionID: 1, Limit: 1})
			var protocolErr *PromotionProtocolError
			if !errors.As(err, &protocolErr) {
				t.Fatalf("error = %v, want PromotionProtocolError", err)
			}
		})
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"products":[],"total":0,"last_id":"","future_field":{"x":1}}`))
	}))
	defer server.Close()
	response, err := NewClient(WithURI(server.URL)).Promotions().ProductsInPromotionV2(context.Background(), &ProductsInPromotionV2Params{ActionID: 1, Limit: 1})
	if err != nil || response.Total != 0 || response.Products == nil {
		t.Fatalf("unknown response field: response=%+v error=%v", response, err)
	}
}

func TestPromotionV2ListBoundsAndOwnership(t *testing.T) {
	t.Parallel()

	var calls atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.URL.Path != "/v2/actions/products" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(promotionCoreFixture(t, "empty-products.json")))
	}))
	defer server.Close()
	api := NewClient(WithURI(server.URL)).Promotions()

	invalid := []*ProductsInPromotionV2Params{
		nil,
		{ActionID: 0, Limit: 1},
		{ActionID: 1, Limit: 0},
		{ActionID: 1, Limit: 101},
	}
	for _, params := range invalid {
		_, err := api.ProductsInPromotionV2(context.Background(), params)
		var inputErr *PromotionInputError
		if !errors.As(err, &inputErr) {
			t.Errorf("params=%+v error=%v, want PromotionInputError", params, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("invalid requests sent %d HTTP calls", calls.Load())
	}

	for _, limit := range []uint64{1, 100} {
		params := &ProductsInPromotionV2Params{ActionID: 1, Limit: limit, LastID: "cursor"}
		before := *params
		if _, err := api.ProductsInPromotionV2(context.Background(), params); err != nil {
			t.Fatalf("limit %d: %v", limit, err)
		}
		if !reflect.DeepEqual(*params, before) {
			t.Fatalf("params mutated: got %+v, want %+v", *params, before)
		}
	}
	if calls.Load() != 2 {
		t.Fatalf("HTTP calls = %d, want 2 (no automatic cursor traversal)", calls.Load())
	}

	_, err := api.ProductsInPromotionV2(nil, &ProductsInPromotionV2Params{ActionID: 1, Limit: 1})
	var nilContextErr *PromotionInputError
	if !errors.As(err, &nilContextErr) {
		t.Fatalf("nil context error = %v", err)
	}
	_, err = (Promotions{}).ProductsInPromotionV2(context.Background(), &ProductsInPromotionV2Params{ActionID: 1, Limit: 1})
	if err == nil || calls.Load() != 2 {
		t.Fatalf("zero-value Promotions error=%v, calls=%d", err, calls.Load())
	}
}

func TestPromotionV2MutationsWire(t *testing.T) {
	t.Parallel()

	updateResponse := promotionCoreFixture(t, "update-mixed.json")
	deactivateResponse := promotionCoreFixture(t, "deactivate-partial.json")
	var calls atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		var responseBody string
		switch r.URL.Path {
		case "/v1/actions/products/update":
			want := `{"action_id":42,"products":[{"product_id":9007199254740993,"action_price":{"amount":"12345678901234567890.123456789","currency":"RUB"}},{"product_id":18446744073709551615,"action_price":{"amount":"0.00","currency":"KZT"},"stock":0}]}`
			if body := promotionRequestBody(t, r); body != want {
				t.Errorf("update body = %s, want %s", body, want)
			}
			responseBody = updateResponse
		case "/v2/actions/products/deactivate":
			want := `{"action_id":42,"product_ids":["9007199254740993","18446744073709551615"]}`
			if body := promotionRequestBody(t, r); body != want {
				t.Errorf("deactivate body = %s, want %s", body, want)
			}
			responseBody = deactivateResponse
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(responseBody))
	}))
	defer server.Close()
	api := NewClient(WithURI(server.URL)).Promotions()

	zero := uint64(0)
	updateParams := &UpdatePromotionProductsParams{
		ActionID: 42,
		Products: []UpdatePromotionProduct{
			{ProductID: 9007199254740993, ActionPrice: PromotionMoney{Amount: "12345678901234567890.123456789", Currency: "RUB"}},
			{ProductID: PromotionProductID(^uint64(0)), ActionPrice: PromotionMoney{Amount: "0.00", Currency: "KZT"}, Stock: &zero},
		},
	}
	beforeUpdate := *updateParams
	beforeUpdate.Products = append([]UpdatePromotionProduct(nil), updateParams.Products...)
	updated, err := api.UpdateProducts(context.Background(), updateParams)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*updateParams, beforeUpdate) {
		t.Fatalf("update params mutated: got %+v, want %+v", *updateParams, beforeUpdate)
	}
	if !reflect.DeepEqual(updated.ActiveProductIDs, []PromotionProductID{1}) || !reflect.DeepEqual(updated.DeactivatedProductIDs, []PromotionProductID{2}) || len(updated.Rejected) != 1 || updated.Rejected[0].ProductID != 3 || len(updated.Warnings) != 1 || updated.Warnings[0].ProductID != 2 {
		t.Fatalf("update categories were not preserved: %+v", updated)
	}

	removeParams := &RemoveProductFromPromotionV2Params{ActionID: 42, ProductIDs: []PromotionProductID{9007199254740993, PromotionProductID(^uint64(0))}}
	beforeRemove := *removeParams
	beforeRemove.ProductIDs = append([]PromotionProductID(nil), removeParams.ProductIDs...)
	removed, err := api.RemoveProductV2(context.Background(), removeParams)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*removeParams, beforeRemove) {
		t.Fatalf("remove params mutated: got %+v, want %+v", *removeParams, beforeRemove)
	}
	if !reflect.DeepEqual(removed.ProductIDs, []PromotionProductID{1}) {
		t.Fatalf("partial deactivation result = %v", removed.ProductIDs)
	}
	if calls.Load() != 2 {
		t.Fatalf("HTTP calls = %d, want 2", calls.Load())
	}
}

func TestPromotionV2MutationBounds(t *testing.T) {
	t.Parallel()

	var calls atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		switch r.URL.Path {
		case "/v1/actions/products/update":
			_, _ = w.Write([]byte(`{"active_product_ids":[]}`))
		case "/v2/actions/products/deactivate":
			_, _ = w.Write([]byte(`{"product_ids":[]}`))
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	api := NewClient(WithURI(server.URL)).Promotions()
	price := PromotionMoney{Amount: "0", Currency: "RUB"}

	for _, count := range []int{1, 1000} {
		products := make([]UpdatePromotionProduct, count)
		for index := range products {
			products[index] = UpdatePromotionProduct{ProductID: PromotionProductID(index + 1), ActionPrice: price}
		}
		if _, err := api.UpdateProducts(context.Background(), &UpdatePromotionProductsParams{ActionID: 1, Products: products}); err != nil {
			t.Errorf("update count %d: %v", count, err)
		}
		ids := make([]PromotionProductID, count)
		for index := range ids {
			ids[index] = PromotionProductID(index + 1)
		}
		if _, err := api.RemoveProductV2(context.Background(), &RemoveProductFromPromotionV2Params{ActionID: 1, ProductIDs: ids}); err != nil {
			t.Errorf("deactivate count %d: %v", count, err)
		}
	}

	updateInvalid := []*UpdatePromotionProductsParams{
		nil,
		{ActionID: 0, Products: []UpdatePromotionProduct{{ProductID: 1, ActionPrice: price}}},
		{ActionID: 1},
		{ActionID: 1, Products: make([]UpdatePromotionProduct, 1001)},
		{ActionID: 1, Products: []UpdatePromotionProduct{{ProductID: 1, ActionPrice: price}, {ProductID: 1, ActionPrice: price}}},
		{ActionID: 1, Products: []UpdatePromotionProduct{{ProductID: 0, ActionPrice: price}}},
		{ActionID: 1, Products: []UpdatePromotionProduct{{ProductID: 1, ActionPrice: PromotionMoney{Amount: "-1", Currency: "RUB"}}}},
	}
	for _, params := range updateInvalid {
		_, err := api.UpdateProducts(context.Background(), params)
		var inputErr *PromotionInputError
		if !errors.As(err, &inputErr) {
			t.Errorf("update params=%+v error=%v, want PromotionInputError", params, err)
		}
	}

	removeInvalid := []*RemoveProductFromPromotionV2Params{
		nil,
		{ActionID: 0, ProductIDs: []PromotionProductID{1}},
		{ActionID: 1},
		{ActionID: 1, ProductIDs: make([]PromotionProductID, 1001)},
		{ActionID: 1, ProductIDs: []PromotionProductID{1, 1}},
		{ActionID: 1, ProductIDs: []PromotionProductID{0}},
	}
	for _, params := range removeInvalid {
		_, err := api.RemoveProductV2(context.Background(), params)
		var inputErr *PromotionInputError
		if !errors.As(err, &inputErr) {
			t.Errorf("remove params=%+v error=%v, want PromotionInputError", params, err)
		}
	}
	if calls.Load() != 4 {
		t.Fatalf("HTTP calls = %d, want only four boundary requests", calls.Load())
	}
}

func TestPromotionV2PartialAndUnknownResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		body         string
		wantProtocol bool
		wantActive   []PromotionProductID
		wantRejected int
	}{
		{name: "partial confirmation", body: `{"active_product_ids":[1]}`, wantActive: []PromotionProductID{1}},
		{name: "rejected is normal HTTP 200 result", body: `{"rejected":[{"product_id":2,"reason":"not eligible"}]}`, wantRejected: 1},
		{name: "recognized empty results", body: `{"active_product_ids":[],"deactivated_product_ids":[],"rejected":[],"warnings":[]}`, wantActive: []PromotionProductID{}},
		{name: "unknown nested result", body: `{"result":{"active_product_ids":[1]}}`, wantProtocol: true},
		{name: "empty response", body: "", wantProtocol: true},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var calls atomic.Int64
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				calls.Add(1)
				_, _ = w.Write([]byte(test.body))
			}))
			defer server.Close()
			response, err := NewClient(WithURI(server.URL)).Promotions().UpdateProducts(context.Background(), &UpdatePromotionProductsParams{
				ActionID: 1,
				Products: []UpdatePromotionProduct{
					{ProductID: 1, ActionPrice: PromotionMoney{Amount: "1", Currency: "RUB"}},
					{ProductID: 2, ActionPrice: PromotionMoney{Amount: "1", Currency: "RUB"}},
				},
			})
			if test.wantProtocol {
				var protocolErr *PromotionProtocolError
				if !errors.As(err, &protocolErr) {
					t.Fatalf("error = %v, want PromotionProtocolError", err)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(response.ActiveProductIDs, test.wantActive) || len(response.Rejected) != test.wantRejected {
					t.Fatalf("response = %+v", response)
				}
			}
			if calls.Load() != 1 {
				t.Fatalf("HTTP attempts = %d, want exactly one", calls.Load())
			}
		})
	}
}

func promotionCoreFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile("testdata/promotions-v2/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func promotionRequestBody(t *testing.T, r *http.Request) string {
	t.Helper()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read request body: %v", err)
	}
	return string(data)
}
