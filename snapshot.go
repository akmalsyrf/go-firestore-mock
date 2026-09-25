package fsmock

import (
	"time"

	"cloud.google.com/go/firestore"
)

// DocumentSnapshot abstracts *firestore.DocumentSnapshot.
type DocumentSnapshot interface {
	Data() map[string]any
	DataTo(p any) error
	DataAt(path string) (any, error)
	DataAtPath(fp firestore.FieldPath) (any, error)
	Exists() bool
	CreateTime() time.Time
	UpdateTime() time.Time
	ReadTime() time.Time
	Ref() DocumentRef
	// Reference returns the underlying SDK snapshot (escape hatch for cursor APIs).
	Reference() *firestore.DocumentSnapshot
}

// QuerySnapshot abstracts *firestore.QuerySnapshot.
type QuerySnapshot interface {
	Documents() DocumentIterator
	Size() int
	Changes() []firestore.DocumentChange
	ReadTime() time.Time
	// Reference returns the underlying SDK query snapshot.
	Reference() *firestore.QuerySnapshot
}

type documentSnapshotWrapper struct {
	snap *firestore.DocumentSnapshot
}

func (w *documentSnapshotWrapper) Data() map[string]any {
	return w.snap.Data()
}

func (w *documentSnapshotWrapper) DataTo(p any) error {
	return w.snap.DataTo(p)
}

func (w *documentSnapshotWrapper) DataAt(path string) (any, error) {
	return w.snap.DataAt(path)
}

func (w *documentSnapshotWrapper) DataAtPath(fp firestore.FieldPath) (any, error) {
	return w.snap.DataAtPath(fp)
}

func (w *documentSnapshotWrapper) Exists() bool {
	return w.snap.Exists()
}

func (w *documentSnapshotWrapper) CreateTime() time.Time {
	return w.snap.CreateTime
}

func (w *documentSnapshotWrapper) UpdateTime() time.Time {
	return w.snap.UpdateTime
}

func (w *documentSnapshotWrapper) ReadTime() time.Time {
	return w.snap.ReadTime
}

func (w *documentSnapshotWrapper) Ref() DocumentRef {
	return newDocumentRef(w.snap.Ref)
}

func (w *documentSnapshotWrapper) Reference() *firestore.DocumentSnapshot {
	return w.snap
}

type querySnapshotWrapper struct {
	snap *firestore.QuerySnapshot
}

func (w *querySnapshotWrapper) Documents() DocumentIterator {
	return newDocumentIterator(w.snap.Documents)
}

func (w *querySnapshotWrapper) Size() int {
	return w.snap.Size
}

func (w *querySnapshotWrapper) Changes() []firestore.DocumentChange {
	return w.snap.Changes
}

func (w *querySnapshotWrapper) ReadTime() time.Time {
	return w.snap.ReadTime
}

func (w *querySnapshotWrapper) Reference() *firestore.QuerySnapshot {
	return w.snap
}
