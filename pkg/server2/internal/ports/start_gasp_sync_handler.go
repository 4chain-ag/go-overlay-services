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

	return c.Status(fiber.StatusOK).JSON(NewStartGASPSyncResponse())
}

// NewStartGASPSyncHandler creates a new StartGASPSyncHandler with the given provider.
// If the provider is nil, it panics.
func NewStartGASPSyncHandler(provider app.StartGASPSyncProvider) *StartGASPSyncHandler {
	if provider == nil {
		panic("start GASP sync provider is nil")
	}

	return &StartGASPSyncHandler{service: app.NewStartGASPSyncService(provider)}
}

// NewStartGASPSyncResponse returns a new StartGASPSync response.
func NewStartGASPSyncResponse() openapi.StartGASPSync {
	return openapi.StartGASPSync{
		Message: "OK",
	}
}
