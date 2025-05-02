package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// LookupServicesListHandlerOption defines a function that configures a LookupServicesListHandler.
type LookupServicesListHandlerOption func(h *LookupServicesListHandler)

// LookupServicesListService defines the contract for retrieving a list of lookup services.
type LookupServicesListService interface {
	GetList(ctx context.Context) app.LookupServicesListResponse
}

// LookupServicesListHandler handles HTTP requests to retrieve a list of lookup services.
type LookupServicesListHandler struct {
	service LookupServicesListService
}

// GetList handles HTTP requests to retrieve a list of available lookup services.
func (h *LookupServicesListHandler) GetList(c *fiber.Ctx) error {
	result := h.service.GetList(c.UserContext())

	// Convert to OpenAPI response format
	response := make(openapi.LookupServicesListResponse, len(result))
	for key, metadata := range result {
		response[key] = openapi.LookupMetadata{
			Name:             metadata.Name,
			ShortDescription: metadata.ShortDescription,
			IconURL:          metadata.IconURL,
			Version:          metadata.Version,
			InformationURL:   metadata.InformationURL,
		}
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// NewLookupServicesListHandler creates a new instance of LookupServicesListHandler.
// It panics if the provider is nil.
func NewLookupServicesListHandler(provider app.LookupServicesListProvider, options ...LookupServicesListHandlerOption) *LookupServicesListHandler {
	if provider == nil {
		panic("lookup services list provider cannot be nil")
	}

	handler := &LookupServicesListHandler{
		service: app.NewLookupServicesListService(provider),
	}

	for _, opt := range options {
		opt(handler)
	}

	return handler
}
