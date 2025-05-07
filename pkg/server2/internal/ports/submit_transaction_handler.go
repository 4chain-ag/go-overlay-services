package ports

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
)

// XTopicsHeader defines the HTTP header key used for specifying transaction topics.
const XTopicsHeader = "x-topics"

// SubmitTransactionService defines the interface for a service responsible for submitting transactions.
type SubmitTransactionService interface {
	SubmitTransaction(ctx context.Context, topics app.TransactionTopics, body ...byte) (*overlay.Steak, error)
}

// SubmitTransactionHandler handles incoming transaction requests.
// It validates the request body, translates the content into a format compatible
// with the submit transaction service, and invokes the appropriate service logic.
type SubmitTransactionHandler struct {
	service         SubmitTransactionService
	responseTimeout time.Duration
}

// ResponseTimeout returns the configured timeout duration for submitting a transaction response.
// This value defines how long the handler will wait before treating the operation as timed out.
func (s *SubmitTransactionHandler) ResponseTimeout() time.Duration {
	return s.responseTimeout
}

// Handle processes an HTTP request to submit a transaction to the submit transaction service.
// It expects the `x-topics` header to be present and valid. On success, it returns
// HTTP 200 OK with a STEAK response (openapi.SubmitTransactionResponse).
// If the header is missing or invalid, it returns HTTP 400 Bad Request.
// If an error occurs during transaction submission, it returns the corresponding application error.
func (s *SubmitTransactionHandler) Handle(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	topics, found := headers[http.CanonicalHeaderKey(XTopicsHeader)]
	if !found {
		return c.Status(fiber.StatusBadRequest).JSON(NewRequestMissingHeaderResponse(XTopicsHeader))
	}

	steak, err := s.service.SubmitTransaction(c.UserContext(), topics, c.Body()...)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(NewSubmitTransactionSuccessResponse(steak))
}

// NewSubmitTransactionHandler creates a new SubmitTransactionHandler with the given provider and timeout.
// If the provider is nil, it panics. The request timeout determines how long the handler will wait
// for a response from the submit transaction service before timing out and responding with a request timeout response.
func NewSubmitTransactionHandler(provider app.SubmitTransactionProvider, timeout time.Duration) *SubmitTransactionHandler {
	if provider == nil {
		panic("submit transaction provider is nil")
	}

	handler := SubmitTransactionHandler{
		service:         app.NewSubmitTransactionService(provider, timeout),
		responseTimeout: timeout,
	}
	return &handler
}

// NewSubmitTransactionSuccessResponse creates a successful response for the transaction submission.
// It maps the Steak data to an OpenAPI response format.
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

// SubmitTransactionServiceInternalError is the internal server error response for transaction submission.
// This error is returned when an internal issue occurs while processing the submitted transaction.
// Typically, this happens when the overlay engine encounters an unexpected error.
var SubmitTransactionServiceInternalError = openapi.InternalServerErrorResponse{
	Message: "Unable to process submitted transaction octet-stream due to an error in the overlay engine.",
}

// SubmitTransactionRequestInvalidTopicsHeaderFormat is the bad request response for invalid topic header format.
// This error occurs when the topics header provided in the request is either missing or incorrectly formatted.
// For example, an empty string or invalid character in the topic header would trigger this error.
var SubmitTransactionRequestInvalidTopicsHeaderFormat = openapi.BadRequestResponse{
	Message: "One or more topic headers are in an invalid format. Empty string values are not allowed.",
}

// RequestTimeout defines the default duration after which a request is considered timed out.
const RequestTimeout = 5 * time.Second

// NewRequestMissingHeaderResponse creates a bad request response indicating that a required HTTP header is missing.
// It takes the name of the missing header and returns an openapi.BadRequestResponse with a descriptive message.
func NewRequestMissingHeaderResponse(header string) openapi.BadRequestResponse {
	return openapi.BadRequestResponse{
		Message: fmt.Sprintf("The submitted request does not include required header: %s.", header),
	}
}
