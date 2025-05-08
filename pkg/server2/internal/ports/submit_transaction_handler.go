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

// RequestTimeout defines the default duration after which a request is considered timed out.
const RequestTimeout = 5 * time.Second

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

// Handle processes an HTTP request to submit a transaction to the submit transaction service.
// It expects the `x-topics` header to be present and valid. On success, it returns
// HTTP 200 OK with a STEAK response (openapi.SubmitTransactionResponse).
// If the header is missing or invalid, it returns HTTP 400 Bad Request.
// If an error occurs during transaction submission, it returns the corresponding application error.
func (s *SubmitTransactionHandler) Handle(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	topics, found := headers[http.CanonicalHeaderKey(XTopicsHeader)]
	if !found {
		return NewMissingXTopicsHeaderError()
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

func NewMissingXTopicsHeaderError() app.Error {
	str := fmt.Sprintf("The submitted request does not include required header: %s.", XTopicsHeader)
	return app.NewIncorrectInputError(str, str)
}
