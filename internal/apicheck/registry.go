package apicheck

import (
	"reflect"

	"cloud.google.com/go/firestore"
	"github.com/akmalsyrf/go-firestore-mock/v2"
)

// pair maps an SDK concrete type to the fsmock interface that should wrap it.
type pair struct {
	name  string
	sdk   reflect.Type
	iface reflect.Type
}

// registry is the single source of truth for wrapped SDK types.
func registry() []pair {
	return []pair{
		{name: "Client", sdk: reflect.TypeOf((*firestore.Client)(nil)), iface: reflect.TypeOf((*fsmock.Client)(nil)).Elem()},
		{name: "Query", sdk: reflect.TypeOf(firestore.Query{}), iface: reflect.TypeOf((*fsmock.Query)(nil)).Elem()},
		{name: "CollectionRef", sdk: reflect.TypeOf((*firestore.CollectionRef)(nil)), iface: reflect.TypeOf((*fsmock.CollectionRef)(nil)).Elem()},
		{name: "CollectionGroupRef", sdk: reflect.TypeOf((*firestore.CollectionGroupRef)(nil)), iface: reflect.TypeOf((*fsmock.CollectionGroupRef)(nil)).Elem()},
		{name: "DocumentRef", sdk: reflect.TypeOf((*firestore.DocumentRef)(nil)), iface: reflect.TypeOf((*fsmock.DocumentRef)(nil)).Elem()},
		{name: "DocumentSnapshot", sdk: reflect.TypeOf((*firestore.DocumentSnapshot)(nil)), iface: reflect.TypeOf((*fsmock.DocumentSnapshot)(nil)).Elem()},
		{name: "QuerySnapshot", sdk: reflect.TypeOf((*firestore.QuerySnapshot)(nil)), iface: reflect.TypeOf((*fsmock.QuerySnapshot)(nil)).Elem()},
		{name: "Transaction", sdk: reflect.TypeOf((*firestore.Transaction)(nil)), iface: reflect.TypeOf((*fsmock.Transaction)(nil)).Elem()},
		{name: "WriteBatch", sdk: reflect.TypeOf((*firestore.WriteBatch)(nil)), iface: reflect.TypeOf((*fsmock.WriteBatch)(nil)).Elem()}, //nolint:staticcheck
		{name: "BulkWriter", sdk: reflect.TypeOf((*firestore.BulkWriter)(nil)), iface: reflect.TypeOf((*fsmock.BulkWriter)(nil)).Elem()},
		{name: "AggregationQuery", sdk: reflect.TypeOf((*firestore.AggregationQuery)(nil)), iface: reflect.TypeOf((*fsmock.AggregationQuery)(nil)).Elem()},
		{name: "AggregationResult", sdk: reflect.TypeOf(firestore.AggregationResult(nil)), iface: reflect.TypeOf((*fsmock.AggregationResult)(nil)).Elem()},
		{name: "VectorQuery", sdk: reflect.TypeOf(firestore.VectorQuery{}), iface: reflect.TypeOf((*fsmock.VectorQuery)(nil)).Elem()},
		{name: "DocumentIterator", sdk: reflect.TypeOf((*firestore.DocumentIterator)(nil)), iface: reflect.TypeOf((*fsmock.DocumentIterator)(nil)).Elem()},
		{name: "DocumentRefIterator", sdk: reflect.TypeOf((*firestore.DocumentRefIterator)(nil)), iface: reflect.TypeOf((*fsmock.DocumentRefIterator)(nil)).Elem()},
		{name: "CollectionIterator", sdk: reflect.TypeOf((*firestore.CollectionIterator)(nil)), iface: reflect.TypeOf((*fsmock.CollectionIterator)(nil)).Elem()},
		{name: "QuerySnapshotIterator", sdk: reflect.TypeOf((*firestore.QuerySnapshotIterator)(nil)), iface: reflect.TypeOf((*fsmock.QuerySnapshotIterator)(nil)).Elem()},
		{name: "DocumentSnapshotIterator", sdk: reflect.TypeOf((*firestore.DocumentSnapshotIterator)(nil)), iface: reflect.TypeOf((*fsmock.DocumentSnapshotIterator)(nil)).Elem()},
		{name: "PipelineSource", sdk: reflect.TypeOf((*firestore.PipelineSource)(nil)), iface: reflect.TypeOf((*fsmock.PipelineSource)(nil)).Elem()},
		{name: "Pipeline", sdk: reflect.TypeOf((*firestore.Pipeline)(nil)), iface: reflect.TypeOf((*fsmock.Pipeline)(nil)).Elem()},
		{name: "PipelineSnapshot", sdk: reflect.TypeOf((*firestore.PipelineSnapshot)(nil)), iface: reflect.TypeOf((*fsmock.PipelineSnapshot)(nil)).Elem()},
		{name: "PipelineResult", sdk: reflect.TypeOf((*firestore.PipelineResult)(nil)), iface: reflect.TypeOf((*fsmock.PipelineResult)(nil)).Elem()},
		{name: "PipelineResultIterator", sdk: reflect.TypeOf((*firestore.PipelineResultIterator)(nil)), iface: reflect.TypeOf((*fsmock.PipelineResultIterator)(nil)).Elem()},
	}
}

// substitutions rewrites SDK types to fsmock interface names for signature comparison.
var substitutions = map[string]string{
	"*firestore.Client":                   "fsmock.Client",
	"firestore.Query":                     "fsmock.Query",
	"*firestore.Query":                    "fsmock.Query",
	"*firestore.CollectionRef":            "fsmock.CollectionRef",
	"*firestore.CollectionGroupRef":       "fsmock.CollectionGroupRef",
	"*firestore.DocumentRef":              "fsmock.DocumentRef",
	"*firestore.DocumentSnapshot":         "fsmock.DocumentSnapshot",
	"*firestore.QuerySnapshot":            "fsmock.QuerySnapshot",
	"*firestore.Transaction":              "fsmock.Transaction",
	"*firestore.WriteBatch":               "fsmock.WriteBatch",
	"*firestore.BulkWriter":               "fsmock.BulkWriter",
	"*firestore.AggregationQuery":         "fsmock.AggregationQuery",
	"firestore.AggregationResult":         "fsmock.AggregationResult",
	"firestore.VectorQuery":               "fsmock.VectorQuery",
	"*firestore.DocumentIterator":         "fsmock.DocumentIterator",
	"*firestore.DocumentRefIterator":      "fsmock.DocumentRefIterator",
	"*firestore.CollectionIterator":       "fsmock.CollectionIterator",
	"*firestore.QuerySnapshotIterator":    "fsmock.QuerySnapshotIterator",
	"*firestore.DocumentSnapshotIterator": "fsmock.DocumentSnapshotIterator",
	"*firestore.PipelineSource":           "fsmock.PipelineSource",
	"*firestore.Pipeline":                 "fsmock.Pipeline",
	"*firestore.PipelineSnapshot":         "fsmock.PipelineSnapshot",
	"*firestore.PipelineResult":           "fsmock.PipelineResult",
	"*firestore.PipelineResultIterator":   "fsmock.PipelineResultIterator",
	"[]*firestore.DocumentRef":            "[]fsmock.DocumentRef",
	"[]*firestore.DocumentSnapshot":       "[]fsmock.DocumentSnapshot",
	"[]firestore.Query":                   "[]fsmock.Query",
}

// typeString returns a stable string for a reflect.Type, applying substitutions.
func typeString(t reflect.Type) string {
	if t == nil {
		return "nil"
	}
	s := stringify(t)
	if sub, ok := substitutions[s]; ok {
		return sub
	}
	// Slice / pointer of substituted elem
	if t.Kind() == reflect.Slice {
		elem := typeString(t.Elem())
		return "[]" + elem
	}
	return s
}

func stringify(t reflect.Type) string {
	if t.Name() != "" {
		pkg := t.PkgPath()
		if pkg == "cloud.google.com/go/firestore" {
			return "firestore." + t.Name()
		}
		if pkg == "github.com/akmalsyrf/go-firestore-mock/v2" {
			return "fsmock." + t.Name()
		}
		if pkg == "" {
			return t.Name()
		}
		return t.String()
	}
	switch t.Kind() {
	case reflect.Pointer:
		return "*" + stringify(t.Elem())
	case reflect.Slice:
		return "[]" + stringify(t.Elem())
	case reflect.Array:
		return "[n]" + stringify(t.Elem())
	case reflect.Map:
		return "map[" + stringify(t.Key()) + "]" + stringify(t.Elem())
	case reflect.Func:
		return "func"
	default:
		return t.String()
	}
}

// methodSig returns a normalized signature string for method i on t.
func methodSig(t reflect.Type, i int) string {
	m := t.Method(i)
	mt := m.Type
	in := make([]string, 0, mt.NumIn())
	start := 0
	if t.Kind() != reflect.Interface {
		start = 1 // skip receiver
	}
	for j := start; j < mt.NumIn(); j++ {
		if mt.IsVariadic() && j == mt.NumIn()-1 {
			in = append(in, "..."+typeString(mt.In(j).Elem()))
			continue
		}
		in = append(in, typeString(mt.In(j)))
	}
	out := make([]string, mt.NumOut())
	for j := 0; j < mt.NumOut(); j++ {
		out[j] = typeString(mt.Out(j))
	}
	return m.Name + "(" + join(in) + ") (" + join(out) + ")"
}

func join(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	out := ss[0]
	for i := 1; i < len(ss); i++ {
		out += ", " + ss[i]
	}
	return out
}

func methodNames(t reflect.Type) map[string]bool {
	out := make(map[string]bool)
	if t == nil {
		return out
	}
	if t.Kind() != reflect.Interface && t.Kind() != reflect.Pointer {
		t = reflect.PointerTo(t)
	}
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if m.PkgPath != "" {
			continue
		}
		out[m.Name] = true
	}
	return out
}

func methodIndex(t reflect.Type) map[string]int {
	out := map[string]int{}
	if t.Kind() != reflect.Interface && t.Kind() != reflect.Pointer {
		t = reflect.PointerTo(t)
	}
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if m.PkgPath != "" {
			continue
		}
		out[m.Name] = i
	}
	return out
}

func normalizeType(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Interface && t.Kind() != reflect.Pointer {
		return reflect.PointerTo(t)
	}
	return t
}
