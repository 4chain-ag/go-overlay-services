package testabilities

import (
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
)

// MockTopicManagerDocumentationProvider is a simple mock implementation for testing
type MockTopicManagerDocumentationProvider struct {
	ShouldFail bool
}

// GetDocumentationForTopicManager simulates a documentation retrieval operation
func (m *MockTopicManagerDocumentationProvider) GetDocumentationForTopicManager(topicManager string) (string, error) {
	if m.ShouldFail {
		return "", errors.New("documentation not found")
	}
	return "# Test Documentation\nThis is a test markdown document.", nil
}

// MockTopicManagerDocumentationProviderExpectations defines the expected behavior of the mock
type MockTopicManagerDocumentationProviderExpectations struct {
	ShouldFail bool
}

// NewMockTopicManagerDocumentationProviderMock creates a new mock provider with the specified expectations
func NewMockTopicManagerDocumentationProviderMock(t *testing.T, expectations MockTopicManagerDocumentationProviderExpectations) *MockTopicManagerDocumentationProvider {
	return &MockTopicManagerDocumentationProvider{
		ShouldFail: expectations.ShouldFail,
	}
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
		s.topicManagerDocumentationProvider = NewMockTopicManagerDocumentationProviderMock(s.t, MockTopicManagerDocumentationProviderExpectations{ShouldFail: false})
	}
}

// WithTopicManagerDocumentationError configures the overlay engine stub to return an error
// when attempting to retrieve topic manager documentation.
func WithTopicManagerDocumentationError() TestOverlayEngineStubOption {
	return func(s *TestOverlayEngineStub) {
		s.topicManagerDocumentationProvider = NewMockTopicManagerDocumentationProviderMock(s.t, MockTopicManagerDocumentationProviderExpectations{ShouldFail: true})
	}
}

// topicManagerDocumentationProviderAlwaysSuccessStub is a mock implementation of TopicManagerDocumentationProvider that always succeeds.
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
