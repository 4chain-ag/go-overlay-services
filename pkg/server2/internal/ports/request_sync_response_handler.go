package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// RequestSyncResponseService defines the interface for the sync response service
type RequestSyncResponseService interface {
	RequestSyncResponse(ctx context.Context, dto *app.RequestSyncResponseDTO, topic string) (*core.GASPInitialResponse, error)
}

// RequestSyncResponseHandler handles requests for sync responses
type RequestSyncResponseHandler struct {
	service RequestSyncResponseService
}

// Handle processes sync response requests
func (h *RequestSyncResponseHandler) Handle(c *fiber.Ctx, params openapi.RequestSyncResponseParams) error {
	var requestBody openapi.RequestSyncResponseJSONRequestBody
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": app.NewRequestSyncResponseInvalidJSONError(),
		})
	}

	dto := app.RequestSyncResponseDTO{
		Version: requestBody.Version,
		Since:   requestBody.Since,
	}

	response, err := h.service.RequestSyncResponse(c.Context(), &dto, params.XBSVTopic)

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

	return &RequestSyncResponseHandler{service: app.NewRequestSyncResponseService(provider)}
}

// NewRequestSyncResponseSuccessResponse creates a successful response for the sync response request
// It maps the GASPInitialResponse data to an OpenAPI response format.
func NewRequestSyncResponseSuccessResponse(response *core.GASPInitialResponse) *openapi.RequestSyncResResponse {

	if response == nil {
		return &openapi.RequestSyncResResponse{
			UTXOList: []openapi.UTXOItem{},
			Since:    0,
		}
	}

	utxoList := make([]openapi.UTXOItem, 0, len(response.UTXOList))

	for _, utxo := range response.UTXOList {
		utxoList = append(utxoList, openapi.UTXOItem{
			Txid: utxo.Txid.String(),
			Vout: int(utxo.OutputIndex),
		})
	}

	return &openapi.RequestSyncResResponse{
		UTXOList: utxoList,
		Since:    int(response.Since),
	}

}
