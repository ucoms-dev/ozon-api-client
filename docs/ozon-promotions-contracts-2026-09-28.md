# Ozon promotions contract notes — 2026-09-28

## Scope and evidence

This field matrix records the eight seller-action operations in the SDK plan. The official announcement and operation pages were reviewed on 2026-09-28. The public Seller API OpenAPI schema was retrieved from Ozon's documentation host on that date; the normalized response-shape extract is saved at [ozon/testdata/promotions-v2/swagger-shapes.json](../ozon/testdata/promotions-v2/swagger-shapes.json). The full 3.68 MB source document was not retained, and no full-source SHA-256 is claimed. The saved normalized extract has SHA-256 415e65b5d60a567b69b6ef2d7453d9dc673af9f79940213e15dfa7e167746af3.

Schema inspection is not a live API response test. Contract fixtures labeled synthetic-spec are synthetic payloads derived from the reviewed operation schemas and this field matrix; they are not captured or redacted seller responses. No authenticated call against the new methods is represented as verified here.

## Official sources

| Source | Contract area |
|---|---|
| [Ozon announcement: changes to Seller API promotion methods](https://dev.ozon.ru/start/563-Obnovlenie-metodov-raboty-s-aktsiiami-v-Seller-API-v2/) | Cutover date, behavior notes, known narrative contradiction |
| [Promos](https://docs.ozon.ru/api/seller/#operation/Promos) | GET /v1/actions, including available auto-add dates |
| [ActionsCandidates](https://docs.ozon.ru/api/seller/#operation/ActionsCandidates) | POST /v2/actions/candidates |
| [ActionsProducts](https://docs.ozon.ru/api/seller/#operation/ActionsProducts) | POST /v2/actions/products |
| [ActionsProductsUpdate](https://docs.ozon.ru/api/seller/#operation/ActionsProductsUpdate) | POST /v1/actions/products/update |
| [ActionsProductsDeactivate](https://docs.ozon.ru/api/seller/#operation/ActionsProductsDeactivate) | POST /v2/actions/products/deactivate |
| [ActionsAutoAddProductsListV2](https://docs.ozon.ru/api/seller/#operation/ActionsAutoAddProductsListV2) | POST /v2/actions/auto-add/products/list |
| [ActionsAutoAddProductsCandidatesV2](https://docs.ozon.ru/api/seller/#operation/ActionsAutoAddProductsCandidatesV2) | POST /v2/actions/auto-add/products/candidates |
| [ActionsAutoAddProductsUpdateV2](https://docs.ozon.ru/api/seller/#operation/ActionsAutoAddProductsUpdateV2) | POST /v2/actions/auto-add/products/update |
| [ActionsAutoAddProductsDeleteV2](https://docs.ozon.ru/api/seller/#operation/ActionsAutoAddProductsDeleteV2) | POST /v2/actions/auto-add/products/delete |
| [Seller API OpenAPI JSON](https://docs.ozon.ru/api/seller/swagger.json?1790601089484) | Retrieved 2026-09-28; source response Date header: Mon, 28 Sep 2026 13:11:45 GMT |

## Endpoint field matrix

All eight operations use POST. Request JSON names below are exact wire names. The SDK response types embed core.CommonResponse and expose endpoint fields at the response top level; the new methods do not decode a result wrapper.

| SDK method | Request fields | Response fields | Paging / notes |
|---|---|---|---|
| ProductsAvailableForPromotionV2 — [ActionsCandidates](https://docs.ozon.ru/api/seller/#operation/ActionsCandidates) | action_id:uint64, limit:uint64, last_id:string | products[], total:uint64, last_id:string | Cursor paging; start with empty last_id, then pass the returned cursor. No automatic page traversal. |
| ProductsInPromotionV2 — [ActionsProducts](https://docs.ozon.ru/api/seller/#operation/ActionsProducts) | action_id:uint64, limit:uint64, last_id:string | products[], total:uint64, last_id:string | Cursor paging; pass the returned cursor to the next call. |
| UpdateProducts — [ActionsProductsUpdate](https://docs.ozon.ru/api/seller/#operation/ActionsProductsUpdate) | action_id:uint64, products[]:{product_id:uint64, action_price:PromotionMoney, stock?:uint64} | active_product_ids[], deactivated_product_ids[], rejected[], warnings[] | stock omission and explicit zero differ. product_id is numeric in this request. |
| RemoveProductV2 — [ActionsProductsDeactivate](https://docs.ozon.ru/api/seller/#operation/ActionsProductsDeactivate) | action_id:uint64, product_ids:string[] | product_ids[] | The request array uses decimal strings for uint64 IDs. Do not invent a rejected response field. |
| ListAutoAddProductsV2 — [ActionsAutoAddProductsListV2](https://docs.ozon.ru/api/seller/#operation/ActionsAutoAddProductsListV2) | action_id:uint64, auto_add_date:date-time, limit:uint64, offset:uint64 | products[], total:uint64 | Offset paging; date is caller-selected from auto_add_dates. |
| ListAutoAddCandidatesV2 — [ActionsAutoAddProductsCandidatesV2](https://docs.ozon.ru/api/seller/#operation/ActionsAutoAddProductsCandidatesV2) | action_id:uint64, auto_add_date:date-time, limit:uint64, offset:uint64 | products[], total:uint64 | Offset paging; no cursor parameter. |
| UpdateAutoAddProductsV2 — [ActionsAutoAddProductsUpdateV2](https://docs.ozon.ru/api/seller/#operation/ActionsAutoAddProductsUpdateV2) | action_id:uint64, auto_add_date:date-time, products[]:{id:uint64, action_price:PromotionMoney, stock?:uint64} | product_ids[], deactivated_ids[], rejected[], warnings[], below_min_price[], extremely_low_price[], failed_price[] | Auto-add product identifier is id, not product_id; diagnostics are separate result categories. |
| DeleteAutoAddProductsV2 — [ActionsAutoAddProductsDeleteV2](https://docs.ozon.ru/api/seller/#operation/ActionsAutoAddProductsDeleteV2) | action_id:uint64, auto_add_date:date-time, product_ids:string[] | product_ids[] | The request array uses decimal strings for uint64 IDs. |

## Product item fields

### Candidates and participants

PromotionCandidateV2 and PromotionParticipantV2 share the following price and product fields. JSON names are snake_case as shown; Go names are the corresponding exported fields.

| JSON field | Wire meaning / SDK representation |
|---|---|
| id | Ozon product ID (PromotionProductID, an unsigned 64-bit integer). |
| price, action_price, max_action_price, alert_max_action_price, price_min_elastic, price_max_elastic, marketplace_seller_price, min_seller_price | Optional price values as *PromotionMoney; each money object has string amount and currency. |
| alert_max_action_price_failed, is_quarantined | Optional booleans; absence remains distinct from false. |
| current_boost, min_boost, max_boost | Numeric boost values (float64 in the transport DTO); they are not money. |
| min_stock, recommended_stock | Unsigned quantities, kept separate from prices. |
| website_prices | PromotionWebsitePrices: optional price and prices_by_schema, whose schema-keyed entries can carry optional black_price and green_price. Schema keys are preserved as strings, not constrained to a hard-coded enum. |

A participant additionally has add_mode:string and optional stock:uint64. Its documented string values are SELLER and AUTO; unknown values are retained. A candidate has neither add_mode nor stock in this schema. The SDK does not synthesize missing fields.

### Auto-add products and candidates

Both auto-add item shapes can carry price, base_price, max_discount_price, action_price_to_auto_add, marketplace_seller_price, and min_seller_price as optional PromotionMoney, plus currency, offer_id, sku, name, min_action_quantity, quantity_to_auto_add, has_expired_min_seller_price, will_be_quarantined, and website_prices.

| Item | Identifier | Auto-add flag |
|---|---|---|
| AutoAddProductV2 | product_id (PromotionProductID) | add_mode?:bool |
| AutoAddCandidateV2 | id (PromotionProductID) | is_manually_added?:bool |

offer_id and sku are not substitutes for product_id or id. The string add_mode on an ordinary promotion participant and the boolean add_mode on an auto-add product are separate wire contracts.

## Encoding and validation rules

### Money and IDs

- PromotionMoney.amount is a decimal string. The SDK validates the decimal form and requires a non-empty, trimmed currency value without parsing into binary floating point, forcing RUB, or assuming a fixed decimal scale. Negative amounts, empty amounts, exponent notation, NaN, and infinity are invalid requests. Zero is a valid transport value, though Ozon may reject it for business reasons. An omitted response Money field stays nil.
- PromotionProductID reads a JSON integer or decimal string exactly into uint64, including values above JavaScript's exact integer range. It rejects fractions, signs, exponent notation, overflow, null, and empty strings. Response decoding accepts either number or string form; ID request arrays for deactivate/delete use strings, while update item IDs use numbers.
- Update auto-add diagnostics use PromotionPriceValidation{key, value} with key as PromotionProductID and value as json.Number; value is an API diagnostic number, not PromotionMoney.
- PromotionIssue has product_id and reason. PromotionInputError and PromotionProtocolError identify the method and concise reason without embedding credentials, sensitive headers, or a full HTTP body.

### SDK input safeguards

The OpenAPI schema does not mark all these constraints as required or impose all of the SDK limits. The SDK intentionally validates at its own boundary:

- A supplied action_id and product identifier must be positive.
- List limits are 1–100. Both ordinary promotion lists use last_id; auto-add lists use offset.
- Mutation batches must be non-empty, no larger than 1000, and contain unique positive IDs. Auto-add update has at least one item. Oversized writes are rejected; the SDK does not chunk them.
- auto_add_date must be non-zero for auto-add list and mutation calls. It comes from the seller's returned schedule; the SDK does not invent a date or compare it with local time.
- Optional stock is a pointer: nil means omit the field, while a pointer to zero explicitly sends zero. The SDK validates wire shape, not whether a promotion's action type permits stock. The caller owns that business choice.

### Response envelopes and partial success

The OpenAPI schemas do not mark every response member as required. The SDK applies a stricter fail-closed check for HTTP 200 to prevent an empty or malformed body from looking like a successful empty result:

- List responses must decode a non-null products array and total. products:[] with total:0 is valid; absent or null products is a protocol error. Both ordinary promotion list responses also carry last_id.
- A mutation response must contain at least one documented result field of the expected type. Deactivate and delete responses must contain product_ids. Correctly typed empty arrays are valid but confirm no IDs by themselves.
- An HTTP 200 rejected or warnings result is an ordinary per-item/partial result. It is not folded into confirmed IDs or turned into a transport error. A partial list of confirmed IDs remains partial.
- Unknown response fields are ignored for forward compatibility. Shape checks apply to the new success responses; existing global transport behavior for other SDK operations is unchanged.
- Valid Ozon error-status response bodies retain CommonResponse status and details. Transport/context failures, malformed JSON, invalid local input, or malformed HTTP 200 envelopes return a Go error.

## Auto-add price diagnostic schema/example discrepancy

The generated example shown on the update operation page has a diagnostics field as a single object and illustrates an empty string for value. The operation's OpenAPI response schema instead defines below_min_price, extremely_low_price, and failed_price as arrays of objects with key integer (uint64) and value number (double). The official OpenAPI response schema was inspected on 2026-09-28, and its normalized relevant shapes are in [swagger-shapes.json](../ozon/testdata/promotions-v2/swagger-shapes.json). The SDK and synthetic fixtures follow the schema: []PromotionPriceValidation, preserving the numeric text with json.Number. The example and schema disagreement remains documented; no captured live operation response is claimed as evidence.

## Date behavior and source ambiguities

- Ozon announced that corresponding legacy methods are disabled on 2026-10-13. The new endpoints use their fixed paths before and after that date; the SDK has no local-clock route switch, hidden v1 fallback, retry, extra action-type lookup, or automatic pagination/chunking.
- The announcement says the post-transition UpdateProducts rule for non-promo-code actions sets the maximum action price; deciding whether a product stays active at that limit is the consumer's business rule. Auto-add changes are applied on auto_add_date. A successful HTTP response is not proof that a future change ran; consumers should read back state after the date.
- Stock is not sent for boost actions after the announced change. Do not treat a stock-discount action and “Maximum boost” as synonyms. The SDK does not infer the action type or decide whether to send stock.
- The announcement's introductory text about deactivate and the detailed update field description disagree about the direction of the price comparison. The SDK does not calculate that comparison; a consumer enabling price automation must resolve the business rule separately.
- discounts-task methods remain separate from seller actions. The docs reviewed here establish no disable date for /v1/actions/discounts-task/list.
