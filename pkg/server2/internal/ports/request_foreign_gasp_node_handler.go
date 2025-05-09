package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// XBSVTopicHeader is the header key for the topic.

const XBSVTopicHeader = "X-BSV-Topic"

// RequestForeignGASPNodeService defines the interface for a service handling foreign GASP node requests.

type RequestForeignGASPNodeService interface {
	RequestForeignGASPNodeWithStrings(ctx context.Context, graphIDStr string, txIDStr string, outputIndex uint32, topic string) (*core.GASPNode, error)
}

// RequestForeignGASPNodeHandler handles requests for foreign GASP nodes.

type RequestForeignGASPNodeHandler struct {
	service RequestForeignGASPNodeService
}

// RequestForeignGASPNodePayload represents the request payload.

type RequestForeignGASPNodePayload struct {
	GraphID string `json:"graphID"`

	TxID string `json:"txID"`

	OutputIndex uint32 `json:"outputIndex"`
}

// Handle processes requests for foreign GASP nodes.

func (h *RequestForeignGASPNodeHandler) Handle(c *fiber.Ctx) error {

	// Get topic from header

	topic := c.Get(XBSVTopicHeader)

	// Parse request body

	var payload RequestForeignGASPNodePayload

	if err := c.BodyParser(&payload); err != nil {

		return app.NewIncorrectInputError("Invalid request body", "The submitted request body is invalid or malformed")

	}

	// Call service with string parameters - service will handle validation and conversion

	node, err := h.service.RequestForeignGASPNodeWithStrings(

		c.Context(),

		payload.GraphID,

		payload.TxID,

		payload.OutputIndex,

		topic,
	)

	if err != nil {

		return err

	}

	return c.Status(fiber.StatusOK).JSON(node)

}

// NewRequestForeignGASPNodeHandler creates a new handler instance.

func NewRequestForeignGASPNodeHandler(provider app.RequestForeignGASPNodeProvider) *RequestForeignGASPNodeHandler {

	if provider == nil {

		panic("request foreign GASP node provider is nil")

	}

	return &RequestForeignGASPNodeHandler{

		service: app.NewRequestForeignGASPNodeService(provider),
	}

}

// RequestForeignGASPNodeInternalErrorResponse is the error response for internal errors.

var RequestForeignGASPNodeInternalErrorResponse = openapi.InternalServerErrorResponse{

	Message: "Unable to process foreign GASP node request due to an error in the overlay engine.",
}
