package ozon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// PromotionMoney is an exact decimal amount and its provider-supplied currency.
// Amount remains a string so values do not lose precision through float64.
type PromotionMoney struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

func (money PromotionMoney) MarshalJSON() ([]byte, error) {
	if err := validatePromotionMoney("", "", money); err != nil {
		return nil, err
	}
	type wire PromotionMoney
	return json.Marshal(wire(money))
}

func (money *PromotionMoney) UnmarshalJSON(data []byte) error {
	type wire PromotionMoney
	var decoded wire
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	value := PromotionMoney(decoded)
	if !validPromotionMoney(value) {
		return errors.New("invalid promotion money")
	}
	*money = value
	return nil
}

// PromotionProductID preserves uint64 IDs received as JSON integers or strings.
// It intentionally does not enforce positivity; callers validate that rule at
// the request or response boundary where the value represents a product.
type PromotionProductID uint64

func (id PromotionProductID) MarshalJSON() ([]byte, error) {
	return strconv.AppendUint(nil, uint64(id), 10), nil
}

func (id *PromotionProductID) UnmarshalJSON(data []byte) error {
	value := bytes.TrimSpace(data)
	if len(value) == 0 {
		return errors.New("invalid promotion product ID")
	}

	var digits []byte
	if value[0] == '"' {
		var text string
		if err := json.Unmarshal(value, &text); err != nil {
			return errors.New("invalid promotion product ID")
		}
		text = strings.TrimSpace(text)
		if text == "" {
			return errors.New("invalid promotion product ID")
		}
		digits = []byte(text)
	} else {
		digits = value
	}
	for _, digit := range digits {
		if digit < '0' || digit > '9' {
			return errors.New("invalid promotion product ID")
		}
	}
	parsed, err := strconv.ParseUint(string(digits), 10, 64)
	if err != nil {
		return errors.New("invalid promotion product ID")
	}
	*id = PromotionProductID(parsed)
	return nil
}

// PromotionIssue is a per-product rejection or warning returned by Ozon.
type PromotionIssue struct {
	ProductID PromotionProductID `json:"product_id"`
	Reason    string             `json:"reason"`
}

// PromotionPriceValidation is an exact numeric diagnostic from auto-add.
type PromotionPriceValidation struct {
	Key   PromotionProductID `json:"key"`
	Value json.Number        `json:"value"`
}

// PromotionWebsitePrices contains prices for the default and named delivery
// schemas. Schema names remain open strings because Ozon may add new values.
type PromotionWebsitePrices struct {
	Price          *PromotionMoney                  `json:"price,omitempty"`
	PricesBySchema map[string]PromotionSchemaPrices `json:"prices_by_schema,omitempty"`
}

// PromotionSchemaPrices contains prices for one delivery schema.
type PromotionSchemaPrices struct {
	BlackPrice *PromotionMoney `json:"black_price,omitempty"`
	GreenPrice *PromotionMoney `json:"green_price,omitempty"`
}

// PromotionInputError reports a rejected caller-supplied value before a request.
type PromotionInputError struct {
	Method string
	Field  string
	Reason string
}

func (err *PromotionInputError) Error() string {
	if err == nil {
		return "invalid promotion input"
	}
	if err.Method != "" && err.Field != "" {
		return fmt.Sprintf("%s: invalid %s: %s", err.Method, err.Field, err.Reason)
	}
	if err.Method != "" {
		return fmt.Sprintf("%s: invalid input: %s", err.Method, err.Reason)
	}
	return "invalid promotion input: " + err.Reason
}

// PromotionProtocolError reports an unexpected Ozon success response shape.
// Cause is available through Unwrap and is kept out of Error text so response
// bodies and decoder details are not accidentally exposed in diagnostics.
type PromotionProtocolError struct {
	Method string
	Reason string
	cause  error
}

func (err *PromotionProtocolError) Error() string {
	if err == nil {
		return "invalid promotion response"
	}
	if err.Method != "" {
		return fmt.Sprintf("%s: invalid response: %s", err.Method, err.Reason)
	}
	return "invalid promotion response: " + err.Reason
}

func (err *PromotionProtocolError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.cause
}

func promotionInputError(method, field, reason string) error {
	return &PromotionInputError{Method: method, Field: field, Reason: reason}
}

func promotionProtocolError(method, reason string) error {
	return &PromotionProtocolError{Method: method, Reason: reason}
}

func promotionProtocolErrorWithCause(method, reason string, cause error) error {
	return &PromotionProtocolError{Method: method, Reason: reason, cause: cause}
}

func validatePromotionMoney(method, field string, money PromotionMoney) error {
	if !validPromotionMoney(money) {
		return promotionInputError(method, field, "amount must be a non-negative decimal string and currency must be non-empty")
	}
	return nil
}

func validPromotionMoney(money PromotionMoney) bool {
	if money.Amount == "" || money.Currency == "" || strings.TrimSpace(money.Currency) != money.Currency {
		return false
	}
	decimalPointSeen := false
	integerDigits := 0
	fractionDigits := 0
	for _, char := range money.Amount {
		switch {
		case char >= '0' && char <= '9':
			if decimalPointSeen {
				fractionDigits++
			} else {
				integerDigits++
			}
		case char == '.' && !decimalPointSeen:
			decimalPointSeen = true
		default:
			return false
		}
	}
	return integerDigits > 0 && (!decimalPointSeen || fractionDigits > 0)
}

func validatePromotionActionID(method string, actionID uint64) error {
	if actionID == 0 {
		return promotionInputError(method, "action_id", "must be greater than zero")
	}
	return nil
}

func validatePromotionProductID(method, field string, id PromotionProductID) error {
	if id == 0 {
		return promotionInputError(method, field, "must be greater than zero")
	}
	return nil
}

func validatePromotionProductIDs(method, field string, ids []PromotionProductID, max int) error {
	if len(ids) == 0 || len(ids) > max {
		return promotionInputError(method, field, fmt.Sprintf("must contain between 1 and %d product IDs", max))
	}
	seen := make(map[PromotionProductID]struct{}, len(ids))
	for _, id := range ids {
		if err := validatePromotionProductID(method, field, id); err != nil {
			return err
		}
		if _, exists := seen[id]; exists {
			return promotionInputError(method, field, "must not contain duplicate product IDs")
		}
		seen[id] = struct{}{}
	}
	return nil
}

func validatePromotionLimit(method string, limit uint64) error {
	if limit < 1 || limit > 100 {
		return promotionInputError(method, "limit", "must be between 1 and 100")
	}
	return nil
}

func validatePromotionCall(method string, ctx context.Context, clientPresent bool) error {
	if ctx == nil {
		return promotionInputError(method, "context", "must not be nil")
	}
	if !clientPresent {
		return promotionInputError(method, "client", "must be configured")
	}
	return nil
}

func validatePromotionResultIDs(method, field string, ids []PromotionProductID) error {
	for _, id := range ids {
		if id == 0 {
			return promotionProtocolError(method, "field "+field+" contains a zero product ID")
		}
	}
	return nil
}

func promotionWireStock(stock *uint64) interface{} {
	if stock == nil {
		return nil
	}
	return *stock
}

type promotionEnvelopeFieldKind uint8

const (
	promotionEnvelopeArray promotionEnvelopeFieldKind = iota
	promotionEnvelopeUint64
	promotionEnvelopeNumber
	promotionEnvelopeString
)

type promotionEnvelopeField struct {
	name string
	kind promotionEnvelopeFieldKind
}

func validatePromotionEnvelope(
	method string,
	data []byte,
	required []promotionEnvelopeField,
	oneOf []promotionEnvelopeField,
) (map[string]json.RawMessage, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, promotionProtocolError(method, "response body is empty")
	}
	if trimmed[0] != '{' {
		return nil, promotionProtocolError(method, "response must be a JSON object")
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(trimmed, &fields); err != nil {
		return nil, promotionProtocolErrorWithCause(method, "response is not valid JSON", err)
	}
	if fields == nil {
		return nil, promotionProtocolError(method, "response must be a JSON object")
	}

	for _, field := range required {
		value, exists := fields[field.name]
		if !exists {
			return nil, promotionProtocolError(method, "missing required field "+field.name)
		}
		if err := validatePromotionEnvelopeField(method, field, value); err != nil {
			return nil, err
		}
	}

	validOneOf := len(oneOf) == 0
	for _, field := range oneOf {
		value, exists := fields[field.name]
		if !exists {
			continue
		}
		if err := validatePromotionEnvelopeField(method, field, value); err != nil {
			return nil, err
		}
		validOneOf = true
	}
	if !validOneOf {
		return nil, promotionProtocolError(method, "response contains no recognized result field")
	}
	return fields, nil
}

func validatePromotionEnvelopeField(method string, field promotionEnvelopeField, value json.RawMessage) error {
	trimmed := bytes.TrimSpace(value)
	valid := false
	switch field.kind {
	case promotionEnvelopeArray:
		valid = len(trimmed) > 0 && trimmed[0] == '['
	case promotionEnvelopeUint64:
		var number uint64
		valid = len(trimmed) > 0 && trimmed[0] != 'n' && json.Unmarshal(trimmed, &number) == nil
	case promotionEnvelopeNumber:
		decoder := json.NewDecoder(bytes.NewReader(trimmed))
		decoder.UseNumber()
		var decoded interface{}
		if err := decoder.Decode(&decoded); err == nil {
			_, valid = decoded.(json.Number)
		}
	case promotionEnvelopeString:
		var decoded string
		valid = len(trimmed) > 0 && trimmed[0] == '"' && json.Unmarshal(trimmed, &decoded) == nil
	}
	if valid {
		return nil
	}
	return promotionProtocolError(method, "field "+field.name+" has an invalid type")
}
