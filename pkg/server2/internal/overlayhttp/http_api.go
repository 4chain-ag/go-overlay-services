package overlayhttp

import (
	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/gofiber/fiber/v2"
)

type API struct {
	submitTransactionHandler  *SubmitTransactionHandler
	advertisementsSyncHandler *AdvertisementsSyncHandler
}

func (a *API) AdvertisementsSync(c *fiber.Ctx) error {
	return a.advertisementsSyncHandler.Handle(c)
}

func (a *API) SubmitTransaction(c *fiber.Ctx, params openapi.SubmitTransactionParams) error {
	return a.submitTransactionHandler.Handle(c, params)
}

func NewAPI(provider engine.OverlayEngineProvider) openapi.ServerInterface {
	if provider == nil {
		panic("overlay engine provider is nil")
	}

	return &API{
		submitTransactionHandler:  NewSubmitTransactionHandler(provider),
		advertisementsSyncHandler: NewAdvertisementsSyncHandler(provider),
	}
}
