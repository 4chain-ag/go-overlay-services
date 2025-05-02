package testabilities

import (
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// TopicManagerMetadataMock represents mock metadata for a topic manager.
type TopicManagerMetadataMock struct {
	Description string
	Icon        string
	Version     string
	InfoUrl     string
}

// WithTopicManagersList returns a TestOverlayEngineStubOption that configures the engine
// to return a specific list of topic managers.
func WithTopicManagersList(metadata map[string]*TopicManagerMetadataMock) TestOverlayEngineStubOption {
	return func(stub *TestOverlayEngineStub) {
		stub.listTopicManagersProvider = &mockTopicManagersListProvider{
			metadata: metadata,
		}
	}
}

// WithEmptyTopicManagersList returns a TestOverlayEngineStubOption that configures the engine
// to return an empty list of topic managers.
func WithEmptyTopicManagersList() TestOverlayEngineStubOption {
	return WithTopicManagersList(map[string]*TopicManagerMetadataMock{})
}

// mockTopicManagersListProvider is a mock implementation of the topic managers list provider.
type mockTopicManagersListProvider struct {
	metadata map[string]*TopicManagerMetadataMock
}

// ListTopicManagers returns a mock list of topic managers.
func (m *mockTopicManagersListProvider) ListTopicManagers() map[string]*overlay.MetaData {
	result := make(map[string]*overlay.MetaData, len(m.metadata))

	for name, mock := range m.metadata {
		if mock == nil {
			result[name] = nil
			continue
		}

		result[name] = &overlay.MetaData{
			Description: mock.Description,
			Icon:        mock.Icon,
			Version:     mock.Version,
			InfoUrl:     mock.InfoUrl,
		}
	}

	return result
}
