package ozon

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const autoAddTestDate = "2026-10-13T03:00:00.123456789+03:00"

func readPromotionAutoAddFixture(t *testing.T, name string) []byte {
	t.Helper()

	data, err := os.ReadFile("testdata/promotions-v2/" + name)
	if err != nil {
		t.Fatalf("read fixture %q: %v", name, err)
	}
	return data
}

func newPromotionAutoAddServer(
	t *testing.T,
	wantPath string,
	status int,
	body []byte,
	inspect func(*testing.T, *http.Request, []byte),
) (*httptest.Server, *atomic.Int32) {
	t.Helper()

	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != http.MethodPost {
			t.Errorf("request method = %q, want POST", r.Method)
		}
		if r.URL.Path != wantPath {
			t.Errorf("request path = %q, want %q", r.URL.Path, wantPath)
		}
		if got := r.Header.Get("Client-Id"); got != "client-42" {
			t.Errorf("Client-Id = %q, want client-42", got)
		}
		if got := r.Header.Get("Api-Key"); got != "key-42" {
			t.Errorf("Api-Key = %q, want key-42", got)
		}
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
		}
		if inspect != nil {
			inspect(t, r, requestBody)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if len(body) > 0 {
			if _, err := w.Write(body); err != nil {
				t.Errorf("write response body: %v", err)
			}
		}
	}))
	return server, &calls
}

func promotionAutoAddClient(serverURL string) *Client {
	return NewClient(WithURI(serverURL), WithClientId("client-42"), WithAPIKey("key-42"))
}

func decodePromotionAutoAddRequest(t *testing.T, body []byte) map[string]interface{} {
	t.Helper()

	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.UseNumber()
	var request map[string]interface{}
	if err := decoder.Decode(&request); err != nil {
		t.Fatalf("decode request JSON: %v", err)
	}
	return request
}

func parsePromotionAutoAddDate(t *testing.T) time.Time {
	t.Helper()

	date, err := time.Parse(time.RFC3339Nano, autoAddTestDate)
	if err != nil {
		t.Fatal(err)
	}
	return date
}

func requirePromotionProtocolError(t *testing.T, err error) {
	t.Helper()

	var protocolErr *PromotionProtocolError
	if !errors.As(err, &protocolErr) {
		t.Fatalf("error = %v, want *PromotionProtocolError", err)
	}
}

func requirePromotionInputError(t *testing.T, err error) {
	t.Helper()

	var inputErr *PromotionInputError
	if !errors.As(err, &inputErr) {
		t.Fatalf("error = %v, want *PromotionInputError", err)
	}
}

func TestPromotionAutoAddV2ListsWire(t *testing.T) {
	t.Parallel()

	t.Run("products", func(t *testing.T) {
		t.Parallel()
		server, calls := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/list", http.StatusOK, readPromotionAutoAddFixture(t, "auto-add-products.json"), func(t *testing.T, _ *http.Request, body []byte) {
			request := decodePromotionAutoAddRequest(t, body)
			if request["action_id"] != json.Number("41") || request["auto_add_date"] != autoAddTestDate || request["offset"] != json.Number("5") || request["limit"] != json.Number("100") {
				t.Errorf("list request = %#v", request)
			}
			if _, ok := request["last_id"]; ok {
				t.Error("offset-paginated auto-add request included last_id")
			}
		})
		defer server.Close()

		params := &ListAutoAddProductsV2Params{ActionID: 41, AutoAddDate: parsePromotionAutoAddDate(t), Offset: 5, Limit: 100}
		wantParams := *params
		response, err := promotionAutoAddClient(server.URL).Promotions().ListAutoAddProductsV2(context.Background(), params)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(params, &wantParams) {
			t.Fatalf("params mutated: got %+v, want %+v", params, wantParams)
		}
		if calls.Load() != 1 || response.StatusCode != http.StatusOK || response.Total != 1 || len(response.Products) != 1 {
			t.Fatalf("calls=%d response=%+v", calls.Load(), response)
		}
		product := response.Products[0]
		if uint64(product.ProductID) != 9007199254740993 || product.SKU != 1234567 || product.OfferID != "offer-auto-add-1" {
			t.Fatalf("product identity = %+v", product)
		}
		if product.Price == nil || product.Price.Amount != "12345678901234567890.123456789" || product.BasePrice == nil || product.BasePrice.Amount != "13000000000000000000.000000001" {
			t.Fatalf("product money lost precision: %+v", product)
		}
		if product.AddMode == nil || !*product.AddMode || product.HasExpiredMinSellerPrice == nil || *product.HasExpiredMinSellerPrice || product.WillBeQuarantined == nil || !*product.WillBeQuarantined {
			t.Fatalf("product optional flags = %+v", product)
		}
		if product.WebsitePrices == nil || product.WebsitePrices.PricesBySchema["fbs"].BlackPrice == nil || product.WebsitePrices.PricesBySchema["fbs"].BlackPrice.Amount != "12000000000000000000.00" {
			t.Fatalf("website prices = %+v", product.WebsitePrices)
		}
	})

	t.Run("candidates", func(t *testing.T) {
		t.Parallel()
		server, _ := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/candidates", http.StatusOK, readPromotionAutoAddFixture(t, "auto-add-candidates.json"), func(t *testing.T, _ *http.Request, body []byte) {
			request := decodePromotionAutoAddRequest(t, body)
			if request["action_id"] != json.Number("41") || request["auto_add_date"] != autoAddTestDate || request["offset"] != json.Number("0") || request["limit"] != json.Number("1") {
				t.Errorf("candidate request = %#v", request)
			}
			if _, ok := request["last_id"]; ok {
				t.Error("offset-paginated auto-add request included last_id")
			}
		})
		defer server.Close()

		response, err := promotionAutoAddClient(server.URL).Promotions().ListAutoAddCandidatesV2(context.Background(), &ListAutoAddCandidatesV2Params{
			ActionID: 41, AutoAddDate: parsePromotionAutoAddDate(t), Limit: 1, Offset: 0,
		})
		if err != nil {
			t.Fatal(err)
		}
		if response.Total != 1 || len(response.Products) != 1 {
			t.Fatalf("response = %+v", response)
		}
		candidate := response.Products[0]
		if uint64(candidate.ID) != ^uint64(0) || candidate.SKU != 7654321 || candidate.OfferID != "offer-auto-candidate-1" {
			t.Fatalf("candidate identity = %+v", candidate)
		}
		if candidate.IsManuallyAdded == nil || *candidate.IsManuallyAdded {
			t.Fatalf("candidate mode fields = %+v", candidate)
		}
		if candidate.Price == nil || candidate.Price.Amount != "99999999999999999999.123456789" || candidate.BasePrice != nil || candidate.ActionPriceToAutoAdd != nil || candidate.WebsitePrices != nil {
			t.Fatalf("candidate optional money = %+v", candidate)
		}
	})
}

func TestPromotionAutoAddV2ListBoundsAndEnvelope(t *testing.T) {
	t.Parallel()

	t.Run("rejects invalid params before HTTP", func(t *testing.T) {
		t.Parallel()
		server, calls := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/list", http.StatusOK, readPromotionAutoAddFixture(t, "auto-add-empty-list.json"), nil)
		defer server.Close()
		client := promotionAutoAddClient(server.URL).Promotions()
		valid := ListAutoAddProductsV2Params{ActionID: 41, AutoAddDate: parsePromotionAutoAddDate(t), Limit: 1}
		cases := []struct {
			name   string
			params *ListAutoAddProductsV2Params
		}{
			{name: "nil params"},
			{name: "zero action", params: &ListAutoAddProductsV2Params{AutoAddDate: valid.AutoAddDate, Limit: 1}},
			{name: "zero date", params: &ListAutoAddProductsV2Params{ActionID: 41, Limit: 1}},
			{name: "zero limit", params: &ListAutoAddProductsV2Params{ActionID: 41, AutoAddDate: valid.AutoAddDate}},
			{name: "limit above maximum", params: &ListAutoAddProductsV2Params{ActionID: 41, AutoAddDate: valid.AutoAddDate, Limit: 101}},
		}
		for _, test := range cases {
			t.Run(test.name, func(t *testing.T) {
				_, err := client.ListAutoAddProductsV2(context.Background(), test.params)
				requirePromotionInputError(t, err)
			})
		}
		if calls.Load() != 0 {
			t.Fatalf("invalid inputs sent %d HTTP requests", calls.Load())
		}
	})

	t.Run("empty arrays are valid and malformed envelopes fail", func(t *testing.T) {
		t.Parallel()
		validBody := readPromotionAutoAddFixture(t, "auto-add-empty-list.json")
		server, _ := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/candidates", http.StatusOK, validBody, nil)
		response, err := promotionAutoAddClient(server.URL).Promotions().ListAutoAddCandidatesV2(context.Background(), &ListAutoAddCandidatesV2Params{
			ActionID: 41, AutoAddDate: parsePromotionAutoAddDate(t), Limit: 1,
		})
		server.Close()
		if err != nil || response == nil || response.Products == nil || len(response.Products) != 0 || response.Total != 0 {
			t.Fatalf("empty list response=%+v err=%v", response, err)
		}

		malformed := []string{
			"",
			"null",
			"{}",
			`{"result":{"products":[],"total":0}}`,
			`{"products":null,"total":0}`,
			`{"products":[],"total":-1}`,
			`{"products":{},"total":0}`,
			`{"products":[{"product_id":0}],"total":1}`,
			`{"products":[{"product_id":1,"price":{"amount":-1,"currency":"RUB"}}],"total":1}`,
		}
		for _, body := range malformed {
			t.Run(body, func(t *testing.T) {
				server, calls := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/list", http.StatusOK, []byte(body), nil)
				defer server.Close()
				_, err := promotionAutoAddClient(server.URL).Promotions().ListAutoAddProductsV2(context.Background(), &ListAutoAddProductsV2Params{
					ActionID: 41, AutoAddDate: parsePromotionAutoAddDate(t), Limit: 1,
				})
				requirePromotionProtocolError(t, err)
				if calls.Load() != 1 {
					t.Fatalf("calls = %d, want one", calls.Load())
				}
			})
		}

		candidateServer, _ := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/candidates", http.StatusOK, []byte(`{"products":[{"id":0}],"total":1}`), nil)
		defer candidateServer.Close()
		_, err = promotionAutoAddClient(candidateServer.URL).Promotions().ListAutoAddCandidatesV2(context.Background(), &ListAutoAddCandidatesV2Params{
			ActionID: 41, AutoAddDate: parsePromotionAutoAddDate(t), Limit: 1,
		})
		requirePromotionProtocolError(t, err)
	})
}

func TestPromotionAutoAddV2MutationsWire(t *testing.T) {
	t.Parallel()

	t.Run("update keeps exact IDs, money and partial result categories", func(t *testing.T) {
		t.Parallel()
		server, calls := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/update", http.StatusOK, readPromotionAutoAddFixture(t, "auto-add-update-mixed.json"), func(t *testing.T, _ *http.Request, body []byte) {
			request := decodePromotionAutoAddRequest(t, body)
			if request["action_id"] != json.Number("41") || request["auto_add_date"] != autoAddTestDate {
				t.Errorf("update envelope = %#v", request)
			}
			if _, ok := request["to_update"]; ok {
				t.Error("request used obsolete to_update field")
			}
			products, ok := request["products"].([]interface{})
			if !ok || len(products) != 3 {
				t.Fatalf("products = %#v", request["products"])
			}
			first := products[0].(map[string]interface{})
			if first["id"] != json.Number("9007199254740993") {
				t.Errorf("first update product = %#v", first)
			}
			if _, ok := first["product_id"]; ok {
				t.Error("update product used product_id instead of id")
			}
			if _, ok := first["stock"]; ok {
				t.Error("nil stock was encoded instead of omitted")
			}
			price := first["action_price"].(map[string]interface{})
			if price["amount"] != "12345678901234567890.123456789" || price["currency"] != "RUB" {
				t.Errorf("action price = %#v", price)
			}
			second := products[1].(map[string]interface{})
			if second["id"] != json.Number("18") || second["stock"] != json.Number("0") {
				t.Errorf("explicit zero stock was not preserved: %#v", second)
			}
			third := products[2].(map[string]interface{})
			if third["id"] != json.Number("19") || third["stock"] != json.Number("5") {
				t.Errorf("nonzero stock was not preserved: %#v", third)
			}
		})
		defer server.Close()

		zero := uint64(0)
		five := uint64(5)
		params := &UpdateAutoAddProductsV2Params{
			ActionID: 41, AutoAddDate: parsePromotionAutoAddDate(t),
			Products: []UpdateAutoAddProductV2{
				{ID: PromotionProductID(9007199254740993), ActionPrice: PromotionMoney{Amount: "12345678901234567890.123456789", Currency: "RUB"}},
				{ID: PromotionProductID(18), ActionPrice: PromotionMoney{Amount: "0", Currency: "RUB"}, Stock: &zero},
				{ID: PromotionProductID(19), ActionPrice: PromotionMoney{Amount: "12.34", Currency: "RUB"}, Stock: &five},
			},
		}
		wantParams := &UpdateAutoAddProductsV2Params{
			ActionID: 41, AutoAddDate: parsePromotionAutoAddDate(t),
			Products: []UpdateAutoAddProductV2{
				{ID: PromotionProductID(9007199254740993), ActionPrice: PromotionMoney{Amount: "12345678901234567890.123456789", Currency: "RUB"}},
				{ID: PromotionProductID(18), ActionPrice: PromotionMoney{Amount: "0", Currency: "RUB"}, Stock: &zero},
				{ID: PromotionProductID(19), ActionPrice: PromotionMoney{Amount: "12.34", Currency: "RUB"}, Stock: &five},
			},
		}
		response, err := promotionAutoAddClient(server.URL).Promotions().UpdateAutoAddProductsV2(context.Background(), params)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(params, wantParams) {
			t.Fatalf("params mutated: got %+v, want %+v", params, wantParams)
		}
		if calls.Load() != 1 || len(response.ProductIDs) != 2 || uint64(response.ProductIDs[0]) != 9007199254740993 || response.DeactivatedIDs[0] != 19 {
			t.Fatalf("update response = %+v", response)
		}
		if len(response.Rejected) != 1 || response.Rejected[0].Reason != "price rejected" || len(response.Warnings) != 1 || response.Warnings[0].Reason != "check schedule" {
			t.Fatalf("partial results = %+v", response)
		}
		if len(response.BelowMinPrice) != 1 || response.BelowMinPrice[0].Value.String() != "12.34567890123456789" || len(response.ExtremelyLowPrice) != 1 || response.ExtremelyLowPrice[0].Value.String() != "0.0000000000000000123" || len(response.FailedPrice) != 1 || response.FailedPrice[0].Value.String() != "99999999999999999999.999999999" {
			t.Fatalf("price diagnostics lost exact JSON numbers: %+v", response)
		}
	})

	t.Run("delete sends decimal string IDs", func(t *testing.T) {
		t.Parallel()
		server, _ := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/delete", http.StatusOK, readPromotionAutoAddFixture(t, "auto-add-delete.json"), func(t *testing.T, _ *http.Request, body []byte) {
			request := decodePromotionAutoAddRequest(t, body)
			if request["action_id"] != json.Number("41") || request["auto_add_date"] != autoAddTestDate {
				t.Errorf("delete envelope = %#v", request)
			}
			if _, ok := request["products"]; ok {
				t.Error("delete request included update products")
			}
			ids, ok := request["product_ids"].([]interface{})
			if !ok || !reflect.DeepEqual(ids, []interface{}{"9007199254740993", "18446744073709551615"}) {
				t.Errorf("product_ids = %#v, want decimal strings", request["product_ids"])
			}
		})
		defer server.Close()

		params := &DeleteAutoAddProductsV2Params{
			ActionID: 41, AutoAddDate: parsePromotionAutoAddDate(t),
			ProductIDs: []PromotionProductID{PromotionProductID(9007199254740993), PromotionProductID(^uint64(0))},
		}
		wantParams := &DeleteAutoAddProductsV2Params{
			ActionID: 41, AutoAddDate: parsePromotionAutoAddDate(t),
			ProductIDs: []PromotionProductID{PromotionProductID(9007199254740993), PromotionProductID(^uint64(0))},
		}
		response, err := promotionAutoAddClient(server.URL).Promotions().DeleteAutoAddProductsV2(context.Background(), params)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(params, wantParams) || len(response.ProductIDs) != 2 || uint64(response.ProductIDs[0]) != 9007199254740993 || uint64(response.ProductIDs[1]) != ^uint64(0) {
			t.Fatalf("params=%+v response=%+v", params, response)
		}
	})
}

func TestPromotionAutoAddV2MutationValidationAndDiagnostics(t *testing.T) {
	t.Parallel()

	server, calls := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/update", http.StatusOK, []byte(`{"product_ids":[]}`), nil)
	defer server.Close()
	client := promotionAutoAddClient(server.URL).Promotions()
	date := parsePromotionAutoAddDate(t)
	validPrice := PromotionMoney{Amount: "12.345", Currency: "RUB"}
	tooManyUpdates := make([]UpdateAutoAddProductV2, 1001)
	for index := range tooManyUpdates {
		tooManyUpdates[index] = UpdateAutoAddProductV2{ID: PromotionProductID(index + 1), ActionPrice: validPrice}
	}

	invalidUpdates := []struct {
		name   string
		params *UpdateAutoAddProductsV2Params
	}{
		{name: "nil params"},
		{name: "zero action", params: &UpdateAutoAddProductsV2Params{AutoAddDate: date, Products: []UpdateAutoAddProductV2{{ID: 1, ActionPrice: validPrice}}}},
		{name: "zero date", params: &UpdateAutoAddProductsV2Params{ActionID: 41, Products: []UpdateAutoAddProductV2{{ID: 1, ActionPrice: validPrice}}}},
		{name: "empty products", params: &UpdateAutoAddProductsV2Params{ActionID: 41, AutoAddDate: date}},
		{name: "zero id", params: &UpdateAutoAddProductsV2Params{ActionID: 41, AutoAddDate: date, Products: []UpdateAutoAddProductV2{{ID: 0, ActionPrice: validPrice}}}},
		{name: "duplicate id", params: &UpdateAutoAddProductsV2Params{ActionID: 41, AutoAddDate: date, Products: []UpdateAutoAddProductV2{{ID: 1, ActionPrice: validPrice}, {ID: 1, ActionPrice: validPrice}}}},
		{name: "more than one thousand products", params: &UpdateAutoAddProductsV2Params{ActionID: 41, AutoAddDate: date, Products: tooManyUpdates}},
		{name: "invalid money", params: &UpdateAutoAddProductsV2Params{ActionID: 41, AutoAddDate: date, Products: []UpdateAutoAddProductV2{{ID: 1, ActionPrice: PromotionMoney{Amount: "1e3", Currency: "RUB"}}}}},
	}
	for _, test := range invalidUpdates {
		t.Run(test.name, func(t *testing.T) {
			_, err := client.UpdateAutoAddProductsV2(context.Background(), test.params)
			requirePromotionInputError(t, err)
		})
	}
	if calls.Load() != 0 {
		t.Fatalf("invalid update inputs sent %d HTTP requests", calls.Load())
	}

	invalidDiagnostics := []string{
		`{"below_min_price":{"key":1,"value":12.3}}`,
		`{"below_min_price":[{"key":1,"value":"12.3"}]}`,
		`{"below_min_price":[{"key":1,"value":""}]}`,
		`{"below_min_price":[{"key":1,"value":null}]}`,
		`{"below_min_price":[{"key":1,"value":{}}]}`,
	}
	for _, body := range invalidDiagnostics {
		t.Run("reject diagnostics "+body, func(t *testing.T) {
			responseServer, _ := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/update", http.StatusOK, []byte(body), nil)
			defer responseServer.Close()
			_, err := promotionAutoAddClient(responseServer.URL).Promotions().UpdateAutoAddProductsV2(context.Background(), &UpdateAutoAddProductsV2Params{
				ActionID: 41, AutoAddDate: date, Products: []UpdateAutoAddProductV2{{ID: 1, ActionPrice: validPrice}},
			})
			requirePromotionProtocolError(t, err)
		})
	}

	diagnosticServer, _ := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/update", http.StatusOK, readPromotionAutoAddFixture(t, "auto-add-price-validation.json"), nil)
	defer diagnosticServer.Close()
	diagnosticResponse, err := promotionAutoAddClient(diagnosticServer.URL).Promotions().UpdateAutoAddProductsV2(context.Background(), &UpdateAutoAddProductsV2Params{
		ActionID: 41, AutoAddDate: date, Products: []UpdateAutoAddProductV2{{ID: 1, ActionPrice: validPrice}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnosticResponse.BelowMinPrice) != 1 || diagnosticResponse.BelowMinPrice[0].Value.String() != "12.34567890123456789" || len(diagnosticResponse.ExtremelyLowPrice) != 1 || diagnosticResponse.ExtremelyLowPrice[0].Value.String() != "0.0000000000000000123" || len(diagnosticResponse.FailedPrice) != 1 || diagnosticResponse.FailedPrice[0].Value.String() != "99999999999999999999.999999999" {
		t.Fatalf("price validation fixture = %+v", diagnosticResponse)
	}

	unknownResultServer, _ := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/update", http.StatusOK, []byte(`{"future_field":[]}`), nil)
	defer unknownResultServer.Close()
	_, err = promotionAutoAddClient(unknownResultServer.URL).Promotions().UpdateAutoAddProductsV2(context.Background(), &UpdateAutoAddProductsV2Params{
		ActionID: 41, AutoAddDate: date, Products: []UpdateAutoAddProductV2{{ID: 1, ActionPrice: validPrice}},
	})
	requirePromotionProtocolError(t, err)

	zeroIDResponseServer, _ := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/update", http.StatusOK, []byte(`{"product_ids":[0]}`), nil)
	defer zeroIDResponseServer.Close()
	_, err = promotionAutoAddClient(zeroIDResponseServer.URL).Promotions().UpdateAutoAddProductsV2(context.Background(), &UpdateAutoAddProductsV2Params{
		ActionID: 41, AutoAddDate: date, Products: []UpdateAutoAddProductV2{{ID: 1, ActionPrice: validPrice}},
	})
	requirePromotionProtocolError(t, err)
}

func TestPromotionAutoAddV2MutationLimitAndCallGuards(t *testing.T) {
	t.Parallel()

	t.Run("accepts one thousand updates and deletions", func(t *testing.T) {
		t.Parallel()
		date := parsePromotionAutoAddDate(t)
		updates := make([]UpdateAutoAddProductV2, 1000)
		ids := make([]PromotionProductID, 1000)
		for index := range updates {
			id := PromotionProductID(index + 1)
			updates[index] = UpdateAutoAddProductV2{ID: id, ActionPrice: PromotionMoney{Amount: "1", Currency: "RUB"}}
			ids[index] = id
		}

		updateServer, _ := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/update", http.StatusOK, []byte(`{"product_ids":[]}`), func(t *testing.T, _ *http.Request, body []byte) {
			request := decodePromotionAutoAddRequest(t, body)
			if got := len(request["products"].([]interface{})); got != 1000 {
				t.Errorf("update count = %d, want 1000", got)
			}
		})
		_, err := promotionAutoAddClient(updateServer.URL).Promotions().UpdateAutoAddProductsV2(context.Background(), &UpdateAutoAddProductsV2Params{ActionID: 41, AutoAddDate: date, Products: updates})
		updateServer.Close()
		if err != nil {
			t.Fatal(err)
		}

		deleteServer, _ := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/delete", http.StatusOK, []byte(`{"product_ids":[]}`), func(t *testing.T, _ *http.Request, body []byte) {
			request := decodePromotionAutoAddRequest(t, body)
			if got := len(request["product_ids"].([]interface{})); got != 1000 {
				t.Errorf("delete count = %d, want 1000", got)
			}
		})
		_, err = promotionAutoAddClient(deleteServer.URL).Promotions().DeleteAutoAddProductsV2(context.Background(), &DeleteAutoAddProductsV2Params{ActionID: 41, AutoAddDate: date, ProductIDs: ids})
		deleteServer.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("rejects more than one thousand before HTTP", func(t *testing.T) {
		t.Parallel()
		server, calls := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/delete", http.StatusOK, []byte(`{"product_ids":[]}`), nil)
		defer server.Close()
		ids := make([]PromotionProductID, 1001)
		for index := range ids {
			ids[index] = PromotionProductID(index + 1)
		}
		_, err := promotionAutoAddClient(server.URL).Promotions().DeleteAutoAddProductsV2(context.Background(), &DeleteAutoAddProductsV2Params{
			ActionID: 41, AutoAddDate: parsePromotionAutoAddDate(t), ProductIDs: ids,
		})
		requirePromotionInputError(t, err)
		if calls.Load() != 0 {
			t.Fatalf("oversized request sent %d HTTP requests", calls.Load())
		}
	})

	t.Run("nil receiver and nil context fail safely", func(t *testing.T) {
		t.Parallel()
		date := parsePromotionAutoAddDate(t)
		listProducts := &ListAutoAddProductsV2Params{ActionID: 41, AutoAddDate: date, Limit: 1}
		listCandidates := &ListAutoAddCandidatesV2Params{ActionID: 41, AutoAddDate: date, Limit: 1}
		update := &UpdateAutoAddProductsV2Params{ActionID: 41, AutoAddDate: date, Products: []UpdateAutoAddProductV2{{ID: 1, ActionPrice: PromotionMoney{Amount: "1", Currency: "RUB"}}}}
		delete := &DeleteAutoAddProductsV2Params{ActionID: 41, AutoAddDate: date, ProductIDs: []PromotionProductID{1}}
		_, err := (Promotions{}).ListAutoAddProductsV2(context.Background(), listProducts)
		requirePromotionInputError(t, err)
		_, err = (Promotions{}).ListAutoAddCandidatesV2(context.Background(), listCandidates)
		requirePromotionInputError(t, err)
		_, err = (Promotions{}).UpdateAutoAddProductsV2(context.Background(), update)
		requirePromotionInputError(t, err)
		_, err = (Promotions{}).DeleteAutoAddProductsV2(context.Background(), delete)
		requirePromotionInputError(t, err)
		_, err = promotionAutoAddClient("http://127.0.0.1:1").Promotions().ListAutoAddProductsV2(nil, listProducts)
		requirePromotionInputError(t, err)
		_, err = promotionAutoAddClient("http://127.0.0.1:1").Promotions().ListAutoAddCandidatesV2(nil, listCandidates)
		requirePromotionInputError(t, err)
		_, err = promotionAutoAddClient("http://127.0.0.1:1").Promotions().UpdateAutoAddProductsV2(nil, update)
		requirePromotionInputError(t, err)
		_, err = promotionAutoAddClient("http://127.0.0.1:1").Promotions().DeleteAutoAddProductsV2(nil, delete)
		requirePromotionInputError(t, err)
	})
}

func TestPromotionAutoAddV2DeleteValidationAndHTTPError(t *testing.T) {
	t.Parallel()

	t.Run("rejects invalid params before HTTP", func(t *testing.T) {
		t.Parallel()
		server, calls := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/delete", http.StatusOK, readPromotionAutoAddFixture(t, "auto-add-delete.json"), nil)
		defer server.Close()
		client := promotionAutoAddClient(server.URL).Promotions()
		date := parsePromotionAutoAddDate(t)
		cases := []struct {
			name   string
			params *DeleteAutoAddProductsV2Params
		}{
			{name: "nil params"},
			{name: "zero action", params: &DeleteAutoAddProductsV2Params{AutoAddDate: date, ProductIDs: []PromotionProductID{1}}},
			{name: "zero date", params: &DeleteAutoAddProductsV2Params{ActionID: 41, ProductIDs: []PromotionProductID{1}}},
			{name: "empty IDs", params: &DeleteAutoAddProductsV2Params{ActionID: 41, AutoAddDate: date}},
			{name: "zero ID", params: &DeleteAutoAddProductsV2Params{ActionID: 41, AutoAddDate: date, ProductIDs: []PromotionProductID{0}}},
			{name: "duplicate IDs", params: &DeleteAutoAddProductsV2Params{ActionID: 41, AutoAddDate: date, ProductIDs: []PromotionProductID{1, 1}}},
		}
		for _, test := range cases {
			t.Run(test.name, func(t *testing.T) {
				_, err := client.DeleteAutoAddProductsV2(context.Background(), test.params)
				requirePromotionInputError(t, err)
			})
		}
		if calls.Load() != 0 {
			t.Fatalf("invalid delete inputs sent %d HTTP requests", calls.Load())
		}
	})

	t.Run("preserves Ozon HTTP error and does not retry", func(t *testing.T) {
		t.Parallel()
		server, calls := newPromotionAutoAddServer(t, "/v2/actions/auto-add/products/delete", http.StatusConflict, []byte(`{"code":42,"message":"conflict","details":[{"typeUrl":"type","value":"bad"}]}`), nil)
		defer server.Close()
		response, err := promotionAutoAddClient(server.URL).Promotions().DeleteAutoAddProductsV2(context.Background(), &DeleteAutoAddProductsV2Params{
			ActionID: 41, AutoAddDate: parsePromotionAutoAddDate(t), ProductIDs: []PromotionProductID{1},
		})
		if err != nil {
			t.Fatal(err)
		}
		if calls.Load() != 1 || response.StatusCode != http.StatusConflict || response.Code != 42 || response.Message != "conflict" || len(response.Details) != 1 {
			t.Fatalf("calls=%d response=%+v", calls.Load(), response)
		}
	})
}
