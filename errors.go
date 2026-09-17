package fsmock

import "errors"

// Sentinel errors for unwrap and construction failures.
// Use errors.Is to match them.
var (
	// ErrNilClient is returned by NewClient when the underlying SDK client is nil.
	ErrNilClient = errors.New("fsmock: nil *firestore.Client")

	// ErrNilArgument is returned when a required wrapper argument is nil.
	ErrNilArgument = errors.New("fsmock: nil argument")

	// ErrForeignImplementation is returned when an interface value is not one of
	// the package wrappers produced by NewClient (and has no usable Reference()).
	ErrForeignImplementation = errors.New("fsmock: foreign interface implementation")
)
