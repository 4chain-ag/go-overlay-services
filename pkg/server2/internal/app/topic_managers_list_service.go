package app

import (
	"context"

	"github.com/bsv-blockchain/go-sdk/overlay"
)

// TopicManagersListProvider defines the contract for retrieving a list of topic managers.
type TopicManagersListProvider interface {
	ListTopicManagers() map[string]*overlay.MetaData
}

// TopicManagerMetadata represents the metadata for a topic manager.
type TopicManagerMetadata struct {
	Name             string  `json:"name"`
	ShortDescription string  `json:"shortDescription"`
	IconURL          *string `json:"iconURL,omitempty"`
	Version          *string `json:"version,omitempty"`
	InformationURL   *string `json:"informationURL,omitempty"`
}

// TopicManagersListResponse contains the response data for the topic managers list operation.
type TopicManagersListResponse map[string]TopicManagerMetadata

// TopicManagersListService provides functionality for retrieving a list of topic managers.
type TopicManagersListService struct {
	provider TopicManagersListProvider
}

// GetList retrieves the list of available topic managers.
func (s *TopicManagersListService) GetList(ctx context.Context) TopicManagersListResponse {
	engineTopicManagers := s.provider.ListTopicManagers()
	result := make(TopicManagersListResponse, len(engineTopicManagers))

	setIfNotEmpty := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

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
			topicManagerMetadata.IconURL = setIfNotEmpty(metadata.Icon)
			topicManagerMetadata.Version = setIfNotEmpty(metadata.Version)
			topicManagerMetadata.InformationURL = setIfNotEmpty(metadata.InfoUrl)
		}

		result[name] = topicManagerMetadata
	}

	return result
}

// NewTopicManagersListService creates a new TopicManagersListService with the given provider.
// It panics if the provider is nil.
func NewTopicManagersListService(provider TopicManagersListProvider) *TopicManagersListService {
	if provider == nil {
		panic("topic managers list provider cannot be nil")
	}

	return &TopicManagersListService{
		provider: provider,
	}
}
