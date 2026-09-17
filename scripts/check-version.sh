#!/usr/bin/env bash
# check-version.sh — ensure version.go, go.mod, and COMPATIBILITY.md agree.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

minor_from_version_go=$(sed -n 's/^const SupportedFirestoreMinor = \([0-9][0-9]*\).*/\1/p' version.go)
sdk_ver=$(sed -n 's/^const FirestoreSDKVersion = "\(.*\)".*/\1/p' version.go)
pkg_ver=$(sed -n 's/^const Version = "\(.*\)".*/\1/p' version.go)
gomod_fs=$(go list -m -f '{{.Version}}' cloud.google.com/go/firestore)

if [[ -z "$minor_from_version_go" || -z "$sdk_ver" || -z "$pkg_ver" ]]; then
  echo "version-check FAIL: could not parse version.go" >&2
  exit 1
fi

# FirestoreSDKVersion like v1.25.0 -> minor 25
sdk_minor=$(echo "$sdk_ver" | sed -n 's/^v1\.\([0-9][0-9]*\)\..*/\1/p')
# Version like 2.25.0 -> minor 25
pkg_minor=$(echo "$pkg_ver" | sed -n 's/^2\.\([0-9][0-9]*\)\..*/\1/p')
# go.mod like v1.25.0
gomod_minor=$(echo "$gomod_fs" | sed -n 's/^v1\.\([0-9][0-9]*\)\..*/\1/p')

compat_line=$(grep -E '^\| `v2\.[0-9]+\.x`' COMPATIBILITY.md | head -1 || true)
compat_minor=$(echo "$compat_line" | sed -n 's/.*v2\.\([0-9][0-9]*\)\.x.*/\1/p')

echo "version.go Version=$pkg_ver SupportedFirestoreMinor=$minor_from_version_go FirestoreSDKVersion=$sdk_ver"
echo "go.mod firestore=$gomod_fs"
echo "COMPATIBILITY.md top row minor=$compat_minor"

fail=0
if [[ "$minor_from_version_go" != "$sdk_minor" ]]; then
  echo "FAIL: SupportedFirestoreMinor ($minor_from_version_go) != minor of FirestoreSDKVersion ($sdk_ver)" >&2
  fail=1
fi
if [[ "$minor_from_version_go" != "$pkg_minor" ]]; then
  echo "FAIL: SupportedFirestoreMinor ($minor_from_version_go) != minor of Version ($pkg_ver)" >&2
  fail=1
fi
if [[ "$minor_from_version_go" != "$gomod_minor" ]]; then
  echo "FAIL: SupportedFirestoreMinor ($minor_from_version_go) != go.mod firestore minor ($gomod_fs)" >&2
  echo "  Fix: go get cloud.google.com/go/firestore@v1.${minor_from_version_go}.0  OR bump version.go" >&2
  fail=1
fi
if [[ "$minor_from_version_go" != "$compat_minor" ]]; then
  echo "FAIL: SupportedFirestoreMinor ($minor_from_version_go) != COMPATIBILITY.md top row ($compat_line)" >&2
  fail=1
fi

# When building a release tag, GITHUB_REF_TYPE=tag / GITHUB_REF_NAME=v2.25.0
if [[ "${GITHUB_REF_TYPE:-}" == "tag" ]]; then
  tag="${GITHUB_REF_NAME}"
  tag_minor=$(echo "$tag" | sed -n 's/^v2\.\([0-9][0-9]*\)\..*/\1/p')
  if [[ "$tag_minor" != "$minor_from_version_go" ]]; then
    echo "FAIL: tag $tag minor != SupportedFirestoreMinor $minor_from_version_go" >&2
    fail=1
  fi
  tag_ver="${tag#v}"
  if [[ "$tag_ver" != "$pkg_ver" ]]; then
    echo "FAIL: tag $tag != version.go Version $pkg_ver" >&2
    fail=1
  fi
fi

if [[ "$fail" -ne 0 ]]; then
  echo "version-check FAIL" >&2
  exit 1
fi
echo "version-check OK"
