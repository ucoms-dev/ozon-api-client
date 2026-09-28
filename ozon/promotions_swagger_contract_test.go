package ozon

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

// This matrix is extracted from the provider schema, not generated from the DTOs.
func TestPromotionSwaggerFieldCoverage(t *testing.T) {
	data, err := os.ReadFile("testdata/promotions-v2/swagger-shapes.json")
	if err != nil {
		t.Fatal(err)
	}
	var contract struct {
		Shapes map[string]map[string]string `json:"shapes"`
	}
	if err := json.Unmarshal(data, &contract); err != nil {
		t.Fatal(err)
	}
	models := []interface{}{
		PromotionMoney{}, PromotionIssue{}, PromotionPriceValidation{}, PromotionWebsitePrices{}, PromotionSchemaPrices{},
		PromotionParticipantV2{}, PromotionCandidateV2{}, AutoAddProductV2{}, AutoAddCandidateV2{},
		ProductsInPromotionV2Response{}, ProductsAvailableForPromotionV2Response{}, UpdatePromotionProductsResponse{}, RemoveProductFromPromotionV2Response{},
		ListAutoAddProductsV2Response{}, ListAutoAddCandidatesV2Response{}, UpdateAutoAddProductsV2Response{}, DeleteAutoAddProductsV2Response{},
		UpdatePromotionProduct{}, UpdateAutoAddProductV2{},
	}
	for _, model := range models {
		typ := reflect.TypeOf(model)
		t.Run(typ.Name(), func(t *testing.T) {
			want, ok := contract.Shapes[typ.Name()]
			if !ok {
				t.Fatal("model absent from schema extract")
			}
			got := promotionJSONFields(typ)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("wire fields = %#v; Swagger fields = %#v", got, want)
			}
		})
	}
	if len(models) != len(contract.Shapes) {
		t.Fatal("not all schema models were checked")
	}
}

func promotionJSONFields(typ reflect.Type) map[string]string {
	result := map[string]string{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" || field.Type.Name() == "CommonResponse" {
			continue
		}
		if field.Anonymous {
			for key, value := range promotionJSONFields(field.Type) {
				result[key] = value
			}
			continue
		}
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "" || name == "-" {
			continue
		}
		ft := field.Type
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		kind := ft.Kind().String()
		switch ft.Kind() {
		case reflect.Struct, reflect.Map:
			kind = "object"
		case reflect.Slice:
			kind = "array"
		case reflect.Bool:
			kind = "boolean"
		case reflect.Float64:
			kind = "number"
		}
		if ft == reflect.TypeOf(json.Number("")) {
			kind = "number"
		}
		result[name] = kind
	}
	return result
}
