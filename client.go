package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

//go:generate mockgen -source=client.go -destination=client_mock.go -package=firestore

type FirestoreClient interface {
	Collection(path string) CollectionRef
	CollectionGroup(collectionID string) Query
	Doc(path string) DocumentRef
	DocFromFullPath(fullPath string) DocumentRef
	Close() error
	BulkWriter(ctx context.Context) BulkWriter
	Batch() WriteBatch
	RunTransaction(ctx context.Context, f func(context.Context, Transaction) error, opts ...firestore.TransactionOption) error
	Collections(ctx context.Context) CollectionIterator
	GetAll(ctx context.Context, docRefs []*firestore.DocumentRef) ([]DocumentSnapshot, error)
}

// firebaseClientWrapper wraps real firestore.Client
type firebaseClientWrapper struct {
	client *firestore.Client
}

func (w *firebaseClientWrapper) Collection(path string) CollectionRef {
	return &collectionRefWrapper{ref: w.client.Collection(path)}
}

func (w *firebaseClientWrapper) CollectionGroup(collectionID string) Query {
	return &queryWrapper{q: w.client.CollectionGroup(collectionID).Query}
}

func (w *firebaseClientWrapper) Doc(path string) DocumentRef {
	return &documentRefWrapper{ref: w.client.Doc(path)}
}

func (w *firebaseClientWrapper) DocFromFullPath(fullPath string) DocumentRef {
	ref := w.client.DocFromFullPath(fullPath)
	if ref == nil {
		return nil
	}
	return &documentRefWrapper{ref: ref}
}

func (w *firebaseClientWrapper) Close() error {
	return w.client.Close()
}

func (w *firebaseClientWrapper) BulkWriter(ctx context.Context) BulkWriter {
	return &bulkWriterWrapper{bw: w.client.BulkWriter(ctx)}
}

func (w *firebaseClientWrapper) Batch() WriteBatch {
	return &writeBatchWrapper{wb: w.client.Batch()}
}

func (w *firebaseClientWrapper) RunTransaction(ctx context.Context, f func(context.Context, Transaction) error, opts ...firestore.TransactionOption) error {
	return w.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return f(ctx, &transactionWrapper{tx: tx})
	}, opts...)
}

func (w *firebaseClientWrapper) Collections(ctx context.Context) CollectionIterator {
	return &collectionIteratorWrapper{iter: w.client.Collections(ctx)}
}

func (w *firebaseClientWrapper) GetAll(ctx context.Context, docRefs []*firestore.DocumentRef) ([]DocumentSnapshot, error) {
	snaps, err := w.client.GetAll(ctx, docRefs)
	if err != nil {
		return nil, err
	}

	result := make([]DocumentSnapshot, len(snaps))
	for i, snap := range snaps {
		result[i] = &documentSnapshotWrapper{snap: snap}
	}
	return result, nil
}

// NewFirestoreClient wraps real client
func NewFirestoreClient(client *firestore.Client) FirestoreClient {
	return &firebaseClientWrapper{client: client}
}

// Query abstracts Firestore query behavior used by repos (Where, Documents, ...).
//
// The interface mirrors the read-side of *firestore.Query and *firestore.CollectionRef;
// CollectionRef embeds Query so any CollectionRef value satisfies Query.
type Query interface {
	Where(path string, op string, value any) Query
	WherePath(fp firestore.FieldPath, op string, value any) Query
	WhereEntity(ef firestore.EntityFilter) Query
	OrderBy(path string, dir firestore.Direction) Query
	OrderByPath(fp firestore.FieldPath, dir firestore.Direction) Query
	Limit(n int) Query
	LimitToLast(n int) Query
	Offset(n int) Query
	StartAt(docSnapshotOrFieldValues ...any) Query
	StartAfter(docSnapshotOrFieldValues ...any) Query
	EndAt(docSnapshotOrFieldValues ...any) Query
	EndBefore(docSnapshotOrFieldValues ...any) Query
	Select(paths ...string) Query
	SelectPaths(fieldPaths ...firestore.FieldPath) Query
	Documents(ctx context.Context) DocumentIterator
	Snapshots(ctx context.Context) QuerySnapshotIterator
	NewAggregationQuery() AggregationQuery
}

type queryWrapper struct{ q firestore.Query }

func (w *queryWrapper) Where(path string, op string, value any) Query {
	return &queryWrapper{q: w.q.Where(path, op, value)}
}

func (w *queryWrapper) WherePath(fp firestore.FieldPath, op string, value any) Query {
	return &queryWrapper{q: w.q.WherePath(fp, op, value)}
}

func (w *queryWrapper) WhereEntity(ef firestore.EntityFilter) Query {
	return &queryWrapper{q: w.q.WhereEntity(ef)}
}

func (w *queryWrapper) OrderBy(path string, dir firestore.Direction) Query {
	return &queryWrapper{q: w.q.OrderBy(path, dir)}
}

func (w *queryWrapper) OrderByPath(fp firestore.FieldPath, dir firestore.Direction) Query {
	return &queryWrapper{q: w.q.OrderByPath(fp, dir)}
}

func (w *queryWrapper) Limit(n int) Query {
	return &queryWrapper{q: w.q.Limit(n)}
}

func (w *queryWrapper) LimitToLast(n int) Query {
	return &queryWrapper{q: w.q.LimitToLast(n)}
}

func (w *queryWrapper) Offset(n int) Query {
	return &queryWrapper{q: w.q.Offset(n)}
}

func (w *queryWrapper) StartAt(docSnapshotOrFieldValues ...any) Query {
	return &queryWrapper{q: w.q.StartAt(docSnapshotOrFieldValues...)}
}

func (w *queryWrapper) StartAfter(docSnapshotOrFieldValues ...any) Query {
	return &queryWrapper{q: w.q.StartAfter(docSnapshotOrFieldValues...)}
}

func (w *queryWrapper) EndAt(docSnapshotOrFieldValues ...any) Query {
	return &queryWrapper{q: w.q.EndAt(docSnapshotOrFieldValues...)}
}

func (w *queryWrapper) EndBefore(docSnapshotOrFieldValues ...any) Query {
	return &queryWrapper{q: w.q.EndBefore(docSnapshotOrFieldValues...)}
}

func (w *queryWrapper) Select(paths ...string) Query {
	return &queryWrapper{q: w.q.Select(paths...)}
}

func (w *queryWrapper) SelectPaths(fieldPaths ...firestore.FieldPath) Query {
	return &queryWrapper{q: w.q.SelectPaths(fieldPaths...)}
}

func (w *queryWrapper) Documents(ctx context.Context) DocumentIterator {
	return &documentIteratorWrapper{iter: w.q.Documents(ctx)}
}

func (w *queryWrapper) Snapshots(ctx context.Context) QuerySnapshotIterator {
	return &querySnapshotIteratorWrapper{iter: w.q.Snapshots(ctx)}
}

func (w *queryWrapper) NewAggregationQuery() AggregationQuery {
	return &aggregationQueryWrapper{aq: w.q.NewAggregationQuery()}
}

// documentIteratorWrapper wraps real firestore.DocumentIterator
type documentIteratorWrapper struct {
	iter *firestore.DocumentIterator
}

func (w *documentIteratorWrapper) Next() (*firestore.DocumentSnapshot, error) {
	return w.iter.Next()
}

func (w *documentIteratorWrapper) Stop() {
	w.iter.Stop()
}

func (w *documentIteratorWrapper) GetAll() ([]*firestore.DocumentSnapshot, error) {
	return w.iter.GetAll()
}

// bulkWriterWrapper wraps real firestore.BulkWriter
type bulkWriterWrapper struct {
	bw *firestore.BulkWriter
}

func (w *bulkWriterWrapper) Create(docRef *firestore.DocumentRef, data interface{}) (*firestore.BulkWriterJob, error) {
	return w.bw.Create(docRef, data)
}

func (w *bulkWriterWrapper) Set(docRef *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) (*firestore.BulkWriterJob, error) {
	return w.bw.Set(docRef, data, opts...)
}

func (w *bulkWriterWrapper) Update(docRef *firestore.DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) (*firestore.BulkWriterJob, error) {
	return w.bw.Update(docRef, updates, preconds...)
}

func (w *bulkWriterWrapper) Delete(docRef *firestore.DocumentRef, preconds ...firestore.Precondition) (*firestore.BulkWriterJob, error) {
	return w.bw.Delete(docRef, preconds...)
}

func (w *bulkWriterWrapper) Flush() {
	w.bw.Flush()
}

func (w *bulkWriterWrapper) End() {
	w.bw.End()
}
