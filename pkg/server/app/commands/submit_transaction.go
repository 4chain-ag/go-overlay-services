package commands

import (
	"context"
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/dto"
	"github.com/gofiber/fiber/v2"
)

// SubmitTransactionProvider defines the contract that must be fulfilled
// to send a transaction request to the overlay engine for further processing.
// Note: The contract definition is still in development and will be updated after
// migrating the engine code.
type SubmitTransactionProvider interface {
	SubmitTransaction(ctx context.Context) error
}

// SubmitTransactionHandler orchestrates the processing flow of a transaction
// request, including the request body validation, converting the request body
// into an overlay-engine-compatible format, and applying any other necessary
// logic before invoking the engine.
type SubmitTransactionHandler struct {
	provider SubmitTransactionProvider
}

// Handle orchestrates the processing flow of a transaction. It prepares and
// sends a JSON response after invoking the engine and returns an HTTP response
// with the appropriate status code based on the engine's response.
func (s *SubmitTransactionHandler) Handle(c *fiber.Ctx) error {
	err := s.provider.SubmitTransaction(c.Context())
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

// NewSubmitTransactionCommandHandler returns an instance of a SubmitTransactionHandler, utilizing
// an implementation of SubmitTransactionProvider. If the provided argument is nil, it triggers a panic.
func NewSubmitTransactionCommandHandler(provider SubmitTransactionProvider) *SubmitTransactionHandler {
	if provider == nil {
		panic("submit transaction provider is nil")
	}
	return &SubmitTransactionHandler{
		provider: provider,
	}
}
