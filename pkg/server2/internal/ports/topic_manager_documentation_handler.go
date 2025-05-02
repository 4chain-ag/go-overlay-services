package ports

import (
	"context"
	"errors"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// TopicManagerDocumentationHandlerOption defines a function that configures a TopicManagerDocumentationHandler.
type TopicManagerDocumentationHandlerOption func(h *TopicManagerDocumentationHandler)

// TopicManagerDocumentationService defines the contract for retrieving topic manager documentation.
type TopicManagerDocumentationService interface {
	GetDocumentation(ctx context.Context, topicManager string) (string, error)
}

// TopicManagerDocumentationHandler handles HTTP requests to retrieve documentation for topic managers.
type TopicManagerDocumentationHandler struct {
	service TopicManagerDocumentationService
}

// GetDocumentation handles HTTP requests to retrieve documentation for a specific topic manager.
// It extracts the topicManager query parameter, invokes the service, and returns the documentation as JSON.
func (h *TopicManagerDocumentationHandler) GetDocumentation(c *fiber.Ctx) error {
	topicManager := c.Query("topicManager")
	if topicManager == "" {
		return c.Status(fiber.StatusBadRequest).JSON(NewMissingTopicManagerParameterResponse())
	}

	documentation, err := h.service.GetDocumentation(c.UserContext(), topicManager)
	switch {
	case errors.Is(err, app.ErrEmptyTopicManagerName):
		return c.Status(fiber.StatusBadRequest).JSON(NewInvalidTopicManagerParameterResponse())

	case errors.Is(err, app.ErrTopicManagerNotFound):
		return c.Status(fiber.StatusInternalServerError).JSON(NewTopicManagerProviderErrorResponse())

	default:
		return c.Status(fiber.StatusOK).JSON(openapi.TopicManagerDocumentationResponse{
			Documentation: documentation,
		})
	}
}

// NewTopicManagerDocumentationHandler creates a new instance of TopicManagerDocumentationHandler.
// It panics if the provider is nil.
func NewTopicManagerDocumentationHandler(provider app.TopicManagerDocumentationProvider, options ...TopicManagerDocumentationHandlerOption) *TopicManagerDocumentationHandler {
	if provider == nil {
		panic("topic manager documentation provider cannot be nil")
	}

	handler := &TopicManagerDocumentationHandler{
		service: app.NewTopicManagerDocumentationService(provider),
	}

	for _, opt := range options {
		opt(handler)
	}

	return handler
}

// NewMissingTopicManagerParameterResponse returns a bad request response for missing topicManager parameter.
func NewMissingTopicManagerParameterResponse() openapi.BadRequestResponse {
	return openapi.Error{
		Message: "topicManager query parameter is required",
	}
}

// NewInvalidTopicManagerParameterResponse returns a bad request response for invalid topicManager parameter.
func NewInvalidTopicManagerParameterResponse() openapi.BadRequestResponse {
	return openapi.Error{
		Message: "topicManager parameter cannot be empty",
	}
}

// NewTopicManagerProviderErrorResponse returns an error response indicating a failure within the overlay engine.
func NewTopicManagerProviderErrorResponse() openapi.InternalServerErrorResponse {
	return openapi.Error{
		Message: "Unable to retrieve documentation for the requested topic manager",
	}
}
