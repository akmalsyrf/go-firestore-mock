package fsmock

import (
	"context"

	"cloud.google.com/go/firestore"
)

// VectorQuery abstracts firestore.VectorQuery.
type VectorQuery interface {
	Documents(ctx context.Context) DocumentIterator
	Serialize() ([]byte, error)
	Deserialize(bytes []byte) (VectorQuery, error)
}

type vectorQueryWrapper struct {
	vq firestore.VectorQuery
}

func (w *vectorQueryWrapper) Documents(ctx context.Context) DocumentIterator {
	return newDocumentIterator(w.vq.Documents(ctx))
}

func (w *vectorQueryWrapper) Serialize() ([]byte, error) {
	return w.vq.Serialize()
}

func (w *vectorQueryWrapper) Deserialize(bytes []byte) (VectorQuery, error) {
	vq, err := w.vq.Deserialize(bytes)
	if err != nil {
		return nil, err
	}
	return &vectorQueryWrapper{vq: vq}, nil
}
