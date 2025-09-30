package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

//go:generate mockgen -source=client.go -destination=client_mock.go -package=firestore

type FirestoreClient interface {
	Collection(path string) CollectionRef
	Doc(path string) DocumentRef
	Close() error
	BulkWriter(ctx context.Context) BulkWriter
	Batch() *firestore.WriteBatch
	RunTransaction(ctx context.Context, f func(context.Context, *firestore.Transaction) error, opts ...firestore.TransactionOption) error
	Collections(ctx context.Context) *firestore.CollectionIterator
	GetAll(ctx context.Context, docRefs []*firestore.DocumentRef) ([]*firestore.DocumentSnapshot, error)
}

// firebaseClientWrapper wraps real firestore.Client
type firebaseClientWrapper struct {
	client *firestore.Client
}

func (w *firebaseClientWrapper) Collection(path string) CollectionRef {
	return &collectionRefWrapper{ref: w.client.Collection(path)}
}

func (w *firebaseClientWrapper) Doc(path string) DocumentRef {
	return &documentRefWrapper{ref: w.client.Doc(path)}
}

func (w *firebaseClientWrapper) Close() error {
	return w.client.Close()
}

func (w *firebaseClientWrapper) BulkWriter(ctx context.Context) BulkWriter {
	return &bulkWriterWrapper{bw: w.client.BulkWriter(ctx)}
}

func (w *firebaseClientWrapper) Batch() *firestore.WriteBatch {
	return w.client.Batch()
}

func (w *firebaseClientWrapper) RunTransaction(ctx context.Context, f func(context.Context, *firestore.Transaction) error, opts ...firestore.TransactionOption) error {
	return w.client.RunTransaction(ctx, f, opts...)
}

func (w *firebaseClientWrapper) Collections(ctx context.Context) *firestore.CollectionIterator {
	return w.client.Collections(ctx)
}

func (w *firebaseClientWrapper) GetAll(ctx context.Context, docRefs []*firestore.DocumentRef) ([]*firestore.DocumentSnapshot, error) {
	return w.client.GetAll(ctx, docRefs)
}

// NewFirestoreClient wraps real client
func NewFirestoreClient(client *firestore.Client) FirestoreClient {
	return &firebaseClientWrapper{client: client}
}

// Query abstracts Firestore query behavior used by repos (Where, Documents).
type Query interface {
	Where(path string, op string, value any) Query
	OrderBy(path string, dir firestore.Direction) Query
	Limit(n int) Query
	LimitToLast(n int) Query
	Offset(n int) Query
	StartAt(docSnapshotOrFieldValues ...any) Query
	StartAfter(docSnapshotOrFieldValues ...any) Query
	EndAt(docSnapshotOrFieldValues ...any) Query
	EndBefore(docSnapshotOrFieldValues ...any) Query
	Select(paths ...string) Query
	Documents(ctx context.Context) DocumentIterator
	Snapshots(ctx context.Context) *firestore.QuerySnapshotIterator
}

type queryWrapper struct{ q firestore.Query }

func (w *queryWrapper) Where(path string, op string, value any) Query {
	return &queryWrapper{q: w.q.Where(path, op, value)}
}

func (w *queryWrapper) OrderBy(path string, dir firestore.Direction) Query {
	return &queryWrapper{q: w.q.OrderBy(path, dir)}
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

func (w *queryWrapper) Documents(ctx context.Context) DocumentIterator {
	return &documentIteratorWrapper{iter: w.q.Documents(ctx)}
}

func (w *queryWrapper) Snapshots(ctx context.Context) *firestore.QuerySnapshotIterator {
	return w.q.Snapshots(ctx)
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
