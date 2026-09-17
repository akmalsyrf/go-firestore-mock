package apicheck

// exceptions documents SDK methods that are intentionally not mirrored 1:1
// on fsmock interfaces. Keys are "Type.Method".
var exceptions = map[string]string{
	"CollectionRef.WithReadOptions": "SDK returns *CollectionRef; fsmock inherits Query.WithReadOptions returning Query (Go embedding cannot overload by return type). Wrapper still applies options to *firestore.CollectionRef.",
}
