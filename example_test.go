package fsmock_test

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/akmalsyrf/go-firestore-mock/v2"
)

// ExampleNewClient shows the production wrap pattern.
func ExampleNewClient() {
	ctx := context.Background()
	fs, err := firestore.NewClient(ctx, "demo-project")
	if err != nil {
		fmt.Println("skip:", err)
		return
	}
	defer func() { _ = fs.Close() }()

	client, err := fsmock.NewClient(fs)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	_ = client
	fmt.Println("ok")
}

// ExampleMustNewClient shows the panic-on-nil constructor.
func ExampleMustNewClient() {
	defer func() {
		if recover() != nil {
			fmt.Println("panicked on nil")
		}
	}()
	_ = fsmock.MustNewClient(nil)
}
