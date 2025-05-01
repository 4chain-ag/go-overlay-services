package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// StartGASPSyncService abstracts the logic for handling GASP sync operations.

type StartGASPSyncService interface {
	StartGASPSync(ctx context.Context) error
}

// StartGASPSyncHandler orchestrates the processing flow of a GASP sync

// request and applies any necessary logic before invoking the engine.

type StartGASPSyncHandler struct {
	service StartGASPSyncService
}

// Handle orchestrates the processing flow of a GASP sync request.

// It prepares and sends a JSON response after invoking the engine and returns an HTTP response

// with the appropriate status code based on the engine's response.

func (h *StartGASPSyncHandler) Handle(c *fiber.Ctx) error {

	err := h.service.StartGASPSync(c.Context())

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(NewStartGASPSyncProviderErrorResponse())

	}

	return c.Status(fiber.StatusOK).JSON(NewStartGASPSyncSuccessResponse())

}

// NewStartGASPSyncHandler returns an instance of a StartGASPSyncHandler,

// utilizing an implementation of StartGASPSyncProvider. If the provided argument is nil, it triggers a panic.

func NewStartGASPSyncHandler(provider app.StartGASPSyncProvider) *StartGASPSyncHandler {

	if provider == nil {

		panic("start GASP sync provider is nil")

	}

	return &StartGASPSyncHandler{service: app.NewStartGASPSyncService(provider)}

}

// NewStartGASPSyncSuccessResponse creates a successful response for GASP synchronization.

// It returns an instance of openapi.StartGASPSyncResponse with a predefined success message.

func NewStartGASPSyncSuccessResponse() openapi.StartGASPSyncResponse {

	return openapi.StartGASPSyncResponse{

		Message: "GASP sync request successfully delegated to overlay engine.",
	}

}

// NewStartGASPSyncProviderErrorResponse creates an error response for GASP synchronization failures.

// It returns an instance of openapi.InternalServerErrorResponse with a predefined error message

// indicating an issue with the overlay engine during the sync process.

func NewStartGASPSyncProviderErrorResponse() openapi.InternalServerErrorResponse {

	return openapi.Error{Message: "Unable to process GASP sync request due to issues with the overlay engine."}

}
