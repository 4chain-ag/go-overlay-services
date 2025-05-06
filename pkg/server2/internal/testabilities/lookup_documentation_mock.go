package testabilities

import (
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
)

// LookupServiceDocumentationProviderMock is a mock implementation of the lookup service documentation provider
// for testing purposes.
type LookupServiceDocumentationProviderMock struct {
	t             *testing.T
	documentation string
	err           error
}

// GetDocumentationForLookupServiceProvider returns the configured documentation or error for testing.
func (m *LookupServiceDocumentationProviderMock) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	return m.documentation, m.err
}

// WithLookupServiceDocumentationProvider allows setting a custom LookupServiceDocumentationProvider in a TestOverlayEngineStub.
// This can be used to mock lookup service documentation behavior during tests.
func WithLookupServiceDocumentationProvider(provider app.LookupServiceDocumentationProvider) TestOverlayEngineStubOption {
	return func(stub *TestOverlayEngineStub) {
		stub.lookupServiceDocumentationProvider = provider
	}
}

// WithLookupServiceDocumentation configures the overlay engine stub to return successful documentation.
func WithLookupServiceDocumentation(doc string) TestOverlayEngineStubOption {
	return func(s *TestOverlayEngineStub) {
		s.lookupServiceDocumentationProvider = &lookupServiceDocumentationProviderAlwaysSuccessStub{documentation: doc}
	}
}

// WithLookupServiceDocumentationError configures the overlay engine stub to return an error
// when attempting to retrieve lookup service documentation.
func WithLookupServiceDocumentationError() TestOverlayEngineStubOption {
	return func(s *TestOverlayEngineStub) {
		s.lookupServiceDocumentationProvider = &lookupServiceDocumentationProviderAlwaysFailureStub{}
	}
}

// lookupServiceDocumentationProviderAlwaysSuccessStub is a mock implementation of LookupServiceDocumentationProvider that always succeeds.
// It is used as the default LookupServiceDocumentationProvider in the TestOverlayEngineStub.
type lookupServiceDocumentationProviderAlwaysSuccessStub struct {
	documentation string
}

// GetDocumentationForLookupServiceProvider simulates a successful documentation retrieval.
func (s *lookupServiceDocumentationProviderAlwaysSuccessStub) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	return s.documentation, nil
}

// lookupServiceDocumentationProviderAlwaysFailureStub is a mock implementation of LookupServiceDocumentationProvider that always fails.
type lookupServiceDocumentationProviderAlwaysFailureStub struct{}

// GetDocumentationForLookupServiceProvider simulates a failed documentation retrieval.
func (s *lookupServiceDocumentationProviderAlwaysFailureStub) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	return "", errors.New("lookup service documentation error")
}

// NewLookupServiceDocumentationProviderMock creates a new mock provider for lookup service documentation.
func NewLookupServiceDocumentationProviderMock(t *testing.T, documentation string, err error) *LookupServiceDocumentationProviderMock {
	return &LookupServiceDocumentationProviderMock{
		t:             t,
		documentation: documentation,
		err:           err,
	}
}
