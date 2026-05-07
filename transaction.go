package firestore

import (
	"fmt"

	"cloud.google.com/go/firestore"
)

//go:generate mockgen -source=transaction.go -destination=transaction_mock.go -package=firestore

// Transaction abstracts Firestore Transaction behavior.
//
// Read methods (Get, GetAll, Documents, DocumentRefs) read inside the transaction.
// Firestore requires all reads to happen before any writes within a transaction.
type Transaction interface {
	Get(docRef *firestore.DocumentRef) (DocumentSnapshot, error)
	GetAll(docRefs []*firestore.DocumentRef) ([]DocumentSnapshot, error)
	// Documents returns a DocumentIterator for the given Query or CollectionRef
	// (CollectionRef is a Query because it embeds Query in its interface),
	// matching the *firestore.Transaction.Documents(q Queryer) signature.
	Documents(q Query) DocumentIterator
	// DocumentRefs returns a DocumentRefIterator for the given CollectionRef,
	// including missing documents (those that have sub-documents but no own data).
	DocumentRefs(coll CollectionRef) DocumentRefIterator
	Create(docRef *firestore.DocumentRef, data interface{}) error
	Set(docRef *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) error
	Update(docRef *firestore.DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) error
	Delete(docRef *firestore.DocumentRef, preconds ...firestore.Precondition) error
}

type transactionWrapper struct {
	tx *firestore.Transaction
}

func (w *transactionWrapper) Get(docRef *firestore.DocumentRef) (DocumentSnapshot, error) {
	snap, err := w.tx.Get(docRef)
	if err != nil {
		return nil, err
	}
	return &documentSnapshotWrapper{snap: snap}, nil
}

func (w *transactionWrapper) GetAll(docRefs []*firestore.DocumentRef) ([]DocumentSnapshot, error) {
	snaps, err := w.tx.GetAll(docRefs)
	if err != nil {
		return nil, err
	}

	result := make([]DocumentSnapshot, len(snaps))
	for i, snap := range snaps {
		result[i] = &documentSnapshotWrapper{snap: snap}
	}
	return result, nil
}

func (w *transactionWrapper) Create(docRef *firestore.DocumentRef, data interface{}) error {
	return w.tx.Create(docRef, data)
}

func (w *transactionWrapper) Set(docRef *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) error {
	return w.tx.Set(docRef, data, opts...)
}

func (w *transactionWrapper) Update(docRef *firestore.DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) error {
	return w.tx.Update(docRef, updates, preconds...)
}

func (w *transactionWrapper) Delete(docRef *firestore.DocumentRef, preconds ...firestore.Precondition) error {
	return w.tx.Delete(docRef, preconds...)
}

// Documents converts q (Query or CollectionRef wrapper) to the underlying
// firestore.Queryer and delegates to *firestore.Transaction.Documents.
//
// q must be one of the wrappers produced by this package (queryWrapper /
// collectionRefWrapper) or a CollectionRef whose Reference() returns a real
// *firestore.CollectionRef. Custom mocks that satisfy Query but do not back a
// real Firestore queryer cannot be used with the production wrapper; in tests
// you should mock the Transaction interface directly.
func (w *transactionWrapper) Documents(q Query) DocumentIterator {
	queryer, err := toFirestoreQueryer(q)
	if err != nil {
		panic(fmt.Sprintf("go-firestore-mock: transactionWrapper.Documents: %v", err))
	}
	return &documentIteratorWrapper{iter: w.tx.Documents(queryer)}
}

// DocumentRefs delegates to *firestore.Transaction.DocumentRefs using the
// underlying *firestore.CollectionRef behind the CollectionRef wrapper.
func (w *transactionWrapper) DocumentRefs(coll CollectionRef) DocumentRefIterator {
	if coll == nil {
		panic("go-firestore-mock: transactionWrapper.DocumentRefs: nil CollectionRef")
	}
	return &documentRefIteratorWrapper{iter: w.tx.DocumentRefs(coll.Reference())}
}

// toFirestoreQueryer extracts the real firestore.Queryer behind a Query wrapper.
// Returns an error for custom Query implementations that do not back a real
// Firestore type.
func toFirestoreQueryer(q Query) (firestore.Queryer, error) {
	switch v := q.(type) {
	case *queryWrapper:
		return v.q, nil
	case *collectionRefWrapper:
		return v.ref, nil
	}
	if cr, ok := q.(CollectionRef); ok {
		if ref := cr.Reference(); ref != nil {
			return ref, nil
		}
	}
	return nil, fmt.Errorf("Query implementation %T cannot be converted to firestore.Queryer (use a wrapper produced by NewFirestoreClient)", q)
}
