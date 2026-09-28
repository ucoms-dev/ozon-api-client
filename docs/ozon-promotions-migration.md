# Migrating to the Ozon promotions SDK methods

This guide covers the eight seller-action methods added for the September 2026 Ozon Seller API update. The scope is seller promotions (seller-actions); customer discount requests (discounts-task) are a separate API family and are not migrated here.

The contracts were reviewed on 2026-09-28 against Ozon's announcement, operation pages, and the public Seller API OpenAPI schema. Schema inspection does not prove that a live Seller API response has been tested. No authenticated live method response is claimed here.

## Method transition

| Existing SDK method | Replacement | Migration note |
|---|---|---|
| ProductsAvailableForPromotion (POST /v1/actions/candidates) | ProductsAvailableForPromotionV2 (POST /v2/actions/candidates) | New money and response shapes; pagination uses last_id rather than offset. |
| ProductsInPromotion (POST /v1/actions/products) | ProductsInPromotionV2 (POST /v2/actions/products) | Legacy offset pagination becomes last_id cursor paging; response fields move to the top level. |
| AddToPromotion (POST /v1/actions/products/activate) | UpdateProducts (POST /v1/actions/products/update) | This is a contract change, not a drop-in call. The new operation can report active and deactivated IDs, warnings, and rejected items. |
| RemoveProduct (POST /v1/actions/products/deactivate) | RemoveProductV2 for promo-code actions; UpdateProducts for price-based participation | Choose the operation for the action's participation rules. The response shape and product-ID encoding differ from the legacy method. |
| GetAvailablePromotions (GET /v1/actions) | No replacement | It remains the source for available auto_add_dates; the SDK adds that optional field. |
| — | ListAutoAddProductsV2, ListAutoAddCandidatesV2 | New reads require an explicit date from auto_add_dates. |
| — | UpdateAutoAddProductsV2, DeleteAutoAddProductsV2 | New scheduled auto-add operations; a successful HTTP response does not prove a later scheduled change was applied. |

The four legacy methods remain in the SDK for source compatibility, but Ozon announced that their corresponding old endpoints are disabled on **2026-10-13**. The new methods use fixed endpoints and do not switch behavior according to the SDK machine's clock. The announcement says the new operations retain the previous business logic until that date and use the updated behavior after it; that does not make the old and new JSON contracts interchangeable.

## New method and wire summary

All eight methods send one POST request per call. Their response embeds CommonResponse and endpoint fields at the top level; there is no result wrapper.

| Method | Endpoint | Request fields | Endpoint response fields |
|---|---|---|---|
| ProductsAvailableForPromotionV2 | /v2/actions/candidates | action_id, limit, last_id | products, total, last_id |
| ProductsInPromotionV2 | /v2/actions/products | action_id, limit, last_id | products, total, last_id |
| UpdateProducts | /v1/actions/products/update | action_id, products[{product_id, action_price, stock?}] | active_product_ids, deactivated_product_ids, rejected, warnings |
| RemoveProductV2 | /v2/actions/products/deactivate | action_id, product_ids | product_ids |
| ListAutoAddProductsV2 | /v2/actions/auto-add/products/list | action_id, auto_add_date, limit, offset | products, total |
| ListAutoAddCandidatesV2 | /v2/actions/auto-add/products/candidates | action_id, auto_add_date, limit, offset | products, total |
| UpdateAutoAddProductsV2 | /v2/actions/auto-add/products/update | action_id, auto_add_date, products[{id, action_price, stock?}] | product_ids, deactivated_ids, rejected, warnings, below_min_price, extremely_low_price, failed_price |
| DeleteAutoAddProductsV2 | /v2/actions/auto-add/products/delete | action_id, auto_add_date, product_ids | product_ids |

Use the Go method signatures and types in the SDK as the source of exact field types. The [contract matrix](ozon-promotions-contracts-2026-09-28.md) lists every field and links each operation to Ozon's documentation.

## IDs, prices, and absent values

- PromotionMoney carries amount as a decimal string and currency as supplied. It preserves decimal precision; the SDK does not convert prices through float64, impose a two-decimal scale, or assume RUB. A request amount must be a valid non-negative decimal without exponent notation, NaN, or infinity. The wire boundary permits zero; Ozon can still reject a business-ineligible price. An absent response price remains nil.
- Product IDs use uint64 without a float64 round trip. product_id, id, offer_id, and sku are distinct fields: do not use an offer ID or SKU in place of the Ozon product identifier.
- Update request items use numeric product_id (UpdateProducts) or numeric id (UpdateAutoAddProductsV2). Deactivate and delete requests encode product_ids as decimal strings, matching the OpenAPI string<uint64> array. Responses accept either numeric or string IDs because the reviewed Ozon examples and schema do not use one representation consistently.
- stock is optional. A nil Go pointer omits the JSON field; a pointer to zero sends "stock": 0. Do not turn absence into zero. The same principle applies to optional money and boolean response fields.
- A regular promotion participant's add_mode is a string (SELLER or AUTO); preserve unrecognized strings. In auto-add, product add_mode is a boolean, while candidate is_manually_added is a boolean. These fields have different meanings and must not share a string enum.

## Validation, paging, and partial success

The SDK rejects invalid local input before making the request: action IDs and supplied product IDs must be positive; list limits are 1–100; mutation batches contain 1–1000 unique positive product IDs; and auto-add calls require a non-zero auto_add_date. These are SDK safeguards and can be stricter than constraints marked required in the published OpenAPI schema. See the contract notes for exact operation-specific rules. The SDK does not chunk oversized writes.

Both ordinary promotion lists use cursor paging: pass the returned last_id to fetch the next page. Auto-add lists use offset. Each call reads only one page; the SDK does not walk pages automatically.

For HTTP 200, list responses must have a recognized, non-null products array and total; an empty array with total 0 is valid. Mutation responses must contain at least one documented result field with the correct type. Deactivate and delete responses require product_ids. Unknown extra response fields are allowed.

rejected and warnings are per-item result categories, not transport failures. The SDK keeps confirmed, deactivated, rejected, warning, and price-diagnostic IDs separate. A partial confirmed-ID list is returned as received; the SDK never fills it with IDs from the request. Empty result arrays do not mean that every requested item succeeded.

Valid Ozon 4xx/5xx response bodies are returned through CommonResponse (StatusCode, Code, Message, and Details). A Go error indicates a transport or context failure, malformed JSON, invalid caller input, or a malformed HTTP 200 response. A timeout or network failure during a write may leave the outcome unknown. The SDK makes no automatic retry; callers should read back state and reconcile before deciding whether to send another write.

## Dates and the 2026-10-13 behavior change

GetAvailablePromotions returns optional auto_add_dates. Choose an available value and pass it as AutoAddDate; the SDK does not invent dates, schedule a timer, or check local time. Date-time values preserve their instant, UTC offset, and fractional-second precision during JSON round trips.

An HTTP 200 for an auto-add update can include rejected items, warnings, and price diagnostics. Inspect all result categories to determine each item's reported outcome; a successful response does not prove that every item was accepted or that the scheduled change will be applied. The consuming application should read the product/participant list after the scheduled date and reconcile the observed state.

The announcement describes a changed rule for non-promo-code actions: UpdateProducts sets the maximum action price. Whether a price at that limit keeps a product active is a caller business rule, not an SDK decision. The announcement also applies auto-add changes on AutoAddDate and says stock is not sent for boost actions after the transition. The SDK does not infer an action type or issue an extra promotions request to decide this. The caller must choose whether stock is appropriate for the action. See the documented source ambiguity about the deactivate price comparison in the [contract notes](ozon-promotions-contracts-2026-09-28.md).

## Separate API family

ListDiscountRequests, ApproveDiscountRequest, and DeclineDiscountRequest belong to customer discounts-task requests. They are not seller actions and are not redirected to any of the eight methods above. The deprecation notice for /v1/actions/discounts-task/list has no cutoff date established in this migration.
