package fsmock

import (
	"context"

	"cloud.google.com/go/firestore"
)

// CollectionRef abstracts *firestore.CollectionRef.
//
// WithReadOptions is inherited from Query (returns Query). The wrapper
// implementation still applies options to the underlying *firestore.CollectionRef
// and returns a CollectionRef value that also satisfies Query.
type CollectionRef interface {
	Query
	Doc(id string) DocumentRef
	Add(ctx context.Context, data any) (DocumentRef, *firestore.WriteResult, error)
	NewDoc() DocumentRef
	DocumentRefs(ctx context.Context) DocumentRefIterator
	Parent() DocumentRef
	Reference() *firestore.CollectionRef
	ID() string
	Path() string
}

// CollectionGroupRef abstracts *firestore.CollectionGroupRef.
type CollectionGroupRef interface {
	Query
	GetPartitionedQueries(ctx context.Context, partitionCount int) ([]Query, error)
	Reference() *firestore.CollectionGroupRef
}

// collectionRefWrapper embeds queryWrapper so Query methods are not duplicated.
type collectionRefWrapper struct {
	queryWrapper
	ref *firestore.CollectionRef
}

func (w *collectionRefWrapper) Doc(id string) DocumentRef {
	return newDocumentRef(w.ref.Doc(id))
}

func (w *collectionRefWrapper) Add(ctx context.Context, data any) (DocumentRef, *firestore.WriteResult, error) {
	ref, wr, err := w.ref.Add(ctx, data)
	if err != nil {
		return nil, wr, err
	}
	return newDocumentRef(ref), wr, nil
}

func (w *collectionRefWrapper) NewDoc() DocumentRef {
	return newDocumentRef(w.ref.NewDoc())
}

func (w *collectionRefWrapper) DocumentRefs(ctx context.Context) DocumentRefIterator {
	return newDocumentRefIterator(w.ref.DocumentRefs(ctx))
}

func (w *collectionRefWrapper) Parent() DocumentRef {
	return newDocumentRef(w.ref.Parent)
}

func (w *collectionRefWrapper) Reference() *firestore.CollectionRef {
	return w.ref
}

func (w *collectionRefWrapper) ID() string {
	return w.ref.ID
}

func (w *collectionRefWrapper) Path() string {
	return w.ref.Path
}

// WithReadOptions overrides the embedded queryWrapper method so options apply
// to a cloned CollectionRef (and its shared Query.readSettings), leaving the
// receiver unmodified.
func (w *collectionRefWrapper) WithReadOptions(opts ...firestore.ReadOption) Query {
	clone := cloneCollectionRefWithFreshReadSettings(w.ref)
	return newCollectionRef(clone.WithReadOptions(opts...))
}

type collectionGroupRefWrapper struct {
	queryWrapper
	ref *firestore.CollectionGroupRef
}

func (w *collectionGroupRefWrapper) GetPartitionedQueries(ctx context.Context, partitionCount int) ([]Query, error) {
	qs, err := w.ref.GetPartitionedQueries(ctx, partitionCount)
	if err != nil {
		return nil, err
	}
	out := make([]Query, len(qs))
	for i := range qs {
		out[i] = newQuery(qs[i])
	}
	return out, nil
}

func (w *collectionGroupRefWrapper) Reference() *firestore.CollectionGroupRef {
	return w.ref
}
