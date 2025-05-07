package ports

import (
	"context"
	"errors"

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

func (h *LookupProviderDocumentationHandler) GetDocumentation(c *fiber.Ctx) error {

	lookupService := c.Query("lookupService")

	if lookupService == "" {

		return c.Status(fiber.StatusBadRequest).JSON(LookupProviderMissingParameter)

	}

	documentation, err := h.service.GetDocumentation(c.UserContext(), lookupService)

	var target app.Error

	if err != nil && !errors.As(err, &target) {

		return c.Status(fiber.StatusInternalServerError).JSON(UnhandledErrorTypeResponse)

	}

	switch target.ErrorType() {

	case app.ErrorTypeIncorrectInput:

		return c.Status(fiber.StatusBadRequest).JSON(LookupProviderInvalidParameter)

	case app.ErrorTypeProviderFailure:

		return c.Status(fiber.StatusInternalServerError).JSON(LookupProviderError)

	default:

		return c.Status(fiber.StatusOK).JSON(openapi.LookupServiceDocumentationResponse{

			Documentation: documentation,
		})

	}

}

// NewLookupProviderDocumentationHandler creates a new instance of LookupProviderDocumentationHandler.

// Panics if the provider is nil.

func NewLookupProviderDocumentationHandler(provider app.LookupProviderDocumentationProvider) *LookupProviderDocumentationHandler {

	if provider == nil {

		panic("lookup service documentation provider cannot be nil")

	}

	return &LookupProviderDocumentationHandler{

		service: app.NewLookupProviderDocumentationService(provider),
	}

}

// LookupProviderMissingParameter is the bad request response for missing lookupService parameter.

var LookupProviderMissingParameter = openapi.BadRequestResponse{

	Message: "lookupService query parameter is required",
}

// LookupProviderInvalidParameter is the bad request response for invalid lookupService parameter.

var LookupProviderInvalidParameter = openapi.BadRequestResponse{

	Message: "lookup service name cannot be empty",
}

// LookupProviderError is the internal server error response for provider failures.

var LookupProviderError = openapi.InternalServerErrorResponse{

	Message: "Unable to retrieve documentation for the requested lookup service",
}
