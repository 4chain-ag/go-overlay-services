package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// RequestForeignGASPNodeService defines the interface for a service responsible for
// requesting foreign GASP nodes. It encapsulates the business logic for resolving node data.
type RequestForeignGASPNodeService interface {
	// RequestForeignGASPNode processes the given request and returns the corresponding GASP node
	// or an error if the request cannot be fulfilled.
	RequestForeignGASPNode(ctx context.Context, dto app.RequestForeignGASPNodeDTO) (*core.GASPNode, error)
}

// RequestForeignGASPNodeHandler handles HTTP requests for foreign GASP nodes.
// It parses input, delegates to the service layer, and formats the response.
type RequestForeignGASPNodeHandler struct {
	service RequestForeignGASPNodeService
}

// Handle processes an incoming request for a foreign GASP node.
// It parses the request body and parameters, delegates the request to the service,
// and returns a formatted JSON response or an appropriate error.
func (h *RequestForeignGASPNodeHandler) Handle(c *fiber.Ctx, params openapi.RequestForeignGASPNodeParams) error {
	var body openapi.RequestForeignGASPNodeJSONBody
	err := c.BodyParser(&body)
	if err != nil {
		return NewRequestBodyParserError(err)
	}

	node, err := h.service.RequestForeignGASPNode(c.Context(), app.RequestForeignGASPNodeDTO{
		GraphID:     body.GraphID,
		TxID:        body.TxID,
		OutputIndex: body.OutputIndex,
		Topic:       params.XBSVTopic,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(NewRequestForeignGASPNodeSuccessResponse(node))
}

// NewRequestForeignGASPNodeHandler creates and returns a new RequestForeignGASPNodeHandler
// using the provided RequestForeignGASPNodeProvider. It panics if the provider is nil.
func NewRequestForeignGASPNodeHandler(provider app.RequestForeignGASPNodeProvider) *RequestForeignGASPNodeHandler {
	if provider == nil {
		panic("request foreign GASP node provider is nil")
	}
	return &RequestForeignGASPNodeHandler{service: app.NewRequestForeignGASPNodeService(provider)}
}

// NewRequestForeignGASPNodeSuccessResponse constructs a success response from the given GASPNode.
// It maps internal types to the OpenAPI response format, ensuring compatibility with the API spec.
func NewRequestForeignGASPNodeSuccessResponse(node *core.GASPNode) openapi.GASPNode {
	var inputs map[string]any
	if len(node.Inputs) > 0 {
		inputs = make(map[string]any, len(node.Inputs))
		for k, v := range node.Inputs {
			inputs[k] = v
		}
	}

	var graphID string
	if node.GraphID != nil {
		graphID = node.GraphID.String()
	}

	var proof string
	if node.Proof != nil {
		proof = *node.Proof
	}

	return openapi.GASPNode{
		GraphID:        graphID,
		RawTx:          node.RawTx,
		OutputIndex:    node.OutputIndex,
		Proof:          proof,
		TxMetadata:     node.TxMetadata,
		OutputMetadata: node.OutputMetadata,
		Inputs:         inputs,
		AncillaryBeef:  node.AncillaryBeef,
	}
}
