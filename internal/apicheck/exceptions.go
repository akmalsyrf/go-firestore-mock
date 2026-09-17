package apicheck

// ignoredTypes lists SDK types that have exported methods but are intentionally
// not wrapped by fsmock (value types, errors, jobs, helpers).
// Keys are SDK type names; values are required reasons.
var ignoredTypes = map[string]string{
	"BSONObjectID":       "value type with helper methods; pass-through SDK value",
	"BulkWriterJob":      "returned by BulkWriter ops; opaque SDK handle, not mocked",
	"CommitResponse":     "value type returned by Pipeline write stages",
	"DocumentChangeKind": "enum-like value type",
	"ExplainStats":       "value/metrics type; pass-through",
	"FieldNotFoundError": "error type; pass-through",
	"Update":             "value type used as write argument; pass-through",
}

// deviations documents intentional signature differences from the SDK after
// substitution. Keys are "Type.Method". Reason is required.
//
// Also documents fsmock-only methods (escape hatches, field→method adapters,
// convenience helpers) so contributors know they are intentional.
var deviations = map[string]string{
	// Return-type / embedding constraints
	"CollectionRef.WithReadOptions":      "SDK returns *CollectionRef; fsmock inherits Query.WithReadOptions returning Query (Go embedding cannot overload by return type)",
	"Query.WithReadOptions":              "SDK returns *Query; fsmock returns Query interface (value semantics via wrapper)",
	"CollectionGroupRef.WithReadOptions": "inherited from Query; same as Query.WithReadOptions",

	// Escape hatches
	"DocumentSnapshot.Reference": "fsmock escape hatch to *firestore.DocumentSnapshot for cursor APIs",
	"QuerySnapshot.Reference":    "fsmock escape hatch to *firestore.QuerySnapshot",
	"Query.SDKQuery":             "fsmock escape hatch; named SDKQuery to avoid collision with CollectionRef.Reference",

	// SDK fields exposed as methods
	"DocumentRef.ID":              "SDK exposes ID as a field; fsmock exposes ID() method",
	"DocumentRef.Path":            "SDK exposes Path as a field; fsmock exposes Path() method",
	"DocumentRef.Parent":          "SDK field; fsmock method returning CollectionRef",
	"CollectionRef.ID":            "SDK exposes ID as a field; fsmock exposes ID() method",
	"CollectionRef.Path":          "SDK exposes Path as a field; fsmock exposes Path() method",
	"CollectionRef.Parent":        "SDK field; fsmock method returning DocumentRef",
	"DocumentSnapshot.CreateTime": "SDK field; fsmock method",
	"DocumentSnapshot.UpdateTime": "SDK field; fsmock method",
	"DocumentSnapshot.ReadTime":   "SDK field; fsmock method",
	"DocumentSnapshot.Ref":        "SDK field Ref; fsmock method Ref()",
	"QuerySnapshot.Size":          "SDK field; fsmock method",
	"QuerySnapshot.Changes":       "SDK field; fsmock method",
	"QuerySnapshot.ReadTime":      "SDK field; fsmock method",
	"QuerySnapshot.Documents":     "SDK field (iterator); fsmock method Documents()",

	// Behavioral improvements: error instead of panic / interface args
	"Transaction.Documents":                     "fsmock returns (DocumentIterator, error); SDK takes Queryer and panics on bad input",
	"Transaction.DocumentRefs":                  "fsmock returns (DocumentRefIterator, error) instead of panicking on nil",
	"Transaction.Execute":                       "fsmock returns error when Pipeline unwrap fails",
	"AggregationQuery.Transaction":              "fsmock returns (AggregationQuery, error) when Transaction unwrap fails",
	"AggregationQuery.Get":                      "fsmock returns AggregationResult interface instead of concrete map type",
	"AggregationQuery.GetResponse":              "fsmock wraps result in *fsmock.AggregationResponse",
	"AggregationResult.Count":                   "fsmock convenience helper not on SDK",
	"AggregationResult.Data":                    "fsmock returns (map, error) and recovers SDK panics",
	"Client.GetAll":                             "fsmock accepts []DocumentRef interfaces",
	"Client.RunTransaction":                     "fsmock callback receives fsmock.Transaction",
	"BulkWriter.Create":                         "fsmock accepts DocumentRef interface",
	"BulkWriter.Set":                            "fsmock accepts DocumentRef interface",
	"BulkWriter.Update":                         "fsmock accepts DocumentRef interface",
	"BulkWriter.Delete":                         "fsmock accepts DocumentRef interface",
	"WriteBatch.Create":                         "fsmock accepts DocumentRef interface",
	"WriteBatch.Set":                            "fsmock accepts DocumentRef interface",
	"WriteBatch.Update":                         "fsmock accepts DocumentRef interface",
	"WriteBatch.Delete":                         "fsmock accepts DocumentRef interface",
	"PipelineSource.Documents":                  "fsmock accepts []DocumentRef and returns error on unwrap",
	"PipelineSource.CreateFromQuery":            "fsmock accepts Query interface and returns error on unwrap",
	"PipelineSource.CreateFromAggregationQuery": "fsmock accepts AggregationQuery interface and returns error on unwrap",
	"Pipeline.Union":                            "fsmock accepts Pipeline interface and returns error on unwrap",
}

// exceptions lists SDK methods intentionally absent from fsmock under the same name.
// Prefer deviations for signature differences. Keep empty unless a method is skipped entirely.
var exceptions = map[string]string{}
