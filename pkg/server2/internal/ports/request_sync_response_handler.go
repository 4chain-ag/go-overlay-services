package ports

import (
	"context"
	"errors"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// XBSVTopicHeader is the HTTP header name for BSV topic
const XBSVTopicHeader = "X-BSV-Topic"

// RequestSyncResponseService abstracts the logic for handling sync response requests.
type RequestSyncResponseService interface {
	RequestSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error)
}

// RequestSyncResponseHandler orchestrates the processing flow of a sync response
// request and applies any necessary logic before invoking the engine.
type RequestSyncResponseHandler struct {
	service RequestSyncResponseService
}

// Handle orchestrates the processing flow of a sync response request.
// It prepares and sends a JSON response after invoking the engine and returns an HTTP response
// with the appropriate status code based on the engine's response.
func (h *RequestSyncResponseHandler) Handle(c *fiber.Ctx) error {
	topic := c.Get(XBSVTopicHeader)
	if topic == "" {
		return c.Status(fiber.StatusBadRequest).JSON(openapi.Error{
			Message: "Missing 'X-BSV-Topic' header",
		})
	}

	var initialRequest core.GASPInitialRequest
	if err := c.BodyParser(&initialRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(openapi.Error{
			Message: "Invalid request body",
		})
	}

	response, err := h.service.RequestSyncResponse(c.Context(), &initialRequest, topic)
	if err != nil {
		if errors.Is(err, app.ErrRequestSyncResponseProvider) {
			return c.Status(fiber.StatusInternalServerError).JSON(openapi.Error{
				Message: "Unable to process sync response request due to issues with the overlay engine",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(openapi.Error{
			Message: "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// NewRequestSyncResponseHandler returns an instance of a RequestSyncResponseHandler,
// utilizing an implementation of RequestSyncResponseProvider. If the provided argument is nil, it triggers a panic.
func NewRequestSyncResponseHandler(provider app.RequestSyncResponseProvider) *RequestSyncResponseHandler {
	if provider == nil {
		panic("request sync response provider is nil")
	}

	return &RequestSyncResponseHandler{service: app.NewRequestSyncResponseService(provider)}
} 
