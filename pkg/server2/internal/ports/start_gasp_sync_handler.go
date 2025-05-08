package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// StartGASPSyncService defines the interface for a service responsible for initiating GASP synchronization.

type StartGASPSyncService interface {
	StartGASPSync(ctx context.Context) error
}

// StartGASPSyncHandler handles the /api/v1/admin/start-gasp-sync endpoint.

type StartGASPSyncHandler struct {
	service StartGASPSyncService
}

// Handle initiates the GASP sync and returns the appropriate status.

func (h *StartGASPSyncHandler) Handle(c *fiber.Ctx) error {

	if err := h.service.StartGASPSync(c.UserContext()); err != nil {

		return err

	}

	return c.Status(fiber.StatusOK).JSON(StartGASPSyncSuccessResponse)

}

// NewStartGASPSyncHandler creates a new StartGASPSyncHandler with the given provider.

// If the provider is nil, it panics.

func NewStartGASPSyncHandler(provider app.StartGASPSyncProvider) *StartGASPSyncHandler {

	if provider == nil {

		panic("start GASP sync provider is nil")

	}

	return &StartGASPSyncHandler{

		service: app.NewStartGASPSyncService(provider),
	}

}

// StartGASPSyncSuccessResponse is the success response for starting GASP sync.

var StartGASPSyncSuccessResponse = openapi.StartGASPSync{

	Message: "OK",
}

// StartGASPSyncInternalErrorResponse is the internal server error response for GASP sync initiation.

// This error is returned when an internal issue occurs while initiating GASP synchronization.

var StartGASPSyncInternalErrorResponse = openapi.InternalServerErrorResponse{

	Message: "Unable to start GASP synchronization due to an error in the overlay engine.",
}
