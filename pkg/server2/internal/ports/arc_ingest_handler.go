package ports

import (
	"context"
	"errors"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// ArcRequestTimeout defines a default timeout for ARC ingest requests.
const ArcRequestTimeout = 10 * time.Second

// ArcIngestHandlerOption defines a function that configures an ArcIngestHandler.
type ArcIngestHandlerOption func(h *ArcIngestHandler)

// WithArcResponseTimeout configures the timeout duration for Merkle proof processing.
func WithArcResponseTimeout(d time.Duration) ArcIngestHandlerOption {
	return func(h *ArcIngestHandler) {
		h.responseTimeout = d
	}
}

// ArcIngestService defines the interface for a service responsible for handling ARC ingest requests.
type ArcIngestService interface {
	HandleArcIngest(ctx context.Context, txID string, merklePath string, blockHeight uint32) error
}

// ArcIngestHandler handles incoming ARC ingest requests, validating the request body,
// and forwarding the request to the service for processing.
type ArcIngestHandler struct {
	service         ArcIngestService
	responseTimeout time.Duration
}

// ArcIngestRequest defines the expected structure for the ARC ingest request body.
type ArcIngestRequest struct {
	TxID        string `json:"txid"`
	MerklePath  string `json:"merklePath"`
	BlockHeight uint32 `json:"blockHeight"`
}

// Validate checks if all required fields are present and valid.
func (r *ArcIngestRequest) Validate() error {
	if r.TxID == "" {
		return errors.New("missing required field: txid")
	}
	if r.MerklePath == "" {
		return errors.New("missing required field: merkle path")
	}
	return nil
}

// ArcIngestResponse represents the response format for the ArcIngestHandler.
type ArcIngestResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HandleArcIngest processes an HTTP request for ARC ingest.
// It validates the request, processes it through the service, and returns
// appropriate responses for success or failure cases.
func (h *ArcIngestHandler) HandleArcIngest(c *fiber.Ctx) error {
	var request ArcIngestRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(openapi.BadRequestResponse{
			Message: "Invalid request body format",
		})
	}

	if err := request.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(openapi.BadRequestResponse{
			Message: err.Error(),
		})
	}

	err := h.service.HandleArcIngest(c.Context(), request.TxID, request.MerklePath, request.BlockHeight)

	// Check for specific error types
	if err != nil {
		if appErr, ok := err.(app.Error); ok {
			// Handle app.Error types based on their ErrorType
			switch appErr.ErrorType() {
			case app.ErrorTypeOperationTimeout:
				return c.Status(fiber.StatusGatewayTimeout).JSON(openapi.Error{
					Message: appErr.Slug(),
				})
			case app.ErrorTypeUnknown:
				return c.Status(fiber.StatusRequestTimeout).JSON(openapi.Error{
					Message: appErr.Slug(),
				})
			case app.ErrorTypeIncorrectInput:
				return c.Status(fiber.StatusBadRequest).JSON(openapi.Error{
					Message: appErr.Slug(),
				})
			case app.ErrorTypeProviderFailure:
				return c.Status(fiber.StatusInternalServerError).JSON(openapi.Error{
					Message: appErr.Slug(),
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(openapi.Error{
					Message: "Internal server error occurred during processing",
				})
			}
		}

		// Default error handler for non-app.Error types
		return c.Status(fiber.StatusInternalServerError).JSON(openapi.Error{
			Message: "Internal server error occurred during processing",
		})
	}

	// Success case
	return c.Status(fiber.StatusOK).JSON(ArcIngestResponse{
		Status:  "success",
		Message: "Transaction status updated",
	})
}

// NewArcIngestHandler creates a new instance of ArcIngestHandler.
func NewArcIngestHandler(provider app.NewMerkleProofProvider, opts ...ArcIngestHandlerOption) *ArcIngestHandler {
	if provider == nil {
		panic("provider is nil")
	}

	handler := ArcIngestHandler{
		service:         app.NewArcIngestService(provider, app.DefaultArcIngestTimeout),
		responseTimeout: ArcRequestTimeout,
	}

	for _, opt := range opts {
		opt(&handler)
	}

	return &handler
}
