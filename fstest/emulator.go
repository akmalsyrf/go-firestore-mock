package fstest

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/akmalsyrf/go-firestore-mock/v2"
)

// EmulatorHostEnv is the standard Firestore emulator environment variable.
const EmulatorHostEnv = "FIRESTORE_EMULATOR_HOST"

// DefaultProjectID used when projectID is empty.
const DefaultProjectID = "demo-fsmock"

// NewEmulatorClient creates an fsmock.Client backed by the Firestore emulator.
// Requires FIRESTORE_EMULATOR_HOST (e.g. "localhost:8080").
func NewEmulatorClient(ctx context.Context, projectID string) (fsmock.Client, *firestore.Client, error) {
	if !EmulatorAvailable() {
		return nil, nil, fmt.Errorf("fstest: %s is not set", EmulatorHostEnv)
	}
	if projectID == "" {
		projectID = DefaultProjectID
	}

	fs, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, nil, err
	}

	client, err := fsmock.NewClient(fs)
	if err != nil {
		_ = fs.Close()
		return nil, nil, err
	}
	return client, fs, nil
}

// EmulatorAvailable reports whether FIRESTORE_EMULATOR_HOST is set.
func EmulatorAvailable() bool {
	return os.Getenv(EmulatorHostEnv) != ""
}

// RequireEmulator skips when FIRESTORE_EMULATOR_HOST is unset.
// The accuracy gate requires the host and treats fstest skips as failures, so this
// only soft-skips local `go test -tags=integration` without an emulator.
func RequireEmulator(t *testing.T) {
	t.Helper()
	if !EmulatorAvailable() {
		t.Skip("FIRESTORE_EMULATOR_HOST not set; start emulator and re-run with -tags=integration")
	}
}

// Harness holds both the wrapped client and the raw SDK client for parity checks.
type Harness struct {
	Ctx    context.Context
	Client fsmock.Client
	Raw    *firestore.Client
	Prefix string // unique collection prefix for this test
}

// NewHarness creates a client pair and a unique prefix derived from the test name.
func NewHarness(t *testing.T) *Harness {
	t.Helper()
	RequireEmulator(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)

	client, raw, err := NewEmulatorClient(ctx, DefaultProjectID)
	if err != nil {
		t.Fatalf("NewEmulatorClient: %v", err)
	}
	t.Cleanup(func() { _ = raw.Close() })

	return &Harness{
		Ctx:    ctx,
		Client: client,
		Raw:    raw,
		Prefix: fmt.Sprintf("t_%d_%s", time.Now().UnixNano(), sanitize(t.Name())),
	}
}

// Coll returns a unique collection path under this harness prefix.
func (h *Harness) Coll(suffix string) string {
	return h.Prefix + "_" + suffix
}

func sanitize(name string) string {
	b := make([]byte, 0, len(name))
	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
			b = append(b, c)
		default:
			b = append(b, '_')
		}
	}
	if len(b) > 40 {
		b = b[:40]
	}
	return string(b)
}
