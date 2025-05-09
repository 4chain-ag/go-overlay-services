package ports

import (
	"context"
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
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
	GraphID     string `json:"graphID"`
	TxID        string `json:"txID"`
	OutputIndex uint32 `json:"outputIndex"`
}

// Handle processes requests for foreign GASP nodes.
func (h *RequestForeignGASPNodeHandler) Handle(c *fiber.Ctx) error {
	topic := c.Get(XBSVTopicHeader)
	if topic == "" {
		return NewMissingXBSVTopicHeaderError()
	}

	var payload RequestForeignGASPNodePayload
	if err := c.BodyParser(&payload); err != nil {
		return NewInvalidRequestBodyError()
	}

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

// NewMissingXBSVTopicHeaderError returns an Error indicating that the X-BSV-Topic header is missing.
func NewMissingXBSVTopicHeaderError() app.Error {
	str := fmt.Sprintf("The submitted request does not include required header: %s.", XBSVTopicHeader)
	return app.NewIncorrectInputError(str, str)
}

// NewInvalidRequestBodyError returns an Error indicating that the request body is invalid.
func NewInvalidRequestBodyError() app.Error {
	const msg = "The submitted request body is invalid or malformed"
	return app.NewIncorrectInputError(msg, msg)
}
