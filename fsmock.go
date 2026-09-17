package fsmock

import "cloud.google.com/go/firestore"

// NewClient wraps a real *firestore.Client. client must not be nil.
func NewClient(client *firestore.Client) (Client, error) {
	if client == nil {
		return nil, ErrNilClient
	}
	return &clientWrapper{client: client}, nil
}

// MustNewClient is like NewClient but panics on error.
func MustNewClient(client *firestore.Client) Client {
	c, err := NewClient(client)
	if err != nil {
		panic(err)
	}
	return c
}
