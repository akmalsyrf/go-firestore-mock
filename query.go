package fsmock

import (
	"context"

	"cloud.google.com/go/firestore"
)

// Query abstracts firestore.Query behavior.
// CollectionRef and CollectionGroupRef embed Query.
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
	FindNearest(vectorField string, queryVector any, limit int, measure firestore.DistanceMeasure, options *firestore.FindNearestOptions) VectorQuery
	FindNearestPath(vectorFieldPath firestore.FieldPath, queryVector any, limit int, measure firestore.DistanceMeasure, options *firestore.FindNearestOptions) VectorQuery
	Serialize() ([]byte, error)
	Deserialize(bytes []byte) (Query, error)
	WithReadOptions(opts ...firestore.ReadOption) Query
	WithRunOptions(opts ...firestore.RunOption) Query
	Pipeline() Pipeline
}

type queryWrapper struct {
	q firestore.Query
}

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
	return newDocumentIterator(w.q.Documents(ctx))
}

func (w *queryWrapper) Snapshots(ctx context.Context) QuerySnapshotIterator {
	return &querySnapshotIteratorWrapper{iter: w.q.Snapshots(ctx)}
}

func (w *queryWrapper) NewAggregationQuery() AggregationQuery {
	q := w.q
	return newAggregationQuery(q.NewAggregationQuery())
}

func (w *queryWrapper) FindNearest(vectorField string, queryVector any, limit int, measure firestore.DistanceMeasure, options *firestore.FindNearestOptions) VectorQuery {
	return newVectorQuery(w.q.FindNearest(vectorField, queryVector, limit, measure, options))
}

func (w *queryWrapper) FindNearestPath(vectorFieldPath firestore.FieldPath, queryVector any, limit int, measure firestore.DistanceMeasure, options *firestore.FindNearestOptions) VectorQuery {
	return newVectorQuery(w.q.FindNearestPath(vectorFieldPath, queryVector, limit, measure, options))
}

func (w *queryWrapper) Serialize() ([]byte, error) {
	return w.q.Serialize()
}

func (w *queryWrapper) Deserialize(bytes []byte) (Query, error) {
	q, err := w.q.Deserialize(bytes)
	if err != nil {
		return nil, err
	}
	return &queryWrapper{q: q}, nil
}

func (w *queryWrapper) WithReadOptions(opts ...firestore.ReadOption) Query {
	q := w.q
	(&q).WithReadOptions(opts...)
	return &queryWrapper{q: q}
}

func (w *queryWrapper) WithRunOptions(opts ...firestore.RunOption) Query {
	return &queryWrapper{q: w.q.WithRunOptions(opts...)}
}

func (w *queryWrapper) Pipeline() Pipeline {
	return newPipeline(w.q.Pipeline())
}
