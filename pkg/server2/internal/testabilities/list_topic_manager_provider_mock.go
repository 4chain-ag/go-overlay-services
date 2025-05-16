package testabilities

import (
	"testing"

	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

// TopicManagerListProviderMockExpectations defines the expected behavior of the TopicManagerListProviderMock during a test.

type TopicManagerListProviderMockExpectations struct {

	// MetadataList is the mock topic managers that will be returned.

	MetadataList map[string]*overlay.MetaData

	// ListTopicManagersCall indicates whether the ListTopicManagers method is expected to be called during the test.

	ListTopicManagersCall bool
}

// TopicManagerListProviderMock is a mock implementation of a topic manager list provider,

// used for testing the behavior of components that depend on topic manager listing.

type TopicManagerListProviderMock struct {
	t *testing.T

	// expectations defines the expected behavior and outcomes for this mock.

	expectations TopicManagerListProviderMockExpectations

	// called is true if the ListTopicManagers method was called.

	called bool
}

// ListTopicManagers returns the predefined list of topic managers.

func (m *TopicManagerListProviderMock) ListTopicManagers() map[string]*overlay.MetaData {

	m.t.Helper()

	m.called = true

	return m.expectations.MetadataList

}

// AssertCalled verifies that the ListTopicManagers method was called if it was expected to be.

func (m *TopicManagerListProviderMock) AssertCalled() {

	m.t.Helper()

	require.Equal(m.t, m.expectations.ListTopicManagersCall, m.called, "Discrepancy between expected and actual ListTopicManagers call")

}

// NewTopicManagerListProviderMock creates a new instance of TopicManagerListProviderMock with the given expectations.

func NewTopicManagerListProviderMock(t *testing.T, expectations TopicManagerListProviderMockExpectations) *TopicManagerListProviderMock {

	mock := &TopicManagerListProviderMock{

		t: t,

		expectations: expectations,
	}

	return mock

}

// TopicManagerListProviderAlwaysEmpty is a mock that always returns an empty topic manager list.

type TopicManagerListProviderAlwaysEmpty struct{}

// ListTopicManagers returns an empty map of topic managers.

func (*TopicManagerListProviderAlwaysEmpty) ListTopicManagers() map[string]*overlay.MetaData {

	return map[string]*overlay.MetaData{}

}

// TopicManagerListProviderAlwaysSuccess is a mock that always returns a predefined list of topic managers.

type TopicManagerListProviderAlwaysSuccess struct{}

// ListTopicManagers returns a predefined map of topic managers.

func (*TopicManagerListProviderAlwaysSuccess) ListTopicManagers() map[string]*overlay.MetaData {

	return map[string]*overlay.MetaData{

		"topic_manager1": {

			Description: "Description 1",

			Icon: "https://example.com/icon.png",

			Version: "1.0.0",

			InfoUrl: "https://example.com/info",
		},

		"topic_manager2": {

			Description: "Description 2",

			Icon: "https://example.com/icon2.png",

			Version: "2.0.0",

			InfoUrl: "https://example.com/info2",
		},
	}

}
