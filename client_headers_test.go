package core

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestResponseHeaders(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusTooManyRequests} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			calls := 0
			client := NewMockClient(func(w http.ResponseWriter, r *http.Request) {
				calls++
				w.Header().Set("Ratelimit-Remaining", "0")
				w.Header().Set("Retry-After", "2")
				w.Header().Set("Item-Retry-After", "60")
				w.Header().Set("Item-Rate-Limit-Remaining", "0")
				w.WriteHeader(status)
				w.Write([]byte(`{"code":8,"message":"quota exceeded"}`))
			})
			response, err := client.Request(context.Background(), http.MethodPost, "/quota", nil, &struct{}{}, nil)
			if err != nil {
				t.Fatal(err)
			}
			var copied CommonResponse
			response.CopyCommonResponse(&copied)
			if copied.StatusCode != status || calls != 1 {
				t.Fatalf("status=%d calls=%d", copied.StatusCode, calls)
			}
			if copied.Headers.Get("Retry-After") != "2" || copied.Headers.Get("Item-Retry-After") != "60" || copied.Headers.Get("Ratelimit-Remaining") != "0" || copied.Headers.Get("Item-Rate-Limit-Remaining") != "0" {
				t.Fatalf("headers=%v", copied.Headers)
			}
			if status == http.StatusTooManyRequests && (copied.Code != 8 || copied.Message != "quota exceeded") {
				t.Fatalf("error=%+v", copied)
			}
			copied.Headers.Set("Retry-After", "999")
			if response.Headers.Get("Retry-After") != "2" {
				t.Fatal("copy aliases source headers")
			}
			data, err := json.Marshal(copied)
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(data, &fields); err != nil {
				t.Fatal(err)
			}
			if _, ok := fields["Headers"]; ok {
				t.Fatal("transport headers leaked into JSON")
			}
		})
	}
}
