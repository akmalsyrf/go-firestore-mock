package firestore

import (
	"time"

	"cloud.google.com/go/firestore"
)

//go:generate mockgen -source=document_snapshot.go -destination=document_snapshot_mock.go -package=firestore

// DocumentSnapshot abstracts Firestore DocumentSnapshot behavior
type DocumentSnapshot interface {
	Data() map[string]interface{}
	DataTo(p interface{}) error
	DataAt(path string) (interface{}, error)
	Exists() bool
	CreateTime() time.Time
	UpdateTime() time.Time
	ReadTime() time.Time
	Ref() *firestore.DocumentRef
}

type documentSnapshotWrapper struct {
	snap *firestore.DocumentSnapshot
}

func (w *documentSnapshotWrapper) Data() map[string]interface{} {
	return w.snap.Data()
}

func (w *documentSnapshotWrapper) DataTo(p interface{}) error {
	return w.snap.DataTo(p)
}

func (w *documentSnapshotWrapper) DataAt(path string) (interface{}, error) {
	return w.snap.DataAt(path)
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

func (w *documentSnapshotWrapper) Ref() *firestore.DocumentRef {
	return w.snap.Ref
}
