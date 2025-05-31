package ports

import (
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
)

// LookupListService defines the interface for a service that provides access to
// lookup service provider metadata.
type LookupListService interface {
	// ListLookupServiceProviders returns a map of all registered lookup service providers.
	ListLookupServiceProviders() map[string]*overlay.MetaData
}

// LookupListHandler is an HTTP handler that serves requests to retrieve
// a list of registered lookup service providers along with their metadata.
// It uses a LookupListService to fetch the provider data and formats the
// response according to the OpenAPI schema.
type LookupListHandler struct {
	service LookupListService
}

// Handle processes a request to list all available lookup service providers.
// It invokes the service layer to fetch metadata about the providers, transforms
// the result into an OpenAPI-compliant response structure, and returns it as a JSON response with a 200 OK status.
func (h *LookupListHandler) Handle(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(NewLookupListSuccessResponse(h.service.ListLookupServiceProviders()))
}

// NewLookupListHandler creates a new LookupListHandler with the given provider.
// If the provider is nil, it panics.
func NewLookupListHandler(provider app.LookupListProvider) *LookupListHandler {
	if provider == nil {
		panic("lookup list provider is nil")
	}
	return &LookupListHandler{service: app.NewLookupListService(provider)}
}

// NewLookupListSuccessResponse transforms a map of provider metadata into an
// OpenAPI-compliant LookupServiceProvidersListResponse. It maps internal metadata fields into their public-facing API counterparts.
func NewLookupListSuccessResponse(lookupList map[string]*overlay.MetaData) openapi.LookupServiceProvidersListResponse {
	response := make(openapi.LookupServiceProvidersList, len(lookupList))
	for name, metadata := range lookupList {
		response[name] = openapi.LookupServiceProviderMetadata{
			Name:             metadata.Name,
			ShortDescription: metadata.Description,
			IconURL:          metadata.Icon,
			Version:          metadata.Version,
			InformationURL:   metadata.InfoUrl,
		}
	}
	return response
}
