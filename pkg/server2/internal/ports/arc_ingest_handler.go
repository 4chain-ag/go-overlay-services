package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// TODO: This is a temporary implementation, we will require to update it using the arc api key middleware implementation.
// We will also need to update the openapi spec to match the new request body as security field.

// ArcIngestHandler handles HTTP requests for ARC (transaction) ingest operations.
// It delegates processing to the ArcIngestService.
type ArcIngestHandler struct {
	service *app.ArcIngestService
}

// HandleArcIngest processes an HTTP request to ingest ARC transaction data.
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

	err := h.service.HandleArcIngest(context.Background(), dto)
	if err != nil {
		return err // Let the global error handler process the application error
	}

	return c.Status(fiber.StatusOK).JSON(openapi.ArcIngest{
		Status:  "success",
		Message: "Transaction successfully ingested",
	})
}

// NewArcIngestHandler creates a new ArcIngestHandler with the given provider.
// It initializes the internal ArcIngestService.
func NewArcIngestHandler(provider engine.OverlayEngineProvider) *ArcIngestHandler {
	return &ArcIngestHandler{
		service: app.NewArcIngestService(provider),
	}
}
