# Promotions SDK verification — 2026-09-28

Base: `efe8b033adb4df1060035c89ef77fa52d8740796` (`dev`, SDK v1.17.0).
This report covers the local SDK release candidate. It does not claim publication,
authenticated Seller API execution, or migration of the consuming UCOMS service.

## Completed checks

| Check | Result |
|---|---|
| `go test ./... -count=1` on Go 1.26.4 | PASS, all five packages |
| `GOTOOLCHAIN=go1.20.14 go test ./... -count=1` | PASS, all five packages |
| `go test -race ./ozon -run 'TestPromotion\|TestAvailablePromotions' -count=1` | PASS |
| `bash scripts/vet-promotions.sh` on Go 1.26.4 | PASS |
| `GOTOOLCHAIN=go1.20.14 bash scripts/vet-promotions.sh` | PASS |
| `git diff --check` | PASS |
| Public Swagger field coverage | PASS, 19 DTO shapes |
| Existing AST route audit | PASS, all eight exact POST method/path pairs |
| Legacy HTTP contract checks | PASS, original routes and numeric price representation retained |
| UCOMS consumer compilation with a temporary SDK replacement | PASS; test binary compiled, not executed |

Consumer base: `2cfcd5f1dfa16a7b3bb0f877b3af092ac364f558`.
The following command ran from `services/ozon-integration-service` in an isolated
checkout of that base:

```sh
go test -c \
  -modfile /tmp/ozon-sdk-consumer-check-20260928/go.mod \
  -o /tmp/ozon-sdk-consumer-check-20260928/service.test \
  ./internal/service
```

The temporary modfile replaced the SDK with the candidate worktree and resolved
the consumer's local module replacements to that isolated checkout. Tracked
consumer `go.mod` and `go.sum` were not changed. Compilation demonstrates source
compatibility of the current consumer; it does not switch its existing calls to
the new endpoints.

## Full vet baseline and CI scope

`go vet ./...` fails on both the unchanged base and this candidate, on both tested
Go versions. The sorted diagnostics are identical per toolchain. Existing findings
include discarded context cancellation functions in legacy tests and
testing-goroutine/unreachable-code findings in notification tests. No new vet
diagnostic was introduced by this update.

The new CI vet step therefore runs `scripts/vet-promotions.sh`: all production Go
files in `ozon` together with the new promotion tests, plus the contract-audit
packages. No analyzer is disabled. Full tests still run across all packages.
Repair of unrelated legacy tests is outside this SDK change; the full vet command
is explicitly not reported as passing.

CI retains its existing action versions and coverage job. The version matrix uses
`1.20.x` and `1.x` with `check-latest: true` and `stable: true`, which are supported
by `actions/setup-go@v2`. No live GitHub Actions run is claimed.

## Contract evidence and boundary fix

The public OpenAPI document loaded by Ozon's documentation page was inspected for
all eight operations. The saved normalized extract, source URL, observation date,
and SHA-256 are recorded in the [contract notes](ozon-promotions-contracts-2026-09-28.md).
The complete Swagger document was not retained on disk, so no new full-repository
Swagger audit report is claimed. Payload fixtures are synthetic, not live seller
responses.

Tests exposed a pre-existing transport reflection assumption: a non-nil pointer
to primitive `stock` caused `reflect.NumField` to panic. Both new update methods
translate public `*uint64` stock into private wire values before invoking the
transport. Omitted stock and explicit zero remain distinct; the shared transport
and legacy operations retain their behavior.

The tests exercise request paths/bodies, all result categories, empty versus
malformed success responses, ID and decimal precision, bounds, cancellation,
deadlines, and one HTTP attempt per invocation. Future auto-add execution and
actual participation still require application-level readback from Ozon.

Independent exact-commit reviews are a release gate after the scoped task commit;
their results and reviewed commit/tree are reported separately in the handoff.
