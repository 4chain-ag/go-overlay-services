package app

import (
	"context"

	"github.com/bsv-blockchain/go-sdk/overlay"
)

// LookupServicesListProvider defines the contract for retrieving a list of lookup services.
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

// LookupServicesListResponse contains the response data for the lookup services list operation.
type LookupServicesListResponse map[string]LookupMetadata

// LookupServicesListService provides functionality for retrieving a list of lookup services.
type LookupServicesListService struct {
	provider LookupServicesListProvider
}

// GetList retrieves the list of available lookup services.
func (s *LookupServicesListService) GetList(ctx context.Context) LookupServicesListResponse {
	engineLookupProviders := s.provider.ListLookupServiceProviders()
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

// NewLookupServicesListService creates a new LookupServicesListService with the given provider.
// It panics if the provider is nil.
func NewLookupServicesListService(provider LookupServicesListProvider) *LookupServicesListService {
	if provider == nil {
		panic("lookup services list provider cannot be nil")
	}

	return &LookupServicesListService{
		provider: provider,
	}
}
