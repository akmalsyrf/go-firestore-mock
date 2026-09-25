// Package fsmock provides thin wrapper interfaces and gomock-friendly mocks
// over cloud.google.com/go/firestore for dependency injection and unit testing.
//
// # Versioning
//
// Module path v2. Tags follow v2.<firestore-minor>.<patch> so the minor matches
// the Firestore SDK minor this release was tested against (see version.go and
// COMPATIBILITY.md). Breaking interface changes require a new major module path.
//
// # Scope
//
// Every instance method on the wrapped SDK types (Client, Query, CollectionRef,
// DocumentRef, Transaction, WriteBatch, BulkWriter, AggregationQuery, VectorQuery,
// Pipeline, and related iterators/snapshots) has a matching method on an fsmock
// interface. Pipeline option/expression types (Expression, BooleanExpression,
// SearchOption, etc.) are pass-through SDK values — they are constructed by callers
// and do not need to be mocked.
//
// # Production usage
//
// Construct a real *firestore.Client, then wrap it:
//
//	fs, err := firestore.NewClient(ctx, projectID)
//	client, err := fsmock.NewClient(fs)
//
// Never pass a nil *firestore.Client to NewClient.
package fsmock
