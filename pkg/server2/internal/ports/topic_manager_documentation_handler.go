package ports

import (
	"context"

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
// Returns 200 OK with documentation on success
func (h *TopicManagerDocumentationHandler) Handle(c *fiber.Ctx, params openapi.GetTopicManagerDocumentationParams) error {
	documentation, err := h.service.GetDocumentation(c.UserContext(), c.Query("topicManager"))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(openapi.TopicManagerDocumentationResponse{
		Documentation: documentation,
	})
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
