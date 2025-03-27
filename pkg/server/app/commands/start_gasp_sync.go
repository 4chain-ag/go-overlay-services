package commands

import (
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/dto"
	"github.com/gofiber/fiber/v2"
)

// StartGASPSyncProvider defines the contract for triggering GASP sync.
type StartGASPSyncProvider interface {
	StartGASPSync() error
}

// StartGASPSyncHandler handles the /admin/start-gasp-sync endpoint.
type StartGASPSyncHandler struct {
	provider StartGASPSyncProvider
}

// Handle initiates the sync and returns appropriate status.
func (h *StartGASPSyncHandler) Handle(c *fiber.Ctx) error {
	if err := h.provider.StartGASPSync(); err != nil {
		if wrapErr := c.Status(fiber.StatusInternalServerError).JSON(dto.HandlerResponseNonOK); wrapErr != nil {
			return fmt.Errorf("failed to write 500 JSON response: %w", wrapErr)
		}
		return nil
	}

	if wrapErr := c.Status(fiber.StatusOK).JSON(dto.HandlerResponseOK); wrapErr != nil {
		return fmt.Errorf("failed to write 200 JSON response: %w", wrapErr)
	}
	return nil
}

// NewStartGASPSyncHandler constructs the handler.
func NewStartGASPSyncHandler(provider StartGASPSyncProvider) *StartGASPSyncHandler {
	if provider == nil {
		panic("start GASP sync provider is nil")
	}
	return &StartGASPSyncHandler{provider: provider}
}
