package testabilities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TopicManagerDocumentationProviderMockExpectations struct {
	DocumentationCall bool
	Error             error
	Documentation     string
}

var DefaultTopicManagerDocumentationProviderMockExpectations = TopicManagerDocumentationProviderMockExpectations{
	DocumentationCall: true,
	Error:             nil,
	Documentation:     "# Topic Manager Documentation\nThis is a test markdown document.",
}

// TopicManagerDocumentationProviderMock is a simple mock implementation for testing
type TopicManagerDocumentationProviderMock struct {
	t            *testing.T
	expectations TopicManagerDocumentationProviderMockExpectations
	called       bool
}

// GetDocumentationForTopicManager simulates a documentation retrieval operation
func (m *TopicManagerDocumentationProviderMock) GetDocumentationForTopicManager(topicManagerName string) (string, error) {
	m.t.Helper()

	m.called = true

	if m.expectations.Error != nil {
		return "", m.expectations.Error
	}

	return m.expectations.Documentation, nil
}

func (m *TopicManagerDocumentationProviderMock) AssertCalled() {
	m.t.Helper()
	require.Equal(m.t, m.expectations.DocumentationCall, m.called, "Discrepancy between expected and actual DocumentationCall")
}

func NewTopicManagerDocumentationProviderMock(t *testing.T, expectations TopicManagerDocumentationProviderMockExpectations) *TopicManagerDocumentationProviderMock {
	return &TopicManagerDocumentationProviderMock{
		t:            t,
		expectations: expectations,
	}
}
