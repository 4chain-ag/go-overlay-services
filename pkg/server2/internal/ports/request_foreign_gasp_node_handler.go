package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
)

// XBSVTopicHeader is the header key for the topic.
const XBSVTopicHeader = "X-BSV-Topic"

// RequestForeignGASPNodeService defines the interface for a service handling foreign GASP node requests.
type RequestForeignGASPNodeService interface {
	RequestForeignGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error)
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
	// Check for topic header
	topic := c.Get(XBSVTopicHeader)
	if topic == "" {
		return NewMissingXBSVTopicHeaderError()
	}

	// Parse request body
	var payload RequestForeignGASPNodePayload
	if err := c.BodyParser(&payload); err != nil {
		return NewInvalidRequestBodyError()
	}

	// Create outpoint
	outpoint := &overlay.Outpoint{
		OutputIndex: payload.OutputIndex,
	}
	txid, err := chainhash.NewHashFromHex(payload.TxID)
	if err != nil {
		return NewInvalidTxIDError()
	}
	outpoint.Txid = *txid

	// Create graphID
	graphID, err := overlay.NewOutpointFromString(payload.GraphID)
	if err != nil {
		return NewInvalidGraphIDError()
	}

	// Call service
	node, err := h.service.RequestForeignGASPNode(c.Context(), graphID, outpoint, topic)
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

// Error types for RequestForeignGASPNode handler
func NewMissingXBSVTopicHeaderError() app.Error {
	const msg = "The submitted request does not include required header: X-BSV-Topic"
	return app.NewIncorrectInputError(msg, msg)
}

func NewInvalidRequestBodyError() app.Error {
	const msg = "The submitted request body is invalid or malformed"
	return app.NewIncorrectInputError(msg, msg)
}

func NewInvalidTxIDError() app.Error {
	const msg = "The submitted txID is not a valid transaction hash"
	return app.NewIncorrectInputError(msg, msg)
}

func NewInvalidGraphIDError() app.Error {
	const msg = "The submitted graphID is not in a valid format (expected: txID.outputIndex)"
	return app.NewIncorrectInputError(msg, msg)
}

// RequestForeignGASPNodeInternalErrorResponse is the error response for internal errors.
var RequestForeignGASPNodeInternalErrorResponse = openapi.InternalServerErrorResponse{
	Message: "Unable to process foreign GASP node request due to an error in the overlay engine.",
}
