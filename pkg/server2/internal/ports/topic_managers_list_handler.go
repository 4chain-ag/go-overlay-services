package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// TopicManagersListHandlerOption defines a function that configures a TopicManagersListHandler.
type TopicManagersListHandlerOption func(h *TopicManagersListHandler)

// TopicManagersListService defines the contract for retrieving a list of topic managers.
type TopicManagersListService interface {
	GetList(ctx context.Context) app.TopicManagersListResponse
}

// TopicManagersListHandler handles HTTP requests to retrieve a list of topic managers.
type TopicManagersListHandler struct {
	service TopicManagersListService
}

// GetList handles HTTP requests to retrieve a list of available topic managers.
func (h *TopicManagersListHandler) GetList(c *fiber.Ctx) error {
	result := h.service.GetList(c.UserContext())

	// Convert to OpenAPI response format
	response := make(openapi.TopicManagersListResponse, len(result))
	for key, metadata := range result {
		response[key] = openapi.TopicManagerMetadata{
			Name:             metadata.Name,
			ShortDescription: metadata.ShortDescription,
			IconURL:          metadata.IconURL,
			Version:          metadata.Version,
			InformationURL:   metadata.InformationURL,
		}
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// NewTopicManagersListHandler creates a new instance of TopicManagersListHandler.
// It panics if the provider is nil.
func NewTopicManagersListHandler(provider app.TopicManagersListProvider, options ...TopicManagersListHandlerOption) *TopicManagersListHandler {
	if provider == nil {
		panic("topic managers list provider cannot be nil")
	}

	handler := &TopicManagersListHandler{
		service: app.NewTopicManagersListService(provider),
	}

	for _, opt := range options {
		opt(handler)
	}

	return handler
}
