package app

import (
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// LookupListProvider defines the interface for retrieving
// a list of lookup service providers from the overlay engine.
type LookupListProvider interface {
	ListLookupServiceProviders() map[string]*overlay.MetaData
}

// LookupListService provides operations for retrieving and formatting
// lookup service provider metadata from the overlay engine.
type LookupListService struct {
	provider LookupListProvider
}

// ListLookupServiceProviders retrieves the list of lookup service providers
// and formats them into a standardized response structure.
func (s *LookupListService) ListLookupServiceProviders() map[string]*overlay.MetaData {
	return s.provider.ListLookupServiceProviders()
}

// NewLookupListService creates a new LookupListService
// initialized with the given provider. It panics if the provider is nil.
func NewLookupListService(provider LookupListProvider) *LookupListService {
	if provider == nil {
		panic("lookup list provider is nil")
	}
	return &LookupListService{provider: provider}
}
