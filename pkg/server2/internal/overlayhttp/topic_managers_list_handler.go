package overlayhttp

import (
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
	"k8s.io/utils/ptr"
)

// TopicManagersListProvider defines the contract that must be fulfilled

// to retrieve a list of available topic managers from the overlay engine.

type TopicManagersListProvider interface {
	ListTopicManagers() map[string]*overlay.MetaData
}

// TopicManagersListHandler orchestrates the retrieval of available topic managers

// and returns their metadata.

type TopicManagersListHandler struct {
	provider TopicManagersListProvider
}

// Handle processes a request for listing topic managers.

// It retrieves the list of available topic managers and returns their metadata as JSON.

func (t *TopicManagersListHandler) Handle(c *fiber.Ctx) error {

	managers := t.provider.ListTopicManagers()

	if managers == nil {

		return c.Status(fiber.StatusOK).JSON(make(map[string]openapi.TopicManagerMetadata))

	}

	response := make(map[string]openapi.TopicManagerMetadata, len(managers))

	setIfNotEmpty := func(s string) *string {

		if s == "" {

			return nil

		}

		return ptr.To(s)

	}

	coalesce := func(primary, fallback string) string {

		if primary != "" {

			return primary

		}

		return fallback

	}

	for name, metadata := range managers {

		managerMetadata := openapi.TopicManagerMetadata{

			Name: name,

			Description: "No description available",
		}

		if metadata != nil {

			managerMetadata.Description = coalesce(metadata.Description, "No description available")

			managerMetadata.IconURL = setIfNotEmpty(metadata.Icon)

			managerMetadata.Version = setIfNotEmpty(metadata.Version)

			managerMetadata.InformationURL = setIfNotEmpty(metadata.InfoUrl)

		}

		response[name] = managerMetadata

	}

	return c.Status(fiber.StatusOK).JSON(response)

}

// NewTopicManagersListHandler returns a new instance of a TopicManagersListHandler.

// If the provider argument is nil, it triggers a panic.

func NewTopicManagersListHandler(provider TopicManagersListProvider) *TopicManagersListHandler {

	if provider == nil {

		panic("topic managers list provider is nil")

	}

	return &TopicManagersListHandler{

		provider: provider,
	}

}
