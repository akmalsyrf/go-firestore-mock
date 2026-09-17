package fsmock

import (
	"reflect"
	"unsafe"

	"cloud.google.com/go/firestore"
)

// cloneQueryWithFreshReadSettings returns a shallow copy of q whose readSettings
// pointer is a deep copy of the original (or nil). This prevents WithReadOptions
// from mutating the receiver's shared readSettings.
func cloneQueryWithFreshReadSettings(q firestore.Query) firestore.Query {
	cloneReadSettingsField(reflect.ValueOf(&q).Elem(), "readSettings")
	return q
}

// cloneCollectionRefWithFreshReadSettings returns a new *CollectionRef that does
// not share readSettings with the original. CollectionRef and its embedded Query
// share the same readSettings pointer in the SDK; both are replaced.
func cloneCollectionRefWithFreshReadSettings(ref *firestore.CollectionRef) *firestore.CollectionRef {
	if ref == nil {
		return nil
	}
	clone := *ref
	rv := reflect.ValueOf(&clone).Elem()
	newRS := cloneReadSettingsField(rv, "readSettings")
	// Keep embedded Query.readSettings in sync (SDK shares one pointer).
	qField := rv.FieldByName("Query")
	if qField.IsValid() {
		setUnexported(qField.FieldByName("readSettings"), newRS)
	}
	return &clone
}

// cloneDocumentRefWithFreshReadSettings returns a new *DocumentRef that does not
// share readSettings with the original.
func cloneDocumentRefWithFreshReadSettings(ref *firestore.DocumentRef) *firestore.DocumentRef {
	if ref == nil {
		return nil
	}
	clone := *ref
	cloneReadSettingsField(reflect.ValueOf(&clone).Elem(), "readSettings")
	return &clone
}

// cloneReadSettingsField deep-copies the named *readSettings field on structVal
// (an addressable struct Value) and returns the new pointer Value (may be invalid if nil).
func cloneReadSettingsField(structVal reflect.Value, field string) reflect.Value {
	f := structVal.FieldByName(field)
	if !f.IsValid() || f.Kind() != reflect.Pointer {
		return reflect.Value{}
	}
	if f.IsNil() {
		return reflect.Value{}
	}
	src := f.Elem()
	dst := reflect.New(src.Type())
	// Byte-copy: src Value is flagged unexported, so Value.Set panics.
	typedmemmove(dst.UnsafePointer(), unsafe.Pointer(src.UnsafeAddr()), src.Type().Size())
	setUnexported(f, dst)
	return dst
}

func typedmemmove(dst, src unsafe.Pointer, size uintptr) {
	if size == 0 || dst == src {
		return
	}
	srcBytes := unsafe.Slice((*byte)(src), size)
	dstBytes := unsafe.Slice((*byte)(dst), size)
	copy(dstBytes, srcBytes)
}

func setUnexported(field, value reflect.Value) {
	if !field.IsValid() {
		return
	}
	addr := unsafe.Pointer(field.UnsafeAddr())
	if !value.IsValid() || (value.Kind() == reflect.Pointer && value.IsNil()) {
		*(*unsafe.Pointer)(addr) = nil
		return
	}
	if field.Kind() == reflect.Pointer {
		// Avoid Value.Set: value may be flagged unexported (foreign package field).
		*(*unsafe.Pointer)(addr) = value.UnsafePointer()
		return
	}
	target := reflect.NewAt(field.Type(), addr).Elem()
	target.Set(value)
}
