# GoFirestore Mock v2

[![Go Reference](https://pkg.go.dev/badge/github.com/akmalsyrf/go-firestore-mock/v2.svg)](https://pkg.go.dev/github.com/akmalsyrf/go-firestore-mock/v2)

Thin wrapper interfaces and gomock stubs over [`cloud.google.com/go/firestore`](https://pkg.go.dev/cloud.google.com/go/firestore) **v1.25.0**.

Package: **`fsmock`**.

## Install

```bash
go get github.com/akmalsyrf/go-firestore-mock/v2
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
| `/internal/apicheck` | SDK method-set completeness |
| `/scripts/check-accuracy.sh` | PR accuracy gate |

## PR gate (required)

CI job **`gate`** fails unless all of these succeed:

1. **lint** / **unit** (includes apicheck + race)
2. **generate-check** — mocks not stale
3. **accuracy** — `scripts/check-accuracy.sh`
   - apicheck must pass
   - required integration + **parity** tests must **pass** (`t.Skip` = fail)

Enable branch protection: require status check **`gate`**.

```bash
export FIRESTORE_EMULATOR_HOST=127.0.0.1:8080
make accuracy   # or: make gate
```

## Development

```bash
make test
make generate-check
make accuracy          # needs emulator
```

Migrating from v1: [MIGRATION-v1-to-v2.md](MIGRATION-v1-to-v2.md).

## License

MIT — see [LICENSE](LICENSE).
