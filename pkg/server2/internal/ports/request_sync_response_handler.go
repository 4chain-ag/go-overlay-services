package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// XBSVTopicHeader is the header key for the BSV topic

const XBSVTopicHeader = "X-BSV-Topic"

// RequestSyncResponseService defines the interface for the sync response service

type RequestSyncResponseService interface {
	RequestSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error)
}

// RequestSyncResponseHandler handles requests for sync responses

type RequestSyncResponseHandler struct {
	service RequestSyncResponseService
}

// Handle processes sync response requests

func (h *RequestSyncResponseHandler) Handle(c *fiber.Ctx) error {

	// Check for topic header

	topic := c.Get(XBSVTopicHeader)

	if topic == "" {

		return NewMissingXBSVTopicHeaderError()

	}

	// Parse request body

	var initialRequest core.GASPInitialRequest

	if err := c.BodyParser(&initialRequest); err != nil {

		return NewInvalidRequestBodyError()

	}

	// Call service

	response, err := h.service.RequestSyncResponse(c.Context(), &initialRequest, topic)

	if err != nil {

		return err

	}

	return c.Status(fiber.StatusOK).JSON(response)

}

// NewRequestSyncResponseHandler creates a new handler

func NewRequestSyncResponseHandler(provider app.RequestSyncResponseProvider) *RequestSyncResponseHandler {

	if provider == nil {

		panic("request sync response provider is nil")

	}

	return &RequestSyncResponseHandler{

		service: app.NewRequestSyncResponseService(provider),
	}

}

// NewMissingXBSVTopicHeaderError creates an error for missing X-BSV-Topic header

func NewMissingXBSVTopicHeaderError() app.Error {

	const msg = "The submitted request does not include required header: X-BSV-Topic"

	return app.NewIncorrectInputError(msg, msg)

}

// NewInvalidRequestBodyError creates an error for invalid request body

func NewInvalidRequestBodyError() app.Error {

	const msg = "The submitted request body is invalid or malformed"

	return app.NewIncorrectInputError(msg, msg)

}

// RequestSyncResponseInternalErrorResponse is the error response for internal errors

var RequestSyncResponseInternalErrorResponse = openapi.InternalServerErrorResponse{

	Message: "Unable to process sync response request due to an error in the overlay engine.",
}
