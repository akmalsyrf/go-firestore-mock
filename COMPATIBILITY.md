# Compatibility

## Version pairing

| fsmock | Firestore SDK (`cloud.google.com/go/firestore`) |
|--------|--------------------------------------------------|
| `v2.25.x` | `v1.25.x` |

Scheme: **`v2.<firestore-minor>.<patch>`**.

- Major `2` = module path `/v2`.
- Minor = Firestore SDK minor this release was **tested against**.
- Patch = fsmock-only fixes (docs, unwrap bugs, waivers, etc.).

### How to choose a version

```bash
# If your go.mod has cloud.google.com/go/firestore v1.25.3:
go get github.com/akmalsyrf/go-firestore-mock/v2@v2.25.0
```

Prefer the latest `v2.25.x` patch while you stay on Firestore 1.25.

### No backport policy

Only the latest `SupportedFirestoreMinor` (see `version.go`) receives fixes.
Older tags stay forever on the module proxy; we do not maintain release branches
per SDK minor. To get a fix, upgrade Firestore + fsmock together.

### MVS caveat

Go has no upper-bound dependency constraints. The Firestore version in this
module's `go.mod` is a **minimum**. Your module's MVS may select a newer
Firestore. Pairing means "tested against", not "locked exclusively".

After upgrading Firestore past the paired minor, upgrade fsmock to the matching
`v2.<new-minor>.x` (or open an SDK-bump issue if that tag does not exist yet).
