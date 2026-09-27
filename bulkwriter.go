package fsmock

import (
	"cloud.google.com/go/firestore"
)

// BulkWriter abstracts *firestore.BulkWriter.
//
// As of firestore v1.23+, BulkWriter enforces backpressure and surfaces write
// errors. Flush and End may block until in-flight writes complete.
type BulkWriter interface {
	Create(docRef DocumentRef, data any) (BulkWriterJob, error)
	Set(docRef DocumentRef, data any, opts ...firestore.SetOption) (BulkWriterJob, error)
	Update(docRef DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) (BulkWriterJob, error)
	Delete(docRef DocumentRef, preconds ...firestore.Precondition) (BulkWriterJob, error)
	Flush()
	End()
}

// BulkWriterJob abstracts *firestore.BulkWriterJob so callers can mock Results().
type BulkWriterJob interface {
	Results() (*firestore.WriteResult, error)
}

type bulkWriterWrapper struct {
	bw *firestore.BulkWriter
}

type bulkWriterJobWrapper struct {
	job *firestore.BulkWriterJob
}

func newBulkWriterJob(j *firestore.BulkWriterJob) BulkWriterJob {
	if j == nil {
		return nil
	}
	return &bulkWriterJobWrapper{job: j}
}

func (w *bulkWriterWrapper) Create(docRef DocumentRef, data any) (BulkWriterJob, error) {
	ref, err := toDocumentRef(docRef)
	if err != nil {
		return nil, err
	}
	job, err := w.bw.Create(ref, data)
	return newBulkWriterJob(job), err
}

func (w *bulkWriterWrapper) Set(docRef DocumentRef, data any, opts ...firestore.SetOption) (BulkWriterJob, error) {
	ref, err := toDocumentRef(docRef)
	if err != nil {
		return nil, err
	}
	job, err := w.bw.Set(ref, data, opts...)
	return newBulkWriterJob(job), err
}

func (w *bulkWriterWrapper) Update(docRef DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) (BulkWriterJob, error) {
	ref, err := toDocumentRef(docRef)
	if err != nil {
		return nil, err
	}
	job, err := w.bw.Update(ref, updates, preconds...)
	return newBulkWriterJob(job), err
}

func (w *bulkWriterWrapper) Delete(docRef DocumentRef, preconds ...firestore.Precondition) (BulkWriterJob, error) {
	ref, err := toDocumentRef(docRef)
	if err != nil {
		return nil, err
	}
	job, err := w.bw.Delete(ref, preconds...)
	return newBulkWriterJob(job), err
}

func (w *bulkWriterWrapper) Flush() {
	w.bw.Flush()
}

func (w *bulkWriterWrapper) End() {
	w.bw.End()
}

func (w *bulkWriterJobWrapper) Results() (*firestore.WriteResult, error) {
	return w.job.Results()
}
