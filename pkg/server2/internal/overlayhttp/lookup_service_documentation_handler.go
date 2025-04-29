package overlayhttp

import (
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/gofiber/fiber/v2"
)

// LookupServiceDocumentationProvider defines the contract that must be fulfilled
// to retrieve documentation for a lookup service provider from the overlay engine.
type LookupServiceDocumentationProvider interface {
	GetDocumentationForLookupServiceProvider(lookupService string) (string, error)
}

// LookupServiceDocumentationHandler orchestrates the processing flow of a lookup documentation
// request, including the request parameter validation and invoking the engine to retrieve
// the documentation.
type LookupServiceDocumentationHandler struct {
	provider LookupServiceDocumentationProvider
}

// Handle processes a request for lookup service documentation.
// It extracts the lookupService query parameter, invokes the engine provider,
// and returns the documentation as JSON with the appropriate status code.
func (l *LookupServiceDocumentationHandler) Handle(c *fiber.Ctx, params openapi.LookupServiceDocumentationParams) error {
	if params.LookupService == "" {
		return c.Status(fiber.StatusBadRequest).SendString("lookupService query parameter is required")
	}

	documentation, err := l.provider.GetDocumentationForLookupServiceProvider(params.LookupService)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(NewLookupServiceDocumentationErrorResponse())
	}

	return c.Status(fiber.StatusOK).JSON(openapi.LookupServiceDocumentation{
		Documentation: documentation,
	})
}

// NewLookupServiceDocumentationHandler returns an instance of a LookupServiceDocumentationHandler.
// If the provider argument is nil, it triggers a panic.
func NewLookupServiceDocumentationHandler(provider LookupServiceDocumentationProvider) *LookupServiceDocumentationHandler {
	if provider == nil {
		panic("lookup service documentation provider is nil")
	}

	return &LookupServiceDocumentationHandler{
		provider: provider,
	}
}

// NewLookupServiceDocumentationErrorResponse creates an error response for documentation retrieval failures.
func NewLookupServiceDocumentationErrorResponse() openapi.InternalServerErrorResponse {
	return openapi.Error{
		Message: "Unable to retrieve documentation for the specified lookup service provider.",
	}
} 
