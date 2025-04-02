package commands

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

const (
	// Error messages
	errMsgMissingTopicsHeader = "Missing x-topics header"
	errMsgInvalidTopicsFormat = "Invalid x-topics header format"
	errMsgFailedToReadBody    = "Failed to read request body"
	errMsgMethodNotAllowed    = "Method not allowed"
	
	// Header keys
	headerTopics = "x-topics"
)

// SubmitTransactionHandlerResponse defines the response body content that
// will be sent in JSON format after successfully processing the handler logic.
type SubmitTransactionHandlerResponse struct {
	overlay.Steak `json:"steak"`
}

// SubmitTransactionProvider defines the contract that must be fulfilled
// to send a transaction request to the overlay engine for further processing.
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
	// NOTE: comment place holders will be removed before code is merged
	// 1. Validate HTTP method
	if r.Method != http.MethodPost {
		http.Error(w, errMsgMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// 2. Extract and validate the x-topics header
	topicsHeader := r.Header.Get(headerTopics)
	if topicsHeader == "" {
		http.Error(w, errMsgMissingTopicsHeader, http.StatusBadRequest)
		return
	}

	// 3. Parse the topics header as JSON
	var topics []string
	if err := json.Unmarshal([]byte(topicsHeader), &topics); err != nil {
		http.Error(w, errMsgInvalidTopicsFormat, http.StatusBadRequest)
		return
	}

	// 4. Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, errMsgFailedToReadBody, http.StatusBadRequest)
		return
	}

	// 5. Create the TaggedBEEF object
	taggedBEEF := overlay.TaggedBEEF{
		Beef:   body,
		Topics: topics,
	}

	// 6. Set up synchronization for handling async callbacks
	responseSent := false
	var responseSync sync.Mutex

	// Define the callback function - will be called when STEAK is ready
	onSteakReady := func(steak *overlay.Steak) {
		responseSync.Lock()
		defer responseSync.Unlock()
		
		if !responseSent {
			responseSent = true
			jsonutil.SendHTTPResponse(w, http.StatusOK, SubmitTransactionHandlerResponse{Steak: *steak})
		}
	}

	// 7. Submit the transaction
	steak, err := s.provider.Submit(r.Context(), taggedBEEF, engine.SubmitModeCurrent, onSteakReady)
	if err != nil {
		responseSync.Lock()
		defer responseSync.Unlock()
		
		if !responseSent {
			responseSent = true
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// 8. If the callback hasn't been triggered yet, send the response immediately
	responseSync.Lock()
	defer responseSync.Unlock()
	
	if !responseSent {
		responseSent = true
		jsonutil.SendHTTPResponse(w, http.StatusOK, SubmitTransactionHandlerResponse{Steak: steak})
	}
}

// NewSubmitTransactionCommandHandler returns an instance of a SubmitTransactionHandler, utilizing
// an implementation of SubmitTransactionProvider. If the provided argument is nil, it triggers a panic.
func NewSubmitTransactionCommandHandler(provider SubmitTransactionProvider) *SubmitTransactionHandler {
	if provider == nil {
		panic("submit transaction provider is nil")
	}
	return &SubmitTransactionHandler{
		provider: provider,
	}
}
