package ozon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	core "github.com/ucoms-dev/ozon-api-client"
)

// AutoAddProductV2 is a product currently listed for scheduled auto-add to a promotion.
type AutoAddProductV2 struct {
	ProductID                PromotionProductID      `json:"product_id"`
	Price                    *PromotionMoney         `json:"price,omitempty"`
	BasePrice                *PromotionMoney         `json:"base_price,omitempty"`
	MaxDiscountPrice         *PromotionMoney         `json:"max_discount_price,omitempty"`
	ActionPriceToAutoAdd     *PromotionMoney         `json:"action_price_to_auto_add,omitempty"`
	MarketplaceSellerPrice   *PromotionMoney         `json:"marketplace_seller_price,omitempty"`
	MinSellerPrice           *PromotionMoney         `json:"min_seller_price,omitempty"`
	Currency                 string                  `json:"currency,omitempty"`
	OfferID                  string                  `json:"offer_id,omitempty"`
	SKU                      uint64                  `json:"sku,omitempty"`
	Name                     string                  `json:"name,omitempty"`
	MinActionQuantity        uint64                  `json:"min_action_quantity,omitempty"`
	QuantityToAutoAdd        uint64                  `json:"quantity_to_auto_add,omitempty"`
	HasExpiredMinSellerPrice *bool                   `json:"has_expired_min_seller_price,omitempty"`
	WillBeQuarantined        *bool                   `json:"will_be_quarantined,omitempty"`
	AddMode                  *bool                   `json:"add_mode,omitempty"`
	WebsitePrices            *PromotionWebsitePrices `json:"website_prices,omitempty"`
}

// AutoAddCandidateV2 is a product that can be added to a scheduled auto-add promotion.
type AutoAddCandidateV2 struct {
	ID                       PromotionProductID      `json:"id"`
	Price                    *PromotionMoney         `json:"price,omitempty"`
	BasePrice                *PromotionMoney         `json:"base_price,omitempty"`
	MaxDiscountPrice         *PromotionMoney         `json:"max_discount_price,omitempty"`
	ActionPriceToAutoAdd     *PromotionMoney         `json:"action_price_to_auto_add,omitempty"`
	MarketplaceSellerPrice   *PromotionMoney         `json:"marketplace_seller_price,omitempty"`
	MinSellerPrice           *PromotionMoney         `json:"min_seller_price,omitempty"`
	Currency                 string                  `json:"currency,omitempty"`
	OfferID                  string                  `json:"offer_id,omitempty"`
	SKU                      uint64                  `json:"sku,omitempty"`
	Name                     string                  `json:"name,omitempty"`
	MinActionQuantity        uint64                  `json:"min_action_quantity,omitempty"`
	QuantityToAutoAdd        uint64                  `json:"quantity_to_auto_add,omitempty"`
	HasExpiredMinSellerPrice *bool                   `json:"has_expired_min_seller_price,omitempty"`
	WillBeQuarantined        *bool                   `json:"will_be_quarantined,omitempty"`
	IsManuallyAdded          *bool                   `json:"is_manually_added,omitempty"`
	WebsitePrices            *PromotionWebsitePrices `json:"website_prices,omitempty"`
}

// ListAutoAddProductsV2Params selects one scheduled auto-add list page.
type ListAutoAddProductsV2Params struct {
	ActionID    uint64    `json:"action_id"`
	AutoAddDate time.Time `json:"auto_add_date"`
	Limit       uint64    `json:"limit"`
	Offset      uint64    `json:"offset"`
}

// ListAutoAddProductsV2Response is one offset-paginated auto-add product page.
type ListAutoAddProductsV2Response struct {
	core.CommonResponse
	Products []AutoAddProductV2 `json:"products"`
	Total    uint64             `json:"total"`
	decoded  bool
}

func (response *ListAutoAddProductsV2Response) UnmarshalJSON(data []byte) error {
	const method = "ListAutoAddProductsV2"
	response.decoded = false
	if _, err := validatePromotionEnvelope(method, data, []promotionEnvelopeField{
		{name: "products", kind: promotionEnvelopeArray},
		{name: "total", kind: promotionEnvelopeUint64},
	}, nil); err != nil {
		return err
	}
	type wire ListAutoAddProductsV2Response
	if err := json.Unmarshal(data, (*wire)(response)); err != nil {
		return promotionProtocolErrorWithCause(method, "response fields are malformed", err)
	}
	for _, product := range response.Products {
		if product.ProductID == 0 {
			return promotionProtocolError(method, "products contains an invalid product_id")
		}
	}
	response.decoded = true
	return nil
}

// ListAutoAddProductsV2 returns one offset-paginated page for a date supplied by Ozon.
// The schedule date comes from GetAvailablePromotions; this method does not infer dates
// from the local clock or walk later pages automatically.
func (c Promotions) ListAutoAddProductsV2(ctx context.Context, params *ListAutoAddProductsV2Params) (*ListAutoAddProductsV2Response, error) {
	const method = "ListAutoAddProductsV2"
	if params == nil {
		return nil, promotionInputError(method, "params", "must not be nil")
	}
	if err := validatePromotionActionID(method, params.ActionID); err != nil {
		return nil, err
	}
	if params.AutoAddDate.IsZero() {
		return nil, promotionInputError(method, "auto_add_date", "must not be zero")
	}
	if err := validatePromotionLimit(method, params.Limit); err != nil {
		return nil, err
	}
	if err := validatePromotionCall(method, ctx, c.client != nil); err != nil {
		return nil, err
	}

	request := *params
	url := "/v2/actions/auto-add/products/list"
	responseBody := &ListAutoAddProductsV2Response{}
	response, err := c.client.Request(ctx, http.MethodPost, url, &request, responseBody, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&responseBody.CommonResponse)
	if response.StatusCode == http.StatusOK && !responseBody.decoded {
		return nil, promotionProtocolError(method, "response body is empty")
	}
	return responseBody, nil
}

// ListAutoAddCandidatesV2Params selects one candidate page for an auto-add date.
type ListAutoAddCandidatesV2Params struct {
	ActionID    uint64    `json:"action_id"`
	AutoAddDate time.Time `json:"auto_add_date"`
	Limit       uint64    `json:"limit"`
	Offset      uint64    `json:"offset"`
}

// ListAutoAddCandidatesV2Response is one offset-paginated auto-add candidate page.
type ListAutoAddCandidatesV2Response struct {
	core.CommonResponse
	Products []AutoAddCandidateV2 `json:"products"`
	Total    uint64               `json:"total"`
	decoded  bool
}

func (response *ListAutoAddCandidatesV2Response) UnmarshalJSON(data []byte) error {
	const method = "ListAutoAddCandidatesV2"
	response.decoded = false
	if _, err := validatePromotionEnvelope(method, data, []promotionEnvelopeField{
		{name: "products", kind: promotionEnvelopeArray},
		{name: "total", kind: promotionEnvelopeUint64},
	}, nil); err != nil {
		return err
	}
	type wire ListAutoAddCandidatesV2Response
	if err := json.Unmarshal(data, (*wire)(response)); err != nil {
		return promotionProtocolErrorWithCause(method, "response fields are malformed", err)
	}
	for _, candidate := range response.Products {
		if candidate.ID == 0 {
			return promotionProtocolError(method, "products contains an invalid id")
		}
	}
	response.decoded = true
	return nil
}

// ListAutoAddCandidatesV2 returns one page of products eligible for scheduled auto-add.
// The date is explicit because available auto-add dates are supplied by Ozon.
func (c Promotions) ListAutoAddCandidatesV2(ctx context.Context, params *ListAutoAddCandidatesV2Params) (*ListAutoAddCandidatesV2Response, error) {
	const method = "ListAutoAddCandidatesV2"
	if params == nil {
		return nil, promotionInputError(method, "params", "must not be nil")
	}
	if err := validatePromotionActionID(method, params.ActionID); err != nil {
		return nil, err
	}
	if params.AutoAddDate.IsZero() {
		return nil, promotionInputError(method, "auto_add_date", "must not be zero")
	}
	if err := validatePromotionLimit(method, params.Limit); err != nil {
		return nil, err
	}
	if err := validatePromotionCall(method, ctx, c.client != nil); err != nil {
		return nil, err
	}

	request := *params
	url := "/v2/actions/auto-add/products/candidates"
	responseBody := &ListAutoAddCandidatesV2Response{}
	response, err := c.client.Request(ctx, http.MethodPost, url, &request, responseBody, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&responseBody.CommonResponse)
	if response.StatusCode == http.StatusOK && !responseBody.decoded {
		return nil, promotionProtocolError(method, "response body is empty")
	}
	return responseBody, nil
}

// UpdateAutoAddProductV2 describes a product price scheduled for an auto-add date.
type UpdateAutoAddProductV2 struct {
	ID          PromotionProductID `json:"id"`
	ActionPrice PromotionMoney     `json:"action_price"`
	Stock       *uint64            `json:"stock,omitempty"`
}

// UpdateAutoAddProductsV2Params schedules product price updates for one Ozon auto-add date.
type UpdateAutoAddProductsV2Params struct {
	ActionID    uint64                   `json:"action_id"`
	AutoAddDate time.Time                `json:"auto_add_date"`
	Products    []UpdateAutoAddProductV2 `json:"products"`
}

type updateAutoAddProductsV2Request struct {
	ActionID    uint64                          `json:"action_id"`
	AutoAddDate time.Time                       `json:"auto_add_date"`
	Products    []updateAutoAddProductV2Request `json:"products"`
}

type updateAutoAddProductV2Request struct {
	ID          PromotionProductID `json:"id"`
	ActionPrice PromotionMoney     `json:"action_price"`
	// Stock uses uint64 or nil via promotionWireStock; the core default walker
	// otherwise dereferences *uint64 and panics before JSON encoding.
	Stock interface{} `json:"stock,omitempty"`
}

// UpdateAutoAddProductsV2Response reports confirmed, deactivated, rejected, warned,
// and price-validation results independently. The product IDs do not imply that the
// scheduled changes were later applied.
type UpdateAutoAddProductsV2Response struct {
	core.CommonResponse
	ProductIDs        []PromotionProductID       `json:"product_ids,omitempty"`
	DeactivatedIDs    []PromotionProductID       `json:"deactivated_ids,omitempty"`
	Rejected          []PromotionIssue           `json:"rejected,omitempty"`
	Warnings          []PromotionIssue           `json:"warnings,omitempty"`
	BelowMinPrice     []PromotionPriceValidation `json:"below_min_price,omitempty"`
	ExtremelyLowPrice []PromotionPriceValidation `json:"extremely_low_price,omitempty"`
	FailedPrice       []PromotionPriceValidation `json:"failed_price,omitempty"`
	decoded           bool
}

func (response *UpdateAutoAddProductsV2Response) UnmarshalJSON(data []byte) error {
	const method = "UpdateAutoAddProductsV2"
	response.decoded = false
	fields, err := validatePromotionEnvelope(method, data, nil, []promotionEnvelopeField{
		{name: "product_ids", kind: promotionEnvelopeArray},
		{name: "deactivated_ids", kind: promotionEnvelopeArray},
		{name: "rejected", kind: promotionEnvelopeArray},
		{name: "warnings", kind: promotionEnvelopeArray},
		{name: "below_min_price", kind: promotionEnvelopeArray},
		{name: "extremely_low_price", kind: promotionEnvelopeArray},
		{name: "failed_price", kind: promotionEnvelopeArray},
	})
	if err != nil {
		return err
	}
	if err := validateAutoAddPriceDiagnostics(method, fields); err != nil {
		return err
	}
	type wire UpdateAutoAddProductsV2Response
	if err := json.Unmarshal(data, (*wire)(response)); err != nil {
		return promotionProtocolErrorWithCause(method, "response fields are malformed", err)
	}
	if err := validatePromotionResultIDs(method, "product_ids", response.ProductIDs); err != nil {
		return err
	}
	if err := validatePromotionResultIDs(method, "deactivated_ids", response.DeactivatedIDs); err != nil {
		return err
	}
	if err := validateAutoAddIssues(method, "rejected", response.Rejected); err != nil {
		return err
	}
	if err := validateAutoAddIssues(method, "warnings", response.Warnings); err != nil {
		return err
	}
	response.decoded = true
	return nil
}

func validateAutoAddPriceDiagnostics(method string, fields map[string]json.RawMessage) error {
	for _, fieldName := range []string{"below_min_price", "extremely_low_price", "failed_price"} {
		data, exists := fields[fieldName]
		if !exists {
			continue
		}
		var items []map[string]json.RawMessage
		if err := json.Unmarshal(data, &items); err != nil {
			return promotionProtocolErrorWithCause(method, "field "+fieldName+" is malformed", err)
		}
		for index, item := range items {
			keyData, hasKey := item["key"]
			if !hasKey {
				return promotionProtocolError(method, fmt.Sprintf("field %s item %d is missing key", fieldName, index))
			}
			var key PromotionProductID
			if err := json.Unmarshal(keyData, &key); err != nil || key == 0 {
				return promotionProtocolErrorWithCause(method, fmt.Sprintf("field %s item %d has an invalid key", fieldName, index), err)
			}
			valueData, hasValue := item["value"]
			if !hasValue {
				return promotionProtocolError(method, fmt.Sprintf("field %s item %d is missing value", fieldName, index))
			}
			trimmedValue := bytes.TrimSpace(valueData)
			if len(trimmedValue) == 0 || trimmedValue[0] == '"' {
				return promotionProtocolError(method, fmt.Sprintf("field %s item %d value must be a JSON number", fieldName, index))
			}
			var value json.Number
			if err := json.Unmarshal(trimmedValue, &value); err != nil || value == "" {
				return promotionProtocolErrorWithCause(method, fmt.Sprintf("field %s item %d has an invalid number", fieldName, index), err)
			}
		}
	}
	return nil
}

func validateAutoAddIssues(method, field string, issues []PromotionIssue) error {
	for _, issue := range issues {
		if issue.ProductID == 0 {
			return promotionProtocolError(method, field+" contains a zero product_id")
		}
	}
	return nil
}

// UpdateAutoAddProductsV2 schedules new prices for an explicitly supplied Ozon date.
// It sends one update request and reports Ozon's schedule response; callers must list
// the promotion again later to confirm what became active.
func (c Promotions) UpdateAutoAddProductsV2(ctx context.Context, params *UpdateAutoAddProductsV2Params) (*UpdateAutoAddProductsV2Response, error) {
	const method = "UpdateAutoAddProductsV2"
	if params == nil {
		return nil, promotionInputError(method, "params", "must not be nil")
	}
	if err := validatePromotionActionID(method, params.ActionID); err != nil {
		return nil, err
	}
	if params.AutoAddDate.IsZero() {
		return nil, promotionInputError(method, "auto_add_date", "must not be zero")
	}
	productIDs := make([]PromotionProductID, len(params.Products))
	for index, product := range params.Products {
		productIDs[index] = product.ID
	}
	if err := validatePromotionProductIDs(method, "products", productIDs, 1000); err != nil {
		return nil, err
	}
	for index, product := range params.Products {
		if err := validatePromotionMoney(method, fmt.Sprintf("products[%d].action_price", index), product.ActionPrice); err != nil {
			return nil, err
		}
	}
	if err := validatePromotionCall(method, ctx, c.client != nil); err != nil {
		return nil, err
	}

	request := updateAutoAddProductsV2Request{
		ActionID:    params.ActionID,
		AutoAddDate: params.AutoAddDate,
		Products:    make([]updateAutoAddProductV2Request, len(params.Products)),
	}
	for index, product := range params.Products {
		request.Products[index] = updateAutoAddProductV2Request{
			ID:          product.ID,
			ActionPrice: product.ActionPrice,
			Stock:       promotionWireStock(product.Stock),
		}
	}
	url := "/v2/actions/auto-add/products/update"
	responseBody := &UpdateAutoAddProductsV2Response{}
	response, err := c.client.Request(ctx, http.MethodPost, url, &request, responseBody, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&responseBody.CommonResponse)
	if response.StatusCode == http.StatusOK && !responseBody.decoded {
		return nil, promotionProtocolError(method, "response body is empty")
	}
	return responseBody, nil
}

// DeleteAutoAddProductsV2Params identifies scheduled products to remove from one date.
type DeleteAutoAddProductsV2Params struct {
	ActionID    uint64               `json:"action_id"`
	AutoAddDate time.Time            `json:"auto_add_date"`
	ProductIDs  []PromotionProductID `json:"product_ids"`
}

// DeleteAutoAddProductsV2Response lists product IDs returned by Ozon for deletion.
type DeleteAutoAddProductsV2Response struct {
	core.CommonResponse
	ProductIDs []PromotionProductID `json:"product_ids"`
	decoded    bool
}

func (response *DeleteAutoAddProductsV2Response) UnmarshalJSON(data []byte) error {
	const method = "DeleteAutoAddProductsV2"
	response.decoded = false
	if _, err := validatePromotionEnvelope(method, data, []promotionEnvelopeField{
		{name: "product_ids", kind: promotionEnvelopeArray},
	}, nil); err != nil {
		return err
	}
	type wire DeleteAutoAddProductsV2Response
	if err := json.Unmarshal(data, (*wire)(response)); err != nil {
		return promotionProtocolErrorWithCause(method, "response fields are malformed", err)
	}
	if err := validatePromotionResultIDs(method, "product_ids", response.ProductIDs); err != nil {
		return err
	}
	response.decoded = true
	return nil
}

type deleteAutoAddProductsV2Request struct {
	ActionID    uint64    `json:"action_id"`
	AutoAddDate time.Time `json:"auto_add_date"`
	ProductIDs  []string  `json:"product_ids"`
}

// DeleteAutoAddProductsV2 removes the supplied products from a scheduled auto-add date.
func (c Promotions) DeleteAutoAddProductsV2(ctx context.Context, params *DeleteAutoAddProductsV2Params) (*DeleteAutoAddProductsV2Response, error) {
	const method = "DeleteAutoAddProductsV2"
	if params == nil {
		return nil, promotionInputError(method, "params", "must not be nil")
	}
	if err := validatePromotionActionID(method, params.ActionID); err != nil {
		return nil, err
	}
	if params.AutoAddDate.IsZero() {
		return nil, promotionInputError(method, "auto_add_date", "must not be zero")
	}
	if err := validatePromotionProductIDs(method, "product_ids", params.ProductIDs, 1000); err != nil {
		return nil, err
	}
	if err := validatePromotionCall(method, ctx, c.client != nil); err != nil {
		return nil, err
	}

	productIDs := make([]string, len(params.ProductIDs))
	for index, id := range params.ProductIDs {
		productIDs[index] = strconv.FormatUint(uint64(id), 10)
	}
	request := deleteAutoAddProductsV2Request{
		ActionID:    params.ActionID,
		AutoAddDate: params.AutoAddDate,
		ProductIDs:  productIDs,
	}
	url := "/v2/actions/auto-add/products/delete"
	responseBody := &DeleteAutoAddProductsV2Response{}
	response, err := c.client.Request(ctx, http.MethodPost, url, &request, responseBody, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&responseBody.CommonResponse)
	if response.StatusCode == http.StatusOK && !responseBody.decoded {
		return nil, promotionProtocolError(method, "response body is empty")
	}
	return responseBody, nil
}
