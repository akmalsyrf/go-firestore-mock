package apicheck

// Package apicheck compares exported method sets of cloud.google.com/go/firestore
// types against fsmock interfaces. When the SDK gains new methods, this package's
// tests fail unless the method is wrapped or explicitly listed as an exception.
