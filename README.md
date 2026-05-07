# Go Firestore Mock [![GoDoc](https://godoc.org/github.com/akmalsyrf/go-firestore-mock?status.svg)](https://godoc.org/github.com/akmalsyrf/go-firestore-mock)

A comprehensive Go library that provides mock implementations and wrapper interfaces for Google Cloud Firestore operations. This library is designed to make testing Firestore-dependent applications easier by providing both mock objects and wrapper interfaces that abstract the Firestore client.

## Features

- **Wrapper interfaces**: Thin abstraction over the official `cloud.google.com/go/firestore` client for testability and dependency injection
- **gomock-generated mocks**: Generated mocks for each interface (use with `go.uber.org/mock/gomock`)
- **Unit tests**: Broad coverage of wrappers and mocks (on the order of ~190 `Test*` entry points in this module)
- **Type safety**: Call sites depend on interfaces, not concrete `*firestore.Client` types

## Compatibility

This module re-exports types from **`cloud.google.com/go/firestore`** (see `go.mod` for the pinned minor version). In production you **must** construct a real client with `firestore.NewClient` / `firestore.NewClientWithDatabase` (or your app’s factory), then wrap it with `NewFirestoreClient`. The wrapper is a **subset** of the full Firestore API; see the interface definitions in the source for what is supported.

## Installation

```bash
go get github.com/akmalsyrf/go-firestore-mock
```

## Quick Start (production)

Use the official Firestore client, then wrap it. **Never** pass `nil` into `NewFirestoreClient`.

```go
package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	gofirestoremock "github.com/akmalsyrf/go-firestore-mock"
)

func main() {
	ctx := context.Background()
	fs, err := firestore.NewClient(ctx, "your-project-id")
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()

	client := gofirestoremock.NewFirestoreClient(fs)

	collection := client.Collection("users")
	doc := collection.Doc("user123")

	_, err = doc.Set(ctx, map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
	})
	if err != nil {
		log.Fatal(err)
	}

	docSnap, err := doc.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Document data: %v", docSnap.Data())
}
```

## Testing with gomock

```go
package mypkg_test

import (
	"testing"

	gofirestoremock "github.com/akmalsyrf/go-firestore-mock"
	"go.uber.org/mock/gomock"
)

func TestWithMockClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := gofirestoremock.NewMockFirestoreClient(ctrl)
	mockColl := gofirestoremock.NewMockCollectionRef(ctrl)

	mockClient.EXPECT().Collection("users").Return(mockColl)

	collection := mockClient.Collection("users")
	if collection == nil {
		t.Fatal("expected non-nil collection ref")
	}
}
```

## API Reference

### Core Interfaces

#### FirestoreClient
```go
type FirestoreClient interface {
    Collection(path string) CollectionRef
    CollectionGroup(collectionID string) Query
    Doc(path string) DocumentRef
    DocFromFullPath(fullPath string) DocumentRef
    Close() error
    BulkWriter(ctx context.Context) BulkWriter
    Batch() WriteBatch
    RunTransaction(ctx context.Context, f func(context.Context, Transaction) error, opts ...firestore.TransactionOption) error
    Collections(ctx context.Context) CollectionIterator
    GetAll(ctx context.Context, docRefs []*firestore.DocumentRef) ([]DocumentSnapshot, error)
}
```

#### CollectionRef
```go
type CollectionRef interface {
    Query  // Embeds all Query methods
    Doc(id string) DocumentRef
    Add(ctx context.Context, data any) (*firestore.DocumentRef, *firestore.WriteResult, error)
    NewDoc() DocumentRef
    DocumentRefs(ctx context.Context) DocumentRefIterator
    Parent() DocumentRef
    Reference() *firestore.CollectionRef
    ID() string
    Path() string
}
```

#### DocumentRef
```go
type DocumentRef interface {
    Set(ctx context.Context, data any, opts ...firestore.SetOption) (*firestore.WriteResult, error)
	Get(ctx context.Context) (DocumentSnapshot, error)
    Delete(ctx context.Context, opts ...firestore.Precondition) (*firestore.WriteResult, error)
    Update(ctx context.Context, updates []firestore.Update, preconds ...firestore.Precondition) (*firestore.WriteResult, error)
    Create(ctx context.Context, data any) (*firestore.WriteResult, error)
    Collection(path string) CollectionRef
	Collections(ctx context.Context) CollectionIterator
	Snapshots(ctx context.Context) DocumentSnapshotIterator
    Reference() *firestore.DocumentRef
    ID() string
    Path() string
    Parent() *firestore.CollectionRef
}
```

#### Query
```go
type Query interface {
    Where(path, op string, value any) Query
    WherePath(fp firestore.FieldPath, op string, value any) Query
    WhereEntity(ef firestore.EntityFilter) Query
    OrderBy(path string, dir firestore.Direction) Query
    OrderByPath(fp firestore.FieldPath, dir firestore.Direction) Query
    Limit(n int) Query
    LimitToLast(n int) Query
    Offset(n int) Query
    StartAt(docSnapshotOrFieldValues ...any) Query
    StartAfter(docSnapshotOrFieldValues ...any) Query
    EndAt(docSnapshotOrFieldValues ...any) Query
    EndBefore(docSnapshotOrFieldValues ...any) Query
    Select(paths ...string) Query
    SelectPaths(fieldPaths ...firestore.FieldPath) Query
    Documents(ctx context.Context) DocumentIterator
    Snapshots(ctx context.Context) QuerySnapshotIterator
    NewAggregationQuery() AggregationQuery
}
```

#### DocumentIterator / DocumentRefIterator
```go
type DocumentIterator interface {
    Next() (*firestore.DocumentSnapshot, error)
    Stop()
    GetAll() ([]*firestore.DocumentSnapshot, error)
}

type DocumentRefIterator interface {
    Next() (*firestore.DocumentRef, error)
    GetAll() ([]*firestore.DocumentRef, error)
}
```

#### BulkWriter
```go
type BulkWriter interface {
    Create(docRef *firestore.DocumentRef, data interface{}) (*firestore.BulkWriterJob, error)
    Set(docRef *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) (*firestore.BulkWriterJob, error)
    Update(docRef *firestore.DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) (*firestore.BulkWriterJob, error)
    Delete(docRef *firestore.DocumentRef, preconds ...firestore.Precondition) (*firestore.BulkWriterJob, error)
    Flush()
    End()
}
```

#### WriteBatch
```go
type WriteBatch interface {
    Create(docRef *firestore.DocumentRef, data interface{}) WriteBatch
    Set(docRef *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) WriteBatch
    Update(docRef *firestore.DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) WriteBatch
    Delete(docRef *firestore.DocumentRef, preconds ...firestore.Precondition) WriteBatch
    Commit(ctx context.Context) ([]*firestore.WriteResult, error)
}
```

#### Transaction
```go
type Transaction interface {
    Get(docRef *firestore.DocumentRef) (DocumentSnapshot, error)
    GetAll(docRefs []*firestore.DocumentRef) ([]DocumentSnapshot, error)
    // Documents accepts a Query or CollectionRef (CollectionRef embeds Query).
    Documents(q Query) DocumentIterator
    DocumentRefs(coll CollectionRef) DocumentRefIterator
    Create(docRef *firestore.DocumentRef, data interface{}) error
    Set(docRef *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) error
    Update(docRef *firestore.DocumentRef, updates []firestore.Update, preconds ...firestore.Precondition) error
    Delete(docRef *firestore.DocumentRef, preconds ...firestore.Precondition) error
}
```

> The wrapper around a real `*firestore.Transaction` extracts the underlying
> `firestore.Queryer` from `*queryWrapper` / `*collectionRefWrapper` (and from
> custom types that implement `CollectionRef.Reference()`). Pure mock
> implementations of `Query` cannot be passed to the production wrapper – mock
> the `Transaction` interface in tests instead.

#### AggregationQuery
```go
type AggregationQuery interface {
    WithCount(alias string) AggregationQuery
    Get(ctx context.Context) (AggregationResult, error)
}

type AggregationResult interface {
    Count(alias string) (*int64, error)
}
```

#### DocumentSnapshot
```go
type DocumentSnapshot interface {
    Data() map[string]interface{}
    DataTo(p interface{}) error
    DataAt(path string) (interface{}, error)
    DataAtPath(fp firestore.FieldPath) (interface{}, error)
    Exists() bool
    CreateTime() time.Time
    UpdateTime() time.Time
    ReadTime() time.Time
    Ref() *firestore.DocumentRef
}
```

#### Iterator Wrappers
```go
type QuerySnapshotIterator interface {
    Next() (*firestore.QuerySnapshot, error)
    Stop()
}

type CollectionIterator interface {
    Next() (*firestore.CollectionRef, error)
    Stop()
}

type DocumentSnapshotIterator interface {
    Next() (DocumentSnapshot, error)
    Stop()
}
```

### Advanced Query Examples

```go
// Collection Group Query (query across all subcollections with same ID)
allProducts := client.CollectionGroup("products").
    Where("price", ">", 100).
    OrderBy("price", firestore.Asc).
    Limit(50)

// Pagination with OrderBy and Limit
results := collection.
    Where("status", "==", "active").
    OrderBy("createdAt", firestore.Desc).
    Limit(10).
    Offset(20)

// Cursor-based pagination
lastDoc := getPreviousLastDocument()
nextPage := collection.
    OrderBy("createdAt", firestore.Desc).
    StartAfter(lastDoc).
    Limit(10)

// Complex filtering and sorting
products := collection.
    Where("category", "==", "electronics").
    Where("price", ">", 100).
    OrderBy("price", firestore.Asc).
    OrderBy("rating", firestore.Desc).
    Limit(50)

// Field selection
users := collection.
    Where("active", "==", true).
    Select("name", "email").
    Limit(100)

// Count aggregation query
countQuery := collection.
    Where("status", "==", "published").
    NewAggregationQuery().
    WithCount("total")

result, err := countQuery.Get(ctx)
if err != nil {
    log.Fatal(err)
}
count, err := result.Count("total")
if err != nil {
	log.Fatal(err)
}
if count != nil {
	fmt.Printf("Total documents: %d\n", *count)
}
```

### Batch Operations

```go
// WriteBatch for atomic writes
batch := client.Batch()
batch.
    Create(docRef1, data1).
    Set(docRef2, data2).
    Update(docRef3, updates).
    Delete(docRef4)

results, err := batch.Commit(ctx)

// Transaction for atomic read-writes
err := client.RunTransaction(ctx, func(ctx context.Context, tx Transaction) error {
    // Read operations
    doc, err := tx.Get(docRef)
    if err != nil {
        return err
    }
    
    // Write operations
    return tx.Set(docRef, newData)
})

// Read-then-write inside a transaction (e.g. delete every document in a
// subcollection). All reads must happen before any writes per Firestore rules.
err := client.RunTransaction(ctx, func(ctx context.Context, tx Transaction) error {
    coll := client.Collection("matches").Doc(matchID).Collection("undo_actions")

    iter := tx.Documents(coll)
    defer iter.Stop()

    var refs []*firestore.DocumentRef
    for {
        snap, err := iter.Next()
        if err == iterator.Done {
            break
        }
        if err != nil {
            return err
        }
        refs = append(refs, snap.Ref)
    }

    for _, ref := range refs {
        if err := tx.Delete(ref); err != nil {
            return err
        }
    }
    return nil
})
```

### Document Operations

```go
// Create a new document (fails if exists)
_, err := doc.Create(ctx, map[string]interface{}{
    "name": "John Doe",
    "email": "john@example.com",
})

// Update specific fields
_, err := doc.Update(ctx, []firestore.Update{
    {Path: "name", Value: "Jane Doe"},
    {Path: "updatedAt", Value: firestore.ServerTimestamp},
})

// Get document metadata
docID := doc.ID()
docPath := doc.Path()
parentCollection := doc.Parent()

// List subcollections
collections := doc.Collections(ctx)

// Create auto-generated ID document
newDoc := collection.NewDoc()
_, err := newDoc.Set(ctx, data)

// Get collection metadata
collectionID := collection.ID()
collectionPath := collection.Path()
parentDoc := collection.Parent()

// Working with DocumentSnapshot
snap, err := doc.Get(ctx)
if err != nil {
    log.Fatal(err)
}

if snap.Exists() {
    // Get all data
    data := snap.Data()
    
    // Get specific field
    name, _ := snap.DataAt("name")
    
    // Decode to struct
    var user User
    snap.DataTo(&user)
    
    // Get metadata
    createdAt := snap.CreateTime()
    updatedAt := snap.UpdateTime()
}
```

## Testing

Run the test suite:

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Generate HTML coverage report
make test-coverage-html

# Run tests with race detection
make test-race
```

## Development

### Prerequisites

- Go 1.25 or later (required by `cloud.google.com/go/firestore v1.22.0`)
- Make (optional, for using Makefile)

### Available Make Commands

```bash
make help              # Show all available commands
make all               # Clean, deps, fmt, vet, test, and build
make build             # Build the application
make test              # Run tests
make test-coverage     # Run tests with coverage
make fmt               # Format code
make vet               # Run go vet
make lint              # Run linter
make clean             # Clean build artifacts
make deps              # Download dependencies
make deps-update       # Update dependencies
```

### Project Structure

```
go-firestore-mock/
├── client.go                    # Core wrapper implementations
├── collection.go                # Collection wrapper
├── document.go                  # Document wrapper
├── document_snapshot.go         # Document snapshot wrapper
├── document_iterator.go         # Document iterator interface
├── snapshot_iterators.go        # Iterator wrappers
├── aggregation.go               # Aggregation query wrapper
├── bulk_writer.go               # Bulk writer interface
├── write_batch.go               # Write batch interface
├── transaction.go               # Transaction interface
├── *_mock.go                   # Mock implementations
├── *_test.go                   # Unit tests
├── Makefile                    # Build automation
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── LICENSE                     # MIT License
└── README.md                   # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Google Cloud Firestore team for the excellent Go client library
- The Go community for testing best practices
- All contributors who help improve this library
