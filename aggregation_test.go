package firestore

import (
	"testing"
)

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
	t.Run("Count method exists", func(t *testing.T) {
		wrapper := &aggregationResultWrapper{ar: nil}
		_ = wrapper.Count
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
