package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

// CollectionRef abstracts Firestore collection behavior used by repos.
// It also behaves like a Query (Where, Documents).
//
//go:generate mockgen -source=collection.go -destination=collection_mock.go -package=firestore
type CollectionRef interface {
	Query
	Doc(id string) DocumentRef
	Add(ctx context.Context, data any) (*firestore.DocumentRef, *firestore.WriteResult, error)
	NewDoc() DocumentRef
	Parent() DocumentRef
	ID() string
	Path() string
}

type collectionRefWrapper struct{ ref *firestore.CollectionRef }

func (w *collectionRefWrapper) Doc(id string) DocumentRef {
	return &documentRefWrapper{ref: w.ref.Doc(id)}
}

func (w *collectionRefWrapper) Add(ctx context.Context, data any) (*firestore.DocumentRef, *firestore.WriteResult, error) {
	return w.ref.Add(ctx, data)
}

func (w *collectionRefWrapper) Where(path string, op string, value any) Query {
	return &queryWrapper{q: w.ref.Where(path, op, value)}
}

func (w *collectionRefWrapper) Documents(ctx context.Context) DocumentIterator {
	return &documentIteratorWrapper{iter: w.ref.Documents(ctx)}
}

func (w *collectionRefWrapper) OrderBy(path string, dir firestore.Direction) Query {
	return &queryWrapper{q: w.ref.OrderBy(path, dir)}
}

func (w *collectionRefWrapper) Limit(n int) Query {
	return &queryWrapper{q: w.ref.Limit(n)}
}

func (w *collectionRefWrapper) LimitToLast(n int) Query {
	return &queryWrapper{q: w.ref.LimitToLast(n)}
}

func (w *collectionRefWrapper) Offset(n int) Query {
	return &queryWrapper{q: w.ref.Offset(n)}
}

func (w *collectionRefWrapper) StartAt(docSnapshotOrFieldValues ...any) Query {
	return &queryWrapper{q: w.ref.StartAt(docSnapshotOrFieldValues...)}
}

func (w *collectionRefWrapper) StartAfter(docSnapshotOrFieldValues ...any) Query {
	return &queryWrapper{q: w.ref.StartAfter(docSnapshotOrFieldValues...)}
}

func (w *collectionRefWrapper) EndAt(docSnapshotOrFieldValues ...any) Query {
	return &queryWrapper{q: w.ref.EndAt(docSnapshotOrFieldValues...)}
}

func (w *collectionRefWrapper) EndBefore(docSnapshotOrFieldValues ...any) Query {
	return &queryWrapper{q: w.ref.EndBefore(docSnapshotOrFieldValues...)}
}

func (w *collectionRefWrapper) Select(paths ...string) Query {
	return &queryWrapper{q: w.ref.Select(paths...)}
}

func (w *collectionRefWrapper) Snapshots(ctx context.Context) *firestore.QuerySnapshotIterator {
	return w.ref.Snapshots(ctx)
}

func (w *collectionRefWrapper) NewDoc() DocumentRef {
	return &documentRefWrapper{ref: w.ref.NewDoc()}
}

func (w *collectionRefWrapper) Parent() DocumentRef {
	if w.ref.Parent == nil {
		return nil
	}
	return &documentRefWrapper{ref: w.ref.Parent}
}

func (w *collectionRefWrapper) ID() string {
	return w.ref.ID
}

func (w *collectionRefWrapper) Path() string {
	return w.ref.Path
}
