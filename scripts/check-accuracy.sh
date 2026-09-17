#!/usr/bin/env bash
# check-accuracy.sh — PR accuracy gate for go-firestore-mock.
#
# Fails unless:
#   1) apicheck (SDK surface completeness) passes
#   2) required integration + parity tests PASS (skip counts as failure)
#
# Usage:
#   FIRESTORE_EMULATOR_HOST=127.0.0.1:8080 ./scripts/check-accuracy.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

REQUIRED_TESTS=(
  TestIntegration_ParityRawVsWrapped
  TestIntegration_ParityBatch
  TestIntegration_DocumentCRUD
  TestIntegration_CollectionAddNewDocParent
  TestIntegration_SubcollectionAndCollectionsIterator
  TestIntegration_QueryAndDocumentIterator
  TestIntegration_DocumentRefsIterator
  TestIntegration_WriteBatch
  TestIntegration_Transaction
  TestIntegration_BulkWriter
  TestIntegration_Aggregation
  TestIntegration_GetAllAndClientCollections
  TestIntegration_CollectionGroup
  TestIntegration_QuerySelectOffsetCursors
  TestIntegration_SerializeDeserialize
)

if [[ -z "${FIRESTORE_EMULATOR_HOST:-}" ]]; then
  echo "accuracy-gate FAIL: FIRESTORE_EMULATOR_HOST is not set" >&2
  exit 1
fi

echo "==> apicheck (SDK method surface)"
go test -count=1 ./internal/apicheck/

echo "==> required integration/parity tests"
RUN_RE="$(IFS='|'; echo "${REQUIRED_TESTS[*]}")"
JSON_OUT="$(mktemp)"
trap 'rm -f "$JSON_OUT"' EXIT

set +e
go test -tags=integration -count=1 -json -run "^(${RUN_RE})$" ./fstest/ >"$JSON_OUT"
TEST_EXIT=$?
set -e

if [[ "$TEST_EXIT" -ne 0 ]]; then
  echo "accuracy-gate FAIL: go test exited $TEST_EXIT" >&2
  # still print which required tests missed below
fi

pass=0
fail=0
skip=0
missing=0

for name in "${REQUIRED_TESTS[@]}"; do
  # Last action for this test name wins.
  action="$(grep -E "\"Test\":\"${name}\"" "$JSON_OUT" | grep -E '"Action":"(pass|fail|skip)"' | tail -1 | sed -n 's/.*"Action":"\([^"]*\)".*/\1/p' || true)"
  case "$action" in
    pass)
      echo "  PASS  $name"
      pass=$((pass + 1))
      ;;
    fail)
      echo "  FAIL  $name" >&2
      grep -E "\"Test\":\"${name}\"" "$JSON_OUT" | grep '"Action":"output"' | tail -5 | sed -n 's/.*"Output":"\(.*\)".*/    \1/p' | sed 's/\\n/\n    /g' >&2 || true
      fail=$((fail + 1))
      ;;
    skip)
      echo "  SKIP  $name (treated as FAIL for accuracy gate)" >&2
      skip=$((skip + 1))
      ;;
    *)
      echo "  MISS  $name (not executed)" >&2
      missing=$((missing + 1))
      ;;
  esac
done

total=${#REQUIRED_TESTS[@]}
ok=$pass
echo "==> accuracy: ${ok}/${total} required tests passed (fail=${fail} skip=${skip} miss=${missing})"

if [[ "$fail" -gt 0 || "$skip" -gt 0 || "$missing" -gt 0 || "$TEST_EXIT" -ne 0 ]]; then
  echo "accuracy-gate FAIL" >&2
  exit 1
fi

echo "accuracy-gate OK"
