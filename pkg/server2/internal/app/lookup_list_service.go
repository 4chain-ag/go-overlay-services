package app

import (
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// LookupServiceMetadataDTO represents metadata information for a single lookup service provider.
// This DTO is intended for client-facing usage, such as API responses.
type LookupServiceMetadataDTO struct {
	Name        string // Human-readable name of the service provider.
	Description string // Short description of what the service does.
	IconURL     string // URL to an icon that visually represents the service.
	Version     string // Version identifier of the service implementation.
	InfoURL     string // Link to documentation or more information about the service.
}

// LookupServicesMetadataDTO maps a unique service identifier (typically a string key) to
// its corresponding service metadata. It provides a lookup-friendly structure for APIs.
type LookupServicesMetadataDTO map[string]LookupServiceMetadataDTO

// LookupListProvider defines the contract for components capable of providing
// metadata about available lookup service providers. It abstracts the underlying
// source of metadata (e.g., hardcoded map, configuration file, or remote source).
type LookupListProvider interface {
	// ListLookupServiceProviders returns a map of service identifiers to their raw metadata.
	ListLookupServiceProviders() map[string]*overlay.MetaData
}

// LookupListService provides a higher-level abstraction over a LookupListProvider.
// It is responsible for converting internal metadata representations into standardized
// DTOs suitable for use in API layers or other external consumers.
type LookupListService struct {
	provider LookupListProvider
}

// ListLookupServiceProviders orchestrates the retrieval of lookup service metadata.
// It delegates the data retrieval to the underlying provider and returns a structured
// format suitable for external use. Returns a collection of service metadata records.
func (s *LookupListService) ListLookupServiceProviders() LookupServicesMetadataDTO {
	services := s.provider.ListLookupServiceProviders()
	return NewLookupServicesMetadataDTO(services)
}

// NewLookupServicesMetadataDTO creates a DTO-compliant structure from raw internal metadata.
// It maps each entry from the raw metadata map to a user-facing LookupServiceMetadataDTO.
func NewLookupServicesMetadataDTO(services map[string]*overlay.MetaData) LookupServicesMetadataDTO {
	dto := make(LookupServicesMetadataDTO, len(services))
	for serviceKey, metadata := range services {
		dto[serviceKey] = LookupServiceMetadataDTO{
			Name:        metadata.Name,
			Description: metadata.Description,
			IconURL:     metadata.Icon,
			Version:     metadata.Version,
			InfoURL:     metadata.InfoUrl,
		}
	}
	return dto
}

// NewLookupListService constructs a new LookupListService using the provided LookupListProvider.
// It ensures the service is initialized with a non-nil provider. A panic is triggered otherwise,
// as a nil provider would result in runtime errors during service operations.
func NewLookupListService(provider LookupListProvider) *LookupListService {
	if provider == nil {
		panic("lookup list provider is nil")
	}
	return &LookupListService{provider: provider}
}
