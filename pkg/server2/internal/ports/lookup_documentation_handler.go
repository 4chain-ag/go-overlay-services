package ports

import (
	"context"
	"errors"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// LookupServiceDocumentationHandlerOption defines a function that configures a LookupServiceDocumentationHandler.
type LookupServiceDocumentationHandlerOption func(h *LookupServiceDocumentationHandler)

// LookupServiceDocumentationService defines the contract for retrieving lookup service documentation.
type LookupServiceDocumentationService interface {
	GetDocumentation(ctx context.Context, lookupService string) (string, error)
}

// LookupServiceDocumentationHandler handles HTTP requests to retrieve documentation for lookup service providers.
type LookupServiceDocumentationHandler struct {
	service LookupServiceDocumentationService
}

// GetDocumentation handles HTTP requests to retrieve documentation for a specific lookup service provider.
// It extracts the lookupService query parameter, invokes the service, and returns the documentation as JSON.
func (h *LookupServiceDocumentationHandler) GetDocumentation(c *fiber.Ctx) error {
	lookupService := c.Query("lookupService")
	if lookupService == "" {
		return c.Status(fiber.StatusBadRequest).JSON(NewMissingLookupServiceParameterResponse())
	}

	documentation, err := h.service.GetDocumentation(c.UserContext(), lookupService)
	if err == nil {
		return c.Status(fiber.StatusOK).JSON(openapi.LookupServiceDocumentationResponse{
			Documentation: documentation,
		})
	}

	var target app.Error
	if !errors.As(err, &target) {
		return c.Status(fiber.StatusInternalServerError).JSON(UnhandledErrorTypeResponse)
	}

	switch target.ErrorType() {
	case app.ErrorTypeIncorrectInput:
		return c.Status(fiber.StatusBadRequest).JSON(NewInvalidLookupServiceParameterResponse())
	case app.ErrorTypeProviderFailure:
		return c.Status(fiber.StatusInternalServerError).JSON(NewLookupServiceProviderErrorResponse())
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(UnhandledErrorTypeResponse)
	}
}

// NewLookupServiceDocumentationHandler creates a new instance of LookupServiceDocumentationHandler.
// It panics if the provider is nil.
func NewLookupServiceDocumentationHandler(provider app.LookupServiceDocumentationProvider, options ...LookupServiceDocumentationHandlerOption) *LookupServiceDocumentationHandler {
	if provider == nil {
		panic("lookup service documentation provider cannot be nil")
	}

	handler := &LookupServiceDocumentationHandler{
		service: app.NewLookupServiceDocumentationService(provider),
	}

	for _, opt := range options {
		opt(handler)
	}

	return handler
}

// NewMissingLookupServiceParameterResponse returns a bad request response for missing lookupService parameter.
func NewMissingLookupServiceParameterResponse() openapi.BadRequestResponse {
	return openapi.Error{
		Message: "lookupService query parameter is required",
	}
}

// NewInvalidLookupServiceParameterResponse returns a bad request response for invalid lookupService parameter.
func NewInvalidLookupServiceParameterResponse() openapi.BadRequestResponse {
	return openapi.Error{
		Message: "lookupService parameter cannot be empty",
	}
}

// NewLookupServiceProviderErrorResponse returns an error response indicating a failure within the overlay engine.
func NewLookupServiceProviderErrorResponse() openapi.InternalServerErrorResponse {
	return openapi.Error{
		Message: "Unable to retrieve documentation for the requested lookup service",
	}
}
