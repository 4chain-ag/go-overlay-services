package overlayhttp

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/gofiber/fiber/v2"
)

// SyncAdvertisementsProvider defines the contract that must be fulfilled
// to send synchronize advertisements request to the overlay engine for further processing.
type SyncAdvertisementsProvider interface {
	SyncAdvertisements(ctx context.Context) error
}

// AdvertisementsSyncHandler orchestrates the processing flow of a synchronize advertisements
// request and applies any necessary logic before invoking the engine.
type AdvertisementsSyncHandler struct {
	provider SyncAdvertisementsProvider
}

// Handle orchestrates the processing flow of a synchronize advertisements request.
// It prepares and sends a JSON response after invoking the engine and returns an HTTP response
// with the appropriate status code based on the engine's response.
func (h *AdvertisementsSyncHandler) Handle(c *fiber.Ctx) error {
	err := h.provider.SyncAdvertisements(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(NewSyncAdvertisementsProviderErrorResponse())
	}
	return c.Status(fiber.StatusOK).JSON(NewAdvertisementsSyncSuccessResponse())
}

// NewAdvertisementsSyncHandler returns an instance of a AdvertisementsSyncHandler,
// utilizing an implementation of SyncAdvertisementsProvider. If the provided argument is nil, it triggers a panic.
func NewAdvertisementsSyncHandler(provider SyncAdvertisementsProvider) *AdvertisementsSyncHandler {
	if provider == nil {
		panic("sync advertisements provider is nil")
	}

	return &AdvertisementsSyncHandler{provider: provider}
}

// NewAdvertisementsSyncSuccessResponse creates a successful response for advertisement synchronization.
// It returns an instance of openapi.AdvertisementsSyncResponse with a predefined success message.
func NewAdvertisementsSyncSuccessResponse() openapi.AdvertisementsSyncResponse {
	return openapi.AdvertisementsSyncResponse{
		Message: "Advertisement sync request successfully delegated to overlay engine.",
	}
}

// NewSyncAdvertisementsProviderErrorResponse creates an error response for advertisement synchronization failures.
// It returns an instance of openapi.InternalServerErrorResponse with a predefined error message
// indicating an issue with the overlay engine during the sync process.
func NewSyncAdvertisementsProviderErrorResponse() openapi.InternalServerErrorResponse {
	return openapi.Error{Message: "Unable to process sync advertisements request due to issues with the overlay engine."}
}
