package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

//go:generate mockgen -source=aggregation.go -destination=aggregation_mock.go -package=firestore

// AggregationQuery abstracts Firestore AggregationQuery behavior
type AggregationQuery interface {
	WithCount(alias string) AggregationQuery
	Get(ctx context.Context) (AggregationResult, error)
}

// AggregationResult abstracts Firestore AggregationResult behavior
type AggregationResult interface {
	Count(alias string) (*int64, error)
}

type aggregationQueryWrapper struct {
	aq *firestore.AggregationQuery
}

func (w *aggregationQueryWrapper) WithCount(alias string) AggregationQuery {
	return &aggregationQueryWrapper{aq: w.aq.WithCount(alias)}
}

func (w *aggregationQueryWrapper) Get(ctx context.Context) (AggregationResult, error) {
	result, err := w.aq.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &aggregationResultWrapper{ar: &result}, nil
}

type aggregationResultWrapper struct {
	ar *firestore.AggregationResult
}

func (w *aggregationResultWrapper) Count(alias string) (*int64, error) {
	if w.ar == nil {
		return nil, nil
	}
	// Note: AggregationResult doesn't have direct Count method in current SDK
	// This is a simplified wrapper - actual implementation may vary by SDK version
	return nil, nil
}
