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
	// SDKQuery returns the underlying firestore.Query value.
	// Named distinctly from CollectionRef.Reference / CollectionGroupRef.Reference
	// because Go interfaces cannot overload by return type.
	SDKQuery() firestore.Query
}

type queryWrapper struct {
	q   firestore.Query
	err error // deferred cursor unwrap (surfaced on execute / serialize)
}

func (w *queryWrapper) withQ(q firestore.Query) *queryWrapper {
	return &queryWrapper{q: q, err: w.err}
}

func (w *queryWrapper) withCursor(args []any, apply func(...any) firestore.Query) Query {
	unwrapped, err := unwrapCursorArgs(args)
	if err != nil {
		if w.err == nil {
			return &queryWrapper{q: w.q, err: err}
		}
		return w
	}
	return w.withQ(apply(unwrapped...))
}

func (w *queryWrapper) Where(path string, op string, value any) Query {
	return w.withQ(w.q.Where(path, op, value))
}

func (w *queryWrapper) WherePath(fp firestore.FieldPath, op string, value any) Query {
	return w.withQ(w.q.WherePath(fp, op, value))
}

func (w *queryWrapper) WhereEntity(ef firestore.EntityFilter) Query {
	return w.withQ(w.q.WhereEntity(ef))
}

func (w *queryWrapper) OrderBy(path string, dir firestore.Direction) Query {
	return w.withQ(w.q.OrderBy(path, dir))
}

func (w *queryWrapper) OrderByPath(fp firestore.FieldPath, dir firestore.Direction) Query {
	return w.withQ(w.q.OrderByPath(fp, dir))
}

func (w *queryWrapper) Limit(n int) Query {
	return w.withQ(w.q.Limit(n))
}

func (w *queryWrapper) LimitToLast(n int) Query {
	return w.withQ(w.q.LimitToLast(n))
}

func (w *queryWrapper) Offset(n int) Query {
	return w.withQ(w.q.Offset(n))
}

func (w *queryWrapper) StartAt(docSnapshotOrFieldValues ...any) Query {
	return w.withCursor(docSnapshotOrFieldValues, w.q.StartAt)
}

func (w *queryWrapper) StartAfter(docSnapshotOrFieldValues ...any) Query {
	return w.withCursor(docSnapshotOrFieldValues, w.q.StartAfter)
}

func (w *queryWrapper) EndAt(docSnapshotOrFieldValues ...any) Query {
	return w.withCursor(docSnapshotOrFieldValues, w.q.EndAt)
}

func (w *queryWrapper) EndBefore(docSnapshotOrFieldValues ...any) Query {
	return w.withCursor(docSnapshotOrFieldValues, w.q.EndBefore)
}

func (w *queryWrapper) Select(paths ...string) Query {
	return w.withQ(w.q.Select(paths...))
}

func (w *queryWrapper) SelectPaths(fieldPaths ...firestore.FieldPath) Query {
	return w.withQ(w.q.SelectPaths(fieldPaths...))
}

func (w *queryWrapper) Documents(ctx context.Context) DocumentIterator {
	if w.err != nil {
		return errDocumentIterator{err: w.err}
	}
	return newDocumentIterator(w.q.Documents(ctx))
}

func (w *queryWrapper) Snapshots(ctx context.Context) QuerySnapshotIterator {
	if w.err != nil {
		return errQuerySnapshotIterator{err: w.err}
	}
	return &querySnapshotIteratorWrapper{iter: w.q.Snapshots(ctx)}
}

func (w *queryWrapper) NewAggregationQuery() AggregationQuery {
	if w.err != nil {
		return &aggregationQueryWrapper{err: w.err}
	}
	q := w.q
	return newAggregationQuery(q.NewAggregationQuery())
}

func (w *queryWrapper) FindNearest(vectorField string, queryVector any, limit int, measure firestore.DistanceMeasure, options *firestore.FindNearestOptions) VectorQuery {
	if w.err != nil {
		return &vectorQueryWrapper{err: w.err}
	}
	return newVectorQuery(w.q.FindNearest(vectorField, queryVector, limit, measure, options))
}

func (w *queryWrapper) FindNearestPath(vectorFieldPath firestore.FieldPath, queryVector any, limit int, measure firestore.DistanceMeasure, options *firestore.FindNearestOptions) VectorQuery {
	if w.err != nil {
		return &vectorQueryWrapper{err: w.err}
	}
	return newVectorQuery(w.q.FindNearestPath(vectorFieldPath, queryVector, limit, measure, options))
}

func (w *queryWrapper) Serialize() ([]byte, error) {
	if w.err != nil {
		return nil, w.err
	}
	return w.q.Serialize()
}

func (w *queryWrapper) Deserialize(bytes []byte) (Query, error) {
	if w.err != nil {
		return nil, w.err
	}
	q, err := w.q.Deserialize(bytes)
	if err != nil {
		return nil, err
	}
	return &queryWrapper{q: q}, nil
}

func (w *queryWrapper) WithReadOptions(opts ...firestore.ReadOption) Query {
	if w.err != nil {
		return w
	}
	q := cloneQueryWithFreshReadSettings(w.q)
	(&q).WithReadOptions(opts...)
	return &queryWrapper{q: q, err: w.err}
}

func (w *queryWrapper) WithRunOptions(opts ...firestore.RunOption) Query {
	return w.withQ(w.q.WithRunOptions(opts...))
}

func (w *queryWrapper) Pipeline() Pipeline {
	if w.err != nil {
		return &pipelineWrapper{err: w.err}
	}
	return newPipeline(w.q.Pipeline())
}

func (w *queryWrapper) SDKQuery() firestore.Query {
	return w.q
}
