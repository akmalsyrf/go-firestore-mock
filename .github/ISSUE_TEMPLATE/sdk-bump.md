---
name: SDK bump
about: Upgrade cloud.google.com/go/firestore and retarget fsmock pairing
title: "chore: bump firestore to v1.XX.0 (fsmock v2.XX.0)"
labels: dependencies
---

## Target

- Firestore: `v1.___.0`
- fsmock tag: `v2.___.0`

## Checklist

- [ ] `go get cloud.google.com/go/firestore@v1.___.0`
- [ ] Update `version.go` (`Version`, `FirestoreSDKVersion`, `SupportedFirestoreMinor`)
- [ ] Update `COMPATIBILITY.md` top row
- [ ] `go generate ./...`
- [ ] Fix `internal/apicheck` (registry / deviations / ignoredTypes)
- [ ] Fix `waivers.go` / add `fstest` coverage for new methods
- [ ] `make gate` green
- [ ] Tag `v2.___.0` after merge to `main`
