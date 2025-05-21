package ports

import (

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/gofiber/fiber/v2"
)

// TopicManagersListService defines the interface for a service responsible for retrieving
// and formatting topic manager metadata.
type TopicManagersListService interface {
	ListTopicManagers() app.TopicManagers
}

// TopicManagersListHandler handles incoming requests for topic manager information.
// It delegates to the TopicManagersListService to retrieve the metadata and formats
// the response according to the API spec.
type TopicManagersListHandler struct {
	service TopicManagersListService
}

// Handle processes an HTTP request to list all topic managers.
// It returns an HTTP 200 OK with a TopicManagersListResponse.
func (h *TopicManagersListHandler) Handle(c *fiber.Ctx) error {

	return c.Status(fiber.StatusOK).JSON(h.service.ListTopicManagers())
}

// NewTopicManagersListHandler creates a new TopicManagersListHandler with the given provider.
// It initializes the internal TopicManagersListService.
// Panics if the provider is nil.
func NewTopicManagersListHandler(provider app.TopicManagersListProvider) *TopicManagersListHandler {
	if provider == nil {
		panic("topic manager list provider is nil")
	}
	
	return &TopicManagersListHandler{ service: app.NewTopicManagersListService(provider) }
}

func NewTopicManagersListSuccessResponse(topicManagers app.TopicManagers) openapi.TopicManagersListResponse {
	response := make(openapi.TopicManagersList, len(topicManagers))

	for name, metadata := range topicManagers {
		response[name] = openapi.TopicManagerMetadata{
			Name:             metadata.Name,
			ShortDescription: metadata.ShortDescription,
			IconURL:          &metadata.IconURL,
			Version:          &metadata.Version,
			InformationURL:   &metadata.InformationURL,
		}
	}

	return response
}
