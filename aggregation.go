package firestore

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	pb "cloud.google.com/go/firestore/apiv1/firestorepb"
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
	if w.ar == nil || *w.ar == nil {
		return nil, errors.New("go-firestore-mock: empty aggregation result")
	}
	raw, ok := (*w.ar)[alias]
	if !ok {
		return nil, fmt.Errorf("go-firestore-mock: aggregation alias %q not in result", alias)
	}
	n, err := aggregationFieldToInt64(raw)
	if err != nil {
		return nil, fmt.Errorf("go-firestore-mock: decode count for alias %q: %w", alias, err)
	}
	return &n, nil
}

// aggregationFieldToInt64 interprets values as returned by cloud.google.com/go/firestore
// AggregationQuery.Get (map entries are typically *firestorepb.Value).
func aggregationFieldToInt64(v interface{}) (int64, error) {
	switch x := v.(type) {
	case *pb.Value:
		if x == nil {
			return 0, errors.New("nil *firestorepb.Value")
		}
		switch t := x.GetValueType().(type) {
		case *pb.Value_IntegerValue:
			return t.IntegerValue, nil
		case *pb.Value_DoubleValue:
			return int64(t.DoubleValue), nil
		default:
			return 0, fmt.Errorf("unsupported protobuf value type %T", t)
		}
	case int64:
		return x, nil
	case int:
		return int64(x), nil
	case int32:
		return int64(x), nil
	case float64:
		return int64(x), nil
	default:
		return 0, fmt.Errorf("unsupported Go type %T", v)
	}
}
