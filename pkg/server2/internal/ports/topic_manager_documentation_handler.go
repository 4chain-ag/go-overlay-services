package ports

import (
	"context"
	"errors"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// TopicManagerDocumentationService defines the interface for retrieving topic manager documentation.

type TopicManagerDocumentationService interface {
	GetDocumentation(ctx context.Context, topicManager string) (string, error)
}

// TopicManagerDocumentationHandler handles HTTP requests to retrieve documentation for topic managers.

type TopicManagerDocumentationHandler struct {
	service TopicManagerDocumentationService
}

// GetDocumentation handles HTTP requests to retrieve documentation for a specific topic manager.

// It extracts the topicManager query parameter, invokes the service, and returns the documentation as JSON.

// Returns:

//   - 200 OK with documentation on success

//   - 400 Bad Request if topicManager parameter is missing or empty

//   - 500 Internal Server Error if the service fails to retrieve documentation

func (h *TopicManagerDocumentationHandler) GetDocumentation(c *fiber.Ctx) error {

	topicManager := c.Query("topicManager")

	if topicManager == "" {

		return c.Status(fiber.StatusBadRequest).JSON(TopicManagerMissingParameter)

	}

	documentation, err := h.service.GetDocumentation(c.UserContext(), topicManager)

	var target app.Error

	if err != nil && !errors.As(err, &target) {

		return c.Status(fiber.StatusInternalServerError).JSON(UnhandledErrorTypeResponse)

	}

	switch target.ErrorType() {

	case app.ErrorTypeIncorrectInput:

		return c.Status(fiber.StatusBadRequest).JSON(TopicManagerInvalidParameter)

	case app.ErrorTypeProviderFailure:

		return c.Status(fiber.StatusInternalServerError).JSON(TopicManagerError)

	default:

		return c.Status(fiber.StatusOK).JSON(openapi.TopicManagerDocumentationResponse{

			Documentation: documentation,
		})

	}

}

// NewTopicManagerDocumentationHandler creates a new instance of TopicManagerDocumentationHandler.

// Panics if the provider is nil.

func NewTopicManagerDocumentationHandler(provider app.TopicManagerDocumentationProvider) *TopicManagerDocumentationHandler {

	if provider == nil {

		panic("topic manager documentation provider cannot be nil")

	}

	return &TopicManagerDocumentationHandler{

		service: app.NewTopicManagerDocumentationService(provider),
	}

}

// TopicManagerMissingParameter is the bad request response for missing topicManager parameter.

var TopicManagerMissingParameter = openapi.BadRequestResponse{

	Message: "topicManager query parameter is required",
}

// TopicManagerInvalidParameter is the bad request response for invalid topicManager parameter.

var TopicManagerInvalidParameter = openapi.BadRequestResponse{

	Message: "topic manager name cannot be empty",
}

// TopicManagerError is the internal server error response for provider failures.

var TopicManagerError = openapi.InternalServerErrorResponse{

	Message: "Unable to retrieve documentation for the requested topic manager",
}
