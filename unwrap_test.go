package fsmock

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
)

type foreignQuery struct{}

func (foreignQuery) Where(string, string, any) Query { return nil }
func (foreignQuery) WherePath(firestore.FieldPath, string, any) Query {
	return nil
}
func (foreignQuery) WhereEntity(firestore.EntityFilter) Query { return nil }
func (foreignQuery) OrderBy(string, firestore.Direction) Query {
	return nil
}
func (foreignQuery) OrderByPath(firestore.FieldPath, firestore.Direction) Query {
	return nil
}
func (foreignQuery) Limit(int) Query                                 { return nil }
func (foreignQuery) LimitToLast(int) Query                           { return nil }
func (foreignQuery) Offset(int) Query                                { return nil }
func (foreignQuery) StartAt(...any) Query                            { return nil }
func (foreignQuery) StartAfter(...any) Query                         { return nil }
func (foreignQuery) EndAt(...any) Query                              { return nil }
func (foreignQuery) EndBefore(...any) Query                          { return nil }
func (foreignQuery) Select(...string) Query                          { return nil }
func (foreignQuery) SelectPaths(...firestore.FieldPath) Query        { return nil }
func (foreignQuery) Documents(context.Context) DocumentIterator      { return nil }
func (foreignQuery) Snapshots(context.Context) QuerySnapshotIterator { return nil }
func (foreignQuery) NewAggregationQuery() AggregationQuery           { return nil }
func (foreignQuery) FindNearest(string, any, int, firestore.DistanceMeasure, *firestore.FindNearestOptions) VectorQuery {
	return nil
}
func (foreignQuery) FindNearestPath(firestore.FieldPath, any, int, firestore.DistanceMeasure, *firestore.FindNearestOptions) VectorQuery {
	return nil
}
func (foreignQuery) Serialize() ([]byte, error)        { return nil, nil }
func (foreignQuery) Deserialize([]byte) (Query, error) { return nil, nil }
func (foreignQuery) WithReadOptions(...firestore.ReadOption) Query {
	return nil
}
func (foreignQuery) WithRunOptions(...firestore.RunOption) Query { return nil }
func (foreignQuery) Pipeline() Pipeline                          { return nil }

func TestToFirestoreQueryer_RejectsForeign(t *testing.T) {
	_, err := toFirestoreQueryer(foreignQuery{})
	if err == nil {
		t.Fatal("expected error for foreign Query")
	}
}

func TestToFirestoreQueryer_Nil(t *testing.T) {
	_, err := toFirestoreQueryer(nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestToFirestoreQueryer_QueryWrapper(t *testing.T) {
	qw := &queryWrapper{q: firestore.Query{}}
	got, err := toFirestoreQueryer(qw)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := got.(firestore.Query); !ok {
		t.Fatalf("got %T", got)
	}
}

func TestToDocumentRef_Nil(t *testing.T) {
	_, err := toDocumentRef(nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWriteBatch_DeferredUnwrapError(t *testing.T) {
	wb := &writeBatchWrapper{wb: &firestore.WriteBatch{}} //nolint:staticcheck
	wb.Set(nil, map[string]any{"a": 1})
	_, err := wb.Commit(context.Background())
	if err == nil {
		t.Fatal("expected deferred unwrap error on Commit")
	}
}

func TestTransaction_DocumentsForeignQuery(t *testing.T) {
	tx := &transactionWrapper{tx: &firestore.Transaction{}}
	_, err := tx.Documents(foreignQuery{})
	if err == nil {
		t.Fatal("expected error, not panic")
	}
}

func TestTransaction_DocumentRefsNil(t *testing.T) {
	tx := &transactionWrapper{tx: &firestore.Transaction{}}
	_, err := tx.DocumentRefs(nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
