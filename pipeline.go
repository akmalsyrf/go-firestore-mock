package fsmock

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

// PipelineSource abstracts *firestore.PipelineSource.
// Option types remain pass-through SDK values.
type PipelineSource interface {
	Collection(path string, opts ...firestore.CollectionOption) Pipeline
	CollectionGroup(collectionID string, opts ...firestore.CollectionGroupOption) Pipeline
	Database(opts ...firestore.DatabaseOption) Pipeline
	Documents(refs []DocumentRef, opts ...firestore.DocumentsOption) (Pipeline, error)
	CreateFromQuery(query Query) (Pipeline, error)
	CreateFromAggregationQuery(query AggregationQuery) (Pipeline, error)
	Literals(documents []map[string]any, opts ...firestore.LiteralsOption) Pipeline
}

// Pipeline abstracts *firestore.Pipeline stage chaining and execution.
type Pipeline interface {
	Execute(ctx context.Context, opts ...firestore.ExecuteOption) PipelineSnapshot
	WithReadOptions(opts ...firestore.ReadOption) Pipeline
	Limit(limit int, opts ...firestore.LimitOption) Pipeline
	Sort(orders []firestore.Ordering, opts ...firestore.SortOption) Pipeline
	Offset(offset int, opts ...firestore.OffsetOption) Pipeline
	Select(fields []any, opts ...firestore.SelectOption) Pipeline
	Distinct(fields []any, opts ...firestore.DistinctOption) Pipeline
	AddFields(fields []firestore.Selectable, opts ...firestore.AddFieldsOption) Pipeline
	RemoveFields(fields []any, opts ...firestore.RemoveFieldsOption) Pipeline
	Where(condition firestore.BooleanExpression, opts ...firestore.WhereOption) Pipeline
	Aggregate(accumulators []*firestore.AliasedAggregate, opts ...firestore.AggregateOption) Pipeline
	Unnest(field firestore.Selectable, opts ...firestore.UnnestOption) Pipeline
	UnnestWithAlias(fieldpath any, alias string, opts ...firestore.UnnestOption) Pipeline
	Union(other Pipeline, opts ...firestore.UnionOption) (Pipeline, error)
	Sample(sampler *firestore.Sampler, opts ...firestore.SampleOption) Pipeline
	ReplaceWith(fieldpathOrExpr any, opts ...firestore.ReplaceWithOption) Pipeline
	FindNearest(vectorField any, queryVector any, measure firestore.PipelineDistanceMeasure, opts ...firestore.FindNearestOption) Pipeline
	Search(opts ...firestore.SearchOption) Pipeline
	RawStage(name string, args []any, opts ...firestore.StageOption) Pipeline
	Update(opts ...firestore.UpdateOption) Pipeline
	Delete(opts ...firestore.DeleteOption) Pipeline
	ToScalarExpression() firestore.Expression
	ToArrayExpression() firestore.Expression
	Define(variables []*firestore.AliasedExpression, opts ...firestore.DefineOption) Pipeline
	Reference() *firestore.Pipeline
}

// PipelineSnapshot abstracts *firestore.PipelineSnapshot.
type PipelineSnapshot interface {
	Results() PipelineResultIterator
	ExecutionTime() (*time.Time, error)
	ExplainStats() *firestore.ExplainStats
}

// PipelineResult abstracts *firestore.PipelineResult.
type PipelineResult interface {
	Ref() DocumentRef
	CreateTime() *time.Time
	UpdateTime() *time.Time
	ExecutionTime() *time.Time
	Exists() bool
	Data() map[string]any
	DataTo(v any) error
}

// PipelineResultIterator abstracts *firestore.PipelineResultIterator.
type PipelineResultIterator interface {
	Next() (PipelineResult, error)
	Stop()
	GetAll() ([]PipelineResult, error)
}

type pipelineSourceWrapper struct {
	ps *firestore.PipelineSource
}

func (w *pipelineSourceWrapper) Collection(path string, opts ...firestore.CollectionOption) Pipeline {
	return newPipeline(w.ps.Collection(path, opts...))
}

func (w *pipelineSourceWrapper) CollectionGroup(collectionID string, opts ...firestore.CollectionGroupOption) Pipeline {
	return newPipeline(w.ps.CollectionGroup(collectionID, opts...))
}

func (w *pipelineSourceWrapper) Database(opts ...firestore.DatabaseOption) Pipeline {
	return newPipeline(w.ps.Database(opts...))
}

func (w *pipelineSourceWrapper) Documents(refs []DocumentRef, opts ...firestore.DocumentsOption) (Pipeline, error) {
	sdkRefs, err := wrapDocumentRefs(refs)
	if err != nil {
		return nil, err
	}
	return newPipeline(w.ps.Documents(sdkRefs, opts...)), nil
}

func (w *pipelineSourceWrapper) CreateFromQuery(query Query) (Pipeline, error) {
	queryer, err := toFirestoreQueryer(query)
	if err != nil {
		return nil, err
	}
	return newPipeline(w.ps.CreateFromQuery(queryer)), nil
}

func (w *pipelineSourceWrapper) CreateFromAggregationQuery(query AggregationQuery) (Pipeline, error) {
	aq, err := toAggregationQuery(query)
	if err != nil {
		return nil, err
	}
	return newPipeline(w.ps.CreateFromAggregationQuery(aq)), nil
}

func (w *pipelineSourceWrapper) Literals(documents []map[string]any, opts ...firestore.LiteralsOption) Pipeline {
	return newPipeline(w.ps.Literals(documents, opts...))
}

type pipelineWrapper struct {
	p   *firestore.Pipeline
	err error
}

func (w *pipelineWrapper) withP(p *firestore.Pipeline) Pipeline {
	return &pipelineWrapper{p: p, err: w.err}
}

func (w *pipelineWrapper) Reference() *firestore.Pipeline { return w.p }

func (w *pipelineWrapper) Execute(ctx context.Context, opts ...firestore.ExecuteOption) PipelineSnapshot {
	if w.err != nil {
		return errPipelineSnapshot{err: w.err}
	}
	return newPipelineSnapshot(w.p.Execute(ctx, opts...))
}

func (w *pipelineWrapper) WithReadOptions(opts ...firestore.ReadOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.WithReadOptions(opts...))
}

func (w *pipelineWrapper) Limit(limit int, opts ...firestore.LimitOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.Limit(limit, opts...))
}

func (w *pipelineWrapper) Sort(orders []firestore.Ordering, opts ...firestore.SortOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.Sort(orders, opts...))
}

func (w *pipelineWrapper) Offset(offset int, opts ...firestore.OffsetOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.Offset(offset, opts...))
}

func (w *pipelineWrapper) Select(fields []any, opts ...firestore.SelectOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.Select(fields, opts...))
}

func (w *pipelineWrapper) Distinct(fields []any, opts ...firestore.DistinctOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.Distinct(fields, opts...))
}

func (w *pipelineWrapper) AddFields(fields []firestore.Selectable, opts ...firestore.AddFieldsOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.AddFields(fields, opts...))
}

func (w *pipelineWrapper) RemoveFields(fields []any, opts ...firestore.RemoveFieldsOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.RemoveFields(fields, opts...))
}

func (w *pipelineWrapper) Where(condition firestore.BooleanExpression, opts ...firestore.WhereOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.Where(condition, opts...))
}

func (w *pipelineWrapper) Aggregate(accumulators []*firestore.AliasedAggregate, opts ...firestore.AggregateOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.Aggregate(accumulators, opts...))
}

func (w *pipelineWrapper) Unnest(field firestore.Selectable, opts ...firestore.UnnestOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.Unnest(field, opts...))
}

func (w *pipelineWrapper) UnnestWithAlias(fieldpath any, alias string, opts ...firestore.UnnestOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.UnnestWithAlias(fieldpath, alias, opts...))
}

func (w *pipelineWrapper) Union(other Pipeline, opts ...firestore.UnionOption) (Pipeline, error) {
	if w.err != nil {
		return nil, w.err
	}
	op, err := toPipeline(other)
	if err != nil {
		return nil, err
	}
	return w.withP(w.p.Union(op, opts...)), nil
}

func (w *pipelineWrapper) Sample(sampler *firestore.Sampler, opts ...firestore.SampleOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.Sample(sampler, opts...))
}

func (w *pipelineWrapper) ReplaceWith(fieldpathOrExpr any, opts ...firestore.ReplaceWithOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.ReplaceWith(fieldpathOrExpr, opts...))
}

func (w *pipelineWrapper) FindNearest(vectorField any, queryVector any, measure firestore.PipelineDistanceMeasure, opts ...firestore.FindNearestOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.FindNearest(vectorField, queryVector, measure, opts...))
}

func (w *pipelineWrapper) Search(opts ...firestore.SearchOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.Search(opts...))
}

func (w *pipelineWrapper) RawStage(name string, args []any, opts ...firestore.StageOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.RawStage(name, args, opts...))
}

func (w *pipelineWrapper) Update(opts ...firestore.UpdateOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.Update(opts...))
}

func (w *pipelineWrapper) Delete(opts ...firestore.DeleteOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.Delete(opts...))
}

func (w *pipelineWrapper) ToScalarExpression() firestore.Expression {
	if w.err != nil || w.p == nil {
		return nil
	}
	return w.p.ToScalarExpression()
}

func (w *pipelineWrapper) ToArrayExpression() firestore.Expression {
	if w.err != nil || w.p == nil {
		return nil
	}
	return w.p.ToArrayExpression()
}

func (w *pipelineWrapper) Define(variables []*firestore.AliasedExpression, opts ...firestore.DefineOption) Pipeline {
	if w.err != nil {
		return w
	}
	return w.withP(w.p.Define(variables, opts...))
}

// errPipelineSnapshot surfaces a deferred cursor unwrap error from Pipeline.Execute.
type errPipelineSnapshot struct{ err error }

func (e errPipelineSnapshot) Results() PipelineResultIterator {
	return errPipelineResultIterator(e)
}
func (e errPipelineSnapshot) ExecutionTime() (*time.Time, error) { return nil, e.err }
func (e errPipelineSnapshot) ExplainStats() *firestore.ExplainStats {
	return nil
}

type errPipelineResultIterator struct{ err error }

func (e errPipelineResultIterator) Next() (PipelineResult, error) { return nil, e.err }
func (e errPipelineResultIterator) Stop()                         {}
func (e errPipelineResultIterator) GetAll() ([]PipelineResult, error) {
	return nil, e.err
}

type pipelineSnapshotWrapper struct {
	ps *firestore.PipelineSnapshot
}

func (w *pipelineSnapshotWrapper) Results() PipelineResultIterator {
	return &pipelineResultIteratorWrapper{iter: w.ps.Results()}
}

func (w *pipelineSnapshotWrapper) ExecutionTime() (*time.Time, error) {
	return w.ps.ExecutionTime()
}

func (w *pipelineSnapshotWrapper) ExplainStats() *firestore.ExplainStats {
	return w.ps.ExplainStats()
}

type pipelineResultWrapper struct {
	pr *firestore.PipelineResult
}

func (w *pipelineResultWrapper) Ref() DocumentRef {
	return newDocumentRef(w.pr.Ref())
}

func (w *pipelineResultWrapper) CreateTime() *time.Time    { return w.pr.CreateTime() }
func (w *pipelineResultWrapper) UpdateTime() *time.Time    { return w.pr.UpdateTime() }
func (w *pipelineResultWrapper) ExecutionTime() *time.Time { return w.pr.ExecutionTime() }
func (w *pipelineResultWrapper) Exists() bool              { return w.pr.Exists() }
func (w *pipelineResultWrapper) Data() map[string]any      { return w.pr.Data() }
func (w *pipelineResultWrapper) DataTo(v any) error        { return w.pr.DataTo(v) }

type pipelineResultIteratorWrapper struct {
	iter *firestore.PipelineResultIterator
}

func (w *pipelineResultIteratorWrapper) Next() (PipelineResult, error) {
	pr, err := w.iter.Next()
	if err != nil {
		return nil, err
	}
	return newPipelineResult(pr), nil
}

func (w *pipelineResultIteratorWrapper) Stop() {
	w.iter.Stop()
}

func (w *pipelineResultIteratorWrapper) GetAll() ([]PipelineResult, error) {
	prs, err := w.iter.GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]PipelineResult, len(prs))
	for i, pr := range prs {
		out[i] = newPipelineResult(pr)
	}
	return out, nil
}
