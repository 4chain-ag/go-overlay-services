package ports

import (
	"context"
	"math"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// RequestSyncResponseService defines the interface for the sync response service

type RequestSyncResponseService interface {
	RequestSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error)
}

// RequestSyncResponseHandler handles requests for sync responses

type RequestSyncResponseHandler struct {
	service RequestSyncResponseService
}

// Handle processes sync response requests

func (h *RequestSyncResponseHandler) Handle(c *fiber.Ctx, params openapi.RequestSyncResponseParams) error {

	var requestBody openapi.RequestSyncResponseJSONRequestBody

	if err := c.BodyParser(&requestBody); err != nil {

		return NewRequestSyncResponseInvalidRequestBodyError()

	}

	if requestBody.Since < 0 || requestBody.Since > math.MaxUint32 {

		return NewRequestSyncResponseInvalidRequestBodyError()

	}

	initialRequest := &core.GASPInitialRequest{

		Version: requestBody.Version,

		Since: uint32(requestBody.Since),
	}

	response, err := h.service.RequestSyncResponse(c.Context(), initialRequest, params.XBSVTopic)

	if err != nil {

		return err

	}

	return c.Status(fiber.StatusOK).JSON(NewRequestSyncResponseSuccessResponse(response))

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

// NewRequestSyncResponseSuccessResponse creates a successful response for the sync response request

// It maps the GASPInitialResponse data to an OpenAPI response format.

func NewRequestSyncResponseSuccessResponse(response *core.GASPInitialResponse) *openapi.RequestSyncResResponse {

	if response == nil {

		return &openapi.RequestSyncResResponse{

			UTXOList: []struct {
				Txid string `json:"txid"`

				Vout int `json:"vout"`
			}{},

			Since: 0,
		}

	}

	utxoList := make([]struct {
		Txid string `json:"txid"`

		Vout int `json:"vout"`
	}, 0, len(response.UTXOList))

	for _, utxo := range response.UTXOList {

		utxoList = append(utxoList, struct {
			Txid string `json:"txid"`

			Vout int `json:"vout"`
		}{

			Txid: utxo.Txid.String(),

			Vout: int(utxo.OutputIndex),
		})

	}

	return &openapi.RequestSyncResResponse{

		UTXOList: utxoList,

		Since: int(response.Since),
	}

}

// NewRequestSyncResponseInvalidRequestBodyError returns an Error indicating that the request body is invalid.

func NewRequestSyncResponseInvalidRequestBodyError() app.Error {

	return app.NewIncorrectInputError(

		"Invalid request body format or content",

		"The request body contains invalid data or format. Please check the JSON structure and field values.",
	)

}
