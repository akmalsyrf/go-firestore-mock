package apicheck

// waivers lists *xxxWrapper methods that are allowed to lack coverage in the
// accuracy gate. Keys are "receiverType.Method" (e.g. "queryWrapper.Snapshots").
// Reason and Issue are required so waivers do not rot silently.
type waiver struct {
	Reason string
	Issue  string // GitHub issue URL or short tracking id; use "n/a: <why>" if none
}

var waivers = map[string]waiver{
	// Realtime listeners need a long-lived stream; emulator coverage is flaky in CI.
	"queryWrapper.Snapshots": {
		Reason: "realtime Query.Snapshots listener; covered indirectly via interface compliance",
		Issue:  "n/a: emulator listen flaky in short CI",
	},
	"documentRefWrapper.Snapshots": {
		Reason: "realtime DocumentRef.Snapshots listener",
		Issue:  "n/a: emulator listen flaky in short CI",
	},
	"querySnapshotIteratorWrapper.Next": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "n/a: depends on Snapshots waiver",
	},
	"querySnapshotIteratorWrapper.Stop": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "n/a: depends on Snapshots waiver",
	},
	"documentSnapshotIteratorWrapper.Next": {
		Reason: "only reachable via DocumentRef.Snapshots",
		Issue:  "n/a: depends on Snapshots waiver",
	},
	"documentSnapshotIteratorWrapper.Stop": {
		Reason: "only reachable via DocumentRef.Snapshots",
		Issue:  "n/a: depends on Snapshots waiver",
	},
	"querySnapshotWrapper.Documents": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "n/a: depends on Snapshots waiver",
	},
	"querySnapshotWrapper.Size": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "n/a: depends on Snapshots waiver",
	},
	"querySnapshotWrapper.Changes": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "n/a: depends on Snapshots waiver",
	},
	"querySnapshotWrapper.ReadTime": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "n/a: depends on Snapshots waiver",
	},
	"querySnapshotWrapper.Reference": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "n/a: depends on Snapshots waiver",
	},

	// Client lifecycle / global options rarely exercised in shared emulator suites.
	"clientWrapper.Close": {
		Reason: "client closed by harness cleanup on raw SDK client",
		Issue:  "n/a: harness closes *firestore.Client directly",
	},
	"clientWrapper.WithReadOptions": {
		Reason: "client-level read time; prefer collection/query coverage",
		Issue:  "n/a: low priority",
	},
	"clientWrapper.WithAlwaysUseImplicitOrderBy": {
		Reason: "global query default; behavioral flag",
		Issue:  "n/a: low priority",
	},
	"clientWrapper.DocFromFullPath": {
		Reason: "path helper equivalent to Doc for emulator paths",
		Issue:  "n/a: low priority",
	},

	// Pipeline / vector / enterprise features often unsupported by emulator.
	"queryWrapper.FindNearest": {
		Reason: "vector search; emulator support incomplete",
		Issue:  "n/a: emulator",
	},
	"queryWrapper.FindNearestPath": {
		Reason: "vector search; emulator support incomplete",
		Issue:  "n/a: emulator",
	},
	"queryWrapper.Pipeline": {
		Reason: "pipeline smoke may skip on emulator",
		Issue:  "n/a: emulator",
	},
	"queryWrapper.WithRunOptions": {
		Reason: "explain/run options; optional SDK feature",
		Issue:  "n/a: low priority",
	},
	"vectorQueryWrapper.Documents": {
		Reason: "vector query execution; emulator support incomplete",
		Issue:  "n/a: emulator",
	},
	"vectorQueryWrapper.Serialize": {
		Reason: "vector query serialize; emulator support incomplete",
		Issue:  "n/a: emulator",
	},
	"vectorQueryWrapper.Deserialize": {
		Reason: "vector query deserialize; emulator support incomplete",
		Issue:  "n/a: emulator",
	},

	// Aggregation extras (WithCount/WithSum/WithAvg/Get/GetResponse/Data/DataTo covered in fstest)
	"aggregationQueryWrapper.WithSumPath": {
		Reason: "sum aggregation path variant",
		Issue:  "n/a: extend fstest later",
	},
	"aggregationQueryWrapper.WithAvgPath": {
		Reason: "avg aggregation path variant",
		Issue:  "n/a: extend fstest later",
	},
	"aggregationQueryWrapper.Transaction": {
		Reason: "aggregation inside transaction",
		Issue:  "n/a: extend fstest later",
	},
	"aggregationQueryWrapper.Pipeline": {
		Reason: "pipeline from aggregation; emulator",
		Issue:  "n/a: emulator",
	},
	// Partitioned queries / collection group extras
	"collectionGroupRefWrapper.GetPartitionedQueries": {
		Reason: "partition API; heavy emulator setup",
		Issue:  "n/a: extend fstest later",
	},

	// Query path / select variants
	"queryWrapper.WherePath": {
		Reason: "Where covered; Path variant thin",
		Issue:  "n/a: low priority",
	},
	"queryWrapper.WhereEntity": {
		Reason: "EntityFilter variant",
		Issue:  "n/a: low priority",
	},
	"queryWrapper.OrderByPath": {
		Reason: "OrderBy covered; Path variant thin",
		Issue:  "n/a: low priority",
	},
	"queryWrapper.LimitToLast": {
		Reason: "Limit covered; LimitToLast needs special ordering",
		Issue:  "n/a: extend fstest later",
	},
	"queryWrapper.StartAt": {
		Reason: "StartAfter covered for cursors",
		Issue:  "n/a: low priority",
	},
	"queryWrapper.EndAt": {
		Reason: "StartAfter covered for cursors",
		Issue:  "n/a: low priority",
	},
	"queryWrapper.EndBefore": {
		Reason: "StartAfter covered for cursors",
		Issue:  "n/a: low priority",
	},
	"queryWrapper.SelectPaths": {
		Reason: "Select covered; Path variant thin",
		Issue:  "n/a: low priority",
	},
	// Transaction / batch extras
	"transactionWrapper.Execute": {
		Reason: "pipeline execute inside transaction; emulator",
		Issue:  "n/a: emulator",
	},
	"transactionWrapper.WithReadOptions": {
		Reason: "transaction read options",
		Issue:  "n/a: low priority",
	},

	// Pipeline surface — smoke may cover only a subset
	"clientWrapper.Pipeline": {
		Reason: "pipeline source entry; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineSourceWrapper.Collection": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineSourceWrapper.CollectionGroup": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineSourceWrapper.Database": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineSourceWrapper.Documents": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineSourceWrapper.CreateFromQuery": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineSourceWrapper.CreateFromAggregationQuery": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineSourceWrapper.Literals": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Reference": {
		Reason: "pipeline escape hatch",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Execute": {
		Reason: "pipeline execute; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.WithReadOptions": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Limit": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Sort": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Offset": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Select": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Distinct": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.AddFields": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.RemoveFields": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Where": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Aggregate": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Unnest": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.UnnestWithAlias": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Union": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Sample": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.ReplaceWith": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.FindNearest": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Search": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.RawStage": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Update": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Delete": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.ToScalarExpression": {
		Reason: "pipeline expression helper",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.ToArrayExpression": {
		Reason: "pipeline expression helper",
		Issue:  "n/a: emulator",
	},
	"pipelineWrapper.Define": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineSnapshotWrapper.Results": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineSnapshotWrapper.ExecutionTime": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineSnapshotWrapper.ExplainStats": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineResultWrapper.Ref": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineResultWrapper.CreateTime": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineResultWrapper.UpdateTime": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineResultWrapper.ExecutionTime": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineResultWrapper.Exists": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineResultWrapper.Data": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineResultWrapper.DataTo": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineResultIteratorWrapper.Next": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineResultIteratorWrapper.Stop": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},
	"pipelineResultIteratorWrapper.GetAll": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "n/a: emulator",
	},

	// Iterator explain / page helpers
	"documentIteratorWrapper.ExplainMetrics": {
		Reason: "explain metrics optional",
		Issue:  "n/a: low priority",
	},

	"documentRefWrapper.WithReadOptions": {
		Reason: "same pattern as CollectionRef; clone helper shared",
		Issue:  "n/a: unit",
	},
}
