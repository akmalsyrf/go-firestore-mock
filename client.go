package fsmock

import (
	"context"

	"cloud.google.com/go/firestore"
)

// Client abstracts *firestore.Client for dependency injection and mocking.
type Client interface {
	Collection(path string) CollectionRef
	CollectionGroup(collectionID string) CollectionGroupRef
	Doc(path string) DocumentRef
	DocFromFullPath(fullPath string) DocumentRef
	Close() error
	BulkWriter(ctx context.Context) BulkWriter
	Batch() WriteBatch
	RunTransaction(ctx context.Context, f func(context.Context, Transaction) error, opts ...firestore.TransactionOption) error
	Collections(ctx context.Context) CollectionIterator
	GetAll(ctx context.Context, docRefs []DocumentRef) ([]DocumentSnapshot, error)
	WithReadOptions(opts ...firestore.ReadOption) Client
	WithAlwaysUseImplicitOrderBy(b bool) Client
	Pipeline() PipelineSource
}

type clientWrapper struct {
	client *firestore.Client
}

func (w *clientWrapper) Collection(path string) CollectionRef {
	return newCollectionRef(w.client.Collection(path))
}

func (w *clientWrapper) CollectionGroup(collectionID string) CollectionGroupRef {
	return newCollectionGroupRef(w.client.CollectionGroup(collectionID))
}

func (w *clientWrapper) Doc(path string) DocumentRef {
	return newDocumentRef(w.client.Doc(path))
}

func (w *clientWrapper) DocFromFullPath(fullPath string) DocumentRef {
	return newDocumentRef(w.client.DocFromFullPath(fullPath))
}

func (w *clientWrapper) Close() error {
	return w.client.Close()
}

func (w *clientWrapper) BulkWriter(ctx context.Context) BulkWriter {
	return &bulkWriterWrapper{bw: w.client.BulkWriter(ctx)}
}

func (w *clientWrapper) Batch() WriteBatch {
	return &writeBatchWrapper{wb: w.client.Batch()}
}

func (w *clientWrapper) RunTransaction(ctx context.Context, f func(context.Context, Transaction) error, opts ...firestore.TransactionOption) error {
	return w.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return f(ctx, &transactionWrapper{tx: tx})
	}, opts...)
}

func (w *clientWrapper) Collections(ctx context.Context) CollectionIterator {
	return newCollectionIterator(w.client.Collections(ctx))
}

func (w *clientWrapper) GetAll(ctx context.Context, docRefs []DocumentRef) ([]DocumentSnapshot, error) {
	refs, err := wrapDocumentRefs(docRefs)
	if err != nil {
		return nil, err
	}
	snaps, err := w.client.GetAll(ctx, refs)
	if err != nil {
		return nil, err
	}
	return wrapSnapshots(snaps), nil
}

func (w *clientWrapper) WithReadOptions(opts ...firestore.ReadOption) Client {
	return &clientWrapper{client: w.client.WithReadOptions(opts...)}
}

func (w *clientWrapper) WithAlwaysUseImplicitOrderBy(b bool) Client {
	return &clientWrapper{client: w.client.WithAlwaysUseImplicitOrderBy(b)}
}

func (w *clientWrapper) Pipeline() PipelineSource {
	return &pipelineSourceWrapper{ps: w.client.Pipeline()}
}
