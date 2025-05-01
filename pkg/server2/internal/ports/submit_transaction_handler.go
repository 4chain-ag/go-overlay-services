package ports

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
)

// XTopicsHeader defines the HTTP header key used for specifying transaction topics.
const XTopicsHeader = "x-topics"

// SubmitTransactionHandlerOption defines a function that configures a SubmitTransactionHandler.
type SubmitTransactionHandlerOption func(h *SubmitTransactionHandler)

// WithResponseTime sets the timeout duration for awaiting a response from the transaction submission provider.
func WithResponseTime(d time.Duration) SubmitTransactionHandlerOption {
	return func(h *SubmitTransactionHandler) {
		h.responseTimeout = d
	}
}

// WithRequestBodyLimit sets the maximum allowed size (in bytes) for incoming request bodies.
func WithRequestBodyLimit(limit int64) SubmitTransactionHandlerOption {
	return func(h *SubmitTransactionHandler) {
		h.requestBodyLimit = limit
	}
}

// SubmitTransactionService abstracts the logic for handling transaction submissions.
type SubmitTransactionService interface {
	SubmitTransaction(ctx context.Context, topics app.Topics, body ...byte) (*overlay.Steak, error)
}

// SubmitTransactionHandler handles incoming transaction requests.
// It validates the request body, translates the content into a format compatible
// with the overlay engine, and invokes the appropriate service logic.
type SubmitTransactionHandler struct {
	service          SubmitTransactionService
	requestBodyLimit int64
	responseTimeout  time.Duration
}

// SubmitTransaction processes a transaction submission request.
// It invokes the submission service and returns an appropriate HTTP response
// based on the outcome.
func (s *SubmitTransactionHandler) SubmitTransaction(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	topics, found := headers[http.CanonicalHeaderKey(XTopicsHeader)]
	if !found {
		return c.Status(fiber.StatusBadRequest).JSON(NewRequestMissingHeaderResponse(XTopicsHeader))
	}

	steak, err := s.service.SubmitTransaction(c.UserContext(), topics, c.Body()...)
	switch {
	case errors.Is(err, app.ErrMissingTopics) || errors.Is(err, app.ErrInvalidTopicFormat):
		return c.Status(fiber.StatusBadRequest).JSON(NewInvalidRequestTopicsFormatResponse())

	case errors.Is(err, app.ErrReaderLimitExceeded):
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(NewRequestBodyTooLargeResponse(s.requestBodyLimit))

	case errors.Is(err, app.ErrReaderBytesRead):
		return c.Status(fiber.StatusInternalServerError).JSON(NewTaggedBEEFCreationErrorResponse())

	case errors.Is(err, app.ErrSubmitTransactionProviderTimeout):
		return c.Status(fiber.StatusRequestTimeout).JSON(NewRequestTimeoutResponse(s.responseTimeout))

	case errors.Is(err, app.ErrSubmitTransactionProvider):
		return c.Status(fiber.StatusInternalServerError).JSON(NewSubmitTransactionProviderErrorResponse())

	default:
		return c.Status(fiber.StatusOK).JSON(NewSubmitTransactionSuccessResponse(steak))
	}
}

// NewSubmitTransactionHandler constructs a new SubmitTransactionHandler.
// If the provider is nil, it panics.
func NewSubmitTransactionHandler(provider app.SubmitTransactionProvider, options ...SubmitTransactionHandlerOption) *SubmitTransactionHandler {
	if provider == nil {
		panic("submit transaction provider is nil")
	}

	handler := SubmitTransactionHandler{
		service:          app.NewSubmitTransactionService(provider),
		requestBodyLimit: RequestBodyLimit1GB,
		responseTimeout:  RequestTimeout,
	}

	for _, opt := range options {
		opt(&handler)
	}

	return &handler
}

// NewSubmitTransactionSuccessResponse constructs a successful transaction submission response.
func NewSubmitTransactionSuccessResponse(steak *overlay.Steak) *openapi.SubmitTransactionResponse {
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

// NewSubmitTransactionProviderErrorResponse returns an error response indicating a failure within the overlay engine.
func NewSubmitTransactionProviderErrorResponse() openapi.InternalServerErrorResponse {
	return openapi.Error{
		Message: "Unable to process submitted transaction octet-stream due to an error in the overlay engine.",
	}
}

// NewInvalidRequestTopicsFormatResponse returns a bad request response for invalid topic headers.
func NewInvalidRequestTopicsFormatResponse() openapi.BadRequestResponse {
	return openapi.Error{
		Message: "One or more topic headers are in an invalid format. Empty string values are not allowed.",
	}
}

// NewTaggedBEEFCreationErrorResponse returns an error response for failures during tagged BEEF creation.
func NewTaggedBEEFCreationErrorResponse() openapi.InternalServerErrorResponse {
	return openapi.Error{
		Message: "Unable to process submitted transaction octet-stream due to an issue with the request body.",
	}
}
