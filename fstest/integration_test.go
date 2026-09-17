//go:build integration

package fstest_test

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/akmalsyrf/go-firestore-mock/v2"
	"github.com/akmalsyrf/go-firestore-mock/v2/fstest"
	"google.golang.org/api/iterator"
)

func TestIntegration_DocumentCRUD(t *testing.T) {
	h := fstest.NewHarness(t)
	coll := h.Coll("crud")
	doc := h.Client.Collection(coll).Doc("u1")

	if _, err := doc.Create(h.Ctx, map[string]any{"name": "Ada", "n": 1}); err != nil {
		t.Fatalf("Create: %v", err)
	}
	snap, err := doc.Get(h.Ctx)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !snap.Exists() || snap.Data()["name"] != "Ada" {
		t.Fatalf("snap: exists=%v data=%v", snap.Exists(), snap.Data())
	}
	if snap.Ref().ID() != "u1" {
		t.Fatalf("Ref.ID=%s", snap.Ref().ID())
	}
	if snap.CreateTime().IsZero() || snap.UpdateTime().IsZero() || snap.ReadTime().IsZero() {
		t.Fatal("expected timestamps")
	}

	if _, err := doc.Set(h.Ctx, map[string]any{"name": "Ada", "n": 2}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if _, err := doc.Update(h.Ctx, []firestore.Update{{Path: "n", Value: 3}}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	snap, err = doc.Get(h.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	if snap.Data()["n"] != int64(3) {
		t.Fatalf("n=%v", snap.Data()["n"])
	}

	v, err := snap.DataAt("name")
	if err != nil || v != "Ada" {
		t.Fatalf("DataAt: %v %v", v, err)
	}
	v, err = snap.DataAtPath(firestore.FieldPath{"name"})
	if err != nil || v != "Ada" {
		t.Fatalf("DataAtPath: %v %v", v, err)
	}
	var out struct {
		Name string `firestore:"name"`
		N    int64  `firestore:"n"`
	}
	if err := snap.DataTo(&out); err != nil || out.Name != "Ada" || out.N != 3 {
		t.Fatalf("DataTo: %+v %v", out, err)
	}

	if _, err := doc.Delete(h.Ctx); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	snap, err = doc.Get(h.Ctx)
	if err == nil && snap.Exists() {
		t.Fatal("expected missing after delete")
	}
}

func TestIntegration_CollectionAddNewDocParent(t *testing.T) {
	h := fstest.NewHarness(t)
	coll := h.Client.Collection(h.Coll("add"))
	ref, _, err := coll.Add(h.Ctx, map[string]any{"x": 1})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if ref.ID() == "" || ref.Parent().ID() == "" {
		t.Fatalf("ref=%s parent=%s", ref.ID(), ref.Parent().ID())
	}
	nd := coll.NewDoc()
	if _, err := nd.Set(h.Ctx, map[string]any{"y": 2}); err != nil {
		t.Fatalf("NewDoc Set: %v", err)
	}
	if coll.ID() == "" || coll.Path() == "" || coll.Reference() == nil {
		t.Fatal("collection metadata")
	}
}

func TestIntegration_SubcollectionAndCollectionsIterator(t *testing.T) {
	h := fstest.NewHarness(t)
	parent := h.Client.Collection(h.Coll("parent")).Doc("p1")
	if _, err := parent.Set(h.Ctx, map[string]any{"ok": true}); err != nil {
		t.Fatal(err)
	}
	sub := parent.Collection("items")
	if _, err := sub.Doc("i1").Set(h.Ctx, map[string]any{"n": 1}); err != nil {
		t.Fatal(err)
	}

	it := parent.Collections(h.Ctx)
	found := false
	for {
		c, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if c.ID() == "items" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected subcollection items")
	}
	if it.PageInfo() == nil {
		t.Fatal("PageInfo nil")
	}
}

func TestIntegration_QueryAndDocumentIterator(t *testing.T) {
	h := fstest.NewHarness(t)
	collName := h.Coll("query")
	coll := h.Client.Collection(collName)
	for i, name := range []string{"a", "b", "c"} {
		if _, err := coll.Doc(name).Set(h.Ctx, map[string]any{"name": name, "n": i}); err != nil {
			t.Fatal(err)
		}
	}

	q := coll.Where("n", ">=", 1).OrderBy("n", firestore.Asc).Limit(2)
	iter := q.Documents(h.Ctx)
	defer iter.Stop()

	var got []string
	for {
		s, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, s.Data()["name"].(string))
	}
	if len(got) != 2 || got[0] != "b" || got[1] != "c" {
		t.Fatalf("got=%v", got)
	}

	all, err := coll.Where("n", "==", 0).Documents(h.Ctx).GetAll()
	if err != nil || len(all) != 1 {
		t.Fatalf("GetAll: %d %v", len(all), err)
	}
}

func TestIntegration_DocumentRefsIterator(t *testing.T) {
	h := fstest.NewHarness(t)
	coll := h.Client.Collection(h.Coll("refs"))
	if _, err := coll.Doc("r1").Set(h.Ctx, map[string]any{"a": 1}); err != nil {
		t.Fatal(err)
	}
	it := coll.DocumentRefs(h.Ctx)
	if it.PageInfo() == nil {
		t.Fatal("DocumentRefIterator.PageInfo nil")
	}
	refs, err := it.GetAll()
	if err != nil || len(refs) < 1 {
		t.Fatalf("DocumentRefs: %d %v", len(refs), err)
	}
}

func TestIntegration_WriteBatch(t *testing.T) {
	h := fstest.NewHarness(t)
	coll := h.Client.Collection(h.Coll("batch"))
	d1, d2 := coll.Doc("b1"), coll.Doc("b2")

	_, err := h.Client.Batch().
		Set(d1, map[string]any{"v": 1}).
		Set(d2, map[string]any{"v": 2}).
		Update(d1, []firestore.Update{{Path: "v", Value: 10}}).
		Commit(h.Ctx)
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}
	s, err := d1.Get(h.Ctx)
	if err != nil || s.Data()["v"] != int64(10) {
		t.Fatalf("d1=%v err=%v", s.Data(), err)
	}
	_, err = h.Client.Batch().Delete(d2).Commit(h.Ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIntegration_Transaction(t *testing.T) {
	h := fstest.NewHarness(t)
	coll := h.Client.Collection(h.Coll("tx"))
	doc := coll.Doc("t1")
	if _, err := doc.Set(h.Ctx, map[string]any{"n": 1}); err != nil {
		t.Fatal(err)
	}

	err := h.Client.RunTransaction(h.Ctx, func(ctx context.Context, tx fsmock.Transaction) error {
		// All reads must happen before writes.
		snap, err := tx.Get(doc)
		if err != nil {
			return err
		}
		iter, err := tx.Documents(coll.Where("n", ">=", 0))
		if err != nil {
			return err
		}
		defer iter.Stop()
		if _, err := iter.GetAll(); err != nil {
			return err
		}
		n := snap.Data()["n"].(int64)
		return tx.Set(doc, map[string]any{"n": n + 1})
	})
	if err != nil {
		t.Fatalf("RunTransaction: %v", err)
	}
	snap, err := doc.Get(h.Ctx)
	if err != nil || snap.Data()["n"] != int64(2) {
		t.Fatalf("after tx: %v %v", snap.Data(), err)
	}
}

func TestIntegration_BulkWriter(t *testing.T) {
	h := fstest.NewHarness(t)
	coll := h.Client.Collection(h.Coll("bw"))
	bw := h.Client.BulkWriter(h.Ctx)

	d1, d2 := coll.Doc("w1"), coll.Doc("w2")
	if _, err := bw.Set(d1, map[string]any{"a": 1}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if _, err := bw.Set(d2, map[string]any{"a": 2}); err != nil {
		t.Fatalf("Set d2: %v", err)
	}
	bw.Flush()
	bw.End()

	s, err := d1.Get(h.Ctx)
	if err != nil {
		t.Fatalf("Get d1: %v", err)
	}
	if !s.Exists() || s.Data()["a"] != int64(1) {
		t.Fatalf("after flush: data=%v", s.Data())
	}
	s2, err := d2.Get(h.Ctx)
	if err != nil || !s2.Exists() {
		t.Fatalf("d2 missing: %v", err)
	}
}

func TestIntegration_Aggregation(t *testing.T) {
	h := fstest.NewHarness(t)
	coll := h.Client.Collection(h.Coll("agg"))
	for i := 1; i <= 3; i++ {
		if _, err := coll.Doc(fmt.Sprintf("d%d", i)).Set(h.Ctx, map[string]any{"n": i}); err != nil {
			t.Fatal(err)
		}
	}

	aq := coll.NewAggregationQuery().WithCount("c").WithSum("n", "s").WithAvg("n", "avg")
	ar, err := aq.Get(h.Ctx)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	c, err := ar.Count("c")
	if err != nil || c == nil || *c != 3 {
		t.Fatalf("count=%v err=%v", c, err)
	}
	data, err := ar.Data()
	if err != nil {
		t.Fatalf("Data: %v", err)
	}
	if data["c"] == nil {
		t.Fatalf("Data map: %#v", data)
	}
	var m map[string]any
	if err := ar.DataTo(&m); err != nil {
		t.Fatalf("DataTo: %v", err)
	}

	resp, err := coll.Where("n", ">", 0).NewAggregationQuery().WithCount("c").GetResponse(h.Ctx)
	if err != nil || resp == nil || resp.Result == nil {
		t.Fatalf("GetResponse: %v %#v", err, resp)
	}
}

func TestIntegration_GetAllAndClientCollections(t *testing.T) {
	h := fstest.NewHarness(t)
	coll := h.Client.Collection(h.Coll("getall"))
	d1, d2 := coll.Doc("g1"), coll.Doc("g2")
	if _, err := d1.Set(h.Ctx, map[string]any{"k": 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := d2.Set(h.Ctx, map[string]any{"k": 2}); err != nil {
		t.Fatal(err)
	}
	snaps, err := h.Client.GetAll(h.Ctx, []fsmock.DocumentRef{d1, d2})
	if err != nil || len(snaps) != 2 {
		t.Fatalf("GetAll: %d %v", len(snaps), err)
	}

	it := h.Client.Collections(h.Ctx)
	_, err = it.Next()
	if err != nil && err != iterator.Done {
		t.Fatalf("Collections: %v", err)
	}
}

func TestIntegration_CollectionGroup(t *testing.T) {
	h := fstest.NewHarness(t)
	// Same collection ID under two parents
	p1 := h.Client.Collection(h.Coll("cg1")).Doc("p")
	p2 := h.Client.Collection(h.Coll("cg2")).Doc("p")
	if _, err := p1.Set(h.Ctx, map[string]any{"x": 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := p2.Set(h.Ctx, map[string]any{"x": 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := p1.Collection("notes").Doc("n1").Set(h.Ctx, map[string]any{"tag": "cg"}); err != nil {
		t.Fatal(err)
	}
	if _, err := p2.Collection("notes").Doc("n2").Set(h.Ctx, map[string]any{"tag": "cg"}); err != nil {
		t.Fatal(err)
	}

	cg := h.Client.CollectionGroup("notes")
	all, err := cg.Where("tag", "==", "cg").Documents(h.Ctx).GetAll()
	if err != nil {
		t.Fatalf("CollectionGroup query: %v", err)
	}
	if len(all) < 2 {
		t.Fatalf("expected >=2 notes, got %d", len(all))
	}
}

func TestIntegration_QuerySelectOffsetCursors(t *testing.T) {
	h := fstest.NewHarness(t)
	coll := h.Client.Collection(h.Coll("cursor"))
	for i := 0; i < 5; i++ {
		if _, err := coll.Doc(fmt.Sprintf("d%d", i)).Set(h.Ctx, map[string]any{"n": i, "extra": "x"}); err != nil {
			t.Fatal(err)
		}
	}
	q := coll.OrderBy("n", firestore.Asc).Offset(1).Limit(2).Select("n")
	snaps, err := q.Documents(h.Ctx).GetAll()
	if err != nil || len(snaps) != 2 {
		t.Fatalf("got %d err=%v", len(snaps), err)
	}
	if _, ok := snaps[0].Data()["extra"]; ok {
		t.Fatal("Select should omit extra")
	}

	first, err := coll.OrderBy("n", firestore.Asc).Limit(1).Documents(h.Ctx).GetAll()
	if err != nil || len(first) != 1 {
		t.Fatal(err)
	}
	// Use field values for cursors: wrapped DocumentSnapshot is not a *firestore.DocumentSnapshot.
	next, err := coll.OrderBy("n", firestore.Asc).StartAfter(first[0].Data()["n"]).Limit(1).Documents(h.Ctx).GetAll()
	if err != nil || len(next) != 1 {
		t.Fatalf("StartAfter: %v", err)
	}
	if next[0].Data()["n"] != int64(1) {
		t.Fatalf("cursor n=%v", next[0].Data()["n"])
	}
}

func TestIntegration_BSONRoundTrip(t *testing.T) {
	h := fstest.NewHarness(t)
	doc := h.Client.Collection(h.Coll("bson")).Doc("b1")
	in := map[string]any{
		"oid": firestore.BSONObjectID("507f1f77bcf86cd799439011"),
		"i32": firestore.BSONInt32(42),
	}
	if _, err := doc.Set(h.Ctx, in); err != nil {
		// Emulator builds often lack BSON enterprise types.
		t.Skipf("BSON types unsupported by emulator: %v", err)
	}
	snap, err := doc.Get(h.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	data := snap.Data()
	if data["oid"] == nil || data["i32"] == nil {
		t.Fatalf("BSON round-trip: %#v", data)
	}
}

func TestIntegration_SerializeDeserialize(t *testing.T) {
	h := fstest.NewHarness(t)
	coll := h.Client.Collection(h.Coll("ser"))
	if _, err := coll.Doc("s1").Set(h.Ctx, map[string]any{"n": 1}); err != nil {
		t.Fatal(err)
	}
	q := coll.Where("n", "==", 1)
	b, err := q.Serialize()
	if err != nil {
		t.Fatalf("Serialize: %v", err)
	}
	q2, err := q.Deserialize(b)
	if err != nil {
		t.Fatalf("Deserialize: %v", err)
	}
	all, err := q2.Documents(h.Ctx).GetAll()
	if err != nil || len(all) != 1 {
		t.Fatalf("after deserialize: %d %v", len(all), err)
	}
}
