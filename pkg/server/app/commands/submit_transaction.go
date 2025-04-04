package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// XTopicsHeader defines the HTTP header key used for specifying transaction topics.
const XTopicsHeader = "x-topics"

// RequestBodyLimit1GB defines the maximum allowed size for request bodies (1GB).
const RequestBodyLimit1GB = 1000 * 1024 * 1024

var (
	// ErrMissingXTopicsHeader is returned when the required x-topics header is missing.
	ErrMissingXTopicsHeader = errors.New("missing x-topics header")

	// ErrInvalidXTopicsHeaderFormat is returned when the x-topics header has an invalid format.
	ErrInvalidXTopicsHeaderFormat = errors.New("invalid x-topics header format")

	// ErrInvalidHTTPMethod is returned when an unsupported HTTP method is used.
	ErrInvalidHTTPMethod = errors.New("invalid HTTP method")

	// ErrRequestBodyRead is returned when there's an error reading the request body.
	ErrRequestBodyRead = errors.New("failed to read request body")

	// ErrRequestBodyTooLarge is returned when the request body exceeds the size limit.
	ErrRequestBodyTooLarge = errors.New("request body too large")
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
	provider         SubmitTransactionProvider
	requestBodyLimit int64
}

// CreateTaggedBEEFFromRequest extracts the topics from the header and reads the body
// to create a TaggedBEEF object, with size limits applied.
func (s *SubmitTransactionHandler) CreateTaggedBEEFFromRequest(r *http.Request) (*overlay.TaggedBEEF, error) {
	header := r.Header.Get(XTopicsHeader)
	if header == "" {
		return nil, ErrMissingXTopicsHeader
	}

	var topics []string
	if err := json.Unmarshal([]byte(header), &topics); err != nil {
		return nil, ErrInvalidXTopicsHeaderFormat
	}

	reader := io.LimitReader(r.Body, s.requestBodyLimit)
	buff := make([]byte, 64*1024) // typically it's a best balance between performance and memory usage (loop and size).
	bytesRead := 0
	for {
		n, err := reader.Read(buff)
		bytesRead += n

		if bytesRead > int(s.requestBodyLimit) {
			return nil, ErrRequestBodyTooLarge
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, ErrRequestBodyRead
		}
	}

	return &overlay.TaggedBEEF{Beef: buff, Topics: topics}, nil
}

// Handle orchestrates the processing flow of a transaction. It prepares and
// sends a JSON response after invoking the engine and returns an HTTP response
// with the appropriate status code based on the engine's response.
func (s *SubmitTransactionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			fmt.Printf("Failed to close request body: %v\n", err)
		}
	}()

	if r.Method != http.MethodPost {
		http.Error(w, ErrInvalidHTTPMethod.Error(), http.StatusMethodNotAllowed)
		return
	}

	taggedBEEF, err := s.CreateTaggedBEEFFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	steakChan := make(chan *overlay.Steak, 1)
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
	case <-time.After(5 * time.Second): // Timeout to prevent hanging requests
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
	}
}

// NewSubmitTransactionCommandHandler returns an instance of a SubmitTransactionHandler, utilizing
// an implementation of SubmitTransactionProvider. If the provided argument is nil, it returns an error.
func NewSubmitTransactionCommandHandler(provider SubmitTransactionProvider) (*SubmitTransactionHandler, error) {
	if provider == nil {
		return nil, fmt.Errorf("submit transaction provider is nil")
	}
	return &SubmitTransactionHandler{
		provider:         provider,
		requestBodyLimit: RequestBodyLimit1GB,
	}, nil
}
