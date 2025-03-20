package commands

import (
	"context"
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/dto"
	"github.com/gofiber/fiber/v2"
)

// SyncAdvertisementsProvider defines the contract that must be fulfilled
// to send synchronize advertisements request to the overlay engine for further processing.
// Note: The contract definition is still in development and will be updated after
// migrating the engine code.
type SyncAdvertisementsProvider interface {
	SyncAdvertisments(ctx context.Context) error
}

// SyncAdvertismentsHandler orchestrates the processing flow of a synchronize advertisements
// request and applies any necessary logic before invoking the engine.
type SyncAdvertismentsHandler struct {
	provider SyncAdvertisementsProvider
}

// Handle orchestrates the processing flow of a synchronize advertisements request.
// It prepares and sends a JSON response after invoking the engine and returns an HTTP response
// with the appropriate status code based on the engine's response.
func (s *SyncAdvertismentsHandler) Handle(c *fiber.Ctx) error {
	err := s.provider.SyncAdvertisments(c.Context())
	if err != nil {
		if inner := c.Status(fiber.StatusInternalServerError).JSON(dto.HandlerResponseNonOK); inner != nil {
			return fmt.Errorf("failed to send JSON response: %w", inner)
		}
		return nil
	}

	if err := c.Status(fiber.StatusOK).JSON(dto.HandlerResponseOK); err != nil {
		return fmt.Errorf("failed to send JSON response: %w", nil)
	}
	return nil
}

// NewSyncAdvertismentsHandler returns an instance of a SyncAdvertismentsHandler, utilizing
// an implementation of SyncAdvertisementsProvider. If the provided argument is nil, it triggers a panic.
func NewSyncAdvertismentsHandler(provider SyncAdvertisementsProvider) *SyncAdvertismentsHandler {
	if provider == nil {
		panic("sync advertisements provider is nil")
	}
	return &SyncAdvertismentsHandler{
		provider: provider,
	}
}
