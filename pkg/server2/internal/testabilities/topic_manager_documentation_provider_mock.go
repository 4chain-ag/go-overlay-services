package testabilities

import (
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
)

// TopicManagerDocumentationProviderMock is a mock implementation of the topic manager documentation provider
// for testing purposes.
type TopicManagerDocumentationProviderMock struct {
	t             *testing.T
	documentation string
	err           error
}

// GetDocumentationForTopicManager returns the configured documentation or error for testing.
func (m *TopicManagerDocumentationProviderMock) GetDocumentationForTopicManager(topicManager string) (string, error) {
	return m.documentation, m.err
}

// WithTopicManagerDocumentationProvider allows setting a custom TopicManagerDocumentationProvider in a TestOverlayEngineStub.
// This can be used to mock topic manager documentation behavior during tests.
func WithTopicManagerDocumentationProvider(provider app.TopicManagerDocumentationProvider) TestOverlayEngineStubOption {
	return func(stub *TestOverlayEngineStub) {
		stub.topicManagerDocumentationProvider = provider
	}
}

// WithTopicManagerDocumentation configures the overlay engine stub to return successful documentation.
func WithTopicManagerDocumentation(doc string) TestOverlayEngineStubOption {
	return func(s *TestOverlayEngineStub) {
		s.topicManagerDocumentationProvider = &topicManagerDocumentationProviderAlwaysSuccessStub{documentation: doc}
	}
}

// WithTopicManagerDocumentationError configures the overlay engine stub to return an error
// when attempting to retrieve topic manager documentation.
func WithTopicManagerDocumentationError() TestOverlayEngineStubOption {
	return func(s *TestOverlayEngineStub) {
		s.topicManagerDocumentationProvider = &topicManagerDocumentationProviderAlwaysFailureStub{}
	}
}

// topicManagerDocumentationProviderAlwaysSuccessStub is a mock implementation of TopicManagerDocumentationProvider that always succeeds.
// It is used as the default TopicManagerDocumentationProvider in the TestOverlayEngineStub.
type topicManagerDocumentationProviderAlwaysSuccessStub struct {
	documentation string
}

// GetDocumentationForTopicManager simulates a successful documentation retrieval.
func (s *topicManagerDocumentationProviderAlwaysSuccessStub) GetDocumentationForTopicManager(topicManager string) (string, error) {
	return s.documentation, nil
}

// topicManagerDocumentationProviderAlwaysFailureStub is a mock implementation of TopicManagerDocumentationProvider that always fails.
type topicManagerDocumentationProviderAlwaysFailureStub struct{}

// GetDocumentationForTopicManager simulates a failed documentation retrieval.
func (s *topicManagerDocumentationProviderAlwaysFailureStub) GetDocumentationForTopicManager(topicManager string) (string, error) {
	return "", errors.New("topic manager documentation error")
}

// NewTopicManagerDocumentationProviderMock creates a new mock provider for topic manager documentation.
func NewTopicManagerDocumentationProviderMock(t *testing.T, documentation string, err error) *TopicManagerDocumentationProviderMock {
	return &TopicManagerDocumentationProviderMock{
		t:             t,
		documentation: documentation,
		err:           err,
	}
}
