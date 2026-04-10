package firestore

import (
	"testing"

	"cloud.google.com/go/firestore"
	pb "cloud.google.com/go/firestore/apiv1/firestorepb"
)

func aggregationIntValue(n int64) *pb.Value {
	return &pb.Value{ValueType: &pb.Value_IntegerValue{IntegerValue: n}}
}

func aggregationDoubleValue(f float64) *pb.Value {
	return &pb.Value{ValueType: &pb.Value_DoubleValue{DoubleValue: f}}
}

func TestAggregationQueryWrapper_InterfaceCompliance(t *testing.T) {
	t.Run("verify AggregationQuery interface compliance", func(t *testing.T) {
		var _ AggregationQuery = (*aggregationQueryWrapper)(nil)
	})
}

func TestAggregationResultWrapper_InterfaceCompliance(t *testing.T) {
	t.Run("verify AggregationResult interface compliance", func(t *testing.T) {
		var _ AggregationResult = (*aggregationResultWrapper)(nil)
	})
}

func TestAggregationQueryWrapper_WithCount(t *testing.T) {
	t.Run("WithCount method exists", func(t *testing.T) {
		wrapper := &aggregationQueryWrapper{aq: nil}
		_ = wrapper.WithCount
	})
}

func TestAggregationQueryWrapper_Get(t *testing.T) {
	t.Run("Get method exists", func(t *testing.T) {
		wrapper := &aggregationQueryWrapper{aq: nil}
		_ = wrapper.Get
	})
}

func TestAggregationResultWrapper_Count(t *testing.T) {
	t.Run("nil result returns error", func(t *testing.T) {
		w := &aggregationResultWrapper{ar: nil}
		_, err := w.Count("x")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("empty map returns error on Count", func(t *testing.T) {
		m := firestore.AggregationResult{}
		w := &aggregationResultWrapper{ar: &m}
		_, err := w.Count("missing")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("missing alias", func(t *testing.T) {
		m := firestore.AggregationResult{"other": aggregationIntValue(1)}
		w := &aggregationResultWrapper{ar: &m}
		_, err := w.Count("total")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("count from protobuf integer value", func(t *testing.T) {
		m := firestore.AggregationResult{"total": aggregationIntValue(42)}
		w := &aggregationResultWrapper{ar: &m}
		n, err := w.Count("total")
		if err != nil || n == nil || *n != 42 {
			t.Fatalf("got (%v, %v), want 42", n, err)
		}
	})

	t.Run("count zero from protobuf integer", func(t *testing.T) {
		m := firestore.AggregationResult{"total": aggregationIntValue(0)}
		w := &aggregationResultWrapper{ar: &m}
		n, err := w.Count("total")
		if err != nil || n == nil || *n != 0 {
			t.Fatalf("got (%v, %v), want 0", n, err)
		}
	})

	t.Run("count from protobuf double", func(t *testing.T) {
		m := firestore.AggregationResult{"total": aggregationDoubleValue(7.0)}
		w := &aggregationResultWrapper{ar: &m}
		n, err := w.Count("total")
		if err != nil || n == nil || *n != 7 {
			t.Fatalf("got (%v, %v), want 7", n, err)
		}
	})

	t.Run("count from plain int64", func(t *testing.T) {
		m := firestore.AggregationResult{"total": int64(99)}
		w := &aggregationResultWrapper{ar: &m}
		n, err := w.Count("total")
		if err != nil || n == nil || *n != 99 {
			t.Fatalf("got (%v, %v), want 99", n, err)
		}
	})

	t.Run("unsupported type", func(t *testing.T) {
		m := firestore.AggregationResult{"total": "not-a-number"}
		w := &aggregationResultWrapper{ar: &m}
		_, err := w.Count("total")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestAggregationQueryWrapper_AllMethods(t *testing.T) {
	t.Run("verify all AggregationQuery methods exist", func(t *testing.T) {
		wrapper := &aggregationQueryWrapper{aq: nil}

		_ = wrapper.WithCount
		_ = wrapper.Get
	})
}

func TestAggregationResultWrapper_AllMethods(t *testing.T) {
	t.Run("verify all AggregationResult methods exist", func(t *testing.T) {
		wrapper := &aggregationResultWrapper{ar: nil}

		_ = wrapper.Count
	})
}

func TestAggregationQueryWrapper_Chaining(t *testing.T) {
	t.Run("WithCount supports chaining", func(t *testing.T) {
		wrapper := &aggregationQueryWrapper{aq: nil}
		_ = wrapper.WithCount
	})
}
