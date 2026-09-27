package apicheck

// waivers lists *xxxWrapper methods that are allowed to lack coverage in the
// accuracy gate. Keys are "receiverType.Method" (e.g. "queryWrapper.Snapshots").
// Reason and Issue are required so waivers do not rot silently.
//
// A waiver means "not exercised on the emulator in CI", not "unused / unwrapped".
// Large clusters (Pipeline*, Vector*, realtime Snapshots*) are waived because the
// emulator cannot reliably prove them — prefer a GitHub issue URL in Issue when
// tracking follow-up coverage (prefer GitHub issue URLs over n/a).
type waiver struct {
	Reason string
	Issue  string // GitHub issue URL or short tracking id; use "n/a: <why>" if none
}

var waivers = map[string]waiver{
	// Realtime listeners need a long-lived stream; emulator coverage is flaky in CI.
	"queryWrapper.Snapshots": {
		Reason: "realtime Query.Snapshots listener; emulator listen flaky in short CI (not behavioral accuracy proof)",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/8",
	},
	"documentRefWrapper.Snapshots": {
		Reason: "realtime DocumentRef.Snapshots listener",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/8",
	},
	"querySnapshotIteratorWrapper.Next": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/8",
	},
	"querySnapshotIteratorWrapper.Stop": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/8",
	},
	"documentSnapshotIteratorWrapper.Next": {
		Reason: "only reachable via DocumentRef.Snapshots",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/8",
	},
	"documentSnapshotIteratorWrapper.Stop": {
		Reason: "only reachable via DocumentRef.Snapshots",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/8",
	},
	"querySnapshotWrapper.Documents": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/8",
	},
	"querySnapshotWrapper.Size": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/8",
	},
	"querySnapshotWrapper.Changes": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/8",
	},
	"querySnapshotWrapper.ReadTime": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/8",
	},
	"querySnapshotWrapper.Reference": {
		Reason: "only reachable via Query.Snapshots",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/8",
	},

	// Client lifecycle / global options rarely exercised in shared emulator suites.
	"clientWrapper.Close": {
		Reason: "client closed by harness cleanup on raw SDK client",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"clientWrapper.WithReadOptions": {
		Reason: "client-level read time; prefer collection/query coverage",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"clientWrapper.WithAlwaysUseImplicitOrderBy": {
		Reason: "global query default; behavioral flag",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"clientWrapper.DocFromFullPath": {
		Reason: "path helper equivalent to Doc for emulator paths",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},

	// Pipeline / vector / enterprise features often unsupported by emulator.
	"queryWrapper.FindNearest": {
		Reason: "vector search; emulator support incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"queryWrapper.FindNearestPath": {
		Reason: "vector search; emulator support incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"queryWrapper.Pipeline": {
		Reason: "pipeline smoke may skip on emulator",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"queryWrapper.WithRunOptions": {
		Reason: "explain/run options; optional SDK feature",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"vectorQueryWrapper.Documents": {
		Reason: "vector query execution; emulator support incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},

	// Aggregation extras (WithCount/WithSum/WithAvg/Get/GetResponse/Data/DataTo covered in fstest)
	"aggregationQueryWrapper.WithSumPath": {
		Reason: "sum aggregation path variant",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"aggregationQueryWrapper.WithAvgPath": {
		Reason: "avg aggregation path variant",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"aggregationQueryWrapper.Transaction": {
		Reason: "aggregation inside transaction",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"aggregationQueryWrapper.Pipeline": {
		Reason: "pipeline from aggregation; emulator",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	// Partitioned queries / collection group extras
	"collectionGroupRefWrapper.GetPartitionedQueries": {
		Reason: "partition API; heavy emulator setup",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},

	// Query path / select variants
	"queryWrapper.WherePath": {
		Reason: "Where covered; Path variant thin",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"queryWrapper.WhereEntity": {
		Reason: "EntityFilter variant",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"queryWrapper.OrderByPath": {
		Reason: "OrderBy covered; Path variant thin",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"queryWrapper.LimitToLast": {
		Reason: "Limit covered; LimitToLast needs special ordering",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"queryWrapper.StartAt": {
		Reason: "StartAfter covered for cursors; StartAt is the inclusive sibling",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"queryWrapper.EndAt": {
		Reason: "StartAfter covered for cursors; EndAt is the inclusive end sibling",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"queryWrapper.EndBefore": {
		Reason: "StartAfter covered for cursors; EndBefore is the exclusive end sibling",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"queryWrapper.SelectPaths": {
		Reason: "Select covered; Path variant thin",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	// Transaction / batch extras
	"transactionWrapper.Execute": {
		Reason: "pipeline execute inside transaction; emulator",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"transactionWrapper.WithReadOptions": {
		Reason: "transaction read options",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},

	// Pipeline surface — smoke may cover only a subset
	"clientWrapper.Pipeline": {
		Reason: "pipeline source entry; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineSourceWrapper.Collection": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineSourceWrapper.CollectionGroup": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineSourceWrapper.Database": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineSourceWrapper.Documents": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineSourceWrapper.CreateFromQuery": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineSourceWrapper.CreateFromAggregationQuery": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineSourceWrapper.Literals": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Reference": {
		Reason: "pipeline escape hatch",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Execute": {
		Reason: "pipeline execute; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.WithReadOptions": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Limit": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Sort": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Offset": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Select": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Distinct": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.AddFields": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.RemoveFields": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Where": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Aggregate": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Unnest": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.UnnestWithAlias": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Union": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Sample": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.ReplaceWith": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.FindNearest": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Search": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.RawStage": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Update": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Delete": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.ToScalarExpression": {
		Reason: "pipeline expression helper",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.ToArrayExpression": {
		Reason: "pipeline expression helper",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineWrapper.Define": {
		Reason: "pipeline stage; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineSnapshotWrapper.Results": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineSnapshotWrapper.ExecutionTime": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineSnapshotWrapper.ExplainStats": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineResultWrapper.Ref": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineResultWrapper.CreateTime": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineResultWrapper.UpdateTime": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineResultWrapper.ExecutionTime": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineResultWrapper.Exists": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineResultWrapper.Data": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineResultWrapper.DataTo": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineResultIteratorWrapper.Next": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineResultIteratorWrapper.Stop": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
	"pipelineResultIteratorWrapper.GetAll": {
		Reason: "pipeline; emulator incomplete",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},

	// Iterator explain / page helpers
	"documentIteratorWrapper.ExplainMetrics": {
		Reason: "explain metrics optional",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},

	"documentRefWrapper.WithReadOptions": {
		Reason: "same pattern as CollectionRef; clone helper shared",
		Issue:  "https://github.com/akmalsyrf/go-firestore-mock/issues/10",
	},
}
