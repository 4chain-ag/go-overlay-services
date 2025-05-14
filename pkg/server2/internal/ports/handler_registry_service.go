package ports

import (
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// HandlerRegistryService defines the main point for registering HTTP handler dependencies.
// It acts as a central registry for mapping API endpoints to their handler implementations.
type HandlerRegistryService struct {
	submitTransaction  *SubmitTransactionHandler
	syncAdvertisements *SyncAdvertisementsHandler
}

// SetSubmitTransactionHandler sets a custom SubmitTransaction handler implementation.
// This allows replacing the default handler with an alternative (e.g., for testing).
func (h *HandlerRegistryService) SetSubmitTransactionHandler(provider app.SubmitTransactionProvider, timeout time.Duration) {
	h.submitTransaction = NewSubmitTransactionHandler(provider, timeout)
}

// AdvertisementsSync handles the /advertisements/sync API endpoint.
// This method delegates the request to the configured SyncAdvertisementsHandler.
func (h *HandlerRegistryService) AdvertisementsSync(c *fiber.Ctx) error {
	return h.syncAdvertisements.Handle(c)
}

// SetSubmitTransactionHandler sets a custom SubmitTransaction handler implementation.
// This allows replacing the default handler with an alternative (e.g., for testing).
func (h *HandlerRegistryService) SubmitTransaction(c *fiber.Ctx, params openapi.SubmitTransactionParams) error {
	return h.submitTransaction.Handle(c, params)
}

// NewHandlerRegistryService creates and returns a new HandlerRegistryService instance.
// It initializes all handler implementations with their required dependencies.
func NewHandlerRegistryService(provider engine.OverlayEngineProvider, timeout time.Duration) *HandlerRegistryService {
	return &HandlerRegistryService{
		syncAdvertisements: NewSyncAdvertisementsHandler(provider),
		submitTransaction:  NewSubmitTransactionHandler(provider, timeout),
	}
}
