package ports

import (
	"fmt"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// TopicManagersListService defines the interface for a service responsible for retrieving
// and formatting topic manager metadata.
type TopicManagersListService interface {
	ListTopicManagers() app.TopicManagersListResponse
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
	response := h.service.ListTopicManagers()
	return c.Status(http.StatusOK).JSON(response)
}

// NewTopicManagersListHandler creates a new TopicManagersListHandler with the given provider.
// It initializes the internal TopicManagersListService.
// Panics if the provider is nil.
func NewTopicManagersListHandler(provider app.TopicManagersListProvider) *TopicManagersListHandler {
	service, err := app.NewTopicManagersListService(provider)
	if err != nil {
		panic(fmt.Sprintf("failed to create topic managers list service: %v", err))
	}
	return &TopicManagersListHandler{
		service: service,
	}
}

// TopicManagersListServiceInternalError is the internal server error response for topic managers list.
// This error is returned when an internal issue occurs while retrieving the topic managers list.
var TopicManagersListServiceInternalError = openapi.InternalServerErrorResponse{
	Message: "Unable to retrieve topic managers list due to an error in the overlay engine.",
}
