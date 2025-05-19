package testabilities

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

// Standard metadata maps that can be used for testing
var (
	// EmptyMetadata is an empty metadata map
	EmptyMetadata = map[string]*overlay.MetaData{}

	// DefaultMetadata contains standard metadata for testing
	DefaultMetadata = map[string]*overlay.MetaData{
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
	// EmptyExpectedResponse is an empty response
	EmptyExpectedResponse = app.TopicManagersListResponse{}

	// DefaultExpectedResponse contains the standard expected response matching DefaultMetadata
	DefaultExpectedResponse = app.TopicManagersListResponse{
		"topic_manager1": app.TopicManagerMetadata{
			Name:             "topic_manager1",
			ShortDescription: "Description 1",
			IconURL:          ptr.To("https://example.com/icon.png"),
			Version:          ptr.To("1.0.0"),
			InformationURL:   ptr.To("https://example.com/info"),
		},
		"topic_manager2": app.TopicManagerMetadata{
			Name:             "topic_manager2",
			ShortDescription: "Description 2",
			IconURL:          ptr.To("https://example.com/icon2.png"),
			Version:          ptr.To("2.0.0"),
			InformationURL:   ptr.To("https://example.com/info2"),
		},
	}
)

// TopicManagersListProviderMockExpectations defines the expected behavior of the TopicManagersListProviderMock during a test.
type TopicManagersListProviderMockExpectations struct {
	// MetadataList is the mock topic managers that will be returned.
	MetadataList map[string]*overlay.MetaData

	// ListTopicManagersCall indicates whether the ListTopicManagers method is expected to be called during the test.
	ListTopicManagersCall bool
}

// TopicManagersListProviderMock is a mock implementation of a topic manager list provider,
// used for testing the behavior of components that depend on topic manager listing.
type TopicManagersListProviderMock struct {
	t *testing.T

	// expectations defines the expected behavior and outcomes for this mock.
	expectations TopicManagersListProviderMockExpectations

	// called is true if the ListTopicManagers method was called.
	called bool
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

// TopicManagersListProviderAlwaysEmpty is a mock that always returns an empty topic manager list.
type TopicManagersListProviderAlwaysEmpty struct{}

// ListTopicManagers returns an empty map of topic managers.
func (*TopicManagersListProviderAlwaysEmpty) ListTopicManagers() map[string]*overlay.MetaData {
	return EmptyMetadata
}

// TopicManagersListProviderAlwaysSuccess is a mock that always returns a predefined list of topic managers.
type TopicManagersListProviderAlwaysSuccess struct{}

// ListTopicManagers returns a predefined map of topic managers.
func (*TopicManagersListProviderAlwaysSuccess) ListTopicManagers() map[string]*overlay.MetaData {
	return DefaultMetadata
}
