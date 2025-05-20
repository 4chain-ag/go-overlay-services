package app

import (
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// LookupListProvider defines the interface for retrieving
// a list of lookup service providers from the overlay engine.
type LookupListProvider interface {
	ListLookupServiceProviders() map[string]*overlay.MetaData
}

// LookupServiceProviderMetadata represents the metadata for a lookup service provider.
type LookupServiceProviderMetadata struct {
	Name             string  `json:"name"`
	ShortDescription string  `json:"shortDescription"`
	IconURL          *string `json:"iconURL,omitempty"`
	Version          *string `json:"version,omitempty"`
	InformationURL   *string `json:"informationURL,omitempty"`
}

// LookupListResponse defines the response data structure for the lookup service providers list.
type LookupListResponse map[string]LookupServiceProviderMetadata

// LookupListService provides operations for retrieving and formatting
// lookup service provider metadata from the overlay engine.
type LookupListService struct {
	provider LookupListProvider
}

// ListLookup retrieves the list of lookup service providers
// and formats them into a standardized response structure.
func (s *LookupListService) ListLookup() LookupListResponse {
	// Retrieve lookup service providers from the engine
	engineLookupServiceProviders := s.provider.ListLookupServiceProviders()

	// If nil is returned, provide an empty map
	if engineLookupServiceProviders == nil {
		return make(LookupListResponse)
	}

	result := make(LookupListResponse, len(engineLookupServiceProviders))

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

	for name, metadata := range engineLookupServiceProviders {
		lookupServiceProviderMetadata := LookupServiceProviderMetadata{
			Name:             name,
			ShortDescription: "No description available",
		}

		if metadata != nil {
			lookupServiceProviderMetadata.ShortDescription = coalesce(metadata.Description, "No description available")
			lookupServiceProviderMetadata.IconURL = setIfNotEmpty(metadata.Icon)
			lookupServiceProviderMetadata.Version = setIfNotEmpty(metadata.Version)
			lookupServiceProviderMetadata.InformationURL = setIfNotEmpty(metadata.InfoUrl)
		}

		result[name] = lookupServiceProviderMetadata
	}

	return result
}

// NewLookupListService creates a new LookupListService
// initialized with the given provider. It returns an error if the provider is nil.
func NewLookupListService(provider LookupListProvider) (*LookupListService, error) {
	if provider == nil {
		return nil, NewLookupNilProviderError("lookup service provider list provider")
	}
	return &LookupListService{provider: provider}, nil
}

// NewLookupNilProviderError returns an Error indicating that a required lookup service provider was nil,
// which is invalid input when creating a lookup service provider service.
func NewLookupNilProviderError(providerName string) Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       providerName + " cannot be nil",
		slug:      "The required provider was not properly initialized",
	}
}
