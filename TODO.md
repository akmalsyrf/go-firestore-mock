# TODO — Progress & roadmap toward official Firestore parity

This file tracks **current progress** for [github.com/akmalsyrf/go-firestore-mock](https://github.com/akmalsyrf/go-firestore-mock) and work toward **full API compatibility** with [`cloud.google.com/go/firestore`](https://pkg.go.dev/cloud.google.com/go/firestore).

## What “100% compatibility” means here

The official package exposes **constructors** (`NewClient`, `NewClientWithDatabase`, `NewRESTClient`), **dozens of option and filter types**, **Pipeline / vector / aggregation expression APIs**, and **package-level helpers** (`ArrayUnion`, `Increment`, etc.). This repo intentionally exposes **narrow interfaces** plus gomock stubs.

**Practical definition of “100%” for this project:** every **instance method** on the concrete types you choose to mirror (`*firestore.Client`, `Query`, `DocumentRef`, …) has a matching method on the wrapper interface, with semantics delegated to the SDK (or documented if intentionally unsupported). Types the project does **not** wrap (e.g. entire `Pipeline` surface) are listed under “Out of scope unless added”.

---

## Current progress

| Area | Status |
|------|--------|
| **Module & SDK** | `go.mod` pins `cloud.google.com/go/firestore v1.22.0` (Go 1.25). |
| **Client wrapper** | `Collection`, `CollectionGroup`, `Doc`, `DocFromFullPath`, `Close`, `BulkWriter`, `Batch`, `RunTransaction`, `Collections`, `GetAll`. |
| **Query / collection** | `Where`, `WherePath`, `WhereEntity`, `OrderBy`, `OrderByPath`, limit/offset, cursors, `Select`, `SelectPaths`, `Documents`, `Snapshots`, `NewAggregationQuery`. |
| **Document** | CRUD, subcollection, `Collections`, `Snapshots`, metadata (`ID`, `Path`, `Reference`, `Parent`). |
| **Batch / bulk / transaction** | Full write batch & bulk writer; transactions support reads (`Get`, `GetAll`, `Documents(q)`, `DocumentRefs(coll)`) and writes (`Create`, `Set`, `Update`, `Delete`). |
| **Snapshot & iterators** | `DocumentSnapshot` (including timestamps, `Ref`, `DataAtPath`); document/query iterators; `DocumentRefIterator`; `CollectionIterator` partial (see gaps). |
| **Aggregation** | `WithCount` + `Get`; **wrapper `Count` reads values from the SDK `AggregationResult`** (`*firestorepb.Value` / Go numbers). |
| **Mocks** | `go:generate mockgen` for every interface (Client, Query, CollectionRef, DocumentRef, DocumentSnapshot, Transaction, BulkWriter, WriteBatch, AggregationQuery / Result, all iterators). |
| **Documentation** | README: production pattern `firestore.NewClient` → `NewFirestoreClient`, gomock example, compatibility section. |
| **Tests** | `go test ./...` (incl. `toFirestoreQueryer` rejection of foreign Query impls, aggregation tests for `Count` with realistic result maps). |

---

## Remaining work (roadmap to full parity)

Check off when done. Suggested order: **Client / read options → Query → Document / snapshot → iterators → Transaction → full aggregation → Pipeline & advanced**.

### `FirestoreClient` (vs `*firestore.Client`)

- [x] `DocFromFullPath(fullPath string) DocumentRef`
- [ ] `Pipeline() *Pipeline` — requires abstract `Pipeline` + source (`PipelineSource`) for parity; see below.
- [ ] `WithReadOptions(opts ...firestore.ReadOption) *Client` — the official API returns a new client; decide whether the interface returns `FirestoreClient` or another option surface.

### `Query`

- [x] `WherePath(fp firestore.FieldPath, op string, value interface{}) Query`
- [x] `WhereEntity(ef firestore.EntityFilter) Query`
- [x] `OrderByPath(fp firestore.FieldPath, dir firestore.Direction) Query`
- [x] `SelectPaths(fieldPaths []firestore.FieldPath) Query`
- [ ] `FindNearest` / `FindNearestPath` (vector)
- [ ] `Serialize` / `Deserialize`
- [ ] `Pipeline() *Pipeline`
- [ ] `WithReadOptions(opts ...firestore.ReadOption)`
- [ ] `WithRunOptions(opts ...firestore.RunOption)`

### `CollectionRef`

- [x] `DocumentRefs(ctx context.Context) DocumentRefIterator`
- [ ] `WithReadOptions(opts ...firestore.ReadOption) CollectionRef`

### `DocumentRef`

- [ ] `WithReadOptions(opts ...firestore.ReadOption) DocumentRef`

### `DocumentSnapshot`

- [x] `DataAtPath(fp firestore.FieldPath) (interface{}, error)` (in addition to `DataAt(string)`)

### `DocumentIterator`

- [ ] `ExplainMetrics() (*firestore.ExplainMetrics, error)`

### `CollectionIterator`

Align with the SDK: official iterator has `Next`, `GetAll`, `PageInfo`; it does **not** have `Stop`.

- [ ] `GetAll() ([]*firestore.CollectionRef, error)`
- [ ] `PageInfo() *iterator.PageInfo`
- [ ] Remove or replace `Stop()` on the mock interface to avoid misleading callers (breaking change — semver major).

### `Transaction`

- [x] `Documents(q Query) DocumentIterator` (delegates to `*firestore.Transaction.Documents`; accepts a `Query` or `CollectionRef` since `CollectionRef` embeds `Query`)
- [x] `DocumentRefs(coll CollectionRef) DocumentRefIterator`
- [ ] `Execute(p *firestore.Pipeline) (*firestore.PipelineResultIterator, error)` (if Pipeline is supported)
- [ ] `WithReadOptions(opts ...firestore.ReadOption)`

### `AggregationQuery` / results

Current: `WithCount` + `Get` + `AggregationResult.Count`. The SDK also provides:

- [ ] `WithSum` / `WithSumPath`
- [ ] `WithAvg` / `WithAvgPath`
- [ ] `GetResponse(ctx) (*firestore.AggregationResponse, error)` (explain metrics, etc.)
- [ ] `Transaction(tx *firestore.Transaction) *AggregationQuery`
- [ ] Extend the `AggregationResult` interface (e.g. `Get(alias) (interface{}, error)` or per-aggregation methods) to match the full `map[string]interface{}` surface.

### `WriteBatch`

- [ ] `Commit` with commit-response options (`firestore.WithCommitResponseTo`) if you need full parity with commit overloads.

### Package constructors & options (optional)

Not required if callers always construct `*firestore.Client` themselves:

- [ ] Document that `NewClient` / database ID / REST client remain in the official package.
- [ ] Or provide a thin helper that forwards to `firestore.NewClient` and returns `FirestoreClient`.

### Pipeline, VectorQuery, Expression, composite filters

For **true** 100% against pkg.go.dev, add interface layers for (at least):

- [ ] `Pipeline`, `PipelineSource`, `PipelineResult`, `PipelineResultIterator`, `PipelineSnapshot`
- [ ] `VectorQuery` if split from `Query`
- [ ] `Expression`, `PropertyFilter`, `EntityFilter`, `BooleanExpression`, etc., as needed by your product

This is a large surface; treat it as a separate phase prioritized by real usage.

---

## Engineering & release

- [ ] CI (e.g. GitHub Actions): `go test`, `go vet`, coverage.
- [ ] `CHANGELOG.md` + semver tags; breaking-change policy for interface edits.
- [ ] Adjust `Makefile` for a **library** module (focus targets on test / generate / lint; avoid implying an application binary).
- [ ] Periodically upgrade `cloud.google.com/go/firestore` and re-run `go generate` + tests.

---

## References

- Official package: [pkg.go.dev/cloud.google.com/go/firestore](https://pkg.go.dev/cloud.google.com/go/firestore)
- Version used by this module: see `require cloud.google.com/go/firestore` in `go.mod`.
