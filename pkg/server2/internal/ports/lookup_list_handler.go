package ports

import (
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// LookupListService defines the interface for a service responsible for retrieving
// and formatting lookup service provider metadata.
type LookupListService interface {
	ListLookupServiceProviders() app.LookupServiceProviders
}

// LookupListHandler handles incoming requests for lookup service provider information.
// It delegates to the LookupListService to retrieve the metadata and formats
// the response according to the API spec.
type LookupListHandler struct {
	service LookupListService
}

// Handle processes an HTTP request to list all lookup service providers.
// It returns an HTTP 200 OK with a LookupServiceProvidersResponse.
func (h *LookupListHandler) Handle(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(h.service.ListLookupServiceProviders())
}

// NewLookupListHandler creates a new LookupListHandler with the given provider.
// It initializes the internal LookupListService.
// Panics if the provider is nil.
func NewLookupListHandler(provider app.LookupListProvider) *LookupListHandler {
	if provider == nil {
		panic("lookup list provider is nil")
	}
	return &LookupListHandler{service: app.NewLookupListService(provider)}
}

func NewLookupListSuccessResponse(lookupList app.LookupServiceProviders) openapi.LookupServiceProvidersListResponse {
	response := make(openapi.LookupServiceProvidersList, len(lookupList))

	for name, metadata := range lookupList {
		response[name] = openapi.LookupServiceProviderMetadata{
			Name:             metadata.Name,
			ShortDescription: metadata.ShortDescription,
			IconURL:          &metadata.IconURL,
			Version:          &metadata.Version,
			InformationURL:   &metadata.InformationURL,
		}
	}

	return response
}
