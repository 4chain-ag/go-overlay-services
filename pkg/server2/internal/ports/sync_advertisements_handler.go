package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// SyncAdvertisementsService abstracts the logic for handling advertisement synchronization.
// It delegates the responsibility of syncing advertisements to an underlying implementation.
type SyncAdvertisementsService interface {
	SyncAdvertisements(ctx context.Context) error
}

// SyncAdvertisementsHandler orchestrates the processing flow of a synchronize advertisements
// request and applies any necessary logic before invoking the synchronization engine.
type SyncAdvertisementsHandler struct {
	service SyncAdvertisementsService
}

// Handle processes an HTTP request to synchronize advertisements.
// It invokes the underlying service and returns an HTTP 200 OK on success.
// If an internal error occurs during synchronization, it returns an HTTP 500 Internal Server Error.
func (s *SyncAdvertisementsHandler) Handle(c *fiber.Ctx) error {
	err := s.service.SyncAdvertisements(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(SyncAdvertisementsSuccessResponse)
}

// NewSyncAdvertisementsHandler returns a new instance of SyncAdvertisementsHandler,
// using the given SyncAdvertisementsProvider implementation.
// If the provider is nil, the function will panic to avoid misconfigured handlers at runtime.
func NewSyncAdvertisementsHandler(provider app.SyncAdvertisementsProvider) *SyncAdvertisementsHandler {
	if provider == nil {
		panic("sync advertisements provider is nil")
	}

	return &SyncAdvertisementsHandler{
		service: app.NewAdvertisementsSyncService(provider),
	}
}

// SyncAdvertisementsSuccessResponse is returned when the advertisement synchronization
// request is successfully delegated to the overlay engine.
var SyncAdvertisementsSuccessResponse = openapi.AdvertisementsSyncResponse{
	Message: "Advertisement sync request successfully delegated to overlay engine.",
}
