package commands

import (
	"context"
	"fmt"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
)

// StartGASPSyncProvider defines the contract for triggering GASP sync.
type StartGASPSyncProvider interface {
	StartGASPSync(ctx context.Context) error
}

// StartGASPSyncHandler handles the /admin/start-gasp-sync endpoint.
type StartGASPSyncHandler struct {
	provider StartGASPSyncProvider
}

// ResponseStartGASPNodeHandler is the standard success response body format.
type ResponseStartGASPNodeHandler struct {
	Message string `json:"message"`
}

// Handle initiates the sync and returns appropriate status.
func (h *StartGASPSyncHandler) Handle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	if r.Method != http.MethodPost {
		jsonutil.SendHTTPFailureResponse(w, http.StatusMethodNotAllowed, jsonutil.ReasonBadRequest, "method not allowed, only POST is supported")
		return
	}

	if err := h.provider.StartGASPSync(r.Context()); err != nil {
		jsonutil.SendHTTPFailureResponse(w, http.StatusInternalServerError, jsonutil.ReasonInternalError, "failed to start GASP sync")
		return
	}

	jsonutil.SendHTTPResponse(w, http.StatusOK, ResponseStartGASPNodeHandler{Message: "OK"})
}

// NewStartGASPSyncHandler constructs the handler.
func NewStartGASPSyncHandler(provider StartGASPSyncProvider) (*StartGASPSyncHandler, error) {
	if provider == nil {
		return nil, fmt.Errorf("StartGASPSyncProvider is nil")
	}
	return &StartGASPSyncHandler{provider: provider}, nil
}
