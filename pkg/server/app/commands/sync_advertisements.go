package commands

import (
	"context"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
)

// SyncAdvertisementsHandlerResponse defines the response body content that
// will be sent in JSON format after successfully processing the handler logic.
type SyncAdvertisementsHandlerResponse struct {
	Message string `json:"message"`
}

// SyncAdvertisementsProvider defines the contract that must be fulfilled
// to send synchronize advertisements request to the overlay engine for further processing.
type SyncAdvertisementsProvider interface {
	SyncAdvertisements(ctx context.Context) error
}

// SyncAdvertisementsHandler orchestrates the processing flow of a synchronize advertisements
// request and applies any necessary logic before invoking the engine.
type SyncAdvertisementsHandler struct {
	provider SyncAdvertisementsProvider
}

// Handle orchestrates the processing flow of a synchronize advertisements request.
// It prepares and sends a JSON response after invoking the engine and returns an HTTP response
// with the appropriate status code based on the engine's response.
func (s *SyncAdvertisementsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	err := s.provider.SyncAdvertisements(r.Context())
	if err != nil {
		jsonutil.SendHTTPInternalServerErrorTextResponse(w)
		return
	}

	jsonutil.SendHTTPResponse(w, http.StatusOK, SyncAdvertisementsHandlerResponse{Message: "OK"})
}

// NewSyncAdvertisementsCommandHandler returns an instance of a SyncAdvertisementsHandler,
// utilizing an implementation of SyncAdvertisementsProvider.
// If the provided argument is nil, it triggers a panic.
func NewSyncAdvertisementsCommandHandler(provider SyncAdvertisementsProvider) *SyncAdvertisementsHandler {
	if provider == nil {
		panic("sync advertisements provider is nil")
	}
	return &SyncAdvertisementsHandler{
		provider: provider,
	}
}
