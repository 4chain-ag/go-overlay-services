package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// ArcIngestService defines the interface for the arc ingest service
type ArcIngestService interface {
	HandleArcIngest(ctx context.Context, dto *app.ArcIngestDTO) error
}

// ArcIngestHandler handles HTTP requests for ARC (transaction) ingest operations.
// It delegates processing to the ArcIngestService.
type ArcIngestHandler struct {
	service ArcIngestService
}

// Handle processes an HTTP request to ingest ARC transaction data.
// It validates the request body, constructs an ArcIngestDTO, and passes it to the service.
func (h *ArcIngestHandler) Handle(c *fiber.Ctx) error {
	var request openapi.ArcIngestBody
	if err := c.BodyParser(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request format")
	}

	dto := &app.ArcIngestDTO{
		TxID:        request.Txid,
		MerklePath:  request.MerklePath,
		BlockHeight: request.BlockHeight,
	}

	err := h.service.HandleArcIngest(c.Context(), dto)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(NewArcIngestSuccessResponse())
}

// NewArcIngestHandler creates a new ArcIngestHandler with the given provider.
// It initializes the internal ArcIngestService.
func NewArcIngestHandler(provider engine.OverlayEngineProvider) *ArcIngestHandler {
	if provider == nil {
		panic("arc ingest provider is nil")
	}

	return &ArcIngestHandler{
		service: app.NewArcIngestService(provider),
	}
}

// NewArcIngestSuccessResponse creates a successful response for the arc ingest request
// It returns a standardized success response format.
func NewArcIngestSuccessResponse() *openapi.ArcIngest {
	return &openapi.ArcIngest{
		Status:  "success",
		Message: "Transaction successfully ingested",
	}
}
