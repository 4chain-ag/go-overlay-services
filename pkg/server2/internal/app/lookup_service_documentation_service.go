package app

import (
	"context"
	"errors"
)

// LookupServiceDocumentationProvider defines the contract that must be fulfilled
// to retrieve documentation for a lookup service provider from the overlay engine.
type LookupServiceDocumentationProvider interface {
	GetDocumentationForLookupServiceProvider(lookupService string) (string, error)
}

// LookupServiceDocumentationService provides functionality for retrieving documentation
// for a specific lookup service provider.
type LookupServiceDocumentationService struct {
	provider LookupServiceDocumentationProvider
}

// ErrLookupServiceProviderNotFound is returned when the requested lookup service provider documentation cannot be found.
var ErrLookupServiceProviderNotFound = errors.New("lookup service provider documentation not found")

// ErrEmptyLookupServiceName is returned when an empty lookup service name is provided.
var ErrEmptyLookupServiceName = errors.New("lookup service name cannot be empty")

// GetDocumentation retrieves the documentation for a given lookup service provider.
// Returns an error if the lookup service name is empty or if the provider fails to retrieve the documentation.
func (s *LookupServiceDocumentationService) GetDocumentation(ctx context.Context, lookupService string) (string, error) {
	if lookupService == "" {
		return "", ErrEmptyLookupServiceName
	}

	documentation, err := s.provider.GetDocumentationForLookupServiceProvider(lookupService)
	if err != nil {
		return "", errors.Join(err, ErrLookupServiceProviderNotFound)
	}

	return documentation, nil
}

// NewLookupServiceDocumentationService creates a new LookupServiceDocumentationService instance.
// It panics if the provider is nil.
func NewLookupServiceDocumentationService(provider LookupServiceDocumentationProvider) *LookupServiceDocumentationService {
	if provider == nil {
		panic("lookup service documentation provider cannot be nil")
	}

	return &LookupServiceDocumentationService{
		provider: provider,
	}
}
