package commands

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/dto"
	"github.com/gofiber/fiber/v2"
)

type SyncAdvertisementsProvider interface {
	SyncAdvertisments(ctx context.Context) error
}

type SyncAdvertismentsHandler struct {
	provider SyncAdvertisementsProvider
}

func (s *SyncAdvertismentsHandler) Handle(c *fiber.Ctx) error {
	// TODO: Add custom validation logic.
	err := s.provider.SyncAdvertisments(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HandlerResponseNonOK)
	}
	return c.Status(fiber.StatusOK).JSON(dto.HandlerResponseOK)
}

func NewSyncAdvertismentsHandler(provider SyncAdvertisementsProvider) *SyncAdvertismentsHandler {
	if provider == nil {
		panic("sync advertisements provider is nil")
	}
	return &SyncAdvertismentsHandler{
		provider: provider,
	}
}
