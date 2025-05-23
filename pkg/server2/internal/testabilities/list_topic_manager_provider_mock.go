package testabilities

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

// Standard metadata maps that can be used for testing
var (
	// TopicManagerEmptyMetadata is an empty metadata map
	TopicManagerEmptyMetadata = map[string]*overlay.MetaData{}

	// TopicManagerDefaultMetadata contains standard metadata for testing
	TopicManagerDefaultMetadata = map[string]*overlay.MetaData{
		"topic_manager1": {
			Description: "Description 1",
			Icon:        "https://example.com/icon.png",
			Version:     "1.0.0",
			InfoUrl:     "https://example.com/info",
		},
		"topic_manager2": {
			Description: "Description 2",
			Icon:        "https://example.com/icon2.png",
			Version:     "2.0.0",
			InfoUrl:     "https://example.com/info2",
		},
	}
)

// Standard expected responses that can be used for testing
var (
	// TopicManagerEmptyExpectedResponse is an empty response
	TopicManagerEmptyExpectedResponse = app.TopicManagers{}

	// TopicManagerDefaultExpectedResponse contains the standard expected response matching TopicManagerDefaultMetadata
	TopicManagerDefaultExpectedResponse = app.TopicManagers{
		"topic_manager1": app.TopicManagerMetadata{
			Name:             "topic_manager1",
			ShortDescription: "Description 1",
			IconURL:          "https://example.com/icon.png",
			Version:          "1.0.0",
			InformationURL:   "https://example.com/info",
		},
		"topic_manager2": app.TopicManagerMetadata{
			Name:             "topic_manager2",
			ShortDescription: "Description 2",
			IconURL:          "https://example.com/icon2.png",
			Version:          "2.0.0",
			InformationURL:   "https://example.com/info2",
		},
	}
)

// TopicManagersListProviderMockExpectations defines the expected behavior of the TopicManagersListProviderMock during a test.
type TopicManagersListProviderMockExpectations struct {
	MetadataList          map[string]*overlay.MetaData
	Error                 error
	ListTopicManagersCall bool
}

// TopicManagersListProviderMock is a mock implementation of a topic manager list provider,
// used for testing the behavior of components that depend on topic manager listing.
type TopicManagersListProviderMock struct {
	t            *testing.T
	expectations TopicManagersListProviderMockExpectations
	called       bool
}

// ListTopicManagers returns the predefined list of topic managers.
func (m *TopicManagersListProviderMock) ListTopicManagers() map[string]*overlay.MetaData {
	m.t.Helper()
	m.called = true
	return m.expectations.MetadataList
}

// AssertCalled verifies that the ListTopicManagers method was called if it was expected to be.
func (m *TopicManagersListProviderMock) AssertCalled() {
	m.t.Helper()
	require.Equal(m.t, m.expectations.ListTopicManagersCall, m.called, "Discrepancy between expected and actual ListTopicManagers call")
}

// NewTopicManagersListProviderMock creates a new instance of TopicManagersListProviderMock with the given expectations.
func NewTopicManagersListProviderMock(t *testing.T, expectations TopicManagersListProviderMockExpectations) *TopicManagersListProviderMock {
	return &TopicManagersListProviderMock{
		t:            t,
		expectations: expectations,
	}
}
