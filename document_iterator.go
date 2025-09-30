package firestore

import (
	"cloud.google.com/go/firestore"
)

//go:generate mockgen -source=document_iterator.go -destination=document_iterator_mock.go

// DocumentIterator abstracts Firestore document iterator behavior
type DocumentIterator interface {
	Next() (*firestore.DocumentSnapshot, error)
	Stop()
	GetAll() ([]*firestore.DocumentSnapshot, error)
}
