# Ozon Seller API client contract plan

Source snapshot: `swagger (1).json`, SHA-256
`48a0c82a411bf08c9769d61be6d444488af3062cc7e97d2d50067e16d4b05bb8`.

## Compatibility decision

Officially deprecated exported methods stay in the v1 release line so existing
consumers continue to compile. Every deprecated method present in this client
has a Go `Deprecated:` comment and a current replacement. Removing the old
surface is reserved for v2 of this module.

The one other deprecated Swagger operation,
`POST /v1/posting/digital/list`, is not implemented by this client and therefore
does not need a compatibility wrapper.

## Completed in this update

- Added exact v2 review contracts for comment deletion, status changes, counts,
  review details, and cursor-based review lists.
- Added exact v2 product certificate options, parameter validation, and JSON
  certificate creation contracts.
- Kept the existing current posting replacements: FBO v3, FBS v4, unfulfilled
  FBS v4, act status, and act PDF.
- Made successful empty response bodies valid, as required by current v2 review
  mutation responses.
- Added a deterministic method/path audit command and checked-in snapshot.
- Corrected README module paths and documented the v1 deprecation policy.

## Client-only endpoint disposition

The current client has 25 unique method/path pairs absent from this Swagger
snapshot. Absence is not treated as an official deprecation marker: these calls
remain source-compatible, but new code should not adopt them without live Ozon
verification.

### Direct or consolidated replacement exists (17)

| Client endpoint | Replacement in Swagger | Action |
|---|---|---|
| `POST /v1/conditional-cancellation/approve` | `POST /v2/conditional-cancellation/approve` | Add a v2 method before migrating any consumer. |
| `POST /v1/conditional-cancellation/list` | `POST /v2/conditional-cancellation/list` | Add a v2 method before migrating any consumer. |
| `POST /v1/conditional-cancellation/reject` | `POST /v2/conditional-cancellation/reject` | Add a v2 method before migrating any consumer. |
| `GET /v1/draft/create/info` | `POST /v2/draft/create/info` | Replace verb and DTOs; do not alias the old request. |
| `GET /v1/draft/supply/create` | `POST /v2/draft/supply/create` | Replace verb and DTOs; add status polling via `/v2/draft/supply/create/status`. |
| `GET /v1/draft/timeslot/info` | `POST /v2/draft/timeslot/info` | Replace verb and DTOs. |
| `POST /v1/product/import/stocks` | `POST /v2/products/stocks` | Use existing `UpdateQuantityStockProducts`. |
| `POST /v1/quant/get` | `POST /v1/product/quant/info` | Add the current product-quant DTO. |
| `POST /v1/quant/list` | `POST /v1/product/quant/list` | Add the current product-quant DTO. |
| `POST /v2/chat/list` | `POST /v3/chat/list` | Add cursor-capable v3 chat DTOs. |
| `POST /v2/returns/rfbs/compensate` | `POST /v1/returns/rfbs/action/set` | Consolidate into the current action request. |
| `POST /v2/returns/rfbs/receive-return` | `POST /v1/returns/rfbs/action/set` | Consolidate into the current action request. |
| `POST /v2/returns/rfbs/reject` | `POST /v1/returns/rfbs/action/set` | Consolidate into the current action request. |
| `POST /v2/returns/rfbs/return-money` | `POST /v1/returns/rfbs/action/set` | Consolidate into the current action request. |
| `POST /v2/returns/rfbs/verify` | `POST /v1/returns/rfbs/action/set` | Consolidate into the current action request. |
| `POST /v2/supply-order/get` | `POST /v3/supply-order/get` | Add v3 DTOs; keep the old method only for v1 compatibility. |
| `POST /v2/supply-order/list` | `POST /v3/supply-order/list` | Add v3 cursor DTOs; keep the old method only for v1 compatibility. |

### Replacement is a changed workflow (1)

| Client endpoint | Current Swagger family | Action |
|---|---|---|
| `GET /v1/draft/create` | `POST /v1/draft/direct/create`, `/v1/draft/crossdock/create`, `/v1/draft/multi-cluster/create` | Expose separate workflow-specific methods; there is no safe one-call compatibility adapter. |

### No direct replacement found in this Swagger snapshot (7)

| Client endpoint | Action |
|---|---|
| `POST /v1/chat/updates` | Keep v1-compatible only; use chat history/push flows for new integrations after product-level validation. |
| `POST /v1/conditional-cancellation/get` | Keep v1-compatible only; v2 list may contain the needed data, but it is not a documented one-to-one replacement. |
| `POST /v1/product/upload_digital_codes` | Keep v1-compatible only; posting-code upload has different semantics. |
| `POST /v1/product/upload_digital_codes/info` | Keep v1-compatible only; no matching task-status operation is present. |
| `POST /v1/quant/ship` | Keep v1-compatible only; this path is shared by `Ship` and `Status` in the client and has no current Swagger match. |
| `POST /v2/fbs/posting/sent-by-seller` | Keep v1-compatible only; no current Swagger match was found. |
| `POST /v2/posting/fbs/product/change` | Keep v1-compatible only; no current bulk-weight replacement was found. |

## Unsupported Swagger backlog

There are 266 Swagger operations not implemented by the client after this
update. They are inventoried in the generated audit report. They should be added
by business domain and consumer need, not as an unreviewed bulk generation:

1. Migrate the 17 client-only endpoints with documented replacements.
2. Add methods required by UCOMS runtime flows, with request/response wire tests.
3. Add remaining domains in bounded releases, preserving source compatibility.
4. Re-run the audit for each new Swagger snapshot and investigate every method
   mismatch or newly client-only endpoint before release.
