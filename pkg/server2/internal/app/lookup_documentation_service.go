package app

import (
	"context"
)

// LookupServiceProviderDocumentationProvider defines the contract for retrieving documentation
// for a lookup service provider.
type LookupServiceProviderDocumentationProvider interface {
	GetDocumentationForLookupServiceProvider(lookupServiceName string) (string, error)
}

// LookupServiceProviderDocumentationService provides functionality for retrieving lookup service provider documentation.
type LookupServiceProviderDocumentationService struct {
	provider LookupServiceProviderDocumentationProvider
}

// GetDocumentation retrieves documentation for a specific lookup service provider.
// Returns the documentation string on success, or an error if:
// - The lookup service name is empty (ErrorTypeIncorrectInput)
// - The provider fails to retrieve documentation (ErrorTypeProviderFailure)
func (s *LookupServiceProviderDocumentationService) GetDocumentation(ctx context.Context, lookupServiceName string) (string, error) {
	if lookupServiceName == "" {
		return "", NewEmptyLookupServiceNameError()
	}

	documentation, err := s.provider.GetDocumentationForLookupServiceProvider(lookupServiceName)
	if err != nil {
		return "", NewLookupServiceProviderDocumentationError(err)
	}

	return documentation, nil
}

// NewLookupServiceProviderDocumentationService creates a new LookupServiceProviderDocumentationService with the given provider.
// Panics if the provider is nil.
func NewLookupServiceProviderDocumentationService(provider LookupServiceProviderDocumentationProvider) *LookupServiceProviderDocumentationService {
	if provider == nil {
		panic("lookup service provider documentation provider cannot be nil")
	}

	return &LookupServiceProviderDocumentationService{
		provider: provider,
	}
}

// NewEmptyLookupServiceNameError returns an Error indicating that the lookup service name is empty,
// which is invalid input when retrieving documentation.
func NewEmptyLookupServiceNameError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "lookup service name cannot be empty",
		slug:      "A valid lookup service name must be provided to retrieve documentation.",
	}
}

// NewLookupServiceProviderDocumentationError returns an Error indicating that the configured provider
// failed to retrieve documentation for the lookup service.
func NewLookupServiceProviderDocumentationError(err error) Error {
	return Error{
		errorType: ErrorTypeProviderFailure,
		err:       "unable to retrieve documentation for lookup service provider",
		slug:      "Unable to retrieve documentation for lookup service provider due to an internal error. Please try again later or contact the support team.",
	}
}
