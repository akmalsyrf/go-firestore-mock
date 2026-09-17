package apicheck

import (
	"reflect"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/akmalsyrf/go-firestore-mock/v2"
)

// Pair maps an SDK concrete type to the fsmock interface that should wrap it.
type pair struct {
	name  string
	sdk   reflect.Type
	iface reflect.Type
}

func TestMethodParity(t *testing.T) {
	pairs := []pair{
		{
			name:  "Client",
			sdk:   reflect.TypeOf((*firestore.Client)(nil)),
			iface: reflect.TypeOf((*fsmock.Client)(nil)).Elem(),
		},
		{
			name:  "Query",
			sdk:   reflect.TypeOf(firestore.Query{}),
			iface: reflect.TypeOf((*fsmock.Query)(nil)).Elem(),
		},
		{
			name:  "CollectionRef",
			sdk:   reflect.TypeOf((*firestore.CollectionRef)(nil)),
			iface: reflect.TypeOf((*fsmock.CollectionRef)(nil)).Elem(),
		},
		{
			name:  "CollectionGroupRef",
			sdk:   reflect.TypeOf((*firestore.CollectionGroupRef)(nil)),
			iface: reflect.TypeOf((*fsmock.CollectionGroupRef)(nil)).Elem(),
		},
		{
			name:  "DocumentRef",
			sdk:   reflect.TypeOf((*firestore.DocumentRef)(nil)),
			iface: reflect.TypeOf((*fsmock.DocumentRef)(nil)).Elem(),
		},
		{
			name:  "DocumentSnapshot",
			sdk:   reflect.TypeOf((*firestore.DocumentSnapshot)(nil)),
			iface: reflect.TypeOf((*fsmock.DocumentSnapshot)(nil)).Elem(),
		},
		{
			name:  "Transaction",
			sdk:   reflect.TypeOf((*firestore.Transaction)(nil)),
			iface: reflect.TypeOf((*fsmock.Transaction)(nil)).Elem(),
		},
		{
			name:  "WriteBatch",
			sdk:   reflect.TypeOf((*firestore.WriteBatch)(nil)), //nolint:staticcheck
			iface: reflect.TypeOf((*fsmock.WriteBatch)(nil)).Elem(),
		},
		{
			name:  "BulkWriter",
			sdk:   reflect.TypeOf((*firestore.BulkWriter)(nil)),
			iface: reflect.TypeOf((*fsmock.BulkWriter)(nil)).Elem(),
		},
		{
			name:  "AggregationQuery",
			sdk:   reflect.TypeOf((*firestore.AggregationQuery)(nil)),
			iface: reflect.TypeOf((*fsmock.AggregationQuery)(nil)).Elem(),
		},
		{
			name:  "AggregationResult",
			sdk:   reflect.TypeOf(firestore.AggregationResult(nil)),
			iface: reflect.TypeOf((*fsmock.AggregationResult)(nil)).Elem(),
		},
		{
			name:  "VectorQuery",
			sdk:   reflect.TypeOf(firestore.VectorQuery{}),
			iface: reflect.TypeOf((*fsmock.VectorQuery)(nil)).Elem(),
		},
		{
			name:  "DocumentIterator",
			sdk:   reflect.TypeOf((*firestore.DocumentIterator)(nil)),
			iface: reflect.TypeOf((*fsmock.DocumentIterator)(nil)).Elem(),
		},
		{
			name:  "DocumentRefIterator",
			sdk:   reflect.TypeOf((*firestore.DocumentRefIterator)(nil)),
			iface: reflect.TypeOf((*fsmock.DocumentRefIterator)(nil)).Elem(),
		},
		{
			name:  "CollectionIterator",
			sdk:   reflect.TypeOf((*firestore.CollectionIterator)(nil)),
			iface: reflect.TypeOf((*fsmock.CollectionIterator)(nil)).Elem(),
		},
		{
			name:  "QuerySnapshotIterator",
			sdk:   reflect.TypeOf((*firestore.QuerySnapshotIterator)(nil)),
			iface: reflect.TypeOf((*fsmock.QuerySnapshotIterator)(nil)).Elem(),
		},
		{
			name:  "DocumentSnapshotIterator",
			sdk:   reflect.TypeOf((*firestore.DocumentSnapshotIterator)(nil)),
			iface: reflect.TypeOf((*fsmock.DocumentSnapshotIterator)(nil)).Elem(),
		},
		{
			name:  "PipelineSource",
			sdk:   reflect.TypeOf((*firestore.PipelineSource)(nil)),
			iface: reflect.TypeOf((*fsmock.PipelineSource)(nil)).Elem(),
		},
		{
			name:  "Pipeline",
			sdk:   reflect.TypeOf((*firestore.Pipeline)(nil)),
			iface: reflect.TypeOf((*fsmock.Pipeline)(nil)).Elem(),
		},
		{
			name:  "PipelineSnapshot",
			sdk:   reflect.TypeOf((*firestore.PipelineSnapshot)(nil)),
			iface: reflect.TypeOf((*fsmock.PipelineSnapshot)(nil)).Elem(),
		},
		{
			name:  "PipelineResult",
			sdk:   reflect.TypeOf((*firestore.PipelineResult)(nil)),
			iface: reflect.TypeOf((*fsmock.PipelineResult)(nil)).Elem(),
		},
		{
			name:  "PipelineResultIterator",
			sdk:   reflect.TypeOf((*firestore.PipelineResultIterator)(nil)),
			iface: reflect.TypeOf((*fsmock.PipelineResultIterator)(nil)).Elem(),
		},
	}

	for _, p := range pairs {
		t.Run(p.name, func(t *testing.T) {
			ifaceMethods := methodNames(p.iface)
			sdkMethods := methodNames(p.sdk)

			// DocumentSnapshot: SDK fields mapped to methods — skip field-only check
			if p.name == "DocumentSnapshot" {
				for _, want := range []string{"Data", "DataTo", "DataAt", "DataAtPath", "Exists"} {
					if !ifaceMethods[want] {
						t.Errorf("fsmock.%s missing method %s", p.name, want)
					}
					if !sdkMethods[want] {
						t.Errorf("SDK %s missing method %s (unexpected)", p.name, want)
					}
				}
				return
			}

			// AggregationResult: require Data/DataTo; Count is extra on fsmock
			if p.name == "AggregationResult" {
				for _, want := range []string{"Data", "DataTo"} {
					if !sdkMethods[want] {
						t.Errorf("SDK AggregationResult missing %s", want)
					}
					if !ifaceMethods[want] {
						t.Errorf("fsmock.AggregationResult missing %s", want)
					}
				}
				return
			}

			for name := range sdkMethods {
				key := p.name + "." + name
				if reason, ok := exceptions[key]; ok {
					t.Logf("skip %s: %s", key, reason)
					continue
				}
				if !ifaceMethods[name] {
					t.Errorf("fsmock.%s missing method %q present on SDK type %s", p.name, name, p.sdk)
				}
			}
		})
	}
}

func methodNames(t reflect.Type) map[string]bool {
	out := make(map[string]bool)
	if t == nil {
		return out
	}
	// For non-interface types, Method() includes pointer and value receiver methods
	// depending on whether t is pointer. Normalize to pointer for structs.
	if t.Kind() != reflect.Interface && t.Kind() != reflect.Pointer {
		t = reflect.PointerTo(t)
	}
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if m.PkgPath != "" {
			continue // unexported
		}
		out[m.Name] = true
	}
	return out
}
