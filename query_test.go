package fsmock

import (
	"reflect"
	"testing"
	"time"
	"unsafe"

	"cloud.google.com/go/firestore"
)

func TestQueryChaining_CompileShape(t *testing.T) {
	var q Query = &queryWrapper{}
	q = q.Where("a", "==", 1).OrderBy("a", firestore.Asc).Limit(1)
	_ = q
}

func TestWithReadOptions_QueryDoesNotAlias(t *testing.T) {
	var seeded firestore.Query
	(&seeded).WithReadOptions(firestore.ReadTime(time.Unix(1, 0)))
	base := &queryWrapper{q: seeded}

	before := readSettingsTime(base.q)
	if before.IsZero() {
		t.Fatal("expected seeded readSettings")
	}

	derived := base.WithReadOptions(firestore.ReadTime(time.Unix(100, 0)))
	afterBase := readSettingsTime(base.q)
	afterDerived := readSettingsTime(derived.SDKQuery())

	if !before.Equal(afterBase) {
		t.Fatalf("base query readSettings mutated: before=%v after=%v", before, afterBase)
	}
	if afterDerived.Equal(afterBase) {
		t.Fatalf("derived should differ from base: derived=%v base=%v", afterDerived, afterBase)
	}
}

func TestWithReadOptions_CollectionRefDoesNotAlias(t *testing.T) {
	ref := &firestore.CollectionRef{}
	// Seed via SDK method after attaching a fresh readSettings through a Query seed trick:
	// CollectionRef.WithReadOptions requires non-nil readSettings (SDK applies in place).
	var q firestore.Query
	(&q).WithReadOptions(firestore.ReadTime(time.Unix(1, 0)))
	rs := queryReadSettingsPointer(q)
	setCollectionReadSettings(ref, rs)

	w := &collectionRefWrapper{
		queryWrapper: queryWrapper{q: ref.Query},
		ref:          ref,
	}

	before := collectionReadSettingsTime(ref)
	_ = w.WithReadOptions(firestore.ReadTime(time.Unix(100, 0)))
	after := collectionReadSettingsTime(ref)
	if !before.Equal(after) {
		t.Fatalf("original CollectionRef readSettings mutated: before=%v after=%v", before, after)
	}
}

func readSettingsTime(q firestore.Query) time.Time {
	f := reflect.ValueOf(&q).Elem().FieldByName("readSettings")
	if !f.IsValid() || f.IsNil() {
		return time.Time{}
	}
	return timeFromReadSettings(f)
}

func collectionReadSettingsTime(ref *firestore.CollectionRef) time.Time {
	f := reflect.ValueOf(ref).Elem().FieldByName("readSettings")
	if !f.IsValid() || f.IsNil() {
		return time.Time{}
	}
	return timeFromReadSettings(f)
}

func queryReadSettingsPointer(q firestore.Query) reflect.Value {
	return reflect.ValueOf(&q).Elem().FieldByName("readSettings")
}

func setCollectionReadSettings(ref *firestore.CollectionRef, rs reflect.Value) {
	rv := reflect.ValueOf(ref).Elem()
	setUnexported(rv.FieldByName("readSettings"), rs)
	setUnexported(rv.FieldByName("Query").FieldByName("readSettings"), rs)
}

func timeFromReadSettings(f reflect.Value) time.Time {
	rt := f.Elem().FieldByName("readTime")
	return *(*time.Time)(unsafe.Pointer(rt.UnsafeAddr()))
}
