package firestore

import (
	"cloud.google.com/go/firestore"
)

//go:generate mockgen -source=bulk_writer.go -destination=bulk_writer_mock.go

// BulkWriter abstracts Firestore bulk writer behavior
type BulkWriter interface {
	Create(docRef *firestore.DocumentRef, data interface{}) (*firestore.BulkWriterJob, error)
	Set(docRef *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) (*firestore.BulkWriterJob, error)
	Update(docRef *firestore.DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) (*firestore.BulkWriterJob, error)
	Delete(docRef *firestore.DocumentRef, preconds ...firestore.Precondition) (*firestore.BulkWriterJob, error)
	Flush()
	End()
}
