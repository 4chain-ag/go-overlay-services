package app

import (
	"context"
)

// TopicManagerDocumentationProvider defines the contract for retrieving documentation
// for a topic manager.
type TopicManagerDocumentationProvider interface {
	GetDocumentationForTopicManager(topicManager string) (string, error)
}

// TopicManagerDocumentationService provides functionality for retrieving topic manager documentation.
type TopicManagerDocumentationService struct {
	provider TopicManagerDocumentationProvider
}

// GetDocumentation retrieves documentation for a specific topic manager.
// Returns the documentation string on success, or an error if:
// - The topic manager name is empty (ErrorTypeIncorrectInput)
// - The provider fails to retrieve documentation (ErrorTypeProviderFailure)
func (s *TopicManagerDocumentationService) GetDocumentation(ctx context.Context, topicManager string) (string, error) {
	if topicManager == "" {
		return "", NewIncorrectInputError("topic manager name cannot be empty")
	}

	documentation, err := s.provider.GetDocumentationForTopicManager(topicManager)
	if err != nil {
		return "", NewProviderFailureError("unable to retrieve documentation for topic manager")
	}

	return documentation, nil
}

// NewTopicManagerDocumentationService creates a new TopicManagerDocumentationService with the given provider.
// Panics if the provider is nil.
func NewTopicManagerDocumentationService(provider TopicManagerDocumentationProvider) *TopicManagerDocumentationService {
	if provider == nil {
		panic("topic manager documentation provider cannot be nil")
	}

	return &TopicManagerDocumentationService{
		provider: provider,
	}
}
