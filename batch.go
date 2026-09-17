package fsmock

import (
	"context"

	"cloud.google.com/go/firestore"
)

// WriteBatch abstracts *firestore.WriteBatch.
// Unwrap errors are deferred and returned from Commit (chaining pattern).
type WriteBatch interface {
	Create(docRef DocumentRef, data any) WriteBatch
	Set(docRef DocumentRef, data any, opts ...firestore.SetOption) WriteBatch
	Update(docRef DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) WriteBatch
	Delete(docRef DocumentRef, preconds ...firestore.Precondition) WriteBatch
	Commit(ctx context.Context) ([]*firestore.WriteResult, error)
}

type writeBatchWrapper struct {
	wb  *firestore.WriteBatch //nolint:staticcheck // SDK still ships WriteBatch; wrap for mocks.
	err error
}

func (w *writeBatchWrapper) note(err error) {
	if w.err == nil && err != nil {
		w.err = err
	}
}

func (w *writeBatchWrapper) Create(docRef DocumentRef, data any) WriteBatch {
	ref, err := toDocumentRef(docRef)
	w.note(err)
	if err == nil {
		w.wb.Create(ref, data)
	}
	return w
}

func (w *writeBatchWrapper) Set(docRef DocumentRef, data any, opts ...firestore.SetOption) WriteBatch {
	ref, err := toDocumentRef(docRef)
	w.note(err)
	if err == nil {
		w.wb.Set(ref, data, opts...)
	}
	return w
}

func (w *writeBatchWrapper) Update(docRef DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) WriteBatch {
	ref, err := toDocumentRef(docRef)
	w.note(err)
	if err == nil {
		w.wb.Update(ref, updates, preconds...)
	}
	return w
}

func (w *writeBatchWrapper) Delete(docRef DocumentRef, preconds ...firestore.Precondition) WriteBatch {
	ref, err := toDocumentRef(docRef)
	w.note(err)
	if err == nil {
		w.wb.Delete(ref, preconds...)
	}
	return w
}

func (w *writeBatchWrapper) Commit(ctx context.Context) ([]*firestore.WriteResult, error) {
	if w.err != nil {
		return nil, w.err
	}
	return w.wb.Commit(ctx)
}
