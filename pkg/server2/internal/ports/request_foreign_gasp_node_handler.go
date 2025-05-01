package ports

import (
	"context"
	"errors"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
)

// RequestForeignGASPNodeService abstracts the logic for handling foreign GASP node requests.

type RequestForeignGASPNodeService interface {
	RequestForeignGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error)
}

// RequestForeignGASPNodePayload represents the request body for requesting a foreign GASP node.

type RequestForeignGASPNodePayload struct {
	GraphID string `json:"graphID"`

	TxID string `json:"txID"`

	OutputIndex uint32 `json:"outputIndex"`
}

// RequestForeignGASPNodeHandler orchestrates the processing flow of a foreign GASP node

// request and applies any necessary logic before invoking the engine.

type RequestForeignGASPNodeHandler struct {
	service RequestForeignGASPNodeService
}

// Handle orchestrates the processing flow of a foreign GASP node request.

// It prepares and sends a JSON response after invoking the engine and returns an HTTP response

// with the appropriate status code based on the engine's response.

func (h *RequestForeignGASPNodeHandler) Handle(c *fiber.Ctx) error {

	topics := c.Get("X-BSV-Topic")

	if topics == "" {

		return c.Status(fiber.StatusBadRequest).JSON(openapi.Error{

			Message: "Missing 'X-BSV-Topic' header",
		})

	}

	var payload RequestForeignGASPNodePayload

	if err := c.BodyParser(&payload); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(openapi.Error{

			Message: "Invalid request body",
		})

	}

	outpoint := &overlay.Outpoint{

		OutputIndex: payload.OutputIndex,
	}

	txid, err := chainhash.NewHashFromHex(payload.TxID)

	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(openapi.Error{

			Message: "Invalid txID format",
		})

	}

	outpoint.Txid = *txid

	graphID, err := overlay.NewOutpointFromString(payload.GraphID)

	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(openapi.Error{

			Message: "Invalid graphID format",
		})

	}

	node, err := h.service.RequestForeignGASPNode(c.Context(), graphID, outpoint, topics)

	if err != nil {

		if errors.Is(err, app.ErrRequestForeignGASPNodeProvider) {

			return c.Status(fiber.StatusInternalServerError).JSON(openapi.Error{

				Message: "Unable to process foreign GASP node request due to issues with the overlay engine",
			})

		}

		return c.Status(fiber.StatusInternalServerError).JSON(openapi.Error{

			Message: "Internal server error",
		})

	}

	return c.Status(fiber.StatusOK).JSON(node)

}

// NewRequestForeignGASPNodeHandler returns an instance of a RequestForeignGASPNodeHandler,

// utilizing an implementation of RequestForeignGASPNodeProvider. If the provided argument is nil, it triggers a panic.

func NewRequestForeignGASPNodeHandler(provider app.RequestForeignGASPNodeProvider) *RequestForeignGASPNodeHandler {

	if provider == nil {

		panic("request foreign GASP node provider is nil")

	}

	return &RequestForeignGASPNodeHandler{service: app.NewRequestForeignGASPNodeService(provider)}

}
