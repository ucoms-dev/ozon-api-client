package ozon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	core "github.com/ucoms-dev/ozon-api-client"
)

type Chats struct {
	client *core.Client
}

type ListChatsParams struct {
	// Chats filter
	Filter *ListChatsFilter `json:"filter,omitempty"`

	// Number of values in the response. The default value is 30. The maximum value is 1000
	Limit int64 `json:"limit" default:"30"`

	// Number of elements that will be skipped in the response.
	// For example, if offset=10, the response will start with the 11th element found
	Offset int64 `json:"offset,omitempty"`
}

type ListChatsFilter struct {
	// Filter by chat status:
	//   - All
	//   - Opened
	//   - Closed
	ChatStatus string `json:"chat_status" default:"ALL"`

	// Filter by chats with unread messages
	UnreadOnly bool `json:"unread_only"`
}

type ListChatsResponse struct {
	core.CommonResponse

	// Chats data
	Chats []ListChatsChatData `json:"chats"`

	// Total number of chats
	TotalChatsCount int64 `json:"total_chats_count"`

	// Total number of unread messages
	TotalUnreadCount int64 `json:"total_unread_count"`
}

type ListChatsChatData struct {
	// Chat identifier
	ChatId string `json:"chat_id"`

	// Chat status:
	//   - All
	//   - Opened
	//   - Closed
	ChatStatus string `json:"chat_status"`

	// Chat type:
	//   - Seller_Support — support chat
	//   - Buyer_Seller — chat with a customer
	ChatType string `json:"chat_type"`

	// Chat creation date
	CreatedAt time.Time `json:"created_at"`

	// Identifier of the first unread chat message
	FirstUnreadMessageId uint64 `json:"first_unread_message_id"`

	// Identifier of the last message in the chat
	LastMessageId uint64 `json:"last_message_id"`

	// Number of unread messages in the chat
	UnreadCount int64 `json:"unread_count"`
}

// Returns information about chats by specified filters.
//
// Deprecated: use ListV3 for the current cursor-based chat contract.
func (c Chats) List(ctx context.Context, params *ListChatsParams) (*ListChatsResponse, error) {
	url := "/v2/chat/list"

	resp := &ListChatsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type ListChatsV3Params struct {
	Filter *ListChatsV3Filter `json:"filter,omitempty"`
	// Page size, up to 100. Pagination remains the caller's responsibility.
	Limit  int64  `json:"limit" default:"30"`
	Cursor string `json:"cursor,omitempty"`
}

type ListChatsV3Filter struct {
	ChatStatus string `json:"chat_status" default:"ALL"`
	UnreadOnly bool   `json:"unread_only"`
}

// ChatMessageID preserves provider message IDs without floating-point conversion.
// Ozon documents uint64 integers, but examples contain decimal strings including
// values larger than uint64. The opaque decimal identifier is retained exactly.
type ChatMessageID string

func (id *ChatMessageID) UnmarshalJSON(data []byte) error {
	value := bytes.TrimSpace(data)
	if bytes.Equal(value, []byte("null")) {
		*id = ""
		return nil
	}
	var text string
	if len(value) > 0 && value[0] == '"' {
		if err := json.Unmarshal(value, &text); err != nil {
			return err
		}
	} else {
		text = string(value)
	}
	for _, digit := range text {
		if digit < '0' || digit > '9' {
			return fmt.Errorf("invalid chat message identifier")
		}
	}
	if text == "" {
		*id = ""
		return nil
	}
	*id = ChatMessageID(text)
	return nil
}

type ListChatsV3ChatDetails struct {
	// CreatedAt retains the provider date verbatim, including legacy date formats.
	CreatedAt  string `json:"created_at"`
	ChatId     string `json:"chat_id"`
	ChatStatus string `json:"chat_status"`
	ChatType   string `json:"chat_type"`
}

type ListChatsV3Chat struct {
	Chat                 ListChatsV3ChatDetails `json:"chat"`
	FirstUnreadMessageId ChatMessageID          `json:"first_unread_message_id"`
	LastMessageId        ChatMessageID          `json:"last_message_id"`
	UnreadCount          int64                  `json:"unread_count"`
}

type ListChatsV3Response struct {
	core.CommonResponse
	Chats            []ListChatsV3Chat `json:"chats"`
	TotalUnreadCount int64             `json:"total_unread_count"`
	Cursor           string            `json:"cursor"`
	HasNext          bool              `json:"has_next"`
	rawJSON          []byte
}

// RawResponseJSON returns an independent snapshot of the provider body, including
// fields outside the typed contract, for consumers with established richer mappings.
func (response ListChatsV3Response) RawResponseJSON() []byte {
	return append([]byte(nil), response.rawJSON...)
}

func (response *ListChatsV3Response) UnmarshalJSON(data []byte) error {
	type wire ListChatsV3Response
	var decoded struct {
		*wire
		HasNext json.RawMessage `json:"has_next"`
	}
	value := wire{}
	decoded.wire = &value
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	flag, err := decodeSellerPaginationFlag(decoded.HasNext)
	if err != nil {
		return err
	}
	value.HasNext = flag
	value.rawJSON = append([]byte(nil), data...)
	*response = ListChatsV3Response(value)
	return nil
}

// decodeSellerPaginationFlag handles the boolean schema and the string examples
// published for the current chat and warehouse cursor pagination contracts.
func decodeSellerPaginationFlag(data []byte) (bool, error) {
	if len(data) == 0 || bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return false, nil
	}
	var flag bool
	if err := json.Unmarshal(data, &flag); err == nil {
		return flag, nil
	}
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		switch strings.ToLower(strings.TrimSpace(text)) {
		case "true", "1":
			return true, nil
		case "false", "0":
			return false, nil
		}
	}
	return false, fmt.Errorf("invalid seller pagination has_next")
}

// ListV3 reads one page from the current chat API. It does not switch the legacy
// List route or automatically fetch subsequent pages.
func (c Chats) ListV3(ctx context.Context, params *ListChatsV3Params) (*ListChatsV3Response, error) {
	url := "/v3/chat/list"
	resp := &ListChatsV3Response{}
	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)
	return resp, nil
}

type SendMessageParams struct {
	// Chat identifier
	ChatId string `json:"chat_id"`

	// Message text in the plain text format
	Text string `json:"text"`
}

type SendMessageResponse struct {
	core.CommonResponse

	// Method result
	Result string `json:"result"`
}

// Sends a message to an existing chat by its identifier
func (c Chats) SendMessage(ctx context.Context, params *SendMessageParams) (*SendMessageResponse, error) {
	url := "/v1/chat/send/message"

	resp := &SendMessageResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type SendFileParams struct {
	// File as a base64 string
	Base64Content string `json:"base64_content"`

	// Chat identifier
	ChatId string `json:"chat_id"`

	// File name with extension
	Name string `json:"name"`
}

type SendFileResponse struct {
	core.CommonResponse

	// Method result
	Result string `json:"result"`
}

// Sends a file to an existing chat by its identifier
func (c Chats) SendFile(ctx context.Context, params *SendFileParams) (*SendFileResponse, error) {
	url := "/v1/chat/send/file"

	resp := &SendFileResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type ChatHistoryParams struct {
	// Chat idenitifier
	ChatId string `json:"chat_id"`

	// Messages sorting direction:
	//   - Forward—from old messages to new ones.
	//   - Backward—from new messages to old ones.
	// The default value is `Backward`. You can set the number of messages in the limit parameter
	Direction string `json:"direction" default:"Backward"`

	Filter *ChatHistoryFilter `json:"filter,omitempty"`

	// Identifier of the message from which the chat history will be displayed.
	// Default value is the last visible message
	FromMessageId string `json:"from_message_id"`

	// Number of messages in the response. The default value is 50. The maximum value is 1000
	Limit int64 `json:"limit" default:"50"`
}

type ChatHistoryFilter struct {
	MessageIds []string `json:"message_ids"`
}

type ChatHistoryResponse struct {
	core.CommonResponse

	// Indicates that the response returned only a part of messages
	HasNext bool `json:"has_next"`

	// An array of messages sorted according to the direction parameter in the request body
	Messages []ChatHistoryMessage `json:"messages"`
}

type ChatHistoryMessage struct {
	Context *ChatHistoryContext `json:"context,omitempty"`

	// Message creation date
	CreatedAt time.Time `json:"created_at"`

	// Array with message content in Markdown format
	Data []string `json:"data"`

	IsImage bool `json:"is_image"`

	// Indication of the read message
	IsRead bool `json:"is_read"`

	// Message identifier
	MessageId string `json:"message_id"`

	ModarateImageStatus string `json:"moderate_image_status"`

	// Chat participant identifier
	User ChatHistoryMessageUser `json:"user"`
}

type ChatHistoryContext struct {
	OrderNumber string `json:"order_number"`
	SKU         string `json:"sku"`
}

type ChatHistoryMessageUser struct {
	// Chat participant identifier
	Id string `json:"id"`

	// Chat participant type:
	//   - customer
	//   - seller
	//   - crm—system messages
	//   - courier
	//   - support
	Type string `json:"type"`
}

// Returns the history of chat messages. By default messages are shown from newest to oldest.
func (c Chats) History(ctx context.Context, params *ChatHistoryParams) (*ChatHistoryResponse, error) {
	url := "/v3/chat/history"

	resp := &ChatHistoryResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type UpdateChatParams struct {
	// Chat identifier
	ChatId string `json:"chat_id"`

	// Message identifier
	FromMessageId uint64 `json:"from_message_id"`

	// Number of messages in the response
	Limit int64 `json:"limit,omitempty"`
}

type UpdateChatResponse struct {
	core.CommonResponse

	// Method result
	Result []UpdateChatResult `json:"result"`
}

type UpdateChatResult struct {
	// An order or a product user wrote about in the chat
	Context UpdateChatResultContext `json:"context"`

	// Creation date and time
	CreatedAt time.Time `json:"created_at"`

	// Information about the file in the chat. Displayed only for `type = file`
	File UpdateChatResultFile `json:"file"`

	// File identifier
	Id uint64 `json:"id"`

	// Message. Displayed only for `type = text`
	Text string `json:"text"`

	// Message type:
	//   - text
	//   - file
	Type string `json:"type"`

	// Chat participant information
	User UpdateChatResultUser `json:"user"`
}

type UpdateChatResultContext struct {
	// Product inforamtion
	Item UpdateChatResultContextItem `json:"item"`

	// Order information
	Order UpdateChatResultContextOrder `json:"order"`
}

type UpdateChatResultContextItem struct {
	// Product identifier in the Ozon system, SKU
	SKU int64 `json:"sku"`
}

type UpdateChatResultContextOrder struct {
	// Order number
	OrderNumber string `json:"order_number"`

	// Shipment information
	Postings []UpdateChatResultContextOrderPosting `json:"postings"`
}

type UpdateChatResultContextOrderPosting struct {
	// Delivery scheme:
	//   - FBO
	//   - FBS
	//   - RFBS
	//   - Crossborder
	DeliverySchema string `json:"delivery_schema"`

	// Shipment number
	PostingNumber string `json:"posting_number"`

	// List of product identifiers in the shipment
	SKUList []int64 `json:"sku_list"`
}

type UpdateChatResultFile struct {
	// File type
	Mime string `json:"mime"`

	// File name
	Name string `json:"name"`

	// File size in bytes
	Size int64 `json:"size"`

	// File URL
	URL string `json:"url"`
}

type UpdateChatResultUser struct {
	// Chat participant identifier
	Id string `json:"id"`

	// Chat participant chat:
	//   - customer
	//   - seller
	//   - crm—system messages
	//   - courier
	Type string `json:"type"`
}

// Update chat
func (c Chats) Update(ctx context.Context, params *UpdateChatParams) (*UpdateChatResponse, error) {
	url := "/v1/chat/updates"

	resp := &UpdateChatResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type CreateNewChatParams struct {
	// Shipment identifier
	PostingNumber string `json:"posting_number"`
}

type CreateNewChatResponse struct {
	core.CommonResponse

	//Method result
	Result CreateNewChatResult `json:"result"`
}

type CreateNewChatResult struct {
	// Chat identifier
	ChatId string `json:"chat_id"`
}

// Creates a new chat on the shipment with the customer. For example, to clarify the address or the product model
func (c Chats) Create(ctx context.Context, params *CreateNewChatParams) (*CreateNewChatResponse, error) {
	url := "/v1/chat/start"

	resp := &CreateNewChatResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type MarkAsReadParams struct {
	// Chat identifier
	ChatId string `json:"chat_id"`

	// Message identifier
	FromMessageId uint64 `json:"from_message_id"`
}

type MarkAsReadResponse struct {
	core.CommonResponse

	// Number of unread messages in the chat
	UnreadCount int64 `json:"unread_count"`
}

// A method for marking the selected message and messages before it as read
func (c Chats) MarkAsRead(ctx context.Context, params *MarkAsReadParams) (*MarkAsReadResponse, error) {
	url := "/v2/chat/read"

	resp := &MarkAsReadResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}
