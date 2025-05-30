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
	RequestForeignGASPNode(ctx context.Context, dto app.RequestForeignGASPNodeDTO) (*core.GASPNode, error)
}

// RequestForeignGASPNodeHandler handles incoming requests for foreign GASP nodes.
// It delegates to the RequestForeignGASPNodeService to process the request and formats
// the response according to the API spec.
type RequestForeignGASPNodeHandler struct {
	service RequestForeignGASPNodeService
}

// Handle processes an HTTP request to request a foreign GASP node.
// It extracts the topic from X-BSV-Topic header and parameters from JSON body,
// then returns the GASP node or an appropriate error response.
func (h *RequestForeignGASPNodeHandler) Handle(c *fiber.Ctx, params openapi.RequestForeignGASPNodeParams) error {
	var payload openapi.RequestForeignGASPNodeJSONBody
	if err := c.BodyParser(&payload); err != nil {
		return app.NewRequestForeignGASPNodeInvalidJSONError()
	}

	dto := app.RequestForeignGASPNodeDTO{
		GraphID:     payload.GraphID,
		TxID:        payload.TxID,
		OutputIndex: payload.OutputIndex,
		Topic:       params.XBSVTopic,
	}

	node, err := h.service.RequestForeignGASPNode(c.Context(), dto)
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
func NewRequestForeignGASPNodeSuccessResponse(node *core.GASPNode) openapi.GASPNode {
	var inputs map[string]any
	if len(node.Inputs) > 0 {
		inputs = make(map[string]any, len(node.Inputs))
		for k, v := range node.Inputs {
			inputs[k] = v
		}
	}
	//TODO: define the inputs type in the openapi spec to match the actual type

	graphID := ""
	if node.GraphID != nil {
		graphID = node.GraphID.String()
	}

	proof := ""
	if node.Proof != nil {
		proof = *node.Proof
	}

	ancillaryBeef := ""
	if node.AncillaryBeef != nil {
		ancillaryBeef = string(node.AncillaryBeef)
	}

	return openapi.GASPNode{
		GraphID:        graphID,
		RawTx:          node.RawTx,
		OutputIndex:    int(node.OutputIndex),
		Proof:          proof,
		TxMetadata:     node.TxMetadata,
		OutputMetadata: node.OutputMetadata,
		Inputs:         inputs,
		AncillaryBeef:  ancillaryBeef,
	}
}
