package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

//go:generate mockgen -source=write_batch.go -destination=write_batch_mock.go -package=firestore

// WriteBatch abstracts Firestore WriteBatch behavior
type WriteBatch interface {
	Create(docRef *firestore.DocumentRef, data interface{}) WriteBatch
	Set(docRef *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) WriteBatch
	Update(docRef *firestore.DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) WriteBatch
	Delete(docRef *firestore.DocumentRef, preconds ...firestore.Precondition) WriteBatch
	Commit(ctx context.Context) ([]*firestore.WriteResult, error)
}

type writeBatchWrapper struct {
	wb *firestore.WriteBatch
}

func (w *writeBatchWrapper) Create(docRef *firestore.DocumentRef, data interface{}) WriteBatch {
	w.wb.Create(docRef, data)
	return w
}

func (w *writeBatchWrapper) Set(docRef *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) WriteBatch {
	w.wb.Set(docRef, data, opts...)
	return w
}

func (w *writeBatchWrapper) Update(docRef *firestore.DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) WriteBatch {
	w.wb.Update(docRef, updates, preconds...)
	return w
}

func (w *writeBatchWrapper) Delete(docRef *firestore.DocumentRef, preconds ...firestore.Precondition) WriteBatch {
	w.wb.Delete(docRef, preconds...)
	return w
}

func (w *writeBatchWrapper) Commit(ctx context.Context) ([]*firestore.WriteResult, error) {
	return w.wb.Commit(ctx)
}
