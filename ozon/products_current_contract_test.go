package ozon

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateProductImagesV2Contract(t *testing.T) {
	server := httptest.NewServer(requestContractHandler(t, http.MethodPost, "/v2/product/pictures/import", `{"items":[{"offer_id":"offer-1","primary_image":"main","color_image":"","images":[]}]}`, `{"task_id":9223372036854775806}`))
	defer server.Close()
	response, err := NewClient(WithURI(server.URL)).Products().UpdateProductImagesV2(context.Background(), &UpdateProductImagesV2Params{Items: []UpdateProductImagesV2Item{{OfferId: "offer-1", PrimaryImage: "main", Images: []string{}}}})
	if err != nil {
		t.Fatal(err)
	}
	if response.TaskId != 9223372036854775806 || response.StatusCode != http.StatusOK {
		t.Fatalf("response=%+v", response)
	}
}

func TestProductRangeOperationLimits(t *testing.T) {
	for _, limits := range []string{`[{"limit":30000,"limit_type":"RATE_LIMIT_PER_MINUTE"}]`, `{"limit":30000,"limit_type":"RATE_LIMIT_PER_MINUTE"}`} {
		t.Run(limits, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/v4/product/info/limit" {
					t.Errorf("request=%s %s", r.Method, r.URL.Path)
				}
				w.Header().Set("Ratelimit-Remaining", "48")
				w.Write([]byte(`{"daily_create":{"limit":-1,"usage":0,"reset_at":"2026-10-07T00:00:00Z"},"total":{"limit":1500000,"usage":1295913},"operation_limits":` + limits + `}`))
			}))
			defer server.Close()
			response, err := NewClient(WithURI(server.URL)).Products().GetProductRangeLimit(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if len(response.OperationLimits) != 1 || response.OperationLimits[0].Limit != 30000 || response.OperationLimits[0].LimitType != "RATE_LIMIT_PER_MINUTE" || response.DailyCreate.Limit != -1 || response.Total.Usage != 1295913 || response.DailyCreate.ResetAt.IsZero() || response.Headers.Get("Ratelimit-Remaining") != "48" {
				t.Fatalf("response=%+v", response)
			}
		})
	}
}

func TestProductOperationLimitsDecode(t *testing.T) {
	for _, body := range []string{`null`, `[]`, `[{"limit":0,"limit_type":"FUTURE_LIMIT"},{"limit":-1,"limit_type":"UNSPECIFIED"}]`} {
		var limits GetProductRangeLimitOperationLimits
		if err := json.Unmarshal([]byte(body), &limits); err != nil {
			t.Fatal(err)
		}
		if body == `null` && limits != nil {
			t.Fatalf("limits=%v", limits)
		}
		if body != `null` && body != `[]` && (len(limits) != 2 || limits[0].LimitType != "FUTURE_LIMIT" || limits[1].Limit != -1) {
			t.Fatalf("limits=%v", limits)
		}
	}
	for _, body := range []string{`"bad"`, `123`, `true`, `[{"limit":"invalid"}]`} {
		limits := GetProductRangeLimitOperationLimits{{Limit: 42}}
		if err := json.Unmarshal([]byte(body), &limits); err == nil {
			t.Fatalf("accepted %s", body)
		}
		if len(limits) != 1 || limits[0].Limit != 42 {
			t.Fatal("failed decode changed previous value")
		}
	}
	var response GetProductRangeLimitResponse
	if err := json.Unmarshal([]byte(`{}`), &response); err != nil || response.OperationLimits != nil {
		t.Fatalf("missing limits: %+v %v", response, err)
	}
}

func TestProductQuota429Metadata(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Item-Retry-After", "7")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"code":8,"message":"too many requests"}`))
	}))
	defer server.Close()
	response, err := NewClient(WithURI(server.URL)).Products().GetProductRangeLimit(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 429 || response.Code != 8 || response.Headers.Get("Item-Retry-After") != "7" || calls != 1 {
		t.Fatalf("response=%+v calls=%d", response, calls)
	}
}
