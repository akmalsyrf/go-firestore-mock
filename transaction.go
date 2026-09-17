package fsmock

import (
	"cloud.google.com/go/firestore"
)

// Transaction abstracts *firestore.Transaction.
//
// Read methods (Get, GetAll, Documents, DocumentRefs) read inside the transaction.
// Firestore requires all reads to happen before any writes within a transaction.
type Transaction interface {
	Get(docRef DocumentRef) (DocumentSnapshot, error)
	GetAll(docRefs []DocumentRef) ([]DocumentSnapshot, error)
	Documents(q Query) (DocumentIterator, error)
	DocumentRefs(coll CollectionRef) (DocumentRefIterator, error)
	Create(docRef DocumentRef, data any) error
	Set(docRef DocumentRef, data any, opts ...firestore.SetOption) error
	Update(docRef DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) error
	Delete(docRef DocumentRef, preconds ...firestore.Precondition) error
	WithReadOptions(opts ...firestore.ReadOption) Transaction
	Execute(p Pipeline, opts ...firestore.ExecuteOption) (PipelineSnapshot, error)
}

type transactionWrapper struct {
	tx *firestore.Transaction
}

func (w *transactionWrapper) Get(docRef DocumentRef) (DocumentSnapshot, error) {
	ref, err := toDocumentRef(docRef)
	if err != nil {
		return nil, err
	}
	snap, err := w.tx.Get(ref)
	if err != nil {
		return nil, err
	}
	return &documentSnapshotWrapper{snap: snap}, nil
}

func (w *transactionWrapper) GetAll(docRefs []DocumentRef) ([]DocumentSnapshot, error) {
	refs, err := wrapDocumentRefs(docRefs)
	if err != nil {
		return nil, err
	}
	snaps, err := w.tx.GetAll(refs)
	if err != nil {
		return nil, err
	}
	return wrapSnapshots(snaps), nil
}

func (w *transactionWrapper) Create(docRef DocumentRef, data any) error {
	ref, err := toDocumentRef(docRef)
	if err != nil {
		return err
	}
	return w.tx.Create(ref, data)
}

func (w *transactionWrapper) Set(docRef DocumentRef, data any, opts ...firestore.SetOption) error {
	ref, err := toDocumentRef(docRef)
	if err != nil {
		return err
	}
	return w.tx.Set(ref, data, opts...)
}

func (w *transactionWrapper) Update(docRef DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) error {
	ref, err := toDocumentRef(docRef)
	if err != nil {
		return err
	}
	return w.tx.Update(ref, updates, preconds...)
}

func (w *transactionWrapper) Delete(docRef DocumentRef, preconds ...firestore.Precondition) error {
	ref, err := toDocumentRef(docRef)
	if err != nil {
		return err
	}
	return w.tx.Delete(ref, preconds...)
}

// Documents converts q to the underlying firestore.Queryer and delegates.
// q must be a package wrapper (or CollectionRef with a real Reference()).
// Custom mocks that only satisfy Query cannot be used with the production wrapper;
// mock Transaction directly in unit tests instead.
func (w *transactionWrapper) Documents(q Query) (DocumentIterator, error) {
	queryer, err := toFirestoreQueryer(q)
	if err != nil {
		return nil, err
	}
	return newDocumentIterator(w.tx.Documents(queryer)), nil
}

func (w *transactionWrapper) DocumentRefs(coll CollectionRef) (DocumentRefIterator, error) {
	ref, err := toCollectionRef(coll)
	if err != nil {
		return nil, err
	}
	return newDocumentRefIterator(w.tx.DocumentRefs(ref)), nil
}

func (w *transactionWrapper) WithReadOptions(opts ...firestore.ReadOption) Transaction {
	return &transactionWrapper{tx: w.tx.WithReadOptions(opts...)}
}

func (w *transactionWrapper) Execute(p Pipeline, opts ...firestore.ExecuteOption) (PipelineSnapshot, error) {
	pl, err := toPipeline(p)
	if err != nil {
		return nil, err
	}
	return newPipelineSnapshot(w.tx.Execute(pl, opts...)), nil
}
