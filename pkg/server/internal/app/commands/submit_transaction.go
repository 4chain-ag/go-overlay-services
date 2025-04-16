package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// XTopicsHeader defines the HTTP header key used for specifying transaction topics.
const XTopicsHeader = "x-topics"

var (
	// ErrMissingXTopicsHeader is returned when the required x-topics header is missing.
	ErrMissingXTopicsHeader = errors.New("missing x-topics header")

	// ErrInvalidXTopicsHeaderFormat is returned when the x-topics header has an invalid format.
	ErrInvalidXTopicsHeaderFormat = errors.New("invalid x-topics header format")

	// ErrInvalidHTTPMethod is returned when an unsupported HTTP method is used.
	ErrInvalidHTTPMethod = errors.New("invalid HTTP method")
)

// SubmitTransactionHandlerResponse defines the response body content that
// will be sent in JSON format after successfully processing the handler logic.
type SubmitTransactionHandlerResponse struct {
	Steak overlay.Steak `json:"steak"`
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
	provider         SubmitTransactionProvider
	requestBodyLimit int64
	responseTimeout  time.Duration
}

func (s *SubmitTransactionHandler) createTaggedBEEF(body io.Reader, header http.Header) (*overlay.TaggedBEEF, error) {
	actual := header.Get(XTopicsHeader)
	if actual == "" {
		return nil, ErrMissingXTopicsHeader
	}

	// Parse topics from comma-separated list
	topics := strings.Split(actual, ",")

	// Basic validation - ensure we have at least one non-empty topic
	hasValidTopic := false
	for i, topic := range topics {
		topics[i] = strings.TrimSpace(topic) // Trim whitespace
		if topics[i] != "" {
			hasValidTopic = true
		}
	}

	if !hasValidTopic {
		return nil, ErrInvalidXTopicsHeaderFormat
	}

	reader := jsonutil.BodyReader{
		Body:      body,
		ReadLimit: s.requestBodyLimit,
	}

	bb, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("body reader failure: %w", err)
	}

	return &overlay.TaggedBEEF{Beef: bb, Topics: topics}, nil
}

// Handle orchestrates the processing flow of a transaction. It prepares and
// sends a JSON response after invoking the engine and returns an HTTP response
// with the appropriate status code based on the engine's response.
func (s *SubmitTransactionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrInvalidHTTPMethod.Error(), http.StatusMethodNotAllowed)
		return
	}

	taggedBEEF, err := s.createTaggedBEEF(r.Body, r.Header)
	if errors.Is(err, jsonutil.ErrRequestBodyTooLarge) {
		http.Error(w, jsonutil.ErrRequestBodyTooLarge.Error(), http.StatusRequestEntityTooLarge)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	steakChan := make(chan *overlay.Steak, 1)
	defer close(steakChan)

	_, err = s.provider.Submit(r.Context(), *taggedBEEF, engine.SubmitModeCurrent, func(steak *overlay.Steak) {
		steakChan <- steak
	})

	if err != nil {
		jsonutil.SendHTTPResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	select {
	case steak := <-steakChan:
		jsonutil.SendHTTPResponse(w, http.StatusOK, SubmitTransactionHandlerResponse{Steak: *steak})
	case <-time.After(s.responseTimeout):
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
	}
}

// SubmitTransactionHandlerOption defines a function that can configure a SubmitTransactionHandler.
type SubmitTransactionHandlerOption func(h *SubmitTransactionHandler)

// WithResponseTime configures the timeout duration for a response from the transaction submission.
func WithResponseTime(d time.Duration) SubmitTransactionHandlerOption {
	return func(h *SubmitTransactionHandler) {
		h.responseTimeout = d
	}
}

// WithRequestBodyLimit configures the maximum allowed size for request bodies.
func WithRequestBodyLimit(limit int64) SubmitTransactionHandlerOption {
	return func(h *SubmitTransactionHandler) {
		h.requestBodyLimit = limit
	}
}

// NewSubmitTransactionCommandHandler returns an instance of a SubmitTransactionHandler, utilizing
// an implementation of SubmitTransactionProvider. If the provided argument is nil, it returns an error.
func NewSubmitTransactionCommandHandler(provider SubmitTransactionProvider, opts ...SubmitTransactionHandlerOption) (*SubmitTransactionHandler, error) {
	if provider == nil {
		return nil, fmt.Errorf("submit transaction provider is nil")
	}

	h := SubmitTransactionHandler{
		provider:         provider,
		requestBodyLimit: jsonutil.RequestBodyLimit1GB,
		responseTimeout:  10 * time.Second,
	}
	for _, o := range opts {
		o(&h)
	}
	return &h, nil
}
