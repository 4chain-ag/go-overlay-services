package ports

import (
	"context"

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
	service SubmitTransactionService
}

// Handle processes an HTTP request to submit a transaction to the submit transaction service.
// It expects the `x-topics` header to be present and valid. On success, it returns
// HTTP 200 OK with a STEAK response (openapi.SubmitTransactionResponse).
// If an error occurs during transaction submission, it returns the corresponding application error.
func (s *SubmitTransactionHandler) Handle(c *fiber.Ctx, params openapi.SubmitTransactionParams) error {
	steak, err := s.service.SubmitTransaction(c.UserContext(), params.XTopics, c.Body()...)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(NewSubmitTransactionSuccessResponse(steak))
}

// NewSubmitTransactionHandler creates a new SubmitTransactionHandler with the given provider.
// If the provider is nil, it panics.
func NewSubmitTransactionHandler(provider app.SubmitTransactionProvider) *SubmitTransactionHandler {
	if provider == nil {
		panic("submit transaction provider is nil")
	}

	return &SubmitTransactionHandler{service: app.NewSubmitTransactionService(provider)}
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
