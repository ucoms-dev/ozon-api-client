package contractaudit

import (
	"path/filepath"
	"testing"
)

func TestPromotionV2RouteInventory(t *testing.T) {
	operations, err := LoadClientOperations(filepath.Join("..", "..", "ozon"))
	if err != nil {
		t.Fatal(err)
	}
	wanted := map[string]string{
		"/v2/actions/products":                     "ProductsInPromotionV2",
		"/v2/actions/candidates":                   "ProductsAvailableForPromotionV2",
		"/v1/actions/products/update":              "UpdateProducts",
		"/v2/actions/products/deactivate":          "RemoveProductV2",
		"/v2/actions/auto-add/products/list":       "ListAutoAddProductsV2",
		"/v2/actions/auto-add/products/candidates": "ListAutoAddCandidatesV2",
		"/v2/actions/auto-add/products/update":     "UpdateAutoAddProductsV2",
		"/v2/actions/auto-add/products/delete":     "DeleteAutoAddProductsV2",
	}
	for _, operation := range operations {
		name, found := wanted[operation.Path]
		if !found {
			continue
		}
		if operation.Method != "POST" || operation.GoDeprecated || operation.GoName != name {
			t.Errorf("unexpected operation: %+v", operation)
		}
		delete(wanted, operation.Path)
	}
	for path, name := range wanted {
		t.Errorf("missing %s: POST %s", name, path)
	}
}
