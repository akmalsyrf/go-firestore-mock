# Migrating from v1 to v2

## Module & import

```diff
- import gofirestoremock "github.com/akmalsyrf/go-firestore-mock"
+ import "github.com/akmalsyrf/go-firestore-mock/v2"
+ import "github.com/akmalsyrf/go-firestore-mock/v2/mocks"
```

Package name is now `fsmock` (no alias required vs `cloud.google.com/go/firestore`).

## Client construction

```diff
- client := gofirestoremock.NewFirestoreClient(fs)
+ client, err := fsmock.NewClient(fs)
+ if err != nil { ... }
  // or: client := fsmock.MustNewClient(fs)
```

## Interface renames

| v1 | v2 |
|----|----|
| `FirestoreClient` | `Client` |
| `CollectionGroup(...) Query` | `CollectionGroup(...) CollectionGroupRef` |
| mocks in root package | `mocks.NewMockClient`, etc. |

## Transaction / iterators

```diff
- iter := tx.Documents(q)           // panicked on bad Query
+ iter, err := tx.Documents(q)

- snap, err := iter.Next()          // *firestore.DocumentSnapshot
+ snap, err := iter.Next()          // fsmock.DocumentSnapshot
```

`CollectionIterator` no longer has `Stop()`; use `GetAll` / `PageInfo` like the SDK.

## Version pairing

v2 tags follow **`v2.<firestore-minor>.<patch>`**. This release line is `v2.25.x` and pairs with Firestore `v1.25.x`. See [COMPATIBILITY.md](COMPATIBILITY.md).

## Writes

Pass `fsmock.DocumentRef` into batch/transaction/bulkwriter. Invalid refs on `WriteBatch` surface as errors from `Commit`.

## Cursors

`StartAt` / `StartAfter` / `EndAt` / `EndBefore` accept `fsmock.DocumentSnapshot` (unwrapped to `*firestore.DocumentSnapshot` for the SDK).

Snapshots that cannot produce a non-nil `Reference()` (e.g. a bare mock) defer `ErrForeignImplementation` onto the query; the error surfaces from `Documents` / `Snapshots` / `Serialize` (not as a silent field-value cursor).

## QuerySnapshot.Changes

`QuerySnapshot.Changes()` returns `[]firestore.DocumentChange` from the SDK. That is intentional: `DocumentChange.Doc` and `OldDoc` are `*firestore.DocumentSnapshot`, not `fsmock.DocumentSnapshot`. A parallel `fsmock.DocumentChange` wrapper would be a signature deviation; realtime/Snapshots paths are also largely waived for emulator accuracy. Re-wrap manually or use `Reference()` if you need the raw SDK snapshot.

## BulkWriter jobs

`BulkWriter` create/set/update/delete return `fsmock.BulkWriterJob` (mockable `Results()`), not `*firestore.BulkWriterJob`.

## Mocks

```diff
- mock := gofirestoremock.NewMockFirestoreClient(ctrl)
+ mock := mocks.NewMockClient(ctrl)
```
