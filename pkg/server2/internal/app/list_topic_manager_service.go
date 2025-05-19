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
	Name             string  `json:"name"`
	ShortDescription string  `json:"shortDescription"`
	IconURL          *string `json:"iconURL,omitempty"`
	Version          *string `json:"version,omitempty"`
	InformationURL   *string `json:"informationURL,omitempty"`
}

// TopicManagersListResponse defines the response data structure for the topic managers list.
type TopicManagersListResponse map[string]TopicManagerMetadata

// TopicManagersListService provides operations for retrieving and formatting
// topic manager metadata from the overlay engine.
type TopicManagersListService struct {
	provider TopicManagersListProvider
}

// ListTopicManagers retrieves the list of topic managers
// and formats them into a standardized response structure.
func (s *TopicManagersListService) ListTopicManagers() TopicManagersListResponse {
	// Retrieve topic managers from the engine
	engineTopicManagers := s.provider.ListTopicManagers()

	// If nil is returned, provide an empty map
	if engineTopicManagers == nil {
		return make(TopicManagersListResponse)
	}

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

// NewTopicManagersListService creates a new TopicManagersListService
// initialized with the given provider. It returns an error if the provider is nil.
func NewTopicManagersListService(provider TopicManagersListProvider) (*TopicManagersListService, error) {
	if provider == nil {
		return nil, NewTopicManagerNilProviderError("topic manager list provider")
	}
	return &TopicManagersListService{provider: provider}, nil
}

// NewTopicManagerNilProviderError returns an Error indicating that a required topic manager provider was nil,
// which is invalid input when creating a topic manager service.
func NewTopicManagerNilProviderError(providerName string) Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       providerName + " cannot be nil",
		slug:      "The required provider was not properly initialized",
	}
}
