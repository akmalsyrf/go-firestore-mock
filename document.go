package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

// DocumentRef abstracts Firestore document behavior used by repos.
type DocumentRef interface {
	Set(ctx context.Context, data any, opts ...firestore.SetOption) (*firestore.WriteResult, error)
	Get(ctx context.Context) (*firestore.DocumentSnapshot, error)
	Delete(ctx context.Context, opts ...firestore.Precondition) (*firestore.WriteResult, error)
	Update(ctx context.Context, updates []firestore.Update, preconds ...firestore.Precondition) (*firestore.WriteResult, error)
	Create(ctx context.Context, data any) (*firestore.WriteResult, error)
	Collection(path string) CollectionRef
	Collections(ctx context.Context) *firestore.CollectionIterator
	Snapshots(ctx context.Context) *firestore.DocumentSnapshotIterator
	Reference() *firestore.DocumentRef
	ID() string
	Path() string
	Parent() *firestore.CollectionRef
}

type documentRefWrapper struct{ ref *firestore.DocumentRef }

func (w *documentRefWrapper) Set(ctx context.Context, data any, opts ...firestore.SetOption) (*firestore.WriteResult, error) {
	return w.ref.Set(ctx, data, opts...)
}

func (w *documentRefWrapper) Get(ctx context.Context) (*firestore.DocumentSnapshot, error) {
	return w.ref.Get(ctx)
}

func (w *documentRefWrapper) Delete(ctx context.Context, opts ...firestore.Precondition) (*firestore.WriteResult, error) {
	return w.ref.Delete(ctx, opts...)
}

func (w *documentRefWrapper) Update(ctx context.Context, updates []firestore.Update, preconds ...firestore.Precondition) (*firestore.WriteResult, error) {
	return w.ref.Update(ctx, updates, preconds...)
}

func (w *documentRefWrapper) Create(ctx context.Context, data any) (*firestore.WriteResult, error) {
	return w.ref.Create(ctx, data)
}

func (w *documentRefWrapper) Collection(path string) CollectionRef {
	return &collectionRefWrapper{ref: w.ref.Collection(path)}
}

func (w *documentRefWrapper) Collections(ctx context.Context) *firestore.CollectionIterator {
	return w.ref.Collections(ctx)
}

func (w *documentRefWrapper) Snapshots(ctx context.Context) *firestore.DocumentSnapshotIterator {
	return w.ref.Snapshots(ctx)
}

func (w *documentRefWrapper) Reference() *firestore.DocumentRef {
	return w.ref
}

func (w *documentRefWrapper) ID() string {
	return w.ref.ID
}

func (w *documentRefWrapper) Path() string {
	return w.ref.Path
}

func (w *documentRefWrapper) Parent() *firestore.CollectionRef {
	return w.ref.Parent
}
