# Go Firestore Mock

A comprehensive Go library that provides mock implementations and wrapper interfaces for Google Cloud Firestore operations. This library is designed to make testing Firestore-dependent applications easier by providing both mock objects and wrapper interfaces that abstract the Firestore client.

## Features

- **Wrapper Interfaces**: Clean abstraction layer over Firestore client operations
- **Mock Implementations**: Complete mock objects for all Firestore operations
- **Comprehensive Testing**: 100+ unit tests covering all functionality
- **Type Safety**: Full type safety with Go interfaces
- **Easy Integration**: Drop-in replacement for Firestore client in tests

## Installation

```bash
go get github.com/akmalsyrf/go-firestore-mock
```

## Quick Start

### Using Wrapper Interfaces

```go
package main

import (
    "context"
    "log"
    
    "github.com/akmalsyrf/go-firestore-mock"
)

func main() {
    // Create a wrapper client (in production, pass real Firestore client)
    client := gofirestoremock.NewFirestoreClient(nil)
    
    // Use the wrapper as you would use Firestore client
    collection := client.Collection("users")
    doc := collection.Doc("user123")
    
    // Set data
    _, err := doc.Set(context.Background(), map[string]interface{}{
        "name": "John Doe",
        "email": "john@example.com",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Get data
    docSnap, err := doc.Get(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Document data: %v", docSnap.Data())
}
```

### Using Mock Objects

```go
package main

import (
    "context"
    "testing"
    
    "github.com/akmalsyrf/go-firestore-mock"
    "github.com/stretchr/testify/assert"
)

func TestUserOperations(t *testing.T) {
    // Create mock client
    mockClient := gofirestoremock.NewMockFirestoreClient(t)
    
    // Set up expectations
    mockClient.EXPECT().
        Collection("users").
        Return(gofirestoremock.NewMockCollectionRef(t))
    
    // Use the mock
    collection := mockClient.Collection("users")
    assert.NotNil(t, collection)
}
```

## API Reference

### Core Interfaces

#### FirestoreClient
```go
type FirestoreClient interface {
    Collection(path string) CollectionRef
    Doc(path string) DocumentRef
    Close() error
    BulkWriter(ctx context.Context) BulkWriter
    Batch() *firestore.WriteBatch
    RunTransaction(ctx context.Context, f func(context.Context, *firestore.Transaction) error, opts ...firestore.TransactionOption) error
    Collections(ctx context.Context) *firestore.CollectionIterator
    GetAll(ctx context.Context, docRefs []*firestore.DocumentRef) ([]*firestore.DocumentSnapshot, error)
}
```

#### CollectionRef
```go
type CollectionRef interface {
    Query  // Embeds all Query methods
    Doc(id string) DocumentRef
    Add(ctx context.Context, data any) (*firestore.DocumentRef, *firestore.WriteResult, error)
    NewDoc() DocumentRef
    Parent() DocumentRef
    ID() string
    Path() string
}
```

#### DocumentRef
```go
type DocumentRef interface {
    Set(ctx context.Context, data any, opts ...firestore.SetOption) (*firestore.WriteResult, error)
    Get(ctx context.Context) (*firestore.DocumentSnapshot, error)
    Delete(ctx context.Context, opts ...firestore.Precondition) (*firestore.WriteResult, error)
    Update(ctx context.Context, updates []firestore.Update, preconds ...firestore.Precondition) (*firestore.WriteResult, error)
    Create(ctx context.Context, data any) (*firestore.WriteResult, error)
    Collection(path string) CollectionRef
    Collections(ctx context.Context) *firestore.CollectionIterator
    Snapshots(ctx context.Context) *firestore.DocumentSnapshotIterator
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
	OrderBy(path string, dir firestore.Direction) Query
	Limit(n int) Query
	LimitToLast(n int) Query
	Offset(n int) Query
	StartAt(docSnapshotOrFieldValues ...any) Query
	StartAfter(docSnapshotOrFieldValues ...any) Query
	EndAt(docSnapshotOrFieldValues ...any) Query
	EndBefore(docSnapshotOrFieldValues ...any) Query
	Select(paths ...string) Query
	Documents(ctx context.Context) DocumentIterator
	Snapshots(ctx context.Context) *firestore.QuerySnapshotIterator
}
```

#### DocumentIterator
```go
type DocumentIterator interface {
    Next() (*DocumentSnapshot, error)
    Stop()
    GetAll() ([]*DocumentSnapshot, error)
}
```

#### BulkWriter
```go
type BulkWriter interface {
    Create(doc DocumentRef, data interface{}) (WriteResult, error)
    Set(doc DocumentRef, data interface{}, opts ...SetOption) (WriteResult, error)
    Update(doc DocumentRef, updates []Update, opts ...Precondition) (WriteResult, error)
    Delete(doc DocumentRef, opts ...Precondition) (WriteResult, error)
    Flush() error
    End() error
}
```

### Advanced Query Examples

```go
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
```

## Testing

The library includes comprehensive test coverage with 400+ unit tests:

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

- Go 1.19 or later
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
├── client.go              # Core wrapper implementations
├── collection.go          # Collection wrapper
├── document.go            # Document wrapper
├── document_iterator.go   # Document iterator interface
├── bulk_writer.go         # Bulk writer interface
├── *_mock.go             # Mock implementations
├── *_test.go             # Unit tests
├── Makefile              # Build automation
├── .gitignore            # Git ignore rules
└── README.md             # This file
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
