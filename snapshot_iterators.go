package firestore

import (
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

//go:generate mockgen -source=snapshot_iterators.go -destination=snapshot_iterators_mock.go -package=firestore

// QuerySnapshotIterator abstracts Firestore QuerySnapshotIterator behavior
type QuerySnapshotIterator interface {
	Next() (*firestore.QuerySnapshot, error)
	Stop()
}

// CollectionIterator abstracts Firestore CollectionIterator behavior
type CollectionIterator interface {
	Next() (*firestore.CollectionRef, error)
	Stop()
}

// DocumentSnapshotIterator abstracts Firestore DocumentSnapshotIterator behavior
type DocumentSnapshotIterator interface {
	Next() (DocumentSnapshot, error)
	Stop()
}

type querySnapshotIteratorWrapper struct {
	iter *firestore.QuerySnapshotIterator
}

func (w *querySnapshotIteratorWrapper) Next() (*firestore.QuerySnapshot, error) {
	return w.iter.Next()
}

func (w *querySnapshotIteratorWrapper) Stop() {
	w.iter.Stop()
}

type collectionIteratorWrapper struct {
	iter *firestore.CollectionIterator
}

func (w *collectionIteratorWrapper) Next() (*firestore.CollectionRef, error) {
	return w.iter.Next()
}

func (w *collectionIteratorWrapper) Stop() {
	// CollectionIterator doesn't have Stop method, it's managed by iterator protocol
}

type documentSnapshotIteratorWrapper struct {
	iter *firestore.DocumentSnapshotIterator
}

func (w *documentSnapshotIteratorWrapper) Next() (DocumentSnapshot, error) {
	snap, err := w.iter.Next()
	if err != nil {
		if err == iterator.Done {
			return nil, err
		}
		return nil, err
	}
	return &documentSnapshotWrapper{snap: snap}, nil
}

func (w *documentSnapshotIteratorWrapper) Stop() {
	w.iter.Stop()
}
