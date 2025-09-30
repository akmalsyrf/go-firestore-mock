package firestore

import (
	"cloud.google.com/go/firestore"
)

//go:generate mockgen -source=transaction.go -destination=transaction_mock.go -package=firestore

// Transaction abstracts Firestore Transaction behavior
type Transaction interface {
	Get(docRef *firestore.DocumentRef) (DocumentSnapshot, error)
	GetAll(docRefs []*firestore.DocumentRef) ([]DocumentSnapshot, error)
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
