package fsmock

import (
	"fmt"

	"cloud.google.com/go/firestore"
)

// toFirestoreQueryer extracts the real firestore.Queryer behind a Query wrapper.
func toFirestoreQueryer(q Query) (firestore.Queryer, error) {
	if q == nil {
		return nil, fmt.Errorf("%w: Query", ErrNilArgument)
	}
	switch v := q.(type) {
	case *queryWrapper:
		if v.err != nil {
			return nil, v.err
		}
		return v.q, nil
	case *collectionRefWrapper:
		if v.err != nil {
			return nil, v.err
		}
		return v.ref, nil
	case *collectionGroupRefWrapper:
		if v.err != nil {
			return nil, v.err
		}
		return v.ref, nil
	}
	if cr, ok := q.(CollectionRef); ok {
		if ref := cr.Reference(); ref != nil {
			return ref, nil
		}
	}
	return nil, fmt.Errorf("%w: Query %T (use a wrapper from NewClient)", ErrForeignImplementation, q)
}

// toDocumentRef extracts *firestore.DocumentRef from a DocumentRef wrapper.
func toDocumentRef(d DocumentRef) (*firestore.DocumentRef, error) {
	if d == nil {
		return nil, fmt.Errorf("%w: DocumentRef", ErrNilArgument)
	}
	if w, ok := d.(*documentRefWrapper); ok {
		return w.ref, nil
	}
	if ref := d.Reference(); ref != nil {
		return ref, nil
	}
	return nil, fmt.Errorf("%w: DocumentRef %T has nil Reference()", ErrForeignImplementation, d)
}

// toCollectionRef extracts *firestore.CollectionRef from a CollectionRef wrapper.
func toCollectionRef(c CollectionRef) (*firestore.CollectionRef, error) {
	if c == nil {
		return nil, fmt.Errorf("%w: CollectionRef", ErrNilArgument)
	}
	if w, ok := c.(*collectionRefWrapper); ok {
		return w.ref, nil
	}
	if ref := c.Reference(); ref != nil {
		return ref, nil
	}
	return nil, fmt.Errorf("%w: CollectionRef %T has nil Reference()", ErrForeignImplementation, c)
}

// toTransaction extracts *firestore.Transaction from a Transaction wrapper.
func toTransaction(t Transaction) (*firestore.Transaction, error) {
	if t == nil {
		return nil, fmt.Errorf("%w: Transaction", ErrNilArgument)
	}
	if w, ok := t.(*transactionWrapper); ok {
		return w.tx, nil
	}
	return nil, fmt.Errorf("%w: Transaction %T", ErrForeignImplementation, t)
}

// toPipeline extracts *firestore.Pipeline from a Pipeline wrapper.
func toPipeline(p Pipeline) (*firestore.Pipeline, error) {
	if p == nil {
		return nil, fmt.Errorf("%w: Pipeline", ErrNilArgument)
	}
	if w, ok := p.(*pipelineWrapper); ok {
		if w.err != nil {
			return nil, w.err
		}
		return w.p, nil
	}
	return nil, fmt.Errorf("%w: Pipeline %T", ErrForeignImplementation, p)
}

// toAggregationQuery extracts *firestore.AggregationQuery from a wrapper.
func toAggregationQuery(aq AggregationQuery) (*firestore.AggregationQuery, error) {
	if aq == nil {
		return nil, fmt.Errorf("%w: AggregationQuery", ErrNilArgument)
	}
	if w, ok := aq.(*aggregationQueryWrapper); ok {
		if w.err != nil {
			return nil, w.err
		}
		return w.aq, nil
	}
	return nil, fmt.Errorf("%w: AggregationQuery %T", ErrForeignImplementation, aq)
}

// unwrapCursorArgs converts fsmock.DocumentSnapshot values to *firestore.DocumentSnapshot
// so SDK cursor helpers (StartAt/StartAfter/EndAt/EndBefore) recognize them.
//
// A DocumentSnapshot that cannot produce a non-nil Reference() returns
// ErrForeignImplementation (or ErrNilArgument). Passing it through would make the
// SDK treat the value as a field cursor and silently mis-paginate.
func unwrapCursorArgs(args []any) ([]any, error) {
	if len(args) == 0 {
		return args, nil
	}
	out := make([]any, len(args))
	for i, a := range args {
		switch v := a.(type) {
		case nil:
			out[i] = nil
		case *documentSnapshotWrapper:
			if v == nil {
				return nil, fmt.Errorf("%w: DocumentSnapshot", ErrNilArgument)
			}
			if v.snap == nil {
				return nil, fmt.Errorf("%w: DocumentSnapshot has nil Reference()", ErrForeignImplementation)
			}
			out[i] = v.snap
		case DocumentSnapshot:
			if v == nil {
				return nil, fmt.Errorf("%w: DocumentSnapshot", ErrNilArgument)
			}
			ref := v.Reference()
			if ref == nil {
				return nil, fmt.Errorf("%w: DocumentSnapshot %T has nil Reference()", ErrForeignImplementation, v)
			}
			out[i] = ref
		case *firestore.DocumentSnapshot:
			out[i] = v
		default:
			out[i] = a
		}
	}
	return out, nil
}

func wrapDocumentRefs(refs []DocumentRef) ([]*firestore.DocumentRef, error) {
	out := make([]*firestore.DocumentRef, len(refs))
	for i, r := range refs {
		dr, err := toDocumentRef(r)
		if err != nil {
			return nil, err
		}
		out[i] = dr
	}
	return out, nil
}

func wrapSnapshots(snaps []*firestore.DocumentSnapshot) []DocumentSnapshot {
	out := make([]DocumentSnapshot, len(snaps))
	for i, s := range snaps {
		out[i] = &documentSnapshotWrapper{snap: s}
	}
	return out
}

func newDocumentRef(ref *firestore.DocumentRef) DocumentRef {
	if ref == nil {
		return nil
	}
	return &documentRefWrapper{ref: ref}
}

func newCollectionRef(ref *firestore.CollectionRef) CollectionRef {
	if ref == nil {
		return nil
	}
	return &collectionRefWrapper{
		queryWrapper: queryWrapper{q: ref.Query},
		ref:          ref,
	}
}

func newCollectionGroupRef(ref *firestore.CollectionGroupRef) CollectionGroupRef {
	if ref == nil {
		return nil
	}
	return &collectionGroupRefWrapper{
		queryWrapper: queryWrapper{q: ref.Query},
		ref:          ref,
	}
}

func newQuery(q firestore.Query) Query {
	return &queryWrapper{q: q}
}

func newPipeline(p *firestore.Pipeline) Pipeline {
	if p == nil {
		return nil
	}
	return &pipelineWrapper{p: p}
}

func newPipelineSnapshot(ps *firestore.PipelineSnapshot) PipelineSnapshot {
	if ps == nil {
		return nil
	}
	return &pipelineSnapshotWrapper{ps: ps}
}

func newPipelineResult(pr *firestore.PipelineResult) PipelineResult {
	if pr == nil {
		return nil
	}
	return &pipelineResultWrapper{pr: pr}
}

func newAggregationQuery(aq *firestore.AggregationQuery) AggregationQuery {
	if aq == nil {
		return nil
	}
	return &aggregationQueryWrapper{aq: aq}
}

func newVectorQuery(vq firestore.VectorQuery) VectorQuery {
	return &vectorQueryWrapper{vq: vq}
}

func newDocumentIterator(iter *firestore.DocumentIterator) DocumentIterator {
	if iter == nil {
		return nil
	}
	return &documentIteratorWrapper{iter: iter}
}

func newDocumentRefIterator(iter *firestore.DocumentRefIterator) DocumentRefIterator {
	if iter == nil {
		return nil
	}
	return &documentRefIteratorWrapper{iter: iter}
}

func newCollectionIterator(iter *firestore.CollectionIterator) CollectionIterator {
	if iter == nil {
		return nil
	}
	return &collectionIteratorWrapper{iter: iter}
}
