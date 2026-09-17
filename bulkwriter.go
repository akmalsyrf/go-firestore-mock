package fsmock

import (
	"cloud.google.com/go/firestore"
)

// BulkWriter abstracts *firestore.BulkWriter.
//
// As of firestore v1.23+, BulkWriter enforces backpressure and surfaces write
// errors. Flush and End may block until in-flight writes complete.
type BulkWriter interface {
	Create(docRef DocumentRef, data any) (*firestore.BulkWriterJob, error)
	Set(docRef DocumentRef, data any, opts ...firestore.SetOption) (*firestore.BulkWriterJob, error)
	Update(docRef DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) (*firestore.BulkWriterJob, error)
	Delete(docRef DocumentRef, preconds ...firestore.Precondition) (*firestore.BulkWriterJob, error)
	Flush()
	End()
}

type bulkWriterWrapper struct {
	bw *firestore.BulkWriter
}

func (w *bulkWriterWrapper) Create(docRef DocumentRef, data any) (*firestore.BulkWriterJob, error) {
	ref, err := toDocumentRef(docRef)
	if err != nil {
		return nil, err
	}
	return w.bw.Create(ref, data)
}

func (w *bulkWriterWrapper) Set(docRef DocumentRef, data any, opts ...firestore.SetOption) (*firestore.BulkWriterJob, error) {
	ref, err := toDocumentRef(docRef)
	if err != nil {
		return nil, err
	}
	return w.bw.Set(ref, data, opts...)
}

func (w *bulkWriterWrapper) Update(docRef DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) (*firestore.BulkWriterJob, error) {
	ref, err := toDocumentRef(docRef)
	if err != nil {
		return nil, err
	}
	return w.bw.Update(ref, updates, preconds...)
}

func (w *bulkWriterWrapper) Delete(docRef DocumentRef, preconds ...firestore.Precondition) (*firestore.BulkWriterJob, error) {
	ref, err := toDocumentRef(docRef)
	if err != nil {
		return nil, err
	}
	return w.bw.Delete(ref, preconds...)
}

func (w *bulkWriterWrapper) Flush() {
	w.bw.Flush()
}

func (w *bulkWriterWrapper) End() {
	w.bw.End()
}
