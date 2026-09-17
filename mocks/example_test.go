package mocks_test

import (
	"testing"

	"github.com/akmalsyrf/go-firestore-mock/v2/mocks"
	"go.uber.org/mock/gomock"
)

func TestMockClient_CollectionExpectation(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockClient(ctrl)
	mockColl := mocks.NewMockCollectionRef(ctrl)

	mockClient.EXPECT().Collection("users").Return(mockColl)

	if got := mockClient.Collection("users"); got != mockColl {
		t.Fatalf("got %v", got)
	}
}
