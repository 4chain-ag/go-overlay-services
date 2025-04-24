package overlayhttp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
)

// RequestBodyLimit1GB defines the maximum allowed size for request bodies (1GB).
const RequestBodyLimit1GB = 1000 * 1024 * 1024

var (
	// ErrRequestBodyRead is returned when there's an error reading the request body.
	ErrRequestBodyRead = errors.New("failed to read request body")

	// ErrRequestBodyTooLarge is returned when the request body exceeds the size limit.
	ErrRequestBodyTooLarge = errors.New("request body too large")

	// ErrInvalidXTopicsHeaderFormat is returned when the x-topics header has an invalid format.
	ErrInvalidXTopicsHeaderFormat = errors.New("invalid x-topics header format")
)

// SubmitTransactionProvider defines the contract that must be fulfilled
// to send a transaction request to the overlay engine for further processing.
type SubmitTransactionProvider interface {
	Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error)
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

// SubmitTransactionHandler orchestrates the processing flow of a transaction
// request, including the request body validation, converting the request body
// into an overlay-engine-compatible format, and applying any other necessary
// logic before invoking the engine.
type SubmitTransactionHandler struct {
	provider         SubmitTransactionProvider
	requestBodyLimit int64
	responseTimeout  time.Duration
}

// Handle orchestrates the processing flow of a transaction. It prepares and
// sends a JSON response after invoking the engine and returns an HTTP response
// with the appropriate status code based on the engine's response.
func (s *SubmitTransactionHandler) Handle(c *fiber.Ctx, params openapi.SubmitTransactionParams) error {
	taggedBEEF, err := s.createTaggedBEEF(c.Request().Body(), params.XTopics)
	if errors.Is(err, ErrRequestBodyTooLarge) {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.StatusRequestEntityTooLarge)
	}
	if err != nil {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(err.Error())
	}

	steakChan := make(chan *overlay.Steak, 1)
	_, err = s.provider.Submit(c.UserContext(), *taggedBEEF, engine.SubmitModeCurrent, func(steak *overlay.Steak) {
		steakChan <- steak
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.ErrInternalServerError)
	}

	select {
	case steak := <-steakChan:
		return c.Status(fiber.StatusOK).JSON(toOpenAPISubmitTransactionResponse(steak))
	case <-time.After(s.responseTimeout):
		return c.Status(fiber.StatusRequestTimeout).JSON(fiber.StatusRequestTimeout)
	}
}

func (s *SubmitTransactionHandler) createTaggedBEEF(body []byte, topics []string) (*overlay.TaggedBEEF, error) {
	for i, topic := range topics {
		topics[i] = strings.TrimSpace(topic)
		if topics[i] == "" {
			return nil, ErrInvalidXTopicsHeaderFormat
		}
	}

	reader := io.LimitReader(bytes.NewBuffer(body), s.requestBodyLimit+1)
	buff := make([]byte, 64*1024)
	var dst bytes.Buffer
	var bytesRead int64

	for {
		n, err := reader.Read(buff)
		if n > 0 {
			bytesRead += int64(n)
			if bytesRead > s.requestBodyLimit {
				return nil, ErrRequestBodyTooLarge
			}

			if _, inner := dst.Write(buff[:n]); inner != nil {
				return nil, ErrRequestBodyRead
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, ErrRequestBodyRead
		}
	}

	return &overlay.TaggedBEEF{Beef: dst.Bytes(), Topics: topics}, nil
}

// NewSubmitTransactionHandler returns an instance of a SubmitTransactionHandler, utilizing
// an implementation of SubmitTransactionProvider. If the provider argument is nil, it triggers a panic.
func NewSubmitTransactionHandler(provider SubmitTransactionProvider, options ...SubmitTransactionHandlerOption) *SubmitTransactionHandler {
	if provider == nil {
		panic("submit transaction provider is nil")
	}

	handler := SubmitTransactionHandler{
		provider:         provider,
		requestBodyLimit: RequestBodyLimit1GB,
		responseTimeout:  10 * time.Second,
	}

	for _, opt := range options {
		opt(&handler)
	}

	return &handler
}

func toOpenAPISubmitTransactionResponse(steak *overlay.Steak) *openapi.SubmitTransactionResponse {
	if steak == nil {
		return &openapi.SubmitTransactionResponse{
			STEAK: make(openapi.STEAK),
		}
	}

	response := openapi.SubmitTransactionResponse{
		STEAK: make(openapi.STEAK, len(*steak)),
	}

	for key, instructions := range *steak {
		ancillaryIDs := make([]string, 0, len(instructions.AncillaryTxids))
		for _, id := range instructions.AncillaryTxids {
			ancillaryIDs = append(ancillaryIDs, id.String())
		}

		response.STEAK[key] = openapi.AdmittanceInstructions{
			AncillaryTxIDs: ancillaryIDs,
			CoinsRemoved:   instructions.CoinsRemoved,
			CoinsToRetain:  instructions.CoinsToRetain,
			OutputsToAdmit: instructions.OutputsToAdmit,
		}
	}
	return &response
}
