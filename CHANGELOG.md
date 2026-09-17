# Changelog

## [2.25.0] — 2026-09-18

First v2 release. Version scheme: **`v2.<firestore-minor>.<patch>`** (pairs with Firestore `v1.25.x`). See [COMPATIBILITY.md](COMPATIBILITY.md).

### Breaking (from v1)

- Module path is now `github.com/akmalsyrf/go-firestore-mock/v2`.
- Package renamed from `firestore` to `fsmock`.
- `NewFirestoreClient` removed; use `NewClient` / `MustNewClient`.
- `FirestoreClient` → `Client`.
- `CollectionGroup` returns `CollectionGroupRef` (not bare `Query`).
- `DocumentRef.Parent()` returns `CollectionRef` (not `*firestore.DocumentRef`).
- `CollectionRef.Add` returns `(DocumentRef, *WriteResult, error)`.
- Write APIs (`Transaction`, `WriteBatch`, `BulkWriter`) and `GetAll` accept `DocumentRef` interfaces.
- `Transaction.Documents` / `DocumentRefs` return `(…, error)` instead of panicking.
- `DocumentIterator` / snapshot iterators return `DocumentSnapshot` / `QuerySnapshot` interfaces.
- `CollectionIterator.Stop()` removed; added `GetAll` and `PageInfo`.
- Generated mocks moved to package `mocks`.
- Escape hatches: `DocumentSnapshot.Reference()`, `QuerySnapshot.Reference()`, `Query.SDKQuery()`.

### Added

- Compatibility with `cloud.google.com/go/firestore v1.25.0`.
- Query: `FindNearest`, `FindNearestPath`, `Serialize`/`Deserialize`, `WithReadOptions`, `WithRunOptions`, `Pipeline`.
- Client: `WithReadOptions`, `WithAlwaysUseImplicitOrderBy`, `Pipeline`.
- Aggregation: `WithSum`/`WithAvg` (+ Path), `GetResponse`, `Transaction`, `Pipeline`; `AggregationResult.Data`/`DataTo`.
- `VectorQuery`, full `Pipeline*` mockable boundaries.
- Cursor args unwrap `fsmock.DocumentSnapshot` → `*firestore.DocumentSnapshot`.
- `WithReadOptions` clones `readSettings` (no aliasing with the parent query/ref).
- `internal/apicheck` discover + signature parity + covercheck; `fstest` emulator harness.
- Accuracy gate: coverage-derived method evidence + waivers.
- Version constants in `version.go`; `COMPATIBILITY.md`.

### Changed

- Single `//go:generate` mockgen invocation (package mode) into `mocks/mocks.go`.
- Aggregation `Data()` recovers SDK panics into errors.
- Sentinel errors: `ErrNilClient`, `ErrNilArgument`, `ErrForeignImplementation`.
