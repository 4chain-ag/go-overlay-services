package commands

import (
	"encoding/json"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
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

// Handle processes the HTTP request and writes the appropriate response.
func (h *RequestForeignGASPNodeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload requestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	node, err := h.provider.ProvideForeignGASPNode(payload.GraphID, payload.TxID, payload.OutputIndex)
	if err != nil {
		jsonutil.SendHTTPInternalServerErrorTextResponse(w)
		return
	}

	jsonutil.SendHTTPResponse(w, http.StatusOK, node)
}

// NewRequestForeignGASPNodeHandler creates a new handler instance.
func NewRequestForeignGASPNodeHandler(provider RequestForeignGASPNodeProvider) *RequestForeignGASPNodeHandler {
	if provider == nil {
		return nil
	}
	return &RequestForeignGASPNodeHandler{provider: provider}
}
