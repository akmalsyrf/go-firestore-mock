package fsmock

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	pb "cloud.google.com/go/firestore/apiv1/firestorepb"
)

// AggregationQuery abstracts *firestore.AggregationQuery.
type AggregationQuery interface {
	WithCount(alias string) AggregationQuery
	WithSum(path string, alias string) AggregationQuery
	WithSumPath(fp firestore.FieldPath, alias string) AggregationQuery
	WithAvg(path string, alias string) AggregationQuery
	WithAvgPath(fp firestore.FieldPath, alias string) AggregationQuery
	Get(ctx context.Context) (AggregationResult, error)
	GetResponse(ctx context.Context) (*AggregationResponse, error)
	Transaction(tx Transaction) (AggregationQuery, error)
	Pipeline() Pipeline
}

// AggregationResult abstracts firestore.AggregationResult.
// Count is a convenience helper; Data/DataTo match the SDK surface added in v1.23+.
type AggregationResult interface {
	Count(alias string) (*int64, error)
	Data() (map[string]any, error)
	DataTo(p any) error
}

// AggregationResponse wraps firestore.AggregationResponse.
type AggregationResponse struct {
	Result         AggregationResult
	ExplainMetrics *firestore.ExplainMetrics
}

type aggregationQueryWrapper struct {
	aq *firestore.AggregationQuery
}

func (w *aggregationQueryWrapper) WithCount(alias string) AggregationQuery {
	return &aggregationQueryWrapper{aq: w.aq.WithCount(alias)}
}

func (w *aggregationQueryWrapper) WithSum(path string, alias string) AggregationQuery {
	return &aggregationQueryWrapper{aq: w.aq.WithSum(path, alias)}
}

func (w *aggregationQueryWrapper) WithSumPath(fp firestore.FieldPath, alias string) AggregationQuery {
	return &aggregationQueryWrapper{aq: w.aq.WithSumPath(fp, alias)}
}

func (w *aggregationQueryWrapper) WithAvg(path string, alias string) AggregationQuery {
	return &aggregationQueryWrapper{aq: w.aq.WithAvg(path, alias)}
}

func (w *aggregationQueryWrapper) WithAvgPath(fp firestore.FieldPath, alias string) AggregationQuery {
	return &aggregationQueryWrapper{aq: w.aq.WithAvgPath(fp, alias)}
}

func (w *aggregationQueryWrapper) Get(ctx context.Context) (AggregationResult, error) {
	result, err := w.aq.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &aggregationResultWrapper{ar: result}, nil
}

func (w *aggregationQueryWrapper) GetResponse(ctx context.Context) (*AggregationResponse, error) {
	resp, err := w.aq.GetResponse(ctx)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	return &AggregationResponse{
		Result:         &aggregationResultWrapper{ar: resp.Result},
		ExplainMetrics: resp.ExplainMetrics,
	}, nil
}

func (w *aggregationQueryWrapper) Transaction(tx Transaction) (AggregationQuery, error) {
	sdkTx, err := toTransaction(tx)
	if err != nil {
		return nil, err
	}
	return &aggregationQueryWrapper{aq: w.aq.Transaction(sdkTx)}, nil
}

func (w *aggregationQueryWrapper) Pipeline() Pipeline {
	return newPipeline(w.aq.Pipeline())
}

type aggregationResultWrapper struct {
	ar firestore.AggregationResult
}

func (w *aggregationResultWrapper) Count(alias string) (*int64, error) {
	if w.ar == nil {
		return nil, errors.New("fsmock: empty aggregation result")
	}
	raw, ok := w.ar[alias]
	if !ok {
		return nil, fmt.Errorf("fsmock: aggregation alias %q not in result", alias)
	}
	n, err := aggregationFieldToInt64(raw)
	if err != nil {
		return nil, fmt.Errorf("fsmock: decode count for alias %q: %w", alias, err)
	}
	return &n, nil
}

func (w *aggregationResultWrapper) Data() (map[string]any, error) {
	return safeAggregationData(w.ar)
}

func (w *aggregationResultWrapper) DataTo(p any) error {
	return w.ar.DataTo(p)
}

// safeAggregationData recovers from SDK AggregationResult.Data panics
// (SDK panics on decode bugs / unexpected value shapes).
func safeAggregationData(ar firestore.AggregationResult) (m map[string]any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("fsmock: AggregationResult.Data panic: %v", r)
		}
	}()
	return ar.Data(), nil
}

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
