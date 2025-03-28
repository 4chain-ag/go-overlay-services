package commands

import (
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
)

// StartGASPSyncProvider defines the contract for triggering GASP sync.
type StartGASPSyncProvider interface {
	StartGASPSync() error
}

// StartGASPSyncHandler handles the /admin/start-gasp-sync endpoint.
type StartGASPSyncHandler struct {
	provider StartGASPSyncProvider
}

// HandlerResponse is the standard response body format.
type HandlerResponse struct {
	Message string `json:"message"`
}

// Handle initiates the sync and returns appropriate status.
func (h *StartGASPSyncHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.provider.StartGASPSync(); err != nil {
		jsonutil.SendHTTPResponse(w, http.StatusInternalServerError, HandlerResponse{Message: "FAILED"})
		return
	}

	jsonutil.SendHTTPResponse(w, http.StatusOK, HandlerResponse{Message: "OK"})
}

// NewStartGASPSyncHandler constructs the handler.
func NewStartGASPSyncHandler(provider StartGASPSyncProvider) *StartGASPSyncHandler {
	if provider == nil {
		return nil
	}
	return &StartGASPSyncHandler{provider: provider}
}
