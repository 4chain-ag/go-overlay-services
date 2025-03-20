package queries

import (
	"context"
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/dto"
	"github.com/gofiber/fiber/v2"
)

// TopicManagerDocumentationProvider defines the contract that must be fulfilled
// to send a topic manager documentation request to the overlay engine for further processing.
// Note: The contract definition is still in development and will be updated after
// migrating the engine code.
type TopicManagerDocumentationProvider interface {
	GetTopicManagerDocumentation(ctx context.Context) error
}

// TopicManagerDocumentationHandler orchestrates the processing flow of a topic documentation
// request, including the request body validation, converting the request body
// into an overlay-engine-compatible format, and applying any other necessary
// logic before invoking the engine. It returns the requested topic manager
// documentation in the text/markdown format.
type TopicManagerDocumentationHandler struct {
	provider TopicManagerDocumentationProvider
}

// Handle orchestrates the processing flow of a topic manager documentation request.
// It prepares and sends a JSON response after invoking the engine and returns an HTTP response
// with the appropriate status code based on the engine's response.
func (t *TopicManagerDocumentationHandler) Handle(c *fiber.Ctx) error {
	// TODO: Add custom validation logic.
	err := t.provider.GetTopicManagerDocumentation(c.Context())
	if err != nil {
		if inner := c.Status(fiber.StatusInternalServerError).JSON(dto.HandlerResponseNonOK); inner != nil {
			return fmt.Errorf("failed to send JSON response: %w", inner)
		}
		return nil
	}

	if err := c.Status(fiber.StatusOK).JSON(dto.HandlerResponseOK); err != nil {
		return fmt.Errorf("failed to send JSON response: %w", nil)
	}
	return nil
}

// NewTopicManagerDocumentationHandler returns an instance of a TopicManagerDocumentationHandler, utilizing
// an implementation of TopicManagerDocumentationProvider. If the provided argument is nil, it triggers a panic.
func NewTopicManagerDocumentationHandler(provider TopicManagerDocumentationProvider) *TopicManagerDocumentationHandler {
	if provider == nil {
		panic("topic manager documentation provider is nil")
	}
	return &TopicManagerDocumentationHandler{
		provider: provider,
	}
}
