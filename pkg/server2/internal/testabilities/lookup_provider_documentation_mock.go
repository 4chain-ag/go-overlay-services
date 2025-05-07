package testabilities

import (
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
)

// MockProviderDocumentationProvider is a simple mock implementation for testing
type MockProviderDocumentationProvider struct {
	ShouldFail bool
}

// GetDocumentationForLookupServiceProvider simulates a documentation retrieval operation
func (m *MockProviderDocumentationProvider) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	if m.ShouldFail {
		return "", errors.New("documentation not found")
	}
	return "# Test Documentation\nThis is a test markdown document.", nil
}

// MockProviderDocumentationProviderExpectations defines the expected behavior of the mock
type MockProviderDocumentationProviderExpectations struct {
	ShouldFail bool
}

// NewMockProviderDocumentationProviderMock creates a new mock provider with the specified expectations
func NewMockProviderDocumentationProviderMock(t *testing.T, expectations MockProviderDocumentationProviderExpectations) *MockProviderDocumentationProvider {
	return &MockProviderDocumentationProvider{
		ShouldFail: expectations.ShouldFail,
	}
}

// WithLookupProviderDocumentationProvider allows setting a custom LookupProviderDocumentationProvider in a TestOverlayEngineStub.
// This can be used to mock lookup service documentation behavior during tests.
func WithLookupProviderDocumentationProvider(provider app.LookupProviderDocumentationProvider) TestOverlayEngineStubOption {
	return func(stub *TestOverlayEngineStub) {
		stub.lookupServiceDocumentationProvider = provider
	}
}

// WithLookupProviderDocumentation configures the overlay engine stub to return successful documentation.
func WithLookupProviderDocumentation(doc string) TestOverlayEngineStubOption {
	return func(s *TestOverlayEngineStub) {
		s.lookupServiceDocumentationProvider = NewMockProviderDocumentationProviderMock(s.t, MockProviderDocumentationProviderExpectations{ShouldFail: false})
	}
}

// WithLookupProviderDocumentationError configures the overlay engine stub to return an error
// when attempting to retrieve lookup service documentation.
func WithLookupProviderDocumentationError() TestOverlayEngineStubOption {
	return func(s *TestOverlayEngineStub) {
		s.lookupServiceDocumentationProvider = NewMockProviderDocumentationProviderMock(s.t, MockProviderDocumentationProviderExpectations{ShouldFail: true})
	}
}

// lookupProviderDocumentationProviderAlwaysSuccessStub is a mock implementation of LookupProviderDocumentationProvider that always succeeds.
type lookupProviderDocumentationProviderAlwaysSuccessStub struct {
	documentation string
}

// GetDocumentationForLookupServiceProvider simulates a successful documentation retrieval.
func (s *lookupProviderDocumentationProviderAlwaysSuccessStub) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	return s.documentation, nil
}

// lookupProviderDocumentationProviderAlwaysFailureStub is a mock implementation of LookupProviderDocumentationProvider that always fails.
type lookupProviderDocumentationProviderAlwaysFailureStub struct{}

// GetDocumentationForLookupServiceProvider simulates a failed documentation retrieval.
func (s *lookupProviderDocumentationProviderAlwaysFailureStub) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	return "", errors.New("lookup service documentation error")
}
