package testabilities

import (
	"errors"
)

// MockTopicManagerDocumentationProvider is a simple mock implementation for testing
type MockTopicManagerDocumentationProvider struct {
	ShouldFail bool
}

// GetDocumentationForTopicManager simulates a documentation retrieval operation
func (m *MockTopicManagerDocumentationProvider) GetDocumentationForTopicManager(topicManagerName string) (string, error) {
	if m.ShouldFail {
		return "", errors.New("documentation not found")
	}
	return "# Topic Manager Documentation\nThis is a test markdown document.", nil
}
