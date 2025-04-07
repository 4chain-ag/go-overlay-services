package commands

import (
	"context"
	"fmt"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// SubmitTransactionHandlerResponse defines the response body content that
// will be sent in JSON format after successfully processing the handler logic.
type SubmitTransactionHandlerResponse struct {
	overlay.Steak `json:"steak"`
}

// SubmitTransactionProvider defines the contract that must be fulfilled
// to send a transaction request to the overlay engine for further processing.
// Note: The contract definition is still in development and will be updated after
// migrating the engine code.
type SubmitTransactionProvider interface {
	Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error)
}

// SubmitTransactionHandler orchestrates the processing flow of a transaction
// request, including the request body validation, converting the request body
// into an overlay-engine-compatible format, and applying any other necessary
// logic before invoking the engine.
type SubmitTransactionHandler struct {
	provider SubmitTransactionProvider
}

// Handle orchestrates the processing flow of a transaction. It prepares and
// sends a JSON response after invoking the engine and returns an HTTP response
// with the appropriate status code based on the engine's response.
func (s *SubmitTransactionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	// TODO: Add custom validation logic.
	steak, err := s.provider.Submit(r.Context(), overlay.TaggedBEEF{}, engine.SubmitModeCurrent, func(steak *overlay.Steak) {})
	if err != nil {
		jsonutil.SendHTTPFailureResponse(
			w,
			http.StatusInternalServerError,
			jsonutil.ReasonInternalError,
			"failed to submit transaction to overlay engine",
		)
		return
	}

	jsonutil.SendHTTPResponse(w, http.StatusCreated, SubmitTransactionHandlerResponse{Steak: steak})
}

// NewSubmitTransactionCommandHandler returns an instance of a SubmitTransactionHandler, utilizing
// an implementation of SubmitTransactionProvider. If the provided argument is nil, it triggers a panic.
func NewSubmitTransactionCommandHandler(provider SubmitTransactionProvider) (*SubmitTransactionHandler, error) {
	if provider == nil {
		return nil, fmt.Errorf("submit transaction provider is nil")
	}
	return &SubmitTransactionHandler{
		provider: provider,
	}, nil
}
