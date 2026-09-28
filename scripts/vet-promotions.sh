#!/usr/bin/env bash
set -euo pipefail

cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.."

# Full-repository vet currently reports pre-existing findings in legacy tests.
# Vet all production files in ozon together with the new promotion contracts,
# without disabling analyzers or changing unrelated tests.
sources=()
for source in ozon/*.go; do
  if [[ "$source" != *_test.go ]]; then
    sources+=("$source")
  fi
done
go vet "${sources[@]}" ozon/promotions_*_test.go
go vet ./internal/contractaudit ./cmd/contract-audit
