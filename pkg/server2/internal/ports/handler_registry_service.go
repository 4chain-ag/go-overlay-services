package ports

import (
	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// HandlerRegistryService defines the main point for registering HTTP handler dependencies.
// It acts as a central registry for mapping API endpoints to their handler implementations.
type HandlerRegistryService struct {
	topicManagersList *TopicManagersListHandler
	submitTransaction  *SubmitTransactionHandler
	syncAdvertisements *SyncAdvertisementsHandler
}

// AdvertisementsSync method delegates the request to the configured sync advertisements handler.
func (h *HandlerRegistryService) AdvertisementsSync(c *fiber.Ctx) error {
	return h.syncAdvertisements.Handle(c)
}

// SubmitTransaction method delegates the request to the configured submit transaction handler.
func (h *HandlerRegistryService) SubmitTransaction(c *fiber.Ctx, params openapi.SubmitTransactionParams) error {
	return h.submitTransaction.Handle(c, params)
}

// ListTopicManagers method delegates the request to the configured topic managers list handler.
func (h *HandlerRegistryService) ListTopicManagers(c *fiber.Ctx) error {
	return h.topicManagersList.Handle(c)
}

// NewHandlerRegistryService creates and returns a new HandlerRegistryService instance.
// It initializes all handler implementations with their required dependencies.
func NewHandlerRegistryService(provider engine.OverlayEngineProvider) *HandlerRegistryService {
	return &HandlerRegistryService{
		submitTransaction: NewSubmitTransactionHandler(provider),
		topicManagersList: NewTopicManagersListHandler(provider),
		syncAdvertisements: NewSyncAdvertisementsHandler(provider),
	}
}
