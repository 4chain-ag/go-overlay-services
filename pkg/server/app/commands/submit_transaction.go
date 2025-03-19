package commands

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/dto"
	"github.com/gofiber/fiber/v2"
)

type SubmitTransactionProvider interface {
	SubmitTransaction(ctx context.Context) error
}

type SubmitTransactionHandler struct {
	provider SubmitTransactionProvider
}

func (s *SubmitTransactionHandler) Handle(c *fiber.Ctx) error {
	// TODO: Add custom validation logic.
	err := s.provider.SubmitTransaction(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HandlerResponseOK)
	}
	return c.Status(fiber.StatusOK).JSON(dto.HandlerResponseOK)
}

func NewSubmitTransactionCommandHandler(provider SubmitTransactionProvider) *SubmitTransactionHandler {
	if provider == nil {
		panic("submit transaction provider is nil")
	}
	return &SubmitTransactionHandler{
		provider: provider,
	}
}
