package app

import (
	"context"
)

// LookupProviderDocumentationProvider defines the contract for retrieving documentation
// for a lookup service provider.
type LookupProviderDocumentationProvider interface {
	GetDocumentationForLookupServiceProvider(lookupService string) (string, error)
}

// LookupProviderDocumentationService provides functionality for retrieving lookup service documentation.
type LookupProviderDocumentationService struct {
	provider LookupProviderDocumentationProvider
}

// GetDocumentation retrieves documentation for a specific lookup service.
// Returns the documentation string on success, or an error if:
// - The lookup service name is empty (ErrorTypeIncorrectInput)
// - The provider fails to retrieve documentation (ErrorTypeProviderFailure)
func (s *LookupProviderDocumentationService) GetDocumentation(ctx context.Context, lookupService string) (string, error) {
	if lookupService == "" {
		return "", NewIncorrectInputError("lookup service name cannot be empty")
	}

	documentation, err := s.provider.GetDocumentationForLookupServiceProvider(lookupService)
	if err != nil {
		return "", NewProviderFailureError("unable to retrieve documentation for lookup service provider")
	}

	return documentation, nil
}

// NewLookupProviderDocumentationService creates a new LookupProviderDocumentationService with the given provider.
// Panics if the provider is nil.
func NewLookupProviderDocumentationService(provider LookupProviderDocumentationProvider) *LookupProviderDocumentationService {
	if provider == nil {
		panic("lookup documentation provider cannot be nil")
	}

	return &LookupProviderDocumentationService{
		provider: provider,
	}
}
