package app

import (
	"context"
	"errors"

	"github.com/bsv-blockchain/go-sdk/overlay"
)

var (

	// ErrEmptyLookupServiceName is returned when an empty lookup service name is provided

	ErrEmptyLookupServiceName = errors.New("lookup service name cannot be empty")

	// ErrLookupServiceProviderNotFound is returned when a lookup service provider cannot be found

	ErrLookupServiceProviderNotFound = errors.New("lookup service provider not found")
)

// LookupServicesListProvider defines the contract for retrieving a list of lookup services.

type LookupServicesListProvider interface {
	ListLookupServiceProviders() map[string]*overlay.MetaData
}

// LookupMetadata represents the metadata for a lookup service provider.

type LookupMetadata struct {
	Name string `json:"name"`

	ShortDescription string `json:"shortDescription"`

	IconURL *string `json:"iconURL,omitempty"`

	Version *string `json:"version,omitempty"`

	InformationURL *string `json:"informationURL,omitempty"`
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

			Name: name,

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

// LookupServiceDocumentationProvider defines the contract for retrieving documentation

// for a lookup service provider.

type LookupServiceDocumentationProvider interface {
	GetDocumentationForLookupServiceProvider(lookupService string) (string, error)
}

// LookupServiceDocumentationService provides functionality for retrieving lookup service documentation.

type LookupServiceDocumentationService struct {
	provider LookupServiceDocumentationProvider
}

// GetDocumentation retrieves documentation for a specific lookup service.

func (s *LookupServiceDocumentationService) GetDocumentation(ctx context.Context, lookupService string) (string, error) {

	if lookupService == "" {

		return "", NewIncorrectInputError("lookup service name cannot be empty")

	}

	documentation, err := s.provider.GetDocumentationForLookupServiceProvider(lookupService)

	if err != nil {

		return "", NewProviderFailureError("unable to retrieve documentation for lookup service provider")

	}

	return documentation, nil

}

// NewLookupServiceDocumentationService creates a new LookupServiceDocumentationService with the given provider.

// It panics if the provider is nil.

func NewLookupServiceDocumentationService(provider LookupServiceDocumentationProvider) *LookupServiceDocumentationService {

	if provider == nil {

		panic("lookup service documentation provider cannot be nil")

	}

	return &LookupServiceDocumentationService{

		provider: provider,
	}

}
