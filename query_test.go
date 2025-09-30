package firestore

import (
	"testing"

	"cloud.google.com/go/firestore"
)

func TestQueryWrapper_OrderBy(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		direction firestore.Direction
	}{
		{
			name:      "order by ascending",
			path:      "createdAt",
			direction: firestore.Asc,
		},
		{
			name:      "order by descending",
			path:      "updatedAt",
			direction: firestore.Desc,
		},
		{
			name:      "order by nested field",
			path:      "user.name",
			direction: firestore.Asc,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := &queryWrapper{
				q: firestore.Query{},
			}

			result := wrapper.OrderBy(tt.path, tt.direction)

			if result == nil {
				t.Error("OrderBy should not return nil")
			}

			if _, ok := result.(*queryWrapper); !ok {
				t.Error("OrderBy should return *queryWrapper")
			}
		})
	}
}

func TestQueryWrapper_Limit(t *testing.T) {
	tests := []struct {
		name  string
		limit int
	}{
		{
			name:  "limit 10",
			limit: 10,
		},
		{
			name:  "limit 100",
			limit: 100,
		},
		{
			name:  "limit 1",
			limit: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := &queryWrapper{
				q: firestore.Query{},
			}

			result := wrapper.Limit(tt.limit)

			if result == nil {
				t.Error("Limit should not return nil")
			}

			if _, ok := result.(*queryWrapper); !ok {
				t.Error("Limit should return *queryWrapper")
			}
		})
	}
}

func TestQueryWrapper_LimitToLast(t *testing.T) {
	tests := []struct {
		name  string
		limit int
	}{
		{
			name:  "limit to last 10",
			limit: 10,
		},
		{
			name:  "limit to last 5",
			limit: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := &queryWrapper{
				q: firestore.Query{},
			}

			result := wrapper.LimitToLast(tt.limit)

			if result == nil {
				t.Error("LimitToLast should not return nil")
			}

			if _, ok := result.(*queryWrapper); !ok {
				t.Error("LimitToLast should return *queryWrapper")
			}
		})
	}
}

func TestQueryWrapper_Offset(t *testing.T) {
	tests := []struct {
		name   string
		offset int
	}{
		{
			name:   "offset 0",
			offset: 0,
		},
		{
			name:   "offset 10",
			offset: 10,
		},
		{
			name:   "offset 100",
			offset: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := &queryWrapper{
				q: firestore.Query{},
			}

			result := wrapper.Offset(tt.offset)

			if result == nil {
				t.Error("Offset should not return nil")
			}

			if _, ok := result.(*queryWrapper); !ok {
				t.Error("Offset should return *queryWrapper")
			}
		})
	}
}

func TestQueryWrapper_StartAt(t *testing.T) {
	tests := []struct {
		name   string
		values []any
	}{
		{
			name:   "start at single value",
			values: []any{"value1"},
		},
		{
			name:   "start at multiple values",
			values: []any{"value1", 123, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := &queryWrapper{
				q: firestore.Query{},
			}

			result := wrapper.StartAt(tt.values...)

			if result == nil {
				t.Error("StartAt should not return nil")
			}

			if _, ok := result.(*queryWrapper); !ok {
				t.Error("StartAt should return *queryWrapper")
			}
		})
	}
}

func TestQueryWrapper_StartAfter(t *testing.T) {
	tests := []struct {
		name   string
		values []any
	}{
		{
			name:   "start after single value",
			values: []any{"value1"},
		},
		{
			name:   "start after multiple values",
			values: []any{"value1", 456},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := &queryWrapper{
				q: firestore.Query{},
			}

			result := wrapper.StartAfter(tt.values...)

			if result == nil {
				t.Error("StartAfter should not return nil")
			}

			if _, ok := result.(*queryWrapper); !ok {
				t.Error("StartAfter should return *queryWrapper")
			}
		})
	}
}

func TestQueryWrapper_EndAt(t *testing.T) {
	tests := []struct {
		name   string
		values []any
	}{
		{
			name:   "end at single value",
			values: []any{"value1"},
		},
		{
			name:   "end at multiple values",
			values: []any{"value1", 789},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := &queryWrapper{
				q: firestore.Query{},
			}

			result := wrapper.EndAt(tt.values...)

			if result == nil {
				t.Error("EndAt should not return nil")
			}

			if _, ok := result.(*queryWrapper); !ok {
				t.Error("EndAt should return *queryWrapper")
			}
		})
	}
}

func TestQueryWrapper_EndBefore(t *testing.T) {
	tests := []struct {
		name   string
		values []any
	}{
		{
			name:   "end before single value",
			values: []any{"value1"},
		},
		{
			name:   "end before multiple values",
			values: []any{"value1", 999},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := &queryWrapper{
				q: firestore.Query{},
			}

			result := wrapper.EndBefore(tt.values...)

			if result == nil {
				t.Error("EndBefore should not return nil")
			}

			if _, ok := result.(*queryWrapper); !ok {
				t.Error("EndBefore should return *queryWrapper")
			}
		})
	}
}

func TestQueryWrapper_Select(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
	}{
		{
			name:  "select single field",
			paths: []string{"name"},
		},
		{
			name:  "select multiple fields",
			paths: []string{"name", "email", "age"},
		},
		{
			name:  "select nested fields",
			paths: []string{"user.name", "user.email"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := &queryWrapper{
				q: firestore.Query{},
			}

			result := wrapper.Select(tt.paths...)

			if result == nil {
				t.Error("Select should not return nil")
			}

			if _, ok := result.(*queryWrapper); !ok {
				t.Error("Select should return *queryWrapper")
			}
		})
	}
}

func TestQueryWrapper_Snapshots(t *testing.T) {
	t.Run("snapshots method exists", func(t *testing.T) {
		wrapper := &queryWrapper{
			q: firestore.Query{},
		}

		// Test that Snapshots method exists
		// Note: Not calling actual method to avoid nil pointer panic
		_ = wrapper.Snapshots
	})
}

func TestQueryWrapper_ComplexChaining(t *testing.T) {
	t.Run("complex query chaining", func(t *testing.T) {
		wrapper := &queryWrapper{
			q: firestore.Query{},
		}

		// Test complex chaining
		result := wrapper.
			Where("status", "==", "active").
			OrderBy("createdAt", firestore.Desc).
			Limit(10).
			Offset(5)

		if result == nil {
			t.Error("Complex chaining should not return nil")
		}

		if _, ok := result.(*queryWrapper); !ok {
			t.Error("Complex chaining should return *queryWrapper")
		}
	})
}

func TestQueryWrapper_PaginationPattern(t *testing.T) {
	t.Run("pagination with limit and offset", func(t *testing.T) {
		wrapper := &queryWrapper{
			q: firestore.Query{},
		}

		page := wrapper.
			OrderBy("createdAt", firestore.Desc).
			Limit(20).
			Offset(40)

		if page == nil {
			t.Error("Pagination query should not return nil")
		}
	})

	t.Run("cursor-based pagination with StartAfter", func(t *testing.T) {
		wrapper := &queryWrapper{
			q: firestore.Query{},
		}

		page := wrapper.
			OrderBy("createdAt", firestore.Desc).
			StartAfter("lastDocID").
			Limit(20)

		if page == nil {
			t.Error("Cursor pagination query should not return nil")
		}
	})
}

func TestQueryWrapper_SortingAndFiltering(t *testing.T) {
	t.Run("multiple order by", func(t *testing.T) {
		wrapper := &queryWrapper{
			q: firestore.Query{},
		}

		result := wrapper.
			OrderBy("category", firestore.Asc).
			OrderBy("price", firestore.Desc)

		if result == nil {
			t.Error("Multiple OrderBy should not return nil")
		}
	})

	t.Run("filtering and sorting combined", func(t *testing.T) {
		wrapper := &queryWrapper{
			q: firestore.Query{},
		}

		result := wrapper.
			Where("category", "==", "electronics").
			Where("price", ">", 100).
			OrderBy("price", firestore.Asc).
			Limit(50)

		if result == nil {
			t.Error("Combined filter and sort should not return nil")
		}
	})
}

func TestQueryWrapper_NewAggregationQuery(t *testing.T) {
	t.Run("NewAggregationQuery method exists", func(t *testing.T) {
		wrapper := &queryWrapper{
			q: firestore.Query{},
		}

		_ = wrapper.NewAggregationQuery
	})
}

func TestQueryWrapper_InterfaceCompliance(t *testing.T) {
	t.Run("verify Query interface compliance", func(t *testing.T) {
		var _ Query = (*queryWrapper)(nil)
	})

	t.Run("verify all query methods exist", func(t *testing.T) {
		wrapper := &queryWrapper{
			q: firestore.Query{},
		}

		// Test that all methods exist
		_ = wrapper.Where
		_ = wrapper.OrderBy
		_ = wrapper.Limit
		_ = wrapper.LimitToLast
		_ = wrapper.Offset
		_ = wrapper.StartAt
		_ = wrapper.StartAfter
		_ = wrapper.EndAt
		_ = wrapper.EndBefore
		_ = wrapper.Select
		_ = wrapper.Documents
		_ = wrapper.Snapshots
		_ = wrapper.NewAggregationQuery
	})
}
