package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// RequestForeignGASPNodeService defines the interface for a service responsible for
// requesting foreign GASP nodes.
type RequestForeignGASPNodeService interface {
	RequestForeignGASPNode(ctx context.Context, graphIDStr string, txIDStr string, outputIndex uint32, topic string) (*core.GASPNode, error)
}

// RequestForeignGASPNodeHandler handles incoming requests for foreign GASP nodes.
// It delegates to the RequestForeignGASPNodeService to process the request and formats
// the response according to the API spec.
type RequestForeignGASPNodeHandler struct {
	service RequestForeignGASPNodeService
}

// RequestForeignGASPNode processes an HTTP request to request a foreign GASP node.
// It extracts the topic from X-BSV-Topic header and parameters from JSON body,
// then returns the GASP node or an appropriate error response.
func (h *RequestForeignGASPNodeHandler) Handle(c *fiber.Ctx, params openapi.RequestForeignGASPNodeParams) error {
	// Parse JSON body
	var payload openapi.RequestForeignGASPNodeJSONBody
	if err := c.BodyParser(&payload); err != nil {
		return NewInvalidRequestBodyError()
	}

	node, err := h.service.RequestForeignGASPNode(
		c.Context(),
		payload.GraphID,
		payload.TxID,
		payload.OutputIndex,
		params.XBSVTopic,
	)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(NewRequestForeignGASPNodeSuccessResponse(node))
}

// NewRequestForeignGASPNodeHandler creates a new RequestForeignGASPNodeHandler with the given provider.
// If the provider is nil, it panics.
func NewRequestForeignGASPNodeHandler(provider app.RequestForeignGASPNodeProvider) *RequestForeignGASPNodeHandler {
	if provider == nil {
		panic("request foreign GASP node provider is nil")
	}
	return &RequestForeignGASPNodeHandler{service: app.NewRequestForeignGASPNodeService(provider)}
}

// NewRequestForeignGASPNodeSuccessResponse creates a success response for a foreign GASP node request.
func NewRequestForeignGASPNodeSuccessResponse(node *core.GASPNode) *core.GASPNode {
	return node
}

// NewInvalidRequestBodyError returns an Error indicating that the request body is invalid.
func NewInvalidRequestBodyError() app.Error {
	const str = "The submitted request body is invalid or malformed"
	return app.NewIncorrectInputError(str, str)
}
