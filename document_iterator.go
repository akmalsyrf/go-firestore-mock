package firestore

import (
	"cloud.google.com/go/firestore"
)

//go:generate mockgen -source=document_iterator.go -destination=document_iterator_mock.go -package=firestore

// DocumentIterator abstracts Firestore document iterator behavior
type DocumentIterator interface {
	Next() (*firestore.DocumentSnapshot, error)
	Stop()
	GetAll() ([]*firestore.DocumentSnapshot, error)
}

// DocumentRefIterator abstracts Firestore DocumentRefIterator behavior
// (returned by *firestore.CollectionRef.DocumentRefs and *firestore.Transaction.DocumentRefs).
type DocumentRefIterator interface {
	Next() (*firestore.DocumentRef, error)
	GetAll() ([]*firestore.DocumentRef, error)
}

type documentRefIteratorWrapper struct {
	iter *firestore.DocumentRefIterator
}

func (w *documentRefIteratorWrapper) Next() (*firestore.DocumentRef, error) {
	return w.iter.Next()
}

func (w *documentRefIteratorWrapper) GetAll() ([]*firestore.DocumentRef, error) {
	return w.iter.GetAll()
}
