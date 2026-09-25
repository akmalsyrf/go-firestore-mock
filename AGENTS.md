# Agent / contributor rules

This file is the source of truth for humans and coding agents. Follow it strictly.

## Hard rules

1. **Never edit `mocks/` by hand.** Run `make generate` (or `go generate ./...`).
2. **Every new `*xxxWrapper` method** needs either:
   - an integration test in `fstest/` that executes it, **or**
   - a waiver in `internal/apicheck/waivers.go` with non-empty `Reason` + `Issue`.
3. **Every new SDK type with methods** must be added to `internal/apicheck/registry.go` **or** `ignoredTypes` with a Reason.
4. **Signature drift** from the SDK requires a `deviations` entry with a Reason (see `exceptions.go`).
5. **Do not add feature `t.Skip` in `fstest`.** Skips fail the accuracy gate. Prefer a waiver if the emulator cannot support the feature.
   - Exception: `RequireEmulator` may `t.Skip` only when `FIRESTORE_EMULATOR_HOST` is unset (harness gate). The accuracy script already fails if the host is missing, so this path never soft-passes CI.
6. **Run `make gate` before push.** It must match CI.
7. **Do not bump Firestore** without updating `version.go`, `COMPATIBILITY.md`, regenerating mocks, and fixing apicheck/waivers.
8. **Version scheme:** `v2.<firestore-minor>.<patch>`. Never tag a mismatched minor.

## Recipe: wrap a new SDK method

1. Add the method to the matching interface (`client.go`, `query.go`, …).
2. Implement it on the `*xxxWrapper` — unwrap interfaces via `unwrap.go`, wrap returns via `new*` helpers.
3. If cursors/snapshots are involved, unwrap with `Reference()` / `unwrapCursorArgs` (cursor unwrap errors are deferred onto the query and surface on `Documents` / `Snapshots` / `Serialize`).
4. `make generate`
5. Add `fstest` coverage **or** a waiver.
6. `make gate`
7. If apicheck fails on signatures, either fix the wrapper or add `deviations["Type.Method"]`.

## Recipe: bump Firestore SDK

Use the issue template `.github/ISSUE_TEMPLATE/sdk-bump.md`, then:

```bash
go get cloud.google.com/go/firestore@v1.<minor>.0
# update version.go: Version, FirestoreSDKVersion, SupportedFirestoreMinor
# update COMPATIBILITY.md top row
go generate ./...
make gate
```

Tag `v2.<minor>.0` only after `gate` is green on `main`.

## Layout cheat sheet

| Path | Touch when |
|------|------------|
| `*.go` (root) | Wrappers / interfaces |
| `unwrap.go` / `clone.go` | Cross-type conversion / readSettings clone |
| `internal/apicheck/` | Parity / coverage policy |
| `fstest/` | Behavioral proof |
| `mocks/` | Generated only |
| `version.go` | SDK pairing |
