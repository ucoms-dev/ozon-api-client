package ozon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	core "github.com/ucoms-dev/ozon-api-client"
)

// PromotionProductPrices contains the shared price, boost, stock and quarantine
// fields returned for participants and candidates.
type PromotionProductPrices struct {
	ID                        PromotionProductID      `json:"id"`
	Price                     *PromotionMoney         `json:"price,omitempty"`
	ActionPrice               *PromotionMoney         `json:"action_price,omitempty"`
	MaxActionPrice            *PromotionMoney         `json:"max_action_price,omitempty"`
	AlertMaxActionPrice       *PromotionMoney         `json:"alert_max_action_price,omitempty"`
	AlertMaxActionPriceFailed *bool                   `json:"alert_max_action_price_failed,omitempty"`
	CurrentBoost              float64                 `json:"current_boost"`
	PriceMinElastic           *PromotionMoney         `json:"price_min_elastic,omitempty"`
	PriceMaxElastic           *PromotionMoney         `json:"price_max_elastic,omitempty"`
	MinBoost                  float64                 `json:"min_boost"`
	MaxBoost                  float64                 `json:"max_boost"`
	MinStock                  *uint64                 `json:"min_stock,omitempty"`
	RecommendedStock          *uint64                 `json:"recommended_stock,omitempty"`
	MarketplaceSellerPrice    *PromotionMoney         `json:"marketplace_seller_price,omitempty"`
	MinSellerPrice            *PromotionMoney         `json:"min_seller_price,omitempty"`
	IsQuarantined             *bool                   `json:"is_quarantined,omitempty"`
	WebsitePrices             *PromotionWebsitePrices `json:"website_prices,omitempty"`
}

// PromotionParticipantV2 is a product currently participating in an action.
// AddMode preserves unrecognized provider values for forward compatibility.
type PromotionParticipantV2 struct {
	PromotionProductPrices
	AddMode string  `json:"add_mode"`
	Stock   *uint64 `json:"stock,omitempty"`
}

// PromotionCandidateV2 is a product that can join an action.
type PromotionCandidateV2 struct {
	PromotionProductPrices
}

type ProductsInPromotionV2Params struct {
	ActionID uint64 `json:"action_id"`
	Limit    uint64 `json:"limit"`
	LastID   string `json:"last_id"`
}

type ProductsInPromotionV2Response struct {
	core.CommonResponse
	Products []PromotionParticipantV2 `json:"products"`
	Total    uint64                   `json:"total"`
	LastID   string                   `json:"last_id,omitempty"`
	decoded  bool
}

func (response *ProductsInPromotionV2Response) UnmarshalJSON(data []byte) error {
	const method = "ProductsInPromotionV2"
	fields, err := validatePromotionEnvelope(method, data, []promotionEnvelopeField{
		{name: "products", kind: promotionEnvelopeArray},
		{name: "total", kind: promotionEnvelopeUint64},
	}, nil)
	if err != nil {
		return err
	}
	if cursor, exists := fields["last_id"]; exists {
		if err := validatePromotionEnvelopeField(method, promotionEnvelopeField{name: "last_id", kind: promotionEnvelopeString}, cursor); err != nil {
			return err
		}
	}
	type wire ProductsInPromotionV2Response
	var decoded wire
	if err := json.Unmarshal(data, &decoded); err != nil {
		return promotionProtocolErrorWithCause(method, "response fields could not be decoded", err)
	}
	for _, product := range decoded.Products {
		if err := validatePromotionResultIDs(method, "products", []PromotionProductID{product.ID}); err != nil {
			return err
		}
	}
	decoded.decoded = true
	*response = ProductsInPromotionV2Response(decoded)
	return nil
}

type ProductsAvailableForPromotionV2Params struct {
	ActionID uint64 `json:"action_id"`
	Limit    uint64 `json:"limit"`
	LastID   string `json:"last_id"`
}

type ProductsAvailableForPromotionV2Response struct {
	core.CommonResponse
	Products []PromotionCandidateV2 `json:"products"`
	Total    uint64                 `json:"total"`
	LastID   string                 `json:"last_id,omitempty"`
	decoded  bool
}

func (response *ProductsAvailableForPromotionV2Response) UnmarshalJSON(data []byte) error {
	const method = "ProductsAvailableForPromotionV2"
	fields, err := validatePromotionEnvelope(method, data, []promotionEnvelopeField{
		{name: "products", kind: promotionEnvelopeArray},
		{name: "total", kind: promotionEnvelopeUint64},
	}, nil)
	if err != nil {
		return err
	}
	if cursor, exists := fields["last_id"]; exists {
		if err := validatePromotionEnvelopeField(method, promotionEnvelopeField{name: "last_id", kind: promotionEnvelopeString}, cursor); err != nil {
			return err
		}
	}
	type wire ProductsAvailableForPromotionV2Response
	var decoded wire
	if err := json.Unmarshal(data, &decoded); err != nil {
		return promotionProtocolErrorWithCause(method, "response fields could not be decoded", err)
	}
	for _, product := range decoded.Products {
		if err := validatePromotionResultIDs(method, "products", []PromotionProductID{product.ID}); err != nil {
			return err
		}
	}
	decoded.decoded = true
	*response = ProductsAvailableForPromotionV2Response(decoded)
	return nil
}

// ProductsInPromotionV2 returns one cursor page of products participating in
// the promotion. It does not follow LastID automatically. After Ozon's 13
// October 2026 change, use UpdateProducts for writes; this read endpoint keeps
// returning participation state independently of future price changes.
func (c Promotions) ProductsInPromotionV2(ctx context.Context, params *ProductsInPromotionV2Params) (*ProductsInPromotionV2Response, error) {
	const method = "ProductsInPromotionV2"
	if params == nil {
		return nil, promotionInputError(method, "params", "must not be nil")
	}
	if err := validatePromotionCall(method, ctx, c.client != nil); err != nil {
		return nil, err
	}
	if err := validatePromotionActionID(method, params.ActionID); err != nil {
		return nil, err
	}
	if err := validatePromotionLimit(method, params.Limit); err != nil {
		return nil, err
	}

	requestParams := *params
	url := "/v2/actions/products"
	resp := &ProductsInPromotionV2Response{}
	response, err := c.client.Request(ctx, http.MethodPost, url, &requestParams, resp, nil)
	if err != nil {
		return nil, fmt.Errorf("%s request failed: %w", method, err)
	}
	response.CopyCommonResponse(&resp.CommonResponse)
	if response.StatusCode == http.StatusOK && !resp.decoded {
		return nil, promotionProtocolError(method, "response body is empty")
	}
	return resp, nil
}

// ProductsAvailableForPromotionV2 returns one cursor page of candidates. It
// leaves cursor traversal and all participation decisions to the caller.
func (c Promotions) ProductsAvailableForPromotionV2(ctx context.Context, params *ProductsAvailableForPromotionV2Params) (*ProductsAvailableForPromotionV2Response, error) {
	const method = "ProductsAvailableForPromotionV2"
	if params == nil {
		return nil, promotionInputError(method, "params", "must not be nil")
	}
	if err := validatePromotionCall(method, ctx, c.client != nil); err != nil {
		return nil, err
	}
	if err := validatePromotionActionID(method, params.ActionID); err != nil {
		return nil, err
	}
	if err := validatePromotionLimit(method, params.Limit); err != nil {
		return nil, err
	}

	requestParams := *params
	url := "/v2/actions/candidates"
	resp := &ProductsAvailableForPromotionV2Response{}
	response, err := c.client.Request(ctx, http.MethodPost, url, &requestParams, resp, nil)
	if err != nil {
		return nil, fmt.Errorf("%s request failed: %w", method, err)
	}
	response.CopyCommonResponse(&resp.CommonResponse)
	if response.StatusCode == http.StatusOK && !resp.decoded {
		return nil, promotionProtocolError(method, "response body is empty")
	}
	return resp, nil
}

type UpdatePromotionProduct struct {
	ProductID   PromotionProductID `json:"product_id"`
	ActionPrice PromotionMoney     `json:"action_price"`
	Stock       *uint64            `json:"stock,omitempty"`
}

type UpdatePromotionProductsParams struct {
	ActionID uint64                   `json:"action_id"`
	Products []UpdatePromotionProduct `json:"products"`
}

type UpdatePromotionProductsResponse struct {
	core.CommonResponse
	ActiveProductIDs      []PromotionProductID `json:"active_product_ids"`
	DeactivatedProductIDs []PromotionProductID `json:"deactivated_product_ids"`
	Rejected              []PromotionIssue     `json:"rejected"`
	Warnings              []PromotionIssue     `json:"warnings"`
	decoded               bool
}

func (response *UpdatePromotionProductsResponse) UnmarshalJSON(data []byte) error {
	const method = "UpdateProducts"
	_, err := validatePromotionEnvelope(method, data, nil, []promotionEnvelopeField{
		{name: "active_product_ids", kind: promotionEnvelopeArray},
		{name: "deactivated_product_ids", kind: promotionEnvelopeArray},
		{name: "rejected", kind: promotionEnvelopeArray},
		{name: "warnings", kind: promotionEnvelopeArray},
	})
	if err != nil {
		return err
	}
	type wire UpdatePromotionProductsResponse
	var decoded wire
	if err := json.Unmarshal(data, &decoded); err != nil {
		return promotionProtocolErrorWithCause(method, "response fields could not be decoded", err)
	}
	for _, field := range []struct {
		name string
		ids  []PromotionProductID
	}{
		{name: "active_product_ids", ids: decoded.ActiveProductIDs},
		{name: "deactivated_product_ids", ids: decoded.DeactivatedProductIDs},
	} {
		if err := validatePromotionResultIDs(method, field.name, field.ids); err != nil {
			return err
		}
	}
	for _, field := range []struct {
		name   string
		issues []PromotionIssue
	}{
		{name: "rejected", issues: decoded.Rejected},
		{name: "warnings", issues: decoded.Warnings},
	} {
		for _, issue := range field.issues {
			if issue.ProductID == 0 {
				return promotionProtocolError(method, "field "+field.name+" contains a zero product ID")
			}
		}
	}
	decoded.decoded = true
	*response = UpdatePromotionProductsResponse(decoded)
	return nil
}

type updatePromotionProductWire struct {
	ProductID   PromotionProductID `json:"product_id"`
	ActionPrice PromotionMoney     `json:"action_price"`
	Stock       interface{}        `json:"stock,omitempty"`
}

type updatePromotionProductsWire struct {
	ActionID uint64                       `json:"action_id"`
	Products []updatePromotionProductWire `json:"products"`
}

// UpdateProducts sets prices and optional stock for one batch. Stock is sent
// only for action types that use stock; selecting it remains the caller's
// responsibility. Since 13 October 2026, Ozon may apply the action price as a
// maximum price outside the promotion. A successful HTTP response reports
// per-product outcomes and does not prove future participation or application.
func (c Promotions) UpdateProducts(ctx context.Context, params *UpdatePromotionProductsParams) (*UpdatePromotionProductsResponse, error) {
	const method = "UpdateProducts"
	if params == nil {
		return nil, promotionInputError(method, "params", "must not be nil")
	}
	if err := validatePromotionCall(method, ctx, c.client != nil); err != nil {
		return nil, err
	}
	if err := validatePromotionActionID(method, params.ActionID); err != nil {
		return nil, err
	}
	if len(params.Products) == 0 || len(params.Products) > 1000 {
		return nil, promotionInputError(method, "products", "must contain between 1 and 1000 products")
	}
	productIDs := make([]PromotionProductID, len(params.Products))
	for index, product := range params.Products {
		productIDs[index] = product.ProductID
		if err := validatePromotionMoney(method, "products.action_price", product.ActionPrice); err != nil {
			return nil, err
		}
	}
	if err := validatePromotionProductIDs(method, "products", productIDs, 1000); err != nil {
		return nil, err
	}

	request := updatePromotionProductsWire{
		ActionID: params.ActionID,
		Products: make([]updatePromotionProductWire, len(params.Products)),
	}
	for index, product := range params.Products {
		request.Products[index] = updatePromotionProductWire{
			ProductID:   product.ProductID,
			ActionPrice: product.ActionPrice,
			Stock:       promotionWireStock(product.Stock),
		}
	}

	url := "/v1/actions/products/update"
	resp := &UpdatePromotionProductsResponse{}
	response, err := c.client.Request(ctx, http.MethodPost, url, &request, resp, nil)
	if err != nil {
		return nil, fmt.Errorf("%s request failed: %w", method, err)
	}
	response.CopyCommonResponse(&resp.CommonResponse)
	if response.StatusCode == http.StatusOK && !resp.decoded {
		return nil, promotionProtocolError(method, "response body is empty")
	}
	return resp, nil
}

type RemoveProductFromPromotionV2Params struct {
	ActionID   uint64               `json:"action_id"`
	ProductIDs []PromotionProductID `json:"product_ids"`
}

type RemoveProductFromPromotionV2Response struct {
	core.CommonResponse
	ProductIDs []PromotionProductID `json:"product_ids"`
	decoded    bool
}

func (response *RemoveProductFromPromotionV2Response) UnmarshalJSON(data []byte) error {
	const method = "RemoveProductV2"
	if _, err := validatePromotionEnvelope(method, data, []promotionEnvelopeField{
		{name: "product_ids", kind: promotionEnvelopeArray},
	}, nil); err != nil {
		return err
	}
	type wire RemoveProductFromPromotionV2Response
	var decoded wire
	if err := json.Unmarshal(data, &decoded); err != nil {
		return promotionProtocolErrorWithCause(method, "response fields could not be decoded", err)
	}
	if err := validatePromotionResultIDs(method, "product_ids", decoded.ProductIDs); err != nil {
		return err
	}
	decoded.decoded = true
	*response = RemoveProductFromPromotionV2Response(decoded)
	return nil
}

type removeProductFromPromotionV2Wire struct {
	ActionID   uint64   `json:"action_id"`
	ProductIDs []string `json:"product_ids"`
}

// RemoveProductV2 deactivates one batch of products from an action. The
// response contains only product IDs confirmed by Ozon; callers must preserve
// any unconfirmed IDs for their own follow-up policy.
func (c Promotions) RemoveProductV2(ctx context.Context, params *RemoveProductFromPromotionV2Params) (*RemoveProductFromPromotionV2Response, error) {
	const method = "RemoveProductV2"
	if params == nil {
		return nil, promotionInputError(method, "params", "must not be nil")
	}
	if err := validatePromotionCall(method, ctx, c.client != nil); err != nil {
		return nil, err
	}
	if err := validatePromotionActionID(method, params.ActionID); err != nil {
		return nil, err
	}
	if err := validatePromotionProductIDs(method, "product_ids", params.ProductIDs, 1000); err != nil {
		return nil, err
	}

	request := removeProductFromPromotionV2Wire{
		ActionID:   params.ActionID,
		ProductIDs: make([]string, len(params.ProductIDs)),
	}
	for index, id := range params.ProductIDs {
		request.ProductIDs[index] = strconv.FormatUint(uint64(id), 10)
	}

	url := "/v2/actions/products/deactivate"
	resp := &RemoveProductFromPromotionV2Response{}
	response, err := c.client.Request(ctx, http.MethodPost, url, &request, resp, nil)
	if err != nil {
		return nil, fmt.Errorf("%s request failed: %w", method, err)
	}
	response.CopyCommonResponse(&resp.CommonResponse)
	if response.StatusCode == http.StatusOK && !resp.decoded {
		return nil, promotionProtocolError(method, "response body is empty")
	}
	return resp, nil
}
