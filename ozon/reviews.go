package ozon

import (
	"context"
	"net/http"
	"time"

	core "github.com/ucoms-dev/ozon-api-client"
)

type Reviews struct {
	client *core.Client
}

// DeleteCommentV2Params is the request body for /v2/review/comment/delete.
type DeleteCommentV2Params struct {
	CommentID string `json:"comment_id"`
	SKU       int64  `json:"sku"`
}

type DeleteCommentV2Response struct {
	core.CommonResponse
}

// DeleteCommentV2 deletes a review comment using the current v2 API.
func (c Reviews) DeleteCommentV2(ctx context.Context, params *DeleteCommentV2Params) (*DeleteCommentV2Response, error) {
	resp := &DeleteCommentV2Response{}
	response, err := c.client.Request(ctx, http.MethodPost, "/v2/review/comment/delete", params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)
	return resp, nil
}

// ChangeReviewStatusV2Params is the request body for /v2/review/change-status.
type ChangeReviewStatusV2Params struct {
	ReviewIDs []string `json:"review_ids"`
	Status    string   `json:"status"`
}

type ChangeReviewStatusV2Response struct {
	core.CommonResponse
}

// ChangeStatusV2 changes review statuses using the current v2 API.
func (c Reviews) ChangeStatusV2(ctx context.Context, params *ChangeReviewStatusV2Params) (*ChangeReviewStatusV2Response, error) {
	resp := &ChangeReviewStatusV2Response{}
	response, err := c.client.Request(ctx, http.MethodPost, "/v2/review/change-status", params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)
	return resp, nil
}

type CountReviewsV2Response struct {
	core.CommonResponse

	New       int32 `json:"new"`
	Processed int32 `json:"processed"`
	Total     int32 `json:"total"`
	Viewed    int32 `json:"viewed"`
}

// CountV2 returns review counts grouped by the current v2 statuses.
func (c Reviews) CountV2(ctx context.Context) (*CountReviewsV2Response, error) {
	resp := &CountReviewsV2Response{}
	response, err := c.client.Request(ctx, http.MethodPost, "/v2/review/count", nil, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)
	return resp, nil
}

type GetReviewV2Params struct {
	ReviewID string `json:"review_id"`
}

type ReviewV2Photo struct {
	Height int32  `json:"height"`
	URL    string `json:"url"`
	Width  int32  `json:"width"`
}

type ReviewV2Video struct {
	Height               int64  `json:"height"`
	PreviewURL           string `json:"preview_url"`
	ShortVideoPreviewURL string `json:"short_video_preview_url"`
	URL                  string `json:"url"`
	Width                int64  `json:"width"`
}

type GetReviewV2Response struct {
	core.CommonResponse

	CommentsAmount      int32           `json:"comments_amount"`
	DislikesAmount      int32           `json:"dislikes_amount"`
	ID                  string          `json:"id"`
	IsRatingParticipant bool            `json:"is_rating_participant"`
	LikesAmount         int32           `json:"likes_amount"`
	OrderStatus         string          `json:"order_status"`
	Photos              []ReviewV2Photo `json:"photos"`
	PhotosAmount        int32           `json:"photos_amount"`
	PublishedAt         time.Time       `json:"published_at"`
	Rating              int32           `json:"rating"`
	SKU                 int64           `json:"sku"`
	Status              string          `json:"status"`
	Text                string          `json:"text"`
	Videos              []ReviewV2Video `json:"videos"`
	VideosAmount        int32           `json:"videos_amount"`
}

// GetV2 returns a review using the current v2 API.
func (c Reviews) GetV2(ctx context.Context, params *GetReviewV2Params) (*GetReviewV2Response, error) {
	resp := &GetReviewV2Response{}
	response, err := c.client.Request(ctx, http.MethodPost, "/v2/review/info", params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)
	return resp, nil
}

type ListReviewsV2Filters struct {
	OrderStatus   string     `json:"order_status,omitempty"`
	PublishedFrom *time.Time `json:"published_from,omitempty"`
	PublishedTo   *time.Time `json:"published_to,omitempty"`
	SKUs          []string   `json:"skus,omitempty"`
	Status        string     `json:"status,omitempty"`
}

type ListReviewsV2Params struct {
	Filters ListReviewsV2Filters `json:"filters"`
	LastID  string               `json:"last_id,omitempty"`
	Limit   int32                `json:"limit"`
	SortDir Order                `json:"sort_dir,omitempty"`
}

type ReviewV2Summary struct {
	CommentsAmount      int32     `json:"comments_amount"`
	ID                  string    `json:"id"`
	IsRatingParticipant bool      `json:"is_rating_participant"`
	OrderStatus         string    `json:"order_status"`
	PhotosAmount        int32     `json:"photos_amount"`
	PublishedAt         time.Time `json:"published_at"`
	Rating              int32     `json:"rating"`
	SKU                 int64     `json:"sku"`
	Status              string    `json:"status"`
	Text                string    `json:"text"`
	VideosAmount        int32     `json:"videos_amount"`
}

type ListReviewsV2Response struct {
	core.CommonResponse

	HasNext bool              `json:"has_next"`
	LastID  string            `json:"last_id"`
	Reviews []ReviewV2Summary `json:"reviews"`
}

// ListV2 returns reviews using cursor pagination and the current v2 API.
func (c Reviews) ListV2(ctx context.Context, params *ListReviewsV2Params) (*ListReviewsV2Response, error) {
	resp := &ListReviewsV2Response{}
	response, err := c.client.Request(ctx, http.MethodPost, "/v2/review/list", params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)
	return resp, nil
}

type LeaveCommentParams struct {
	// Review status update
	MarkReviewAsProcesses bool `json:"mark_review_as_processed"`

	// Identifier of the parent comment you're replying to
	ParentCommentId string `json:"parent_comment_id"`

	// Review identifier
	ReviewId string `json:"review_id"`

	// Comment text
	Text string `json:"text"`
}

type LeaveCommentResponse struct {
	core.CommonResponse

	// Comment identifier
	CommentId string `json:"comment_id"`
}

// Only available to sellers with the Premium Plus subscription
func (c Reviews) LeaveComment(ctx context.Context, params *LeaveCommentParams) (*LeaveCommentResponse, error) {
	url := "/v1/review/comment/create"

	resp := &LeaveCommentResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type DeleteCommentParams struct {
	// Comment identifier
	CommentId string `json:"comment_id"`
}

type DeleteCommentResponse struct {
	core.CommonResponse
}

// Only available to sellers with the Premium Plus subscription.
//
// Deprecated: Use DeleteCommentV2.
func (c Reviews) DeleteComment(ctx context.Context, params *DeleteCommentParams) (*DeleteCommentResponse, error) {
	url := "/v1/review/comment/delete"

	resp := &DeleteCommentResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type ListCommentsParams struct {
	// Limit of values in the response. Minimum is 20. Maximum is 100
	Limit int32 `json:"limit"`

	// Number of elements that is skipped in the response.
	// For example, if offset = 10, the response starts with the 11th element found
	Offset int32 `json:"offset"`

	// Review identifier
	ReviewId string `json:"review_id"`

	// 	Sorting direction
	SortDir Order `json:"sort_dir"`
}

type ListCommentsResponse struct {
	core.CommonResponse

	// Number of elements in the response
	Offset int32 `json:"offset"`

	// Comment details
	Comments []Comment `json:"comments"`
}

type Comment struct {
	// Comment identifier
	Id string `json:"id"`

	// true, if the comment was left by an official, false if a customer left it
	IsOfficial bool `json:"is_official"`

	// true, if the comment was left by a seller, false if a customer left it
	IsOwner bool `json:"is_owner"`

	// Identifier of the parent comment to reply to
	ParentCommentId string `json:"parent_comment_id"`

	// Date the comment was published
	PublishedAt time.Time `json:"published_at"`

	// Comment text
	Text string `json:"text"`
}

// Only available to sellers with the Premium Plus subscription
//
// Method returns information about comments on reviews that have passed moderation
func (c Reviews) ListComments(ctx context.Context, params *ListCommentsParams) (*ListCommentsResponse, error) {
	url := "/v1/review/comment/list"

	resp := &ListCommentsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// Only available to sellers with the Premium Plus subscription
type ChangeStatusParams struct {
	// Array with review identifiers from 1 to 100
	ReviewIds []string `json:"review_ids"`

	// Review status
	Status string `json:"status"`
}

type ChangeStatusResponse struct {
	core.CommonResponse
}

// Only available to sellers with the Premium Plus subscription.
//
// Deprecated: Use ChangeStatusV2.
func (c Reviews) ChangeStatus(ctx context.Context, params *ChangeStatusParams) (*ChangeStatusResponse, error) {
	url := "/v1/review/change-status"

	resp := &ChangeStatusResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type CountReviewsResponse struct {
	core.CommonResponse

	// Number of processed review
	Processed int32 `json:"processed"`

	// Number of all reviews
	Total int32 `json:"total"`

	// Number of unprocessed reviews
	Unprocessed int32 `json:"unprocessed"`
}

// Only available to sellers with the Premium Plus subscription.
//
// Deprecated: Use CountV2.
func (c Reviews) Count(ctx context.Context) (*CountReviewsResponse, error) {
	url := "/v1/review/count"

	resp := &CountReviewsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, nil, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetReviewParams struct {
	// Review identifier
	ReviewId string `json:"review_id"`
}

type GetReviewResponse struct {
	core.CommonResponse

	ReviewDetails

	// Number of dislikes on the review
	DislikesAmount int32 `json:"dislikes_amount"`

	// Number of likes on the review
	LikesAmount int32 `json:"likes_amount"`

	// Image details
	Photos []ReviewPhoto `json:"photos"`

	// Video details
	Videos []ReviewVideo `json:"videos"`
}

type ReviewDetails struct {
	// Number of comments on the review
	CommentsAmount int32 `json:"comments_amount"`

	// Review identifier
	Id string `json:"id"`

	// true, if the review affects the rating calculation
	IsRatingParticipant bool `json:"is_rating_participant"`

	// Status of the order for which the customer left a review
	OrderStatus string `json:"order_status"`

	// Number of images in the review
	PhotosAmount int32 `json:"photos_amount"`

	// Review publication date
	PublishedAt time.Time `json:"published_at"`

	// Review rating
	Rating int32 `json:"rating"`

	// Product identifier in the Ozon system, SKU
	SKU int64 `json:"sku"`

	// Review status
	Status string `json:"status"`

	// Review text
	Text string `json:"text"`

	// Number of videos for the review
	VideosAmount int32 `json:"videos_amount"`
}

type ReviewPhoto struct {
	// Height
	Height int32 `json:"height"`

	// Link to image
	URL string `json:"url"`

	// Width
	Width int32 `json:"width"`
}

type ReviewVideo struct {
	// Height
	Height int64 `json:"height"`

	// Link to video preview
	PreviewURL string `json:"preview_url"`

	// Link to short video
	ShortVideoPreviewURL string `json:"short_video_preview_url"`

	// Video link
	URL string `json:"url"`

	// Width
	Width int64 `json:"width"`
}

// Only available to sellers with the Premium Plus subscription.
//
// Deprecated: Use GetV2.
func (c Reviews) Get(ctx context.Context, params *GetReviewParams) (*GetReviewResponse, error) {
	url := "/v1/review/info"

	resp := &GetReviewResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, nil, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type ListReviewsParams struct {
	// Identifier of the last review on the page
	LastId string `json:"last_id"`

	// Number of reviews in the response. Minimum is 20, maximum is 100
	Limit int32 `json:"limit"`

	// Sorting direction
	SortDir Order `json:"sort_dir"`

	// Review statuses
	Status string `json:"status"`
}

type ListReviewsResponse struct {
	core.CommonResponse

	// true, if not all reviews were returned in the response
	HasNext bool `json:"has_next"`

	// Identifier of the last review on the page
	LastId string `json:"last_id"`

	// Review details
	Reviews []ReviewDetails `json:"reviews"`
}

// Only available to sellers with the Premium Plus subscription.
//
// Deprecated: Use ListV2.
func (c Reviews) List(ctx context.Context, params *ListReviewsParams) (*ListReviewsResponse, error) {
	url := "/v1/review/list"

	resp := &ListReviewsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, nil, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}
