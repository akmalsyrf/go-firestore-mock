# Changelog

## [2.0.0] — 2026-09-17

### Breaking

- Module path is now `github.com/akmalsyrf/go-firestore-mock/v2`.
- Package renamed from `firestore` to `fsmock`.
- `NewFirestoreClient` removed; use `NewClient` / `MustNewClient`.
- `FirestoreClient` → `Client`.
- `CollectionGroup` returns `CollectionGroupRef` (not bare `Query`).
- `DocumentRef.Parent()` returns `CollectionRef` (not `*firestore.CollectionRef`).
- `CollectionRef.Add` returns `(DocumentRef, *WriteResult, error)`.
- Write APIs (`Transaction`, `WriteBatch`, `BulkWriter`) and `GetAll` accept `DocumentRef` interfaces.
- `Transaction.Documents` / `DocumentRefs` return `(…, error)` instead of panicking.
- `DocumentIterator` / snapshot iterators return `DocumentSnapshot` / `QuerySnapshot` interfaces.
- `CollectionIterator.Stop()` removed; added `GetAll` and `PageInfo`.
- Generated mocks moved to package `mocks`.

### Added

- Compatibility with `cloud.google.com/go/firestore v1.25.0`.
- Query: `FindNearest`, `FindNearestPath`, `Serialize`/`Deserialize`, `WithReadOptions`, `WithRunOptions`, `Pipeline`.
- Client: `WithReadOptions`, `WithAlwaysUseImplicitOrderBy`, `Pipeline`.
- Aggregation: `WithSum`/`WithAvg` (+ Path), `GetResponse`, `Transaction`, `Pipeline`; `AggregationResult.Data`/`DataTo`.
- `VectorQuery`, full `Pipeline*` mockable boundaries.
- `DocumentIterator.ExplainMetrics`, iterator `PageInfo`.
- `internal/apicheck` parity tests, `fstest` emulator harness, GitHub Actions CI, Dependabot.
- DRY Query forwarding via embedded `queryWrapper` on collection wrappers.

### Changed

- Single `//go:generate` mockgen invocation (package mode) into `mocks/mocks.go`.
- Aggregation `Data()` recovers SDK panics into errors.
- BulkWriter docs note backpressure / blocking `Flush`/`End` (SDK ≥ v1.23).
