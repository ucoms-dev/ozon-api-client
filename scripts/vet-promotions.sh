#!/usr/bin/env bash
set -euo pipefail

cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.."

# Full-repository vet currently reports pre-existing findings in legacy tests.
# Vet all production files in ozon together with current SDK contracts,
# without disabling analyzers or changing unrelated tests.
sources=()
for source in ozon/*.go; do
  if [[ "$source" != *_test.go ]]; then
    sources+=("$source")
  fi
done
go vet "${sources[@]}" ozon/promotions_*_test.go   ozon/request_contract_test.go ozon/products_current_contract_test.go   ozon/warehouses_v2_contract_test.go ozon/seller_swagger_contract_test.go ozon/chats_v3_contract_test.go

root_sources=()
for source in *.go; do
  if [[ "$source" != *_test.go ]]; then
    root_sources+=("$source")
  fi
done
go vet "${root_sources[@]}" client_headers_test.go
go vet ./internal/contractaudit ./cmd/contract-audit
