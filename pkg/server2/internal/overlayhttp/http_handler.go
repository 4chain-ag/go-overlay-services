package overlayhttp

import (
	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/gofiber/fiber/v2"
)

type HTTPHandler struct {
	submitTransactionHandler  *SubmitTransactionHandler
	advertisementsSyncHandler *AdvertisementsSyncHandler
}

func (h *HTTPHandler) AdvertisementsSync(c *fiber.Ctx) error {
	return h.advertisementsSyncHandler.Handle(c)
}

func (h *HTTPHandler) SubmitTransaction(c *fiber.Ctx, params openapi.SubmitTransactionParams) error {
	return h.submitTransactionHandler.Handle(c, params)
}

func NewHTTPHandler(provider engine.OverlayEngineProvider) openapi.ServerInterface {
	return &HTTPHandler{
		submitTransactionHandler:  NewSubmitTransactionHandler(provider),
		advertisementsSyncHandler: NewAdvertisementsSyncHandler(provider),
	}
}
