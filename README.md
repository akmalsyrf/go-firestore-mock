# GoFirestore Mock v2

[![Go Reference](https://pkg.go.dev/badge/github.com/akmalsyrf/go-firestore-mock/v2.svg)](https://pkg.go.dev/github.com/akmalsyrf/go-firestore-mock/v2)

Thin wrapper interfaces and gomock stubs over [`cloud.google.com/go/firestore`](https://pkg.go.dev/cloud.google.com/go/firestore).

Package: **`fsmock`**. Current pairing: **fsmock `v2.25.x` ↔ Firestore `v1.25.x`** — see [COMPATIBILITY.md](COMPATIBILITY.md).

## Install

```bash
go get github.com/akmalsyrf/go-firestore-mock/v2@v2.25.0
```

Requires **Go 1.25+**.

## Production

```go
fs, err := firestore.NewClient(ctx, "your-project-id")
client, err := fsmock.NewClient(fs) // never pass nil
```

## Unit tests with gomock

```go
mockClient := mocks.NewMockClient(ctrl)
mockColl := mocks.NewMockCollectionRef(ctrl)
mockClient.EXPECT().Collection("users").Return(mockColl)
```

## Layout

| Path | Role |
|------|------|
| `/` (`fsmock`) | Interfaces + wrappers |
| `/mocks` | Generated gomock stubs |
| `/fstest` | Emulator integration + parity |
| `/internal/apicheck` | Discover + signature parity + covercheck |
| `/scripts/check-accuracy.sh` | PR accuracy gate |
| `COMPATIBILITY.md` | Version pairing table |

## PR gate (required)

CI job **`gate`** fails unless all of these succeed:

1. **lint** / **unit** (includes apicheck + race) / **version-check**
2. **generate-check** — mocks not stale
3. **accuracy** — `scripts/check-accuracy.sh`
   - integration tests must pass (`t.Skip` in `fstest` = fail, except harness when emulator host unset)
   - every `*xxxWrapper` method covered **or** waived (Pipeline / vector / realtime listeners are largely waived when the emulator cannot exercise them)

```bash
make bootstrap
make gate          # starts emulator via docker if needed
```

## Development

Read [AGENTS.md](AGENTS.md) / [CONTRIBUTING.md](CONTRIBUTING.md) before opening a PR.

Migrating from v1: [MIGRATION-v1-to-v2.md](MIGRATION-v1-to-v2.md).

## License

MIT — see [LICENSE](LICENSE).
