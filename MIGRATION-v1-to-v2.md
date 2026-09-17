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

## Writes

Pass `fsmock.DocumentRef` into batch/transaction/bulkwriter. Invalid refs on `WriteBatch` surface as errors from `Commit`.

## Mocks

```diff
- mock := gofirestoremock.NewMockFirestoreClient(ctrl)
+ mock := mocks.NewMockClient(ctrl)
```
