package overlayhttp

import (
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/gofiber/fiber/v2"
)

// TopicManagerDocumentationProvider defines the contract that must be fulfilled
// to retrieve documentation for a topic manager from the overlay engine.
type TopicManagerDocumentationProvider interface {
	GetDocumentationForTopicManager(topicManager string) (string, error)
}

// TopicManagerDocumentationHandler orchestrates the processing flow of a topic manager documentation
// request, including the request parameter validation and invoking the engine to retrieve
// the documentation.
type TopicManagerDocumentationHandler struct {
	provider TopicManagerDocumentationProvider
}

// Handle processes a request for topic manager documentation.
// It extracts the topicManager query parameter, invokes the engine provider,
// and returns the documentation as JSON with the appropriate status code.
func (t *TopicManagerDocumentationHandler) Handle(c *fiber.Ctx, params openapi.TopicManagerDocumentationParams) error {
	if params.TopicManager == "" {
		return c.Status(fiber.StatusBadRequest).SendString("topicManager query parameter is required")
	}

	documentation, err := t.provider.GetDocumentationForTopicManager(params.TopicManager)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(NewTopicManagerDocumentationErrorResponse())
	}

	return c.Status(fiber.StatusOK).JSON(openapi.TopicManagerDocumentation{
		Documentation: documentation,
	})
}

// NewTopicManagerDocumentationHandler returns an instance of a TopicManagerDocumentationHandler.
// If the provider argument is nil, it triggers a panic.
func NewTopicManagerDocumentationHandler(provider TopicManagerDocumentationProvider) *TopicManagerDocumentationHandler {
	if provider == nil {
		panic("topic manager documentation provider is nil")
	}

	return &TopicManagerDocumentationHandler{
		provider: provider,
	}
}

// NewTopicManagerDocumentationErrorResponse creates an error response for documentation retrieval failures.
func NewTopicManagerDocumentationErrorResponse() openapi.InternalServerErrorResponse {
	return openapi.Error{
		Message: "Unable to retrieve documentation for the specified topic manager.",
	}
} 
