package commands

import (
	"context"
	"fmt"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
)

// ForeignSyncResponseProvider defines the contract for providing a foreign sync response.
type ForeignSyncResponseProvider interface {
	ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error)
}

// RequestSyncResponseHandler orchestrates the /request-sync-response flow.
type RequestSyncResponseHandler struct {
	provider ForeignSyncResponseProvider
}

// NewRequestSyncResponseHandler creates a new instance of the handler.
func NewRequestSyncResponseHandler(provider ForeignSyncResponseProvider) (*RequestSyncResponseHandler, error) {
	if provider == nil {
		return nil, fmt.Errorf("foreign sync response provider cannot be nil")
	}
	return &RequestSyncResponseHandler{provider: provider}, nil
}

// Handle processes the HTTP POST /request-sync-response request.
func (h *RequestSyncResponseHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	topics := r.Header["X-Bsv-Topic"]
	if len(topics) == 0 {
		http.Error(w, "missing 'topic' header", http.StatusBadRequest)
		return
	}
	topic := topics[0]

	var initialRequest core.GASPInitialRequest
	if err := jsonutil.DecodeRequestBody(r, &initialRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.provider.ProvideForeignSyncResponse(r.Context(), &initialRequest, topic)
	if err != nil {
		jsonutil.SendHTTPInternalServerErrorTextResponse(w)
		return
	}

	jsonutil.SendHTTPResponse(w, http.StatusOK, resp)
}
