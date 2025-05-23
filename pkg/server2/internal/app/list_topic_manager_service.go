package app

import (
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// TopicManagersListProvider defines the interface for retrieving
// a list of topic managers from the overlay engine.
type TopicManagersListProvider interface {
	ListTopicManagers() map[string]*overlay.MetaData
}

// TopicManagerMetadata represents the metadata for a topic manager.
type TopicManagerMetadata struct {
	Name             string
	ShortDescription string
	IconURL          string
	Version          string
	InformationURL   string
}

// TopicManagersListService provides operations for retrieving and formatting
// topic manager metadata from the overlay engine.
type TopicManagersListService struct {
	provider TopicManagersListProvider
}

type TopicManagers map[string]TopicManagerMetadata

// ListTopicManagers retrieves the list of topic managers
// and formats them into a standardized response structure.
func (s *TopicManagersListService) ListTopicManagers() TopicManagers {
	engineTopicManagers := s.provider.ListTopicManagers()
	if engineTopicManagers == nil {
		return make(TopicManagers)
	}

	result := make(TopicManagers, len(engineTopicManagers))
	coalesce := func(primary, fallback string) string {
		if primary != "" {
			return primary
		}
		return fallback
	}

	for name, metadata := range engineTopicManagers {
		topicManagerMetadata := TopicManagerMetadata{
			Name:             name,
			ShortDescription: "No description available",
		}

		if metadata != nil {
			topicManagerMetadata.ShortDescription = coalesce(metadata.Description, "No description available")
			topicManagerMetadata.IconURL = metadata.Icon
			topicManagerMetadata.Version = metadata.Version
			topicManagerMetadata.InformationURL = metadata.InfoUrl
		}

		result[name] = topicManagerMetadata
	}

	return result
}

// NewTopicManagersListService creates a new TopicManagersListService
// initialized with the given provider. It panics if the provider is nil.
func NewTopicManagersListService(provider TopicManagersListProvider) (*TopicManagersListService) {
	if provider == nil {
		panic("topic manager list provider is nil")
	}
	return &TopicManagersListService{provider: provider}
}
