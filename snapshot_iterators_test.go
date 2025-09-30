package firestore

import (
	"testing"
)

func TestQuerySnapshotIteratorWrapper_InterfaceCompliance(t *testing.T) {
	t.Run("verify QuerySnapshotIterator interface compliance", func(t *testing.T) {
		var _ QuerySnapshotIterator = (*querySnapshotIteratorWrapper)(nil)
	})
}

func TestCollectionIteratorWrapper_InterfaceCompliance(t *testing.T) {
	t.Run("verify CollectionIterator interface compliance", func(t *testing.T) {
		var _ CollectionIterator = (*collectionIteratorWrapper)(nil)
	})
}

func TestDocumentSnapshotIteratorWrapper_InterfaceCompliance(t *testing.T) {
	t.Run("verify DocumentSnapshotIterator interface compliance", func(t *testing.T) {
		var _ DocumentSnapshotIterator = (*documentSnapshotIteratorWrapper)(nil)
	})
}

func TestQuerySnapshotIteratorWrapper_Next(t *testing.T) {
	t.Run("Next method exists", func(t *testing.T) {
		wrapper := &querySnapshotIteratorWrapper{iter: nil}
		_ = wrapper.Next
	})
}

func TestQuerySnapshotIteratorWrapper_Stop(t *testing.T) {
	t.Run("Stop method exists", func(t *testing.T) {
		wrapper := &querySnapshotIteratorWrapper{iter: nil}
		_ = wrapper.Stop
	})
}

func TestCollectionIteratorWrapper_Next(t *testing.T) {
	t.Run("Next method exists", func(t *testing.T) {
		wrapper := &collectionIteratorWrapper{iter: nil}
		_ = wrapper.Next
	})
}

func TestCollectionIteratorWrapper_Stop(t *testing.T) {
	t.Run("Stop method exists", func(t *testing.T) {
		wrapper := &collectionIteratorWrapper{iter: nil}
		_ = wrapper.Stop
	})
}

func TestDocumentSnapshotIteratorWrapper_Next(t *testing.T) {
	t.Run("Next method exists", func(t *testing.T) {
		wrapper := &documentSnapshotIteratorWrapper{iter: nil}
		_ = wrapper.Next
	})
}

func TestDocumentSnapshotIteratorWrapper_Stop(t *testing.T) {
	t.Run("Stop method exists", func(t *testing.T) {
		wrapper := &documentSnapshotIteratorWrapper{iter: nil}
		_ = wrapper.Stop
	})
}

func TestQuerySnapshotIteratorWrapper_AllMethods(t *testing.T) {
	t.Run("verify all QuerySnapshotIterator methods exist", func(t *testing.T) {
		wrapper := &querySnapshotIteratorWrapper{iter: nil}

		_ = wrapper.Next
		_ = wrapper.Stop
	})
}

func TestCollectionIteratorWrapper_AllMethods(t *testing.T) {
	t.Run("verify all CollectionIterator methods exist", func(t *testing.T) {
		wrapper := &collectionIteratorWrapper{iter: nil}

		_ = wrapper.Next
		_ = wrapper.Stop
	})
}

func TestDocumentSnapshotIteratorWrapper_AllMethods(t *testing.T) {
	t.Run("verify all DocumentSnapshotIterator methods exist", func(t *testing.T) {
		wrapper := &documentSnapshotIteratorWrapper{iter: nil}

		_ = wrapper.Next
		_ = wrapper.Stop
	})
}

func TestIteratorWrappers_MethodSignatures(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "QuerySnapshotIterator signatures",
			test: func(t *testing.T) {
				wrapper := &querySnapshotIteratorWrapper{iter: nil}
				_ = wrapper.Next
				_ = wrapper.Stop
			},
		},
		{
			name: "CollectionIterator signatures",
			test: func(t *testing.T) {
				wrapper := &collectionIteratorWrapper{iter: nil}
				_ = wrapper.Next
				_ = wrapper.Stop
			},
		},
		{
			name: "DocumentSnapshotIterator signatures",
			test: func(t *testing.T) {
				wrapper := &documentSnapshotIteratorWrapper{iter: nil}
				_ = wrapper.Next
				_ = wrapper.Stop
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
