package commands

import (
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp"
	"github.com/gofiber/fiber/v2"
)

// RequestForeignGASPNodeProvider defines the contract that must be fulfilled to send a requestForeignGASPNode to the overlay engine.
type RequestForeignGASPNodeProvider interface {
	ProvideForeignGASPNode(graphID string, txID string, outputIndex uint32) (*gasp.GASPNode, error)
}

// RequestForeignGASPNodeHandler orchestrates the requestForeignGASPNode flow.
type RequestForeignGASPNodeHandler struct {
	provider RequestForeignGASPNodeProvider
}

// requestPayload models the incoming request body.
type requestPayload struct {
	GraphID     string `json:"graphID"`
	TxID        string `json:"txid"`
	OutputIndex uint32 `json:"outputIndex"`
}

// Handle handles the request, validates input, calls the engine, and returns the GASP node.
func (h *RequestForeignGASPNodeHandler) Handle(c *fiber.Ctx) error {
	var payload requestPayload
	if err := c.BodyParser(&payload); err != nil {
		if err := c.Status(fiber.StatusBadRequest).JSON(nil); err != nil {
			return fmt.Errorf("failed to send response: %w", err)
		}
	}

	node, err := h.provider.ProvideForeignGASPNode(payload.GraphID, payload.TxID, payload.OutputIndex)
	if err != nil {
		if err := c.Status(fiber.StatusInternalServerError).JSON(nil); err != nil {
			return fmt.Errorf("failed to send response: %w", err)
		}
	}

	if err := c.Status(fiber.StatusOK).JSON(node); err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}
	return nil
}

// NewRequestForeignGASPNodeHandler creates a new handler instance.
func NewRequestForeignGASPNodeHandler(provider RequestForeignGASPNodeProvider) *RequestForeignGASPNodeHandler {
	if provider == nil {
		return nil
	}
	return &RequestForeignGASPNodeHandler{provider: provider}
}
