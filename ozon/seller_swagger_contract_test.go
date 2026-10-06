package ozon

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

type sellerSwaggerShape struct {
	Type       string                        `json:"type"`
	Format     string                        `json:"format"`
	Ref        string                        `json:"$ref"`
	Items      *sellerSwaggerShape           `json:"items"`
	Properties map[string]sellerSwaggerShape `json:"properties"`
}

// The fixture comes from the official Swagger captured in Brave, independently
// of the Go DTOs. Check nested field coverage and wire types as well as routes.
func TestSellerCurrentSwaggerFieldCoverage(t *testing.T) {
	data, err := os.ReadFile("testdata/seller-2026-10-06/swagger-shapes.json")
	if err != nil {
		t.Fatal(err)
	}
	var snapshot struct {
		Schemas map[string]sellerSwaggerShape `json:"schemas"`
	}
	if err := json.Unmarshal(data, &snapshot); err != nil {
		t.Fatal(err)
	}
	models := map[string]interface{}{
		"product.v2.ProductImportPicturesV2Request":         UpdateProductImagesV2Params{},
		"product.v2.ProductImportPicturesV2Request.Item":    UpdateProductImagesV2Item{},
		"product.v2.ProductInfoPicturesV2Response":          UpdateProductImagesV2Response{},
		"v4GetUploadQuotaResponse":                          GetProductRangeLimitResponse{},
		"product.v4.GetUploadQuotaResponse.OperationLimits": GetProductRangeLimitOperationLimit{},
		"v2WarehouseListV2Request":                          GetListOfWarehousesV2Params{},
		"v2WarehouseListV2Response":                         GetListOfWarehousesV2Response{},
		"WarehouseListV2ResponseWarehouse":                  GetListOfWarehousesV2Warehouse{},
		"WarehouseAddressInfo":                              GetListOfWarehousesV2AddressInfo{},
		"WarehouseFirstMile":                                GetListOfWarehousesV2FirstMile{},
		"WarehouseTimetable":                                GetListOfWarehousesV2Timetable{},
		"TimetableWorkingHours":                             GetListOfWarehousesV2WorkingHours{},
		"GetUploadQuotaResponseDailyCreate":                 GetProductRangeLimitUploadQuota{},
		"GetUploadQuotaResponseDailyUpdate":                 GetProductRangeLimitUploadQuota{},
		"GetUploadQuotaResponseTotal":                       GetProductRangeLimitTotal{},
	}
	for name, schema := range snapshot.Schemas {
		t.Run(name, func(t *testing.T) {
			model, ok := models[name]
			if !ok {
				t.Fatalf("unmapped schema %s", name)
			}
			typ := reflect.TypeOf(model)
			fields := map[string]reflect.Type{}
			for i := 0; i < typ.NumField(); i++ {
				f := typ.Field(i)
				if f.Anonymous {
					continue
				}
				key := strings.Split(f.Tag.Get("json"), ",")[0]
				if key != "" && key != "-" {
					fields[key] = f.Type
				}
			}
			if len(fields) != len(schema.Properties) {
				t.Fatalf("DTO has %d fields; schema has %d", len(fields), len(schema.Properties))
			}
			for key, prop := range schema.Properties {
				ft, ok := fields[key]
				if !ok {
					t.Errorf("missing field %s", key)
					continue
				}
				if key == "operation_limits" {
					// Explicit provider discrepancy: documented object, observed array.
					ft = ft.Elem()
				}
				checkSellerSwaggerType(t, key, ft, prop, models)
			}
		})
	}
}

func checkSellerSwaggerType(t *testing.T, field string, typ reflect.Type, shape sellerSwaggerShape, models map[string]interface{}) {
	t.Helper()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if shape.Ref != "" {
		ref := strings.TrimPrefix(shape.Ref, "#/components/schemas/")
		if model, ok := models[ref]; ok {
			if typ != reflect.TypeOf(model) {
				t.Errorf("%s: %v does not match %s", field, typ, ref)
			}
		} else if typ.Kind() != reflect.String {
			t.Errorf("%s: enum %s must remain a string", field, ref)
		}
		return
	}
	switch shape.Type {
	case "array":
		if typ.Kind() != reflect.Slice {
			t.Errorf("%s: expected slice, got %v", field, typ)
			return
		}
		checkSellerSwaggerType(t, field+"[]", typ.Elem(), *shape.Items, models)
	case "string":
		if shape.Format == "date-time" {
			if typ != reflect.TypeOf(time.Time{}) {
				t.Errorf("%s: expected time.Time, got %v", field, typ)
			}
		} else if typ.Kind() != reflect.String {
			t.Errorf("%s: expected string, got %v", field, typ)
		}
	case "integer":
		if typ.Kind() != reflect.Int64 && typ.Kind() != reflect.Int32 && typ.Kind() != reflect.Int {
			t.Errorf("%s: expected integer, got %v", field, typ)
		}
		if shape.Format == "int64" && typ.Kind() != reflect.Int64 {
			t.Errorf("%s: expected int64, got %v", field, typ)
		}
		if shape.Format == "int32" && typ.Kind() != reflect.Int32 {
			t.Errorf("%s: expected int32, got %v", field, typ)
		}
	case "number":
		if typ.Kind() != reflect.Float64 {
			t.Errorf("%s: expected float64, got %v", field, typ)
		}
	case "boolean":
		if typ.Kind() != reflect.Bool {
			t.Errorf("%s: expected bool, got %v", field, typ)
		}
	default:
		t.Fatalf("%s: unsupported schema type %s", field, shape.Type)
	}
}
