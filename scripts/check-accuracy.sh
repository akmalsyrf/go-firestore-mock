#!/usr/bin/env bash
# check-accuracy.sh — PR accuracy gate for go-firestore-mock.
#
# Fails unless:
#   1) integration + unit tests pass with no skips in fstest
#   2) apicheck (discover + signature parity) passes
#   3) covercheck: every *xxxWrapper method is covered or waived
#
# Usage:
#   FIRESTORE_EMULATOR_HOST=127.0.0.1:8080 ./scripts/check-accuracy.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [[ -z "${FIRESTORE_EMULATOR_HOST:-}" ]]; then
  echo "accuracy-gate FAIL: FIRESTORE_EMULATOR_HOST is not set" >&2
  echo "  Fix: make emulator-up   # or export FIRESTORE_EMULATOR_HOST=127.0.0.1:8080" >&2
  exit 1
fi

COVER_OUT="${COVER_OUT:-$ROOT/cover.out}"
JSON_OUT="$(mktemp)"
trap 'rm -f "$JSON_OUT"' EXIT

MOD="github.com/akmalsyrf/go-firestore-mock/v2"

echo "==> tests with coverage (coverpkg=$MOD)"
set +e
go test -tags=integration -count=1 -json \
  -coverpkg="$MOD" \
  -coverprofile="$COVER_OUT" \
  ./... >"$JSON_OUT"
TEST_EXIT=$?
set -e

# Parse skip/fail from JSON (fstest package only for skip=fail policy).
skip=0
fail=0
while IFS= read -r line; do
  action=$(echo "$line" | sed -n 's/.*"Action":"\([^"]*\)".*/\1/p')
  pkg=$(echo "$line" | sed -n 's/.*"Package":"\([^"]*\)".*/\1/p')
  test=$(echo "$line" | sed -n 's/.*"Test":"\([^"]*\)".*/\1/p')
  [[ -z "$test" ]] && continue
  case "$action" in
    fail)
      echo "  FAIL  $pkg $test" >&2
      fail=$((fail + 1))
      ;;
    skip)
      if [[ "$pkg" == *"/fstest" ]]; then
        echo "  SKIP  $pkg $test (treated as FAIL for accuracy gate)" >&2
        echo "    Fix: make the test pass on the emulator, or move optional coverage to a waiver in internal/apicheck/waivers.go and delete the skip." >&2
        skip=$((skip + 1))
      else
        echo "  skip  $pkg $test (allowed outside fstest)"
      fi
      ;;
  esac
done < <(grep -E '"Action":"(pass|fail|skip)"' "$JSON_OUT" || true)

if [[ "$TEST_EXIT" -ne 0 || "$fail" -gt 0 || "$skip" -gt 0 ]]; then
  echo "accuracy-gate FAIL: tests exit=$TEST_EXIT fail=$fail skip=$skip" >&2
  exit 1
fi

echo "==> apicheck (discover + signatures)"
go test -count=1 ./internal/apicheck/ -run 'TestDiscover|TestMethodParity|TestAggregation'

echo "==> covercheck (wrapper method coverage)"
FSMOCK_COVERPROFILE="$COVER_OUT" go test -count=1 ./internal/apicheck/ -run TestWrapperMethodCoverage

echo "accuracy-gate OK"
