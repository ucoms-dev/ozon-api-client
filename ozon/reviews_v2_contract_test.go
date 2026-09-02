package ozon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestReviewV2MutationContracts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		body string
		call func(*Reviews) error
	}{
		{
			name: "delete comment",
			path: "/v2/review/comment/delete",
			body: `{"comment_id":"comment-1","sku":123}`,
			call: func(client *Reviews) error {
				_, err := client.DeleteCommentV2(context.Background(), &DeleteCommentV2Params{CommentID: "comment-1", SKU: 123})
				return err
			},
		},
		{
			name: "change status",
			path: "/v2/review/change-status",
			body: `{"review_ids":["review-1"],"status":"PROCESSED"}`,
			call: func(client *Reviews) error {
				_, err := client.ChangeStatusV2(context.Background(), &ChangeReviewStatusV2Params{ReviewIDs: []string{"review-1"}, Status: "PROCESSED"})
				return err
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(requestContractHandler(t, http.MethodPost, test.path, test.body, ""))
			defer server.Close()
			if err := test.call(NewClient(WithURI(server.URL)).Reviews()); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestReviewV2ReadContracts(t *testing.T) {
	t.Parallel()

	t.Run("count", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(requestContractHandler(t, http.MethodPost, "/v2/review/count", "", `{"new":2,"processed":3,"total":6,"viewed":1}`))
		defer server.Close()
		response, err := NewClient(WithURI(server.URL)).Reviews().CountV2(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if response.New != 2 || response.Viewed != 1 || response.Processed != 3 || response.Total != 6 {
			t.Fatalf("response = %+v", response)
		}
	})

	t.Run("info", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(requestContractHandler(t, http.MethodPost, "/v2/review/info", `{"review_id":"review-1"}`, `{"id":"review-1","sku":123,"published_at":"2026-09-01T10:00:00Z","status":"NEW","order_status":"DELIVERED","photos":[{"height":10,"url":"photo","width":20}],"videos":[{"height":30,"preview_url":"preview","short_video_preview_url":"short","url":"video","width":40}]}`))
		defer server.Close()
		response, err := NewClient(WithURI(server.URL)).Reviews().GetV2(context.Background(), &GetReviewV2Params{ReviewID: "review-1"})
		if err != nil {
			t.Fatal(err)
		}
		if response.ID != "review-1" || response.SKU != 123 || response.PublishedAt.IsZero() || len(response.Photos) != 1 || len(response.Videos) != 1 {
			t.Fatalf("response = %+v", response)
		}
	})

	t.Run("list", func(t *testing.T) {
		t.Parallel()
		publishedFrom := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
		server := httptest.NewServer(requestContractHandler(t, http.MethodPost, "/v2/review/list", `{"filters":{"order_status":"DELIVERED","published_from":"2026-09-01T00:00:00Z","skus":["123"],"status":"NEW"},"last_id":"page-1","limit":20,"sort_dir":"ASC"}`, `{"has_next":true,"last_id":"page-2","reviews":[{"id":"review-1","sku":123,"published_at":"2026-09-01T10:00:00Z","status":"NEW","order_status":"DELIVERED"}]}`))
		defer server.Close()
		response, err := NewClient(WithURI(server.URL)).Reviews().ListV2(context.Background(), &ListReviewsV2Params{
			Filters: ListReviewsV2Filters{OrderStatus: "DELIVERED", PublishedFrom: &publishedFrom, SKUs: []string{"123"}, Status: "NEW"},
			LastID:  "page-1",
			Limit:   20,
			SortDir: Ascending,
		})
		if err != nil {
			t.Fatal(err)
		}
		if !response.HasNext || response.LastID != "page-2" || len(response.Reviews) != 1 || response.Reviews[0].SKU != 123 {
			t.Fatalf("response = %+v", response)
		}
	})
}
