package ozon_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ucoms-dev/ozon-api-client/ozon"
)

func ExamplePromotions_ProductsInPromotionV2() {
	// A local fixture transport keeps this example independent of seller credentials.
	client := ozon.NewMockClient(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"products":[{"id":42,"price":{"amount":"1250.50","currency":"RUB"}}],"total":2,"last_id":"next-page"}`)
	})
	page, err := client.Promotions().ProductsInPromotionV2(context.Background(), &ozon.ProductsInPromotionV2Params{ActionID: 10, Limit: 100})
	if err != nil {
		panic(err)
	}
	if page.StatusCode != http.StatusOK {
		panic(page.Message)
	}
	// Pass LastID unchanged to the next explicit call; the SDK does not fetch it automatically.
	fmt.Println(page.Products[0].Price.Amount, page.Products[0].Price.Currency, page.LastID)
	// Output: 1250.50 RUB next-page
}

func ExamplePromotions_UpdateProducts() {
	client := ozon.NewMockClient(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"active_product_ids":["42"],"deactivated_product_ids":[],"rejected":[{"product_id":43,"reason":"price rejected"}],"warnings":[]}`)
	})
	result, err := client.Promotions().UpdateProducts(context.Background(), &ozon.UpdatePromotionProductsParams{
		ActionID: 10,
		Products: []ozon.UpdatePromotionProduct{
			{ProductID: 42, ActionPrice: ozon.PromotionMoney{Amount: "1250.50", Currency: "RUB"}},
			{ProductID: 43, ActionPrice: ozon.PromotionMoney{Amount: "1500", Currency: "RUB"}},
		},
	})
	if err != nil {
		panic(err)
	}
	if result.StatusCode != http.StatusOK {
		panic(result.Message)
	}
	// HTTP 200 is not confirmation for every input product. Read back the resulting state.
	fmt.Printf("active=%d removed=%d rejected=%d warnings=%d\n", len(result.ActiveProductIDs), len(result.DeactivatedProductIDs), len(result.Rejected), len(result.Warnings))
	// Output: active=1 removed=0 rejected=1 warnings=0
}

func ExamplePromotions_ListAutoAddProductsV2() {
	client := ozon.NewMockClient(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"products":[{"product_id":42,"sku":777,"add_mode":true}],"total":1}`)
	})
	// In production this value comes from GetAvailablePromotions().Result[n].AutoAddDates.
	date := time.Date(2026, 10, 13, 0, 0, 0, 0, time.UTC)
	page, err := client.Promotions().ListAutoAddProductsV2(context.Background(), &ozon.ListAutoAddProductsV2Params{ActionID: 10, AutoAddDate: date, Limit: 100, Offset: 0})
	if err != nil {
		panic(err)
	}
	if page.StatusCode != http.StatusOK {
		panic(page.Message)
	}
	// Auto-add uses offset pagination, not LastID; a scheduled row is not proof of future execution.
	fmt.Printf("scheduled=%d total=%d\n", len(page.Products), page.Total)
	// Output: scheduled=1 total=1
}
