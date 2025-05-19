package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// LookupProviderDocumentationService defines the interface for retrieving lookup service documentation.
type LookupProviderDocumentationService interface {
	GetDocumentation(ctx context.Context, lookupService string) (string, error)
}

// LookupProviderDocumentationHandler handles HTTP requests to retrieve documentation for lookup service providers.
type LookupProviderDocumentationHandler struct {
	service LookupProviderDocumentationService
}

// GetDocumentation handles HTTP requests to retrieve documentation for a specific lookup service provider.
// It extracts the lookupService query parameter, invokes the service, and returns the documentation as JSON.
// Returns:
//   - 200 OK with documentation on success
//   - 400 Bad Request if lookupService parameter is missing or empty
//   - 500 Internal Server Error if the service fails to retrieve documentation
func (h *LookupProviderDocumentationHandler) Handle(c *fiber.Ctx, params openapi.GetLookupServiceProviderDocumentationParams) error {
	documentation, err := h.service.GetDocumentation(c.UserContext(), c.Query("lookupService"))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(openapi.LookupServiceDocumentationResponse{
		Documentation: documentation,
	})
}

// NewLookupProviderDocumentationHandler creates a new instance of LookupProviderDocumentationHandler.
// Panics if the provider is nil.
func NewLookupProviderDocumentationHandler(provider app.LookupServiceDocumentationProvider) *LookupProviderDocumentationHandler {
	if provider == nil {
		panic("lookup service documentation provider cannot be nil")
	}

	return &LookupProviderDocumentationHandler{
		service: app.NewLookupDocumentationService(provider),
	}
}
