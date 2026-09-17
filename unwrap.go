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
		return v.q, nil
	case *collectionRefWrapper:
		return v.ref, nil
	case *collectionGroupRefWrapper:
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
		return w.aq, nil
	}
	return nil, fmt.Errorf("%w: AggregationQuery %T", ErrForeignImplementation, aq)
}

// unwrapCursorArgs converts fsmock.DocumentSnapshot values to *firestore.DocumentSnapshot
// so SDK cursor helpers (StartAt/StartAfter/EndAt/EndBefore) recognize them.
func unwrapCursorArgs(args []any) []any {
	if len(args) == 0 {
		return args
	}
	out := make([]any, len(args))
	for i, a := range args {
		switch v := a.(type) {
		case DocumentSnapshot:
			if v == nil {
				out[i] = a
				continue
			}
			if ref := v.Reference(); ref != nil {
				out[i] = ref
				continue
			}
			out[i] = a
		case *firestore.DocumentSnapshot:
			out[i] = v
		default:
			out[i] = a
		}
	}
	return out
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
