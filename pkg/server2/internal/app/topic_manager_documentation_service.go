package app

import (
	"context"
	"errors"
)

// TopicManagerDocumentationProvider defines the contract that must be fulfilled
// to retrieve documentation for a topic manager from the overlay engine.
type TopicManagerDocumentationProvider interface {
	GetDocumentationForTopicManager(topicManager string) (string, error)
}

// TopicManagerDocumentationService provides functionality for retrieving documentation
// for a specific topic manager.
type TopicManagerDocumentationService struct {
	provider TopicManagerDocumentationProvider
}

// ErrTopicManagerNotFound is returned when the requested topic manager documentation cannot be found.
var ErrTopicManagerNotFound = errors.New("topic manager documentation not found")

// ErrEmptyTopicManagerName is returned when an empty topic manager name is provided.
var ErrEmptyTopicManagerName = errors.New("topic manager name cannot be empty")

// GetDocumentation retrieves the documentation for a given topic manager.
// Returns an error if the topic manager name is empty or if the provider fails to retrieve the documentation.
func (s *TopicManagerDocumentationService) GetDocumentation(ctx context.Context, topicManager string) (string, error) {
	if topicManager == "" {
		return "", ErrEmptyTopicManagerName
	}

	documentation, err := s.provider.GetDocumentationForTopicManager(topicManager)
	if err != nil {
		return "", errors.Join(err, ErrTopicManagerNotFound)
	}

	return documentation, nil
}

// NewTopicManagerDocumentationService creates a new TopicManagerDocumentationService instance.
// It panics if the provider is nil.
func NewTopicManagerDocumentationService(provider TopicManagerDocumentationProvider) *TopicManagerDocumentationService {
	if provider == nil {
		panic("topic manager documentation provider cannot be nil")
	}

	return &TopicManagerDocumentationService{
		provider: provider,
	}
}
