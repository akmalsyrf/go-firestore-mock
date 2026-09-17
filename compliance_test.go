package fsmock

import (
	"testing"

	"cloud.google.com/go/firestore"
	pb "cloud.google.com/go/firestore/apiv1/firestorepb"
)

// Compile-time interface compliance.
var (
	_ Client                   = (*clientWrapper)(nil)
	_ Query                    = (*queryWrapper)(nil)
	_ Query                    = (*collectionRefWrapper)(nil)
	_ CollectionRef            = (*collectionRefWrapper)(nil)
	_ Query                    = (*collectionGroupRefWrapper)(nil)
	_ CollectionGroupRef       = (*collectionGroupRefWrapper)(nil)
	_ DocumentRef              = (*documentRefWrapper)(nil)
	_ DocumentSnapshot         = (*documentSnapshotWrapper)(nil)
	_ QuerySnapshot            = (*querySnapshotWrapper)(nil)
	_ Transaction              = (*transactionWrapper)(nil)
	_ WriteBatch               = (*writeBatchWrapper)(nil)
	_ BulkWriter               = (*bulkWriterWrapper)(nil)
	_ AggregationQuery         = (*aggregationQueryWrapper)(nil)
	_ AggregationResult        = (*aggregationResultWrapper)(nil)
	_ VectorQuery              = (*vectorQueryWrapper)(nil)
	_ PipelineSource           = (*pipelineSourceWrapper)(nil)
	_ Pipeline                 = (*pipelineWrapper)(nil)
	_ PipelineSnapshot         = (*pipelineSnapshotWrapper)(nil)
	_ PipelineResult           = (*pipelineResultWrapper)(nil)
	_ DocumentIterator         = (*documentIteratorWrapper)(nil)
	_ DocumentRefIterator      = (*documentRefIteratorWrapper)(nil)
	_ CollectionIterator       = (*collectionIteratorWrapper)(nil)
	_ QuerySnapshotIterator    = (*querySnapshotIteratorWrapper)(nil)
	_ DocumentSnapshotIterator = (*documentSnapshotIteratorWrapper)(nil)
	_ PipelineResultIterator   = (*pipelineResultIteratorWrapper)(nil)
)

func TestNewClient_Nil(t *testing.T) {
	_, err := NewClient(nil)
	if err != ErrNilClient {
		t.Fatalf("got %v, want ErrNilClient", err)
	}
}

func TestMustNewClient_NilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	_ = MustNewClient(nil)
}

func TestAggregationFieldToInt64(t *testing.T) {
	n, err := aggregationFieldToInt64(&pb.Value{ValueType: &pb.Value_IntegerValue{IntegerValue: 42}})
	if err != nil || n != 42 {
		t.Fatalf("got %d %v", n, err)
	}
	n, err = aggregationFieldToInt64(int64(7))
	if err != nil || n != 7 {
		t.Fatalf("got %d %v", n, err)
	}
	_, err = aggregationFieldToInt64("bad")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAggregationResult_CountAndData(t *testing.T) {
	ar := firestore.AggregationResult{
		"total": &pb.Value{ValueType: &pb.Value_IntegerValue{IntegerValue: 3}},
	}
	w := &aggregationResultWrapper{ar: ar}
	n, err := w.Count("total")
	if err != nil || n == nil || *n != 3 {
		t.Fatalf("Count: %v %v", n, err)
	}
	_, err = w.Count("missing")
	if err == nil {
		t.Fatal("expected missing alias error")
	}
	// Data may panic on non-pb values; our Count path is the primary helper.
	// DataTo with proper pb values:
	var m map[string]any
	if err := w.DataTo(&m); err != nil {
		t.Fatalf("DataTo: %v", err)
	}
	if m["total"] != int64(3) {
		t.Fatalf("DataTo map: %#v", m)
	}
}

func TestSafeAggregationData_RecoversPanic(t *testing.T) {
	// Non-pb values cause Data() to panic in the SDK.
	ar := firestore.AggregationResult{"x": "not-pb"}
	_, err := safeAggregationData(ar)
	if err == nil {
		t.Fatal("expected error from recovered panic")
	}
}
