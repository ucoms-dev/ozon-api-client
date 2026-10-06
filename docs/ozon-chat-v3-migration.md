# Chat v3 and quota presence — v1.20.0

Source: official Seller API Swagger captured through Brave on 2026-10-06.
The selected provider schemas are preserved in
`ozon/testdata/seller-2026-10-06/chat-v3-schemas.json`.

`Chats().ListV3` reads one POST `/v3/chat/list` page. Use `ListChatsV3Params`
with a limit up to 100, an optional filter and a cursor. The response contains
nested typed chats, total unread count, cursor and has-next. Pagination,
retry policy and transport scheduling remain with the caller. `List` keeps its
legacy v2 route, parameters and response.

Documented discrepancies are handled narrowly:

- `has_next` is a boolean in the schema and a string in examples. Current chat
  and warehouse responses accept booleans and string true/false or 1/0 values;
  invalid flags return decoding errors rather than silently dropping pages.
- Message IDs are uint64 integers in the schema and strings in examples. One
  official example is larger than uint64. `ChatMessageID` retains the opaque
  decimal identifier exactly, accepting numeric and string JSON values without
  float64 conversion. Invalid non-decimal IDs are rejected.
- `created_at` is retained as a raw string so consumers can preserve their
  existing tolerant parsing of empty, date-only and legacy date formats.
- Chat list `chats` declares items without an explicit array type. The SDK models
  it as a slice as shown in the official examples.

`RawResponseJSON()` returns a separate copy of the received chat-list body so
established consumers can retain names, variants and unknown fields outside the
typed surface. Mutating the returned bytes does not change the retained body.

`GetProductRangeLimitResponse.QuotaDataPresent()` is true only when daily-create,
daily-update and total quotas all contain non-null limit and usage fields. Explicit
zero is present, and -1 remains the provider's unlimited sentinel. Consumers must
still validate HTTP status and supported values; missing data is not exhaustion.
Null entries in an operation-limit array are rejected rather than becoming zero
quotas.

No automatic waiting, write retry, API route switching or paging is introduced.
The existing SDK file structure, constructors and provider-error response
contract are preserved.
