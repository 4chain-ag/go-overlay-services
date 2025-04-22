package commands

import (
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
)

// ARCIngestHandler is a temporary test handler implementation used as a stub
// for the overlay HTTP server. It will be removed after mergning PR #119.
type ARCIngestHandler struct {
}

// Handle is a no-op call that returns a StatusInternalServerError response if the ARC API key was not set
// during initialization. Otherwise, it returns a plain text response with HTTP status OK (200).
func (a *ARCIngestHandler) Handle(w http.ResponseWriter, r *http.Request) {
	jsonutil.SendHTTPResponse(w, http.StatusOK, http.StatusText(http.StatusOK))
}

// NewARCIngestHandler returns a new instance of ARCIngestHandler.
func NewARCIngestHandler() (*ARCIngestHandler, error) {
	return &ARCIngestHandler{}, nil
}
