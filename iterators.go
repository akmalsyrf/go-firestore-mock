package fsmock

import (
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// DocumentIterator abstracts *firestore.DocumentIterator.
type DocumentIterator interface {
	Next() (DocumentSnapshot, error)
	Stop()
	GetAll() ([]DocumentSnapshot, error)
	ExplainMetrics() (*firestore.ExplainMetrics, error)
}

// DocumentRefIterator abstracts *firestore.DocumentRefIterator.
type DocumentRefIterator interface {
	Next() (DocumentRef, error)
	GetAll() ([]DocumentRef, error)
	PageInfo() *iterator.PageInfo
}

// CollectionIterator abstracts *firestore.CollectionIterator.
// Unlike DocumentIterator, the SDK CollectionIterator has no Stop method.
type CollectionIterator interface {
	Next() (CollectionRef, error)
	GetAll() ([]CollectionRef, error)
	PageInfo() *iterator.PageInfo
}

// QuerySnapshotIterator abstracts *firestore.QuerySnapshotIterator.
type QuerySnapshotIterator interface {
	Next() (QuerySnapshot, error)
	Stop()
}

// DocumentSnapshotIterator abstracts *firestore.DocumentSnapshotIterator.
type DocumentSnapshotIterator interface {
	Next() (DocumentSnapshot, error)
	Stop()
}

type documentIteratorWrapper struct {
	iter *firestore.DocumentIterator
}

func (w *documentIteratorWrapper) Next() (DocumentSnapshot, error) {
	snap, err := w.iter.Next()
	if err != nil {
		return nil, err
	}
	return &documentSnapshotWrapper{snap: snap}, nil
}

func (w *documentIteratorWrapper) Stop() {
	w.iter.Stop()
}

func (w *documentIteratorWrapper) GetAll() ([]DocumentSnapshot, error) {
	snaps, err := w.iter.GetAll()
	if err != nil {
		return nil, err
	}
	return wrapSnapshots(snaps), nil
}

func (w *documentIteratorWrapper) ExplainMetrics() (*firestore.ExplainMetrics, error) {
	return w.iter.ExplainMetrics()
}

type documentRefIteratorWrapper struct {
	iter *firestore.DocumentRefIterator
}

func (w *documentRefIteratorWrapper) Next() (DocumentRef, error) {
	ref, err := w.iter.Next()
	if err != nil {
		return nil, err
	}
	return newDocumentRef(ref), nil
}

func (w *documentRefIteratorWrapper) GetAll() ([]DocumentRef, error) {
	refs, err := w.iter.GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]DocumentRef, len(refs))
	for i, r := range refs {
		out[i] = newDocumentRef(r)
	}
	return out, nil
}

func (w *documentRefIteratorWrapper) PageInfo() *iterator.PageInfo {
	return w.iter.PageInfo()
}

type collectionIteratorWrapper struct {
	iter *firestore.CollectionIterator
}

func (w *collectionIteratorWrapper) Next() (CollectionRef, error) {
	ref, err := w.iter.Next()
	if err != nil {
		return nil, err
	}
	return newCollectionRef(ref), nil
}

func (w *collectionIteratorWrapper) GetAll() ([]CollectionRef, error) {
	refs, err := w.iter.GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]CollectionRef, len(refs))
	for i, r := range refs {
		out[i] = newCollectionRef(r)
	}
	return out, nil
}

func (w *collectionIteratorWrapper) PageInfo() *iterator.PageInfo {
	return w.iter.PageInfo()
}

type querySnapshotIteratorWrapper struct {
	iter *firestore.QuerySnapshotIterator
}

func (w *querySnapshotIteratorWrapper) Next() (QuerySnapshot, error) {
	snap, err := w.iter.Next()
	if err != nil {
		return nil, err
	}
	return &querySnapshotWrapper{snap: snap}, nil
}

func (w *querySnapshotIteratorWrapper) Stop() {
	w.iter.Stop()
}

type documentSnapshotIteratorWrapper struct {
	iter *firestore.DocumentSnapshotIterator
}

func (w *documentSnapshotIteratorWrapper) Next() (DocumentSnapshot, error) {
	snap, err := w.iter.Next()
	if err != nil {
		return nil, err
	}
	return &documentSnapshotWrapper{snap: snap}, nil
}

func (w *documentSnapshotIteratorWrapper) Stop() {
	w.iter.Stop()
}
