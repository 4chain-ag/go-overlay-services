package commands

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
)

var (
	// ErrMissingService is returned when the required service field is missing.
	ErrMissingService = errors.New("missing required field: service")

	// ErrInvalidRequestBody is returned when the request body cannot be parsed.
	ErrInvalidRequestBody = errors.New("invalid request body")

	// ErrMethodNotAllowed is returned when an unsupported HTTP method is used.
	ErrMethodNotAllowed = errors.New("method not allowed")
)

// LookupHandlerResponse defines the response body content that
// will be sent in JSON format after successfully processing the handler logic.
type LookupHandlerResponse struct {
	*lookup.LookupAnswer `json:"answer"`
}

// LookupQuestionProvider defines the contract that must be fulfilled
// to process lookup questions in the overlay engine.
type LookupQuestionProvider interface {
	Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error)
}

// LookupHandler orchestrates the processing flow of a lookup question,
// including the request body validation, converting the request body
// into an overlay-engine-compatible format, and applying any other necessary
// logic before invoking the engine.
type LookupHandler struct {
	provider LookupQuestionProvider
}

// Handle orchestrates the processing flow of a lookup request. It validates the
// request body, passes it to the engine's lookup method, and returns the result.
func (h *LookupHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	var question lookup.LookupQuestion
	if err := jsonutil.DecodeRequestBody(r, &question); err != nil {
		http.Error(w, fmt.Sprintf("%s: %s", ErrInvalidRequestBody.Error(), err.Error()), http.StatusBadRequest)
		return
	}

	if question.Service == "" {
		http.Error(w, ErrMissingService.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.provider.Lookup(r.Context(), &question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonutil.SendHTTPResponse(w, http.StatusOK, LookupHandlerResponse{LookupAnswer: result})
}

// NewLookupHandler returns an instance of a LookupHandler, utilizing
// an implementation of LookupQuestionProvider. If the provided argument is nil, it returns an error.
func NewLookupHandler(provider LookupQuestionProvider) (*LookupHandler, error) {
	if provider == nil {
		return nil, fmt.Errorf("lookup question provider is nil")
	}
	return &LookupHandler{
		provider: provider,
	}, nil
}
