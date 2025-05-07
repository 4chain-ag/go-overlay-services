package app

import (
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// LookupServicesListProvider defines the interface for retrieving
// a list of lookup service providers from the overlay engine.
type LookupServicesListProvider interface {
	ListLookupServiceProviders() map[string]*overlay.MetaData
}

// LookupMetadata represents the metadata for a lookup service provider.
type LookupMetadata struct {
	Name             string  `json:"name"`
	ShortDescription string  `json:"shortDescription"`
	IconURL          *string `json:"iconURL,omitempty"`
	Version          *string `json:"version,omitempty"`
	InformationURL   *string `json:"informationURL,omitempty"`
}

// LookupServicesListResponse defines the response data structure for the lookup services list.
type LookupServicesListResponse map[string]LookupMetadata

// LookupServicesListService provides operations for retrieving and formatting
// lookup service provider metadata from the overlay engine.
type LookupServicesListService struct {
	provider LookupServicesListProvider
}

// ListLookupServiceProviders retrieves the list of lookup service providers
// and formats them into a standardized response structure.
func (s *LookupServicesListService) ListLookupServiceProviders() LookupServicesListResponse {
	// Retrieve providers from the engine
	engineLookupProviders := s.provider.ListLookupServiceProviders()

	// If nil is returned, provide an empty map
	if engineLookupProviders == nil {
		return make(LookupServicesListResponse)
	}

	result := make(LookupServicesListResponse, len(engineLookupProviders))

	setIfNotEmpty := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	coalesce := func(primary, fallback string) string {
		if primary != "" {
			return primary
		}
		return fallback
	}

	for name, metadata := range engineLookupProviders {
		lookupMetadata := LookupMetadata{
			Name:             name,
			ShortDescription: "No description available",
		}

		if metadata != nil {
			lookupMetadata.ShortDescription = coalesce(metadata.Description, "No description available")
			lookupMetadata.IconURL = setIfNotEmpty(metadata.Icon)
			lookupMetadata.Version = setIfNotEmpty(metadata.Version)
			lookupMetadata.InformationURL = setIfNotEmpty(metadata.InfoUrl)
		}

		result[name] = lookupMetadata
	}

	return result
}

// NewLookupServicesListService creates a new LookupServicesListService
// initialized with the given provider. It returns an error if the provider is nil.
func NewLookupServicesListService(provider LookupServicesListProvider) (*LookupServicesListService, error) {
	if provider == nil {
		return nil, NewIncorrectInputError("lookup services list provider cannot be nil")
	}
	return &LookupServicesListService{provider: provider}, nil
}
