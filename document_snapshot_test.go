package firestore

import (
	"testing"
)

func TestDocumentSnapshotWrapper_InterfaceCompliance(t *testing.T) {
	t.Run("verify DocumentSnapshot interface compliance", func(t *testing.T) {
		var _ DocumentSnapshot = (*documentSnapshotWrapper)(nil)
	})
}

func TestDocumentSnapshotWrapper_Data(t *testing.T) {
	t.Run("Data method exists", func(t *testing.T) {
		wrapper := &documentSnapshotWrapper{snap: nil}
		_ = wrapper.Data
	})
}

func TestDocumentSnapshotWrapper_DataTo(t *testing.T) {
	t.Run("DataTo method exists", func(t *testing.T) {
		wrapper := &documentSnapshotWrapper{snap: nil}
		_ = wrapper.DataTo
	})
}

func TestDocumentSnapshotWrapper_DataAt(t *testing.T) {
	t.Run("DataAt method exists", func(t *testing.T) {
		wrapper := &documentSnapshotWrapper{snap: nil}
		_ = wrapper.DataAt
	})
}

func TestDocumentSnapshotWrapper_Exists(t *testing.T) {
	t.Run("Exists method exists", func(t *testing.T) {
		wrapper := &documentSnapshotWrapper{snap: nil}
		_ = wrapper.Exists
	})
}

func TestDocumentSnapshotWrapper_CreateTime(t *testing.T) {
	t.Run("CreateTime method exists", func(t *testing.T) {
		wrapper := &documentSnapshotWrapper{snap: nil}
		_ = wrapper.CreateTime
	})
}

func TestDocumentSnapshotWrapper_UpdateTime(t *testing.T) {
	t.Run("UpdateTime method exists", func(t *testing.T) {
		wrapper := &documentSnapshotWrapper{snap: nil}
		_ = wrapper.UpdateTime
	})
}

func TestDocumentSnapshotWrapper_ReadTime(t *testing.T) {
	t.Run("ReadTime method exists", func(t *testing.T) {
		wrapper := &documentSnapshotWrapper{snap: nil}
		_ = wrapper.ReadTime
	})
}

func TestDocumentSnapshotWrapper_Ref(t *testing.T) {
	t.Run("Ref method exists", func(t *testing.T) {
		wrapper := &documentSnapshotWrapper{snap: nil}
		_ = wrapper.Ref
	})
}

func TestDocumentSnapshotWrapper_AllMethods(t *testing.T) {
	t.Run("verify all DocumentSnapshot methods exist", func(t *testing.T) {
		wrapper := &documentSnapshotWrapper{snap: nil}

		_ = wrapper.Data
		_ = wrapper.DataTo
		_ = wrapper.DataAt
		_ = wrapper.Exists
		_ = wrapper.CreateTime
		_ = wrapper.UpdateTime
		_ = wrapper.ReadTime
		_ = wrapper.Ref
	})
}

func TestDocumentSnapshotWrapper_MethodSignatures(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "Data signature",
			test: func(t *testing.T) {
				wrapper := &documentSnapshotWrapper{snap: nil}
				_ = wrapper.Data
			},
		},
		{
			name: "DataTo signature",
			test: func(t *testing.T) {
				wrapper := &documentSnapshotWrapper{snap: nil}
				_ = wrapper.DataTo
			},
		},
		{
			name: "DataAt signature",
			test: func(t *testing.T) {
				wrapper := &documentSnapshotWrapper{snap: nil}
				_ = wrapper.DataAt
			},
		},
		{
			name: "Exists signature",
			test: func(t *testing.T) {
				wrapper := &documentSnapshotWrapper{snap: nil}
				_ = wrapper.Exists
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
