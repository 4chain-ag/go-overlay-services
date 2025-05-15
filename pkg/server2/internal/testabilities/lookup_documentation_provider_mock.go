package testabilities

import (
	"errors"
)

// MockLookupServiceProviderDocumentationProvider is a simple mock implementation for testing
type MockLookupServiceProviderDocumentationProvider struct {
	ShouldFail bool
}

// GetDocumentationForLookupServiceProvider simulates a documentation retrieval operation
func (m *MockLookupServiceProviderDocumentationProvider) GetDocumentationForLookupServiceProvider(lookupServiceName string) (string, error) {
	if m.ShouldFail {
		return "", errors.New("documentation not found")
	}
	return "# Test Documentation\nThis is a test markdown document.", nil
}
