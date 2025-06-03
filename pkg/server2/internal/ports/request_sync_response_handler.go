package ports

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// RequestSyncResponseService defines the application-level interface responsible for
// handling foreign sync response requests.
type RequestSyncResponseService interface {
	RequestSyncResponse(ctx context.Context, topic app.Topic, version app.Version, since app.Since) (*app.RequestSyncResponseDTO, error)
}

// RequestSyncResponseHandler handles incoming HTTP requests related to sync response processing.
// It coordinates parsing, validation, service delegation, and response formatting.
type RequestSyncResponseHandler struct {
	service RequestSyncResponseService
}

// Handle process a request to fetch sync response data for a given topic.
// It parses the request body, transforms input into domain models, delegates to the service layer,
// and returns a serialized success response or an appropriate error.
func (h *RequestSyncResponseHandler) Handle(c *fiber.Ctx, params openapi.RequestSyncResponseParams) error {
	var body openapi.RequestSyncResponseJSONRequestBody
	err := c.BodyParser(&body)
	if err != nil {
		return NewRequestBodyParserError(err)
	}

	dto, err := h.service.RequestSyncResponse(
		c.Context(),
		app.NewTopic(params.XBSVTopic),
		app.Version(body.Version),
		app.Since(body.Since),
	)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(NewRequestSyncResponseSuccessResponse(dto))
}

// NewRequestSyncResponseHandler creates a new instance of RequestSyncResponseHandler,
// wiring it with the provided application-level service provider.
func NewRequestSyncResponseHandler(provider app.RequestSyncResponseProvider) *RequestSyncResponseHandler {
	return &RequestSyncResponseHandler{
		service: app.NewRequestSyncResponseService(provider),
	}
}

// NewRequestSyncResponseSuccessResponse transforms a RequestSyncResponseDTO into OpenAPI-compliant
// response format suitable for HTTP transmission.
func NewRequestSyncResponseSuccessResponse(response *app.RequestSyncResponseDTO) *openapi.RequestSyncResResponse {
	if response == nil {
		return &openapi.RequestSyncResResponse{UTXOList: []openapi.UTXOItem{}, Since: 0}
	}

	utxos := make([]openapi.UTXOItem, 0, len(response.UTXOList))
	for _, utxo := range response.UTXOList {
		utxos = append(utxos, openapi.UTXOItem{Txid: utxo.TxID, Vout: int(utxo.OutputIndex)})
	}

	return &openapi.RequestSyncResResponse{
		UTXOList: utxos,
		Since:    int(response.Since),
	}
}
