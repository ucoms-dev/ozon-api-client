package ozon

import (
	"context"
	"net/http"
	"time"

	core "github.com/ucoms-dev/ozon-api-client"
)

type FBO struct {
	client *core.Client
}

type GetFBOShipmentsListParams struct {
	// Sorting direction
	Direction Order `json:"dir,omitempty"`

	// Shipment search filter
	Filter GetFBOShipmentsListFilter `json:"filter"`

	// Number of values in the response. Maximum is 1000, minimum is 1
	Limit int64 `json:"limit"`

	// Number of elements that will be skipped in the response. For example, if offset=10, the response will start with the 11th element found
	Offset int64 `json:"offset,omitempty"`

	// true if the address transliteration from Cyrillic to Latin is enabled
	Translit bool `json:"translit,omitempty"`

	// Additional fields to add to the response
	With *GetFBOShipmentsListWith `json:"with,omitempty"`
}

// Shipment search filter
type GetFBOShipmentsListFilter struct {
	// Period start in YYYY-MM-DD format
	Since *core.TimeFormat `json:"since,omitempty"`

	// Shipment status
	Status string `json:"status,omitempty"`

	// Period end in YYYY-MM-DD format
	To *core.TimeFormat `json:"to,omitempty"`
}

// Additional fields to add to the response
type GetFBOShipmentsListWith struct {
	// Specify true to add analytics data to the response
	AnalyticsData bool `json:"analytics_data"`

	// Specify true to add financial data to the response
	FinancialData bool `json:"financial_data"`
}

type GetFBOShipmentsListResponse struct {
	core.CommonResponse

	// Shipments list
	Result []GetFBOShipmentsListResult `json:"result"`
}

type GetFBOShipmentsListResult struct {
	// Additional data for shipment list
	AdditionalData []GetFBOShipmentsListResultAdditionalData `json:"additional_data"`

	// Analytical data
	AnalyticsData GetFBOShipmentsListResultAnalyticsData `json:"analytics_data"`

	// Shipment cancellation reason identifier
	CancelReasonId int64 `json:"cancel_reason_id"`

	// Shipment cancellation details
	Cancellation FBOPostingCancellation `json:"cancellation"`

	// Date and time of shipment creation
	CreatedAt time.Time `json:"created_at"`

	// Information about an order created on an external platform
	ExternalOrder PostingExternalOrder `json:"external_order"`

	// Financial data
	FinancialData FBOFinancialData `json:"financial_data"`

	// Date and time of shipment processing start
	InProccessAt time.Time `json:"in_process_at"`

	// Buyer legal information
	LegalInfo PostingLegalInfo `json:"legal_info"`

	// Identifier of the order to which the shipment belongs
	OrderId int64 `json:"order_id"`

	// Number of the order to which the shipment belongs
	OrderNumber string `json:"order_number"`

	// Shipment number
	PostingNumber string `json:"posting_number"`

	// Number of products in the shipment
	Products []FBOPostingProduct `json:"products"`

	// Shipment status
	Status string `json:"status"`

	// Shipment substatus
	Substatus string `json:"substatus"`
}

type GetFBOShipmentsListResultAdditionalData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetFBOShipmentsListResultAnalyticsData struct {
	// Delivery method
	DeliveryType string `json:"delivery_type"`

	// Indication that the recipient is a legal person
	//   * true — a legal person,
	//   * false — a natural person.
	IsLegal bool `json:"is_legal"`

	// Premium subscription
	IsPremium bool `json:"is_premium"`

	// Payment method
	PaymentTypeGroupName PaymentTypeGroupName `json:"payment_type_group_name"`

	// Warehouse identifier
	WarehouseId int64 `json:"warehouse_id"`

	// Name of the warehouse from which the order is shipped
	WarehouseName string `json:"warehouse_name"`
}

type FBOPostingProduct struct {
	// Activation codes for services and digital products
	DigitalCodes []string `json:"digital_codes"`

	// Currency of your prices. It matches the currency set in the personal account settings
	CurrencyCode string `json:"currency_code"`

	// Product name
	Name string `json:"name"`

	// Product identifier in the seller's system
	OfferId string `json:"offer_id"`

	// Product price
	Price string `json:"price"`

	// Quantity of products in the shipment
	Quantity int64 `json:"quantity"`

	// Product identifier in the Ozon system, SKU
	SKU int64 `json:"sku"`
}

type FBOFinancialData struct {
	// Identifier of the cluster, where the shipment is sent from
	ClusterFrom string `json:"cluster_from"`

	// Identifier of the cluster, where the shipment is delivered to
	ClusterTo string `json:"cluster_to"`

	// Services
	PostingServices MarketplaceServices `json:"posting_services"`

	// Products list
	Products []FinancialDataProduct `json:"products"`
}

// Returns a list of shipments for a specified period of time. You can additionally filter the shipments by their status.
//
// Deprecated: Ozon discontinued /v2/posting/fbo/list on August 31, 2026.
// Use GetShipmentsListV3.
func (c FBO) GetShipmentsList(ctx context.Context, params *GetFBOShipmentsListParams) (*GetFBOShipmentsListResponse, error) {
	url := "/v2/posting/fbo/list"

	resp := &GetFBOShipmentsListResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// GetFBOShipmentsListV3Params contains cursor-based filters for the current FBO
// posting list API.
type GetFBOShipmentsListV3Params struct {
	Cursor   string                      `json:"cursor,omitempty"`
	Filter   GetFBOShipmentsListV3Filter `json:"filter"`
	Limit    int64                       `json:"limit"`
	SortDir  Order                       `json:"sort_dir,omitempty"`
	Translit bool                        `json:"translit,omitempty"`
	With     *GetFBOShipmentsListV3With  `json:"with,omitempty"`
}

type GetFBOShipmentsListV3Filter struct {
	OrderNumbers   []string  `json:"order_numbers,omitempty"`
	PostingNumbers []string  `json:"posting_numbers,omitempty"`
	Since          time.Time `json:"since"`
	Statuses       []string  `json:"statuses,omitempty"`
	To             time.Time `json:"to"`
}

type GetFBOShipmentsListV3With struct {
	AnalyticsData bool `json:"analytics_data,omitempty"`
	FinancialData bool `json:"financial_data,omitempty"`
	LegalInfo     bool `json:"legal_info,omitempty"`
}

type GetFBOShipmentsListV3Response struct {
	core.CommonResponse
	Cursor   string                      `json:"cursor"`
	HasNext  bool                        `json:"has_next"`
	Postings []GetFBOShipmentsListResult `json:"postings"`
}

// GetFBOShipmentsListV3CurrentResponse contains the current Ozon wire model.
// It is separate from GetFBOShipmentsListV3Response to preserve the public
// source contract released in v1.16.0.
type GetFBOShipmentsListV3CurrentResponse struct {
	core.CommonResponse
	Cursor   string                         `json:"cursor"`
	HasNext  bool                           `json:"has_next"`
	Postings []GetFBOShipmentsListV3Posting `json:"postings"`
}

type GetFBOShipmentsListV3Posting struct {
	GetFBOShipmentsListResult
	AnalyticsData GetFBOShipmentsListV3AnalyticsData `json:"analytics_data"`
	FinancialData GetFBOShipmentsListV3FinancialData `json:"financial_data"`
	Products      []GetFBOShipmentsListV3Product     `json:"products"`
}

type GetFBOShipmentsListV3AnalyticsData struct {
	City                    string               `json:"city"`
	ClientDeliveryDateBegin time.Time            `json:"client_delivery_date_begin"`
	ClientDeliveryDateEnd   time.Time            `json:"client_delivery_date_end"`
	DeliveryType            string               `json:"delivery_type"`
	IsLegal                 bool                 `json:"is_legal"`
	IsPremium               bool                 `json:"is_premium"`
	PaymentTypeGroupName    PaymentTypeGroupName `json:"payment_type_group_name"`
	WarehouseId             int64                `json:"warehouse_id"`
	WarehouseName           string               `json:"warehouse_name"`
}

type GetFBOShipmentsListV3Product struct {
	DigitalCodes        []string     `json:"digital_codes"`
	IsMarketplaceBuyout bool         `json:"is_marketplace_buyout"`
	Name                string       `json:"name"`
	OfferId             string       `json:"offer_id"`
	Price               PostingMoney `json:"price"`
	Quantity            int64        `json:"quantity"`
	SKU                 int64        `json:"sku"`
}

type GetFBOShipmentsListV3FinancialData struct {
	ClusterFrom string                                      `json:"cluster_from"`
	ClusterTo   string                                      `json:"cluster_to"`
	Products    []GetFBOShipmentsListV3FinancialDataProduct `json:"products"`
}

type GetFBOShipmentsListV3FinancialDataProduct struct {
	Actions              []string                                 `json:"actions"`
	Commission           GetFBOShipmentsListV3FinancialCommission `json:"commission"`
	OldPrice             float64                                  `json:"old_price"`
	Payout               float64                                  `json:"payout"`
	Price                float64                                  `json:"price"`
	ProductId            int64                                    `json:"product_id"`
	Quantity             int64                                    `json:"quantity"`
	TotalDiscountPercent float64                                  `json:"total_discount_percent"`
	TotalDiscountValue   float64                                  `json:"total_discount_value"`
}

type GetFBOShipmentsListV3FinancialCommission struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Percent  int64   `json:"percent"`
}

type FBOPostingCancellation struct {
	CancelReason          string `json:"cancel_reason"`
	CancellationInitiator string `json:"cancellation_initiator"`
	CancellationType      string `json:"cancellation_type"`
}

type PostingExternalOrder struct {
	IsExternal   bool   `json:"is_external"`
	PlatformName string `json:"platform_name"`
}

type PostingLegalInfo struct {
	CompanyName string `json:"company_name"`
	INN         string `json:"inn"`
	KPP         string `json:"kpp"`
}

// GetShipmentsListV3 returns FBO postings using cursor pagination and adapts
// the current wire response to the source-compatible v1.16.0 response model.
func (c FBO) GetShipmentsListV3(ctx context.Context, params *GetFBOShipmentsListV3Params) (*GetFBOShipmentsListV3Response, error) {
	current, err := c.GetShipmentsListV3Current(ctx, params)
	if err != nil {
		return nil, err
	}

	resp := &GetFBOShipmentsListV3Response{
		CommonResponse: current.CommonResponse,
		Cursor:         current.Cursor,
		HasNext:        current.HasNext,
		Postings:       make([]GetFBOShipmentsListResult, len(current.Postings)),
	}
	for i, posting := range current.Postings {
		resp.Postings[i] = adaptFBOShipmentV3(posting)
	}
	return resp, nil
}

// GetShipmentsListV3Current returns FBO postings using the exact current Ozon
// response contract, including money objects and nested commissions.
func (c FBO) GetShipmentsListV3Current(ctx context.Context, params *GetFBOShipmentsListV3Params) (*GetFBOShipmentsListV3CurrentResponse, error) {
	resp := &GetFBOShipmentsListV3CurrentResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, "/v3/posting/fbo/list", params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

func adaptFBOShipmentV3(posting GetFBOShipmentsListV3Posting) GetFBOShipmentsListResult {
	legacy := posting.GetFBOShipmentsListResult
	legacy.AnalyticsData = GetFBOShipmentsListResultAnalyticsData{
		DeliveryType:         posting.AnalyticsData.DeliveryType,
		IsLegal:              posting.AnalyticsData.IsLegal,
		IsPremium:            posting.AnalyticsData.IsPremium,
		PaymentTypeGroupName: posting.AnalyticsData.PaymentTypeGroupName,
		WarehouseId:          posting.AnalyticsData.WarehouseId,
		WarehouseName:        posting.AnalyticsData.WarehouseName,
	}
	legacy.FinancialData = FBOFinancialData{
		ClusterFrom: posting.FinancialData.ClusterFrom,
		ClusterTo:   posting.FinancialData.ClusterTo,
		Products:    make([]FinancialDataProduct, len(posting.FinancialData.Products)),
	}
	for i, product := range posting.FinancialData.Products {
		legacy.FinancialData.Products[i] = FinancialDataProduct{
			Actions:                 product.Actions,
			CommissionAmount:        product.Commission.Amount,
			CommissionPercent:       product.Commission.Percent,
			CommissionsCurrencyCode: product.Commission.Currency,
			OldPrice:                product.OldPrice,
			Payout:                  product.Payout,
			Price:                   product.Price,
			ProductId:               product.ProductId,
			Quantity:                product.Quantity,
			TotalDiscountPercent:    product.TotalDiscountPercent,
			TotalDiscountValue:      product.TotalDiscountValue,
		}
	}
	legacy.Products = make([]FBOPostingProduct, len(posting.Products))
	for i, product := range posting.Products {
		legacy.Products[i] = FBOPostingProduct{
			DigitalCodes: product.DigitalCodes,
			CurrencyCode: product.Price.Currency,
			Name:         product.Name,
			OfferId:      product.OfferId,
			Price:        product.Price.Amount,
			Quantity:     product.Quantity,
			SKU:          product.SKU,
		}
	}
	return legacy
}

type GetShipmentDetailsParams struct {
	// Shipment number
	PostingNumber string `json:"posting_number"`

	// true if the address transliteration from Cyrillic to Latin is enabled
	Translit bool `json:"translit,omitempty"`

	// Additional fields to add to the response
	With *GetShipmentDetailsWith `json:"with,omitempty"`
}

type GetShipmentDetailsWith struct {
	// Specify true to add analytics data to the response
	AnalyticsData bool `json:"analytics_data"`

	// Specify true to add financial data to the response
	FinancialData bool `json:"financial_data"`
}

type GetShipmentDetailsResponse struct {
	core.CommonResponse

	// Method result
	Result GetShipmentDetailsResult `json:"result"`
}

type GetShipmentDetailsResult struct {
	// Additional data
	AdditionalData []GetShipmentDetailsResultAdditionalData `json:"additional_data"`

	// Analytical data
	AnalyticsData GetShipmentDetailsResultAnalyticsData `json:"analytics_data"`

	// Shipment cancellation reason identifier
	CancelReasonId int64 `json:"cancel_reason_id"`

	// Date and time of shipment creation
	CreatedAt time.Time `json:"created_at"`

	// Financial data
	FinancialData FBOFinancialData `json:"financial_data"`

	// Date and time of shipment processing start
	InProcessAt time.Time `json:"in_process_at"`

	// Identifier of the order to which the shipment belongs
	OrderId int64 `json:"order_id"`

	// Number of the order to which the shipment belongs
	OrderNumber string `json:"order_number"`

	// Shipment number
	PostingNumber string `json:"posting_number"`

	// Number of products in the shipment
	Products []FBOPostingProduct `json:"products"`

	// Shipment status
	Status string `json:"status"`
}

type GetShipmentDetailsResultAdditionalData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetShipmentDetailsResultAnalyticsData struct {
	// Delivery method
	DeliveryType string `json:"delivery_type"`

	// Indication that the recipient is a legal person:
	//   - true — a legal person
	//   - false — a natural person
	IsLegal bool `json:"is_legal"`

	// Premium subscription
	IsPremium bool `json:"is_premium"`

	// Payment method
	PaymentTypeGroupName string `json:"payment_type_group_name"`

	// Warehouse identifier
	WarehouseId int64 `json:"warehouse_id"`

	// Name of the warehouse from which the order is shipped
	WarehouseName string `json:"warehouse_name"`
}

// Returns information about the shipment by its identifier
func (c FBO) GetShipmentDetails(ctx context.Context, params *GetShipmentDetailsParams) (*GetShipmentDetailsResponse, error) {
	url := "/v2/posting/fbo/get"

	resp := &GetShipmentDetailsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type ListSupplyRequestsParams struct {
	// Filter
	Filter *ListSupplyRequestsFilter `json:"filter"`

	// Customizing the display of the requests list
	Paging *ListSupplyRequestsPaging `json:"paging"`
}

type ListSupplyRequestsFilter struct {
	States []string `json:"states"`
}

type ListSupplyRequestsPaging struct {
	// Supply number from which the list of requests will start
	FromOrderId int64 `json:"from_supply_order_id"`

	// Number of requests in the response
	Limit int32 `json:"limit"`
}

type ListSupplyRequestsResponse struct {
	core.CommonResponse

	// Supply request identifier you last requested
	LastSupplyOrderId int64 `json:"last_supply_order_id"`

	// Supply request identifier
	SupplyOrderId []string `json:"supply_order_id"`
}

// Requests with supply to a specific warehouse and through a virtual distribution center (vDC) are taken into account
func (c FBO) ListSupplyRequests(ctx context.Context, params *ListSupplyRequestsParams) (*ListSupplyRequestsResponse, error) {
	url := "/v2/supply-order/list"

	resp := &ListSupplyRequestsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetSupplyRequestInfoParams struct {
	// Supply request identifier in the Ozon system
	OrderIds []string `json:"order_ids"`
}

type GetSupplyRequestInfoResponse struct {
	core.CommonResponse

	// Supply request details
	Orders []SupplyOrder `json:"orders"`

	// Warehouse details
	Warehouses []SupplyWarehouse `json:"warehouses"`
}

type SupplyOrder struct {
	// true if the supply request can be canceled
	CanCancel bool `json:"can_cancel"`

	// Date of supply request creation
	CreationDate string `json:"creation_date"`

	// Request source
	CreationFlow string `json:"creation_flow"`

	// Time remaining in seconds to fill in the supply details. Only for requests from the vDC
	DataFillingDeadline time.Time `json:"data_filling_deadline_utc"`

	// Supply warehouse identifier
	DropoffWarehouseId int64 `json:"dropoff_warehouse_id"`

	// true if the supply request contains Super Economy products
	IsEconom bool `json:"is_econom"`

	// true if the seller has Super supplies enabled
	IsSuperFBO bool `json:"is_super_fbo"`

	// true if the supply request is virtual
	IsVirtual bool `json:"is_virtual"`

	// true if the supply request contains Super products
	ProductSuperFBO bool `json:"product_super_fbo"`

	// Filter by supply status
	State string `json:"state"`

	// Supply request contents
	Supplies []Supply `json:"supplies"`

	// Supply request identifier
	Id int64 `json:"supply_order_id"`

	// Request number
	OrderNumber string `json:"supply_order_number"`

	// Supply time slot
	Timeslot []SupplyTimeslot `json:"timeslot"`

	// Driver and vehicle details
	Vehicle []SupplyVehicle `json:"vehicle"`
}

type Supply struct {
	// Supply contents identifier. Used in the /v1/supply-order/bundle method
	BundleId string `json:"bundle_id"`

	// Filter by supply status
	SupplyState string `json:"supply_state"`

	// Storage warehouse identifier
	StorageWarehouseId int64 `json:"storage_warehouse_id"`

	// Product tags in the supply request
	SupplyTags []SupplyTag `json:"supply_tags"`

	// Supply identifier
	Id int64 `json:"supply_id"`
}

type SupplyTag struct {
	// true if the supply request contains products certified in the Mercury system
	IsEVSDRequired bool `json:"is_evsd_required"`

	// true if the supply request contains jewelry
	IsJewelry bool `json:"is_jewelry"`

	// true if the supply request contains products for which labeling is possible
	IsMarkingPossible bool `json:"is_marking_possible"`

	// true if the supply request contains products for which labeling is mandatory
	IsMarkingRequired bool `json:"is_marking_required"`
}

type SupplyTimeslot struct {
	// Reason why you can't select the supply time slot
	Reasons []string `json:"can_not_set_reasons"`

	// true, if you can select or edit the supply time slot
	CanSet bool `json:"can_set"`

	// true, if the characteristic is required
	IsRequired bool `json:"is_required"`

	Value SupplyTimeslotValue `json:"value"`
}

type SupplyVehicle struct {
	// Reason why you can't select the supply time slot
	Reasons []string `json:"can_not_set_reasons"`

	// true, if you can select or edit the supply time slot
	CanSet bool `json:"can_set"`

	// true, if the characteristic is required
	IsRequired bool `json:"is_required"`

	Value []GetSupplyRequestInfoVehicle `json:"value"`
}

type SupplyTimeslotValue struct {
	// Supply time slot in local time
	Timeslot []SupplyTimeslotValueTimeslot `json:"timeslot"`

	// Time zone
	Timezone []SupplyTimeslotValueTimezone `json:"timezone_info"`
}

type SupplyTimeslotValueTimeslot struct {
	// Supply time slot start
	From time.Time `json:"from"`

	// Supply time slot end
	To time.Time `json:"to"`
}

type SupplyTimeslotValueTimezone struct {
	// Time zone name
	Name string `json:"iana_name"`

	// Time zone offset from UTC-0 in seconds
	Offset string `json:"offset"`
}

type GetSupplyRequestInfoVehicle struct {
	// Driver name
	DriverName string `json:"driver_name"`

	// Driver phone number
	DriverPhone string `json:"driver_phone"`

	// Car model
	VehicleModel string `json:"vehicle_model"`

	// Car number
	VehicleNumber string `json:"vehicle_number"`
}

type SupplyWarehouse struct {
	// Warehouse address
	Address string `json:"address"`

	// Warehouse name
	Name string `json:"name"`

	// Warehouse identifier
	Id int64 `json:"warehouse_id"`
}

// Method for getting detailed information on a supply request.
// Requests with supply both to a specific warehouse and via a
// virtual distribution center (vDC) are taken into account
func (c FBO) GetSupplyRequestInfo(ctx context.Context, params *GetSupplyRequestInfoParams) (*GetSupplyRequestInfoResponse, error) {
	url := "/v2/supply-order/get"

	resp := &GetSupplyRequestInfoResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetWarehouseWorkloadResponse struct {
	core.CommonResponse

	// Method result
	Result []GetWarehouseWorkloadResult `json:"result"`
}

type GetWarehouseWorkloadResult struct {
	// Workload
	Schedule GetWarehouseWorkloadResultSchedule `json:"schedule"`

	// Warehouse
	Warehouse GetWarehouseWorkloadResultWarehouse `json:"warehouse"`
}

type GetWarehouseWorkloadResultSchedule struct {
	// Data on the products quantity supplied to the warehouse
	Capacity []GetWarehouseWorkloadResultScheduleCapacity `json:"capacity"`

	// The closest available date for supply, local time
	Date time.Time `json:"date"`
}

type GetWarehouseWorkloadResultScheduleCapacity struct {
	// Period start, local time
	Start time.Time `json:"start"`

	// Period end, local time
	End time.Time `json:"end"`

	// Average number of products that the warehouse can accept per day for the period
	Value int32 `json:"value"`
}

type GetWarehouseWorkloadResultWarehouse struct {
	// Warehouse identifier
	Id string `json:"id"`

	// Warehouse name
	Name string `json:"name"`
}

// Method returns a list of active Ozon warehouses with information about their average workload in the nearest future
func (c FBO) GetWarehouseWorkload(ctx context.Context) (*GetWarehouseWorkloadResponse, error) {
	url := "/v1/supplier/available_warehouses"

	resp := &GetWarehouseWorkloadResponse{}

	response, err := c.client.Request(ctx, http.MethodGet, url, nil, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetSupplyOrdersByStatusParams struct {
}

type GetSupplyOrdersByStatusResponse struct {
	core.CommonResponse

	Items []SupplyOrdersByStatus `json:"items"`
}

type SupplyOrdersByStatus struct {
	// Number of supply requests in this status
	Count int32 `json:"count"`

	// Supply status
	OrderState string `json:"order_state"`
}

// Returns the number of supply requests in a specific status.
func (c FBO) GetSupplyOrdersByStatus(ctx context.Context) (*GetSupplyOrdersByStatusResponse, error) {
	url := "/v1/supply-order/status/counter"

	resp := &GetSupplyOrdersByStatusResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, nil, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetSupplyTimeslotsParams struct {
	// Supply request identifier
	SupplyOrderId int64 `json:"supply_order_id"`
}

type GetSupplyTimeslotsResponse struct {
	core.CommonResponse

	// Supply time slot
	Timeslots []SupplyTimeslotValueTimeslot `json:"timeslots"`

	// Time zone
	Timezones []SupplyTimeslotValueTimezone `json:"timezone"`
}

func (c FBO) GetSupplyTimeslots(ctx context.Context, params *GetSupplyTimeslotsParams) (*GetSupplyTimeslotsResponse, error) {
	url := "/v1/supply-order/timeslot/get"

	resp := &GetSupplyTimeslotsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type UpdateSupplyTimeslotParams struct {
	// Supply request identifier
	SupplyOrderId int64 `json:"supply_order_id"`

	// Supply time slot details
	Timeslot SupplyTimeslotValueTimeslot `json:"timeslot"`
}

type UpdateSupplyTimeslotResponse struct {
	core.CommonResponse

	// Possible errors
	Errors []string `json:"errors"`

	// Operation identifier
	OperationId string `json:"operation_id"`
}

func (c FBO) UpdateSupplyTimeslot(ctx context.Context, params *UpdateSupplyTimeslotParams) (*UpdateSupplyTimeslotResponse, error) {
	url := "/v1/supply-order/timeslot/update"

	resp := &UpdateSupplyTimeslotResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetSupplyTimeslotStatusParams struct {
	// Operation identifier
	OperationId string `json:"operation_id"`
}

type GetSupplyTimeslotStatusResponse struct {
	core.CommonResponse

	// Possible errors
	Errors []string `json:"errors"`

	// Data status
	Status string `json:"status"`
}

func (c FBO) GetSupplyTimeslotStatus(ctx context.Context, params *GetSupplyTimeslotStatusParams) (*GetSupplyTimeslotStatusResponse, error) {
	url := "/v1/supply-order/timeslot/status"

	resp := &GetSupplyTimeslotStatusResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type CreatePassParams struct {
	// Supply request identifier
	SupplyOrderId int64 `json:"supply_order_id"`

	// Driver and car information
	Vehicle GetSupplyRequestInfoVehicle `json:"vehicle"`
}

type CreatePassResponse struct {
	core.CommonResponse

	// Possible errors
	Errors []string `json:"error_reasons"`

	// Operation identifier
	OperationId string `json:"operation_id"`
}

func (c FBO) CreatePass(ctx context.Context, params *CreatePassParams) (*CreatePassResponse, error) {
	url := "/v1/supply-order/pass/create"

	resp := &CreatePassResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetPassParams struct {
	// Operation identifier
	OperationId string `json:"operation_id"`
}

type GetPassResponse struct {
	core.CommonResponse

	// Possible errors
	Errors []string `json:"errors"`

	// Status of driver and vehicle data entry
	Result string `json:"result"`
}

func (c FBO) GetPass(ctx context.Context, params *GetPassParams) (*GetPassResponse, error) {
	url := "/v1/supply-order/pass/status"

	resp := &GetPassResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetSupplyContentParams struct {
	// Identifiers of supply contents. You can get them using the /v2/supply-order/get method.
	BundleIds []string `json:"bundle_ids"`

	// true, to sort in ascending order
	IsAsc bool `json:"is_asc"`

	// Identifier of the last SKU value on the page.
	LastId string `json:"last_id"`

	// Number of products on the page.
	Limit int32 `json:"limit"`

	// Search query, for example: by name, article code, or SKU
	Query string `json:"query"`

	// Sorting by parameters
	SortField string `json:"sort_field"`

	// Parameters for calculating placement tags.
	ItemTagsCalculation *GetSupplyContentItemTagsCalculation `json:"item_tags_calculation,omitempty"`
}

type GetSupplyContentItemTagsCalculation struct {
	DropoffWarehouseId  string   `json:"dropoff_warehouse_id"`
	StorageWarehouseIds []string `json:"storage_warehouse_ids"`
}

type GetSupplyContentResponse struct {
	core.CommonResponse

	// List of products in the supply request
	Items []SupplyContentItem `json:"items"`

	// Quantity of products in the request
	TotalCount int32 `json:"total_count"`

	// Indication that the response hasn't returned all products
	HasNext bool `json:"has_next"`

	// Identifier of the last value on the page
	LastId string `json:"last_id"`
}

type SupplyContentItem struct {
	// Link to product image
	IconPath string `json:"icon_path"`

	// Product identifier in the Ozon system, SKU
	SKU int64 `json:"sku"`

	// Product name
	Name string `json:"name"`

	// Product items quantity
	Quantity int32 `json:"quantity"`

	// Barcode
	Barcode string `json:"barcode"`

	// Product identifier
	ProductId int64 `json:"product_id"`

	// Quantity of products in one package
	Quant int32 `json:"quant"`

	// true if the quantity of products in one package can be edited
	IsQuantEditable bool `json:"is_quant_editable"`

	// Volume of products in liters
	VolumeInLiters float64 `json:"volume_in_litres"`

	// Volume of all products in liters
	TotalVolumeInLiters float64 `json:"total_volume_in_litres"`

	// Product article code
	ContractorItemCode string `json:"contractor_item_code"`

	// Super product label
	SFBOAttribute string `json:"sfbo_attribute"`

	// Type of wrapper
	ShipmentType string `json:"shipment_type"`

	// Product identifier in the seller's system.
	OfferId string `json:"offer_id"`

	// Calculated product placement tags.
	Tags []string `json:"tags"`

	// Product placement zone.
	PlacementZone string `json:"placement_zone"`
}

func (c FBO) GetSupplyContent(ctx context.Context, params *GetSupplyContentParams) (*GetSupplyContentResponse, error) {
	url := "/v1/supply-order/bundle"

	resp := &GetSupplyContentResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type CreateSupplyDraftParams struct {
	// Clusters identifiers
	ClusterIds []string `json:"cluster_ids"`

	// Shipping point identifier: pick-up point or sorting center. Only for the type = CREATE_TYPE_CROSSDOCK supply type.
	DropoffWarehouseId int64 `json:"drop_off_point_warehouse_id"`

	// Products
	Items []CreateSupplyDraftItem `json:"items"`

	// Supply type
	Type string `json:"type"`
}

type CreateSupplyDraftItem struct {
	// Product quantity
	Quantity int32 `json:"quantity"`

	// Product identifier
	SKU int64 `json:"sku"`
}

type CreateSupplyDraftResponse struct {
	core.CommonResponse

	// Identifier of the supply request draft
	OperationId string `json:"operation_id"`
}

// Create a direct or cross-docking supply request draft and specify the products to supply.
//
// You can leave feedback on this method in the comments section to the discussion in the Ozon for dev community
func (c FBO) CreateSupplyDraft(ctx context.Context, params *CreateSupplyDraftParams) (*CreateSupplyDraftResponse, error) {
	url := "/v1/draft/create"

	resp := &CreateSupplyDraftResponse{}

	response, err := c.client.Request(ctx, http.MethodGet, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetSupplyDraftInfoParams struct {
	// Identifier of the supply request draft
	OperationId string `json:"operation_id"`
}

type GetSupplyDraftInfoResponse struct {
	core.CommonResponse

	// Clusters
	Clusters []SupplyDraftCluster `json:"clusters"`

	// Identifier of the supply request draft
	DraftId int64 `json:"draft_id"`

	// Errors
	Errors []GetSupplyDraftInfoError `json:"errors"`

	// Creation status of the supply request draft
	Status string `json:"status"`
}

type SupplyDraftCluster struct {
	// Cluster identifier
	Id int64 `json:"cluster_id"`

	// Cluster name
	Name string `json:"cluster_name"`

	// Warehouses
	Warehouses []SupplyDraftWarehouse `json:"warehouses"`
}

type SupplyDraftWarehouse struct {
	// Warehouse address
	Address string `json:"address"`

	// Product list bundle
	BundleIds []SupplyDraftWarehouseBundle `json:"bundle_ids"`

	// Warehouse name
	Name string `json:"name"`

	// Bundle of products that don't come in a shipment. Use the parameter in the /v1/supply-order/bundle method to get details.
	RestrictedBundleId string `json:"restricted_bundle_id"`

	// Warehouse availability
	Status SupplyDraftWarehouseStatus `json:"status"`

	// Supply warehouses
	SupplyWarehouse SupplyWarehouse `json:"supply_warehouse"`

	// Warehouse rank in the cluster
	TotalRank int32 `json:"total_rank"`

	// Warehouse rating
	TotalScore float64 `json:"total_score"`

	// Estimated delivery time
	//
	// Nullable
	TravelTimeDays *int64 `json:"travel_time_days"`

	// Warehouse identifier
	Id int64 `json:"warehouse_id"`
}

type SupplyDraftWarehouseBundle struct {
	// Bundle identifier. Use the parameter in the /v1/supply-order/bundle method to get details
	Id string `json:"bundle_id"`

	// Indicates that the UTD is to be passed
	IsDocless bool `json:"is_docless"`
}

type SupplyDraftWarehouseStatus struct {
	// Reason why the warehouse isn't available
	InvalidReason string `json:"invalid_reason"`

	// Warehouse availability
	IsAvailable bool `json:"is_available"`

	// Warehouse status
	State string `json:"state"`
}

type GetSupplyDraftInfoError struct {
	// Possible errors
	Message string `json:"error_message"`

	// Validation errors
	ItemsValidation []GetSupplyDraftInfoValidationError `json:"items_validation"`

	// Unknown clusters identifiers
	UnknownClusterIds []string `json:"unknown_cluster_ids"`
}

type GetSupplyDraftInfoValidationError struct {
	// Error reasons
	Reasons []string `json:"reasons"`

	// Product identifier in the Ozon system, SKU
	SKU int64 `json:"sku"`
}

func (c FBO) GetSupplyDraftInfo(ctx context.Context, params *GetSupplyDraftInfoParams) (*GetSupplyDraftInfoResponse, error) {
	url := "/v1/draft/create/info"

	resp := &GetSupplyDraftInfoResponse{}

	response, err := c.client.Request(ctx, http.MethodGet, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type CreateSupplyFromDraftParams struct {
	// Identifier of the supply request draft
	DraftId int64 `json:"draft_id"`

	// Supply time slot
	Timeslot CreateSupplyFromDraftTimeslot `json:"timeslot"`

	// Shipping warehouse identifier
	WarehouseId int64 `json:"warehouse_id"`
}

type CreateSupplyFromDraftTimeslot struct {
	// Supply time slot start date
	FromInTimezone time.Time `json:"from_in_timezone"`

	// Supply time slot end date
	ToInTimezone time.Time `json:"to_in_timezone"`
}

type CreateSupplyFromDraftResponse struct {
	core.CommonResponse

	// Supply request identifier
	OperationId string `json:"operation_id"`
}

func (c FBO) CreateSupplyFromDraft(ctx context.Context, params *CreateSupplyFromDraftParams) (*CreateSupplyFromDraftResponse, error) {
	url := "/v1/draft/supply/create"

	resp := &CreateSupplyFromDraftResponse{}

	response, err := c.client.Request(ctx, http.MethodGet, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetDraftTimeslotsParams struct {
	// Start date of the available supply time slots period
	DateFrom time.Time `json:"date_from"`

	// End date of the available supply time slots period
	//
	// The maximum period is 28 days from the current date
	DateTo time.Time `json:"date_to"`

	// Identifier of the supply request draft
	DraftId int64 `json:"draft_id"`

	// The warehouses identifiers for which supply time slots are required
	WarehouseIds []string `json:"warehouse_ids"`
}

type GetDraftTimeslotsResponse struct {
	core.CommonResponse

	// Warehouses supply time slots
	DropoffWarehouseTimeslots []DraftTimeslot `json:"drop_off_warehouse_timeslots"`

	// Start date of the necessary period
	RequestedDateFrom time.Time `json:"requested_date_from"`

	// End date of the necessary period
	RequestedDateTo time.Time `json:"requested_date_to"`
}

type DraftTimeslot struct {
	// Current time in the warehouse time zone
	CurrentTimeInTimezone time.Time `json:"current_time_in_timezone"`

	// Supply time slots by dates
	Days []DraftTimeslotDay `json:"days"`

	// Warehouse identifier
	DropoffWarehouseId int64 `json:"drop_off_warehouse_id"`

	// Warehouse time zone
	WarehouseTimezone string `json:"warehouse_timezone"`
}

type DraftTimeslotDay struct {
	// Supply time slots date
	DateInTimezone time.Time `json:"date_in_timezone"`

	// Supply time slots details
	Timeslots []DraftTimeslotDayTimeslot `json:"timeslots"`
}

type DraftTimeslotDayTimeslot struct {
	// Supply time slot start date
	FromInTimezone time.Time `json:"from_in_timezone"`

	// Supply time slot end date
	ToInTimezone time.Time `json:"to_in_timezone"`
}

// Available supply time slots at final shipping warehouses
func (c FBO) GetDraftTimeslots(ctx context.Context, params *GetDraftTimeslotsParams) (*GetDraftTimeslotsResponse, error) {
	url := "/v1/draft/timeslot/info"

	resp := &GetDraftTimeslotsResponse{}

	response, err := c.client.Request(ctx, http.MethodGet, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type CancelSuppyOrderParams struct {
	// Supply request identifier
	OrderId int64 `json:"order_id"`
}

type CancelSuppyOrderResponse struct {
	core.CommonResponse

	// Operation identifier for canceling the request
	OperationId string `json:"operation_id"`
}

// Cancel supply request
func (c FBO) CancelSuppyOrder(ctx context.Context, params *CancelSuppyOrderParams) (*CancelSuppyOrderResponse, error) {
	url := "/v1/supply-order/cancel"

	resp := &CancelSuppyOrderResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type StatusCancelledSupplyOrderParams struct {
	// Operation identifier for canceling the supply request
	OperationId string `json:"operation_id"`
}

type StatusCancelledSupplyOrderResponse struct {
	core.CommonResponse

	// Reason why the supply request can't be canceled
	ErrorReasons []string `json:"error_reasons"`

	// Details on supply request cancellation
	Result StatusCancelledSupplyOrderResult `json:"result"`

	// Cancellation status of the supply request
	Status string `json:"status"`
}

type StatusCancelledSupplyOrderResult struct {
	// true, if supply request is canceled
	IsOrderCancelled bool `json:"is_order_cancelled"`

	// List of canceled supplies
	Supplies []StatusCancelledSupplyOrderSupply `json:"supplies"`
}

type StatusCancelledSupplyOrderSupply struct {
	// Reason why supplies can't be canceled
	ErrorReasons []string `json:"error_reasons"`

	// true, if supply is canceled
	IsSupplyCancelled bool `json:"is_supply_cancelled"`

	// Supply identifier
	SupplyId int64 `json:"supply_id"`
}

// Get status of canceled supply request
func (c FBO) StatusCancelledSupplyOrder(ctx context.Context, params *StatusCancelledSupplyOrderParams) (*StatusCancelledSupplyOrderResponse, error) {
	url := "/v1/supply-order/cancel/status"

	resp := &StatusCancelledSupplyOrderResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}
