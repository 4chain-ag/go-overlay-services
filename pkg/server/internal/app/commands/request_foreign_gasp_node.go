package commands

import (
	"context"
	"fmt"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// RequestForeignGASPNodeProvider defines the contract that must be fulfilled to send a requestForeignGASPNode to the overlay engine.
type RequestForeignGASPNodeProvider interface {
	ProvideForeignGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint) (*core.GASPNode, error)
}

// RequestForeignGASPNodeHandler orchestrates the requestForeignGASPNode flow.
type RequestForeignGASPNodeHandler struct {
	provider RequestForeignGASPNodeProvider
}

// RequestForeignGASPNodeHandlerPayload models the incoming request body.
type RequestForeignGASPNodeHandlerPayload struct {
	GraphID     string `json:"graphID"`
	TxID        string `json:"txID"`
	OutputIndex uint32 `json:"outputIndex"`
}

// Handle processes the HTTP request and writes the appropriate response.
func (h *RequestForeignGASPNodeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload RequestForeignGASPNodeHandlerPayload
	if err := jsonutil.DecodeRequestBody(r, &payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	outpoint := &overlay.Outpoint{
		OutputIndex: payload.OutputIndex,
	}
	txid, err := chainhash.NewHashFromHex(payload.TxID)
	if err != nil {
		http.Error(w, "invalid txid", http.StatusBadRequest)
		return
	} else {
		outpoint.Txid = *txid
	}
	graphId, err := overlay.NewOutpointFromString(payload.GraphID)
	if err != nil {
		http.Error(w, "invalid graphID", http.StatusBadRequest)
		return
	}
	node, err := h.provider.ProvideForeignGASPNode(r.Context(), graphId, outpoint)
	if err != nil {
		jsonutil.SendHTTPInternalServerErrorTextResponse(w)
		return
	}

	jsonutil.SendHTTPResponse(w, http.StatusOK, node)
}

// NewRequestForeignGASPNodeHandler creates a new handler instance.
func NewRequestForeignGASPNodeHandler(provider RequestForeignGASPNodeProvider) (*RequestForeignGASPNodeHandler, error) {
	if provider == nil {
		return nil, fmt.Errorf("request foreign gasp node provider is nil")
	}
	return &RequestForeignGASPNodeHandler{provider: provider}, nil
}
