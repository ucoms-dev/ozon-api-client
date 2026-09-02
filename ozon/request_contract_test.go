package ozon

import (
	"io"
	"net/http"
	"testing"
)

func requestContractHandler(
	t *testing.T,
	wantMethod string,
	wantPath string,
	wantBody string,
	responseBody string,
) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		t.Helper()

		if r.Method != wantMethod {
			t.Errorf("request method = %q, want %q", r.Method, wantMethod)
		}
		if r.URL.Path != wantPath {
			t.Errorf("request path = %q, want %q", r.URL.Path, wantPath)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}
		if string(body) != wantBody {
			t.Errorf("request body = %s, want %s", body, wantBody)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(responseBody)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}
}
