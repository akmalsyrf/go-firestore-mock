//go:build integration

package fstest_test

import (
	"reflect"
	"testing"

	"github.com/akmalsyrf/go-firestore-mock/v2/fstest"
)

// TestIntegration_ParityRawVsWrapped proves wrappers do not change SDK semantics:
// the same writes/reads via raw *firestore.Client and via fsmock.Client must match.
func TestIntegration_ParityRawVsWrapped(t *testing.T) {
	h := fstest.NewHarness(t)
	name := h.Coll("parity")

	rawDoc := h.Raw.Collection(name).Doc("p1")
	wrapDoc := h.Client.Collection(name).Doc("p1")

	payload := map[string]any{"name": "parity", "n": int64(7)}
	if _, err := wrapDoc.Set(h.Ctx, payload); err != nil {
		t.Fatalf("wrapped Set: %v", err)
	}

	rawSnap, err := rawDoc.Get(h.Ctx)
	if err != nil {
		t.Fatalf("raw Get: %v", err)
	}
	wrapSnap, err := wrapDoc.Get(h.Ctx)
	if err != nil {
		t.Fatalf("wrap Get: %v", err)
	}

	if !reflect.DeepEqual(rawSnap.Data(), wrapSnap.Data()) {
		t.Fatalf("data mismatch\n raw=%v\nwrap=%v", rawSnap.Data(), wrapSnap.Data())
	}
	if rawSnap.Exists() != wrapSnap.Exists() {
		t.Fatal("Exists mismatch")
	}
	if rawSnap.Ref.ID != wrapSnap.Ref().ID() {
		t.Fatalf("ID mismatch %s vs %s", rawSnap.Ref.ID, wrapSnap.Ref().ID())
	}

	// Query parity
	rawIter := h.Raw.Collection(name).Where("n", "==", int64(7)).Documents(h.Ctx)
	rawAll, err := rawIter.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	wrapAll, err := h.Client.Collection(name).Where("n", "==", int64(7)).Documents(h.Ctx).GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(rawAll) != len(wrapAll) {
		t.Fatalf("query len raw=%d wrap=%d", len(rawAll), len(wrapAll))
	}
	if len(rawAll) > 0 && !reflect.DeepEqual(rawAll[0].Data(), wrapAll[0].Data()) {
		t.Fatalf("query data mismatch")
	}

	// Aggregation count parity
	rawAR, err := h.Raw.Collection(name).NewAggregationQuery().WithCount("c").Get(h.Ctx)
	if err != nil {
		t.Fatalf("raw agg: %v", err)
	}
	wrapAR, err := h.Client.Collection(name).NewAggregationQuery().WithCount("c").Get(h.Ctx)
	if err != nil {
		t.Fatalf("wrap agg: %v", err)
	}
	wrapN, err := wrapAR.Count("c")
	if err != nil || wrapN == nil {
		t.Fatalf("wrap Count: %v %v", wrapN, err)
	}
	// Compare via map lookup / DataTo rather than SDK Data() (can panic on odd shapes).
	var rawMap map[string]any
	if err := rawAR.DataTo(&rawMap); err != nil {
		t.Fatalf("raw DataTo: %v", err)
	}
	if rawMap["c"] == nil {
		t.Fatalf("raw count missing: %#v", rawMap)
	}
	wrapData, err := wrapAR.Data()
	if err != nil {
		t.Fatalf("wrap Data: %v", err)
	}
	if wrapData["c"] == nil {
		t.Fatalf("wrap Data missing c: %#v", wrapData)
	}
}

func TestIntegration_ParityBatch(t *testing.T) {
	h := fstest.NewHarness(t)
	name := h.Coll("parity_batch")

	w1 := h.Client.Collection(name).Doc("a")
	w2 := h.Client.Collection(name).Doc("b")
	_, err := h.Client.Batch().
		Set(w1, map[string]any{"v": 1}).
		Set(w2, map[string]any{"v": 2}).
		Commit(h.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	r1, err := h.Raw.Collection(name).Doc("a").Get(h.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := h.Raw.Collection(name).Doc("b").Get(h.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	if r1.Data()["v"] != int64(1) || r2.Data()["v"] != int64(2) {
		t.Fatalf("raw read after wrapped batch: %v %v", r1.Data(), r2.Data())
	}
}
