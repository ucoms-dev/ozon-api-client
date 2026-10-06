# Seller API SDK update — 2026-10-06

Base: `3e6999b` (SDK v1.18.0). Official Seller API Swagger captured through Brave
at 2026-10-06T21:47:40.636Z. The SDK remains organized by existing API groups;
constructors, HTTP transport, error semantics and the v1 compatibility policy
remain in place. All additions use the existing typed DTO / request / copied
`CommonResponse` pattern. No dependency was added.

## Changes and migration

- `Products().UpdateProductImagesV2`: POST `/v2/product/pictures/import`,
  up to 100 `items`, up to 50 image URLs per item, `offer_id`, `primary_image`,
  `color_image`, `images`; response `TaskId` is int64. It replaces the complete
  image set per item. Poll `GetProductImportStatus` with the task ID for the
  processing result. The legacy v1 call is retained and marked deprecated.
- `Warehouses().GetListOfWarehousesV2`: POST `/v2/warehouse/list`, `Limit`,
  `Cursor`, optional `WarehouseIds`; response `Cursor`, `HasNext`, `Warehouses`,
  including all current nested address, first-mile and timetable fields.
  Nullable `PauseAt` is a pointer; enums remain strings. The legacy v1 call is
  retained and marked deprecated. The documented v1 1/minute restriction is
  not transferred to v2.
- `GetProductRangeLimitResponse.OperationLimits`: additive typed quotas with
  int64 limits and open-ended limit types. Missing/null stays nil. Successful
  decoding replaces the list; malformed values return an error without
  overwriting the previous value.
- `CommonResponse.Headers`: cloned raw HTTP headers available on success and
  provider JSON errors; `CopyCommonResponse` makes an independent clone.
  Headers are excluded from DTO JSON. No sleeping, background work or retries
  are added. As before, callers must check both Go error and HTTP status.
- Corrected stock update documentation to the current 30-second pair window,
  80 requests/minute/account, and 100 offer/warehouse pairs/request; documented
  the 1000-price batch limit. All 23 specially constrained routes and global
  limits are recorded in the accompanying rate-limit document.

## Documented provider discrepancies

1. Swagger models `operation_limits` as an object, but a successful provider
   response observed on 2026-10-06 contains an array. The SDK accepts both;
   marshaling uses an array. Unknown `limit_type` values are preserved.
2. Warehouse v2 `warehouse_ids` items are strings with format int64 in the
   schema, while the example shows numeric IDs. The request uses decimal
   strings according to the schema. Response warehouse IDs remain int64.
3. The warehouse example shows a string `has_next`, while the schema defines
   boolean. The SDK follows the boolean schema.
4. `quota_by_category` occurs in a provider response but has no current schema.
   It is not assigned an invented DTO contract by this update.

## Evidence and remaining scope

The route inventory covers 482 operations; after this update the SDK has 234
unique route operations, 207 exact matches, 27 client-only paths, zero HTTP
method mismatches and 275 unsupported Swagger operations. Absence from Swagger
alone is not proof of retirement. Unsupported operations and the historical
replacement backlog are retained for bounded updates by consumer need, as the
existing SDK contract plan requires. This update does not claim full schema
verification or implementation of all 482 operations.

The normalized operation inventory and 12 selected schema shapes are checked
in under `ozon/testdata/seller-2026-10-06`. Wire tests verify v2 routes, payloads,
large int64 IDs, cursor pagination, nested decoding, nullable pause timestamps,
quota shapes, unknown limit types, and provider error/HTTP header propagation.
A separate field/type coverage test compares all 12 DTO shapes with the
independently extracted Swagger fixture.

## Validation

| Check | Result |
|---|---|
| `go test ./... -count=1`, Go 1.26.4 | PASS, all five packages |
| `GOTOOLCHAIN=go1.20.14 go test ./... -count=1` | PASS, all five packages, including the final schema coverage test |
| Targeted `go test -race . ./ozon` for all added transport/product/warehouse/schema contracts | PASS |
| Vet of production files and added contract tests, now included in `scripts/vet-promotions.sh` and CI | PASS |
| Full `go vet ./...` | Existing 224-line diagnostic output, identical to a clean archive of base `3e6999b` |
| `git diff --check` | PASS |
| Fresh route audit | 207 exact matches; zero method mismatches; zero Swagger-deprecated methods without Go deprecation docs |

The targeted contracts were first compiled against the unchanged implementation
and failed for absent fields/methods, then passed after implementation. Full
vet failures are existing context-cancellation leaks, testing fatal calls in
non-test goroutines and unreachable legacy test code. They are not concealed
by disabling analyzers. Production files and new tests were separately vetted.
The full suite on Go 1.26.4 preceded the schema coverage test addition; that
additional test passed on Go 1.26.4 separately and in the targeted race run.

Publication of this additive update is designated SDK v1.19.0. Consumer
dependency changes and deployment remain separate tasks. Write endpoints were tested through local HTTP servers; this is
not evidence of authenticated provider write acceptance.
