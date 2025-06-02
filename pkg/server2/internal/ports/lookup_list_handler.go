package ports

import (
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// LookupListService defines the contract for retrieving metadata
// about registered lookup service providers.
//
// It is designed to abstract away the underlying implementation details of how
// the metadata is collected—whether it be static, dynamic, or retrieved from an external system.
type LookupListService interface {
	// ListLookupServiceProviders returns a collection of lookup service metadata,
	// where each entry is keyed by the service name.
	//
	// This metadata typically includes information such as the service's name,
	// description, version, icon, and reference documentation URL.
	ListLookupServiceProviders() app.LookupServicesMetadataDTO
}

// LookupListHandler is an HTTP handler responsible for processing incoming requests
// to enumerate all available lookup service providers.
//
// The handler retrieves the metadata through the associated LookupListService
// and responds with a JSON object that complies with the OpenAPI schema.
type LookupListHandler struct {
	service LookupListService
}

// Handle handles an HTTP GET request to retrieve the list of registered lookup service providers.
//
// The response includes metadata for each service such as name, description,
// version, icon URL, and informational URL. The handler returns this metadata
// as a JSON object with a 200 OK HTTP status.
func (h *LookupListHandler) Handle(c *fiber.Ctx) error {
	metadata := h.service.ListLookupServiceProviders()
	return c.
		Status(fiber.StatusOK).
		JSON(NewLookupServicesMetadataSuccessResponse(metadata))
}

// NewLookupListHandler creates a new LookupListHandler with the given provider.
// If the provider is nil, it panics.
func NewLookupListHandler(provider app.LookupListProvider) *LookupListHandler {
	if provider == nil {
		panic("lookup list provider is nil")
	}
	return &LookupListHandler{service: app.NewLookupListService(provider)}
}

// NewLookupServicesMetadataSuccessResponse converts a service metadata DTO map (internal representation)
// into a response structure that complies with the OpenAPI schema.
//
// It creates an OpenAPI-compliant response containing metadata such as service name, description,
// icon URL, version, and documentation URL for each registered lookup service.
func NewLookupServicesMetadataSuccessResponse(dto app.LookupServicesMetadataDTO) openapi.LookupServiceProvidersListResponse {
	response := make(openapi.LookupServiceProvidersList, len(dto))
	for serviceID, metadata := range dto {
		response[serviceID] = openapi.LookupServiceProviderMetadata{
			Name:             metadata.Name,
			ShortDescription: metadata.Description,
			IconURL:          metadata.IconURL,
			Version:          metadata.Version,
			InformationURL:   metadata.InfoURL,
		}
	}
	return response
}
