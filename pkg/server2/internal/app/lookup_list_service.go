package app

import (
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// LookupListProvider defines the interface for retrieving
// a list of lookup service providers from the overlay engine.
type LookupListProvider interface {
	ListLookupServiceProviders() map[string]*overlay.MetaData
}

// LookupMetadata represents the metadata for a lookup service provider.
type LookupMetadata struct {
	Name             string
	ShortDescription string
	IconURL          string
	Version          string
	InformationURL   string
}

// LookupListService provides operations for retrieving and formatting
// lookup service provider metadata from the overlay engine.
type LookupListService struct {
	provider LookupListProvider
}

type LookupServiceProviders map[string]LookupMetadata

// ListLookupServiceProviders retrieves the list of lookup service providers
// and formats them into a standardized response structure.
func (s *LookupListService) ListLookupServiceProviders() LookupServiceProviders {
	engineLookupList := s.provider.ListLookupServiceProviders()
	if engineLookupList == nil {
		return make(LookupServiceProviders)
	}

	result := make(LookupServiceProviders, len(engineLookupList))
	coalesce := func(primary, fallback string) string {
		if primary != "" {
			return primary
		}
		return fallback
	}

	for name, metadata := range engineLookupList {
		lookupMetadata := LookupMetadata{
			Name:             name,
			ShortDescription: "No description available",
		}

		if metadata != nil {
			lookupMetadata.ShortDescription = coalesce(metadata.Description, "No description available")
			lookupMetadata.IconURL = metadata.Icon
			lookupMetadata.Version = metadata.Version
			lookupMetadata.InformationURL = metadata.InfoUrl
		}

		result[name] = lookupMetadata
	}

	return result
}

// NewLookupListService creates a new LookupListService
// initialized with the given provider. It panics if the provider is nil.
func NewLookupListService(provider LookupListProvider) *LookupListService {
	if provider == nil {
		panic("lookup list provider is nil")
	}
	return &LookupListService{provider: provider}
}
