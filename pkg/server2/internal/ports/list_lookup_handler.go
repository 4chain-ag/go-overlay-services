package ports

import (
	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/gofiber/fiber/v2"
)

// LookupListHandler handles HTTP requests for listing lookup service providers.

type LookupListHandler struct {
	service *app.LookupListService
}

// Handle processes HTTP requests to list lookup service providers.

// It uses the LookupListService to retrieve the list and returns it as JSON.

func (h *LookupListHandler) Handle(c *fiber.Ctx) error {

	response := h.service.ListLookup()

	return c.JSON(response)

}

// NewLookupListHandler creates a new LookupListHandler with the given engine provider.

// It initializes the underlying service and returns an error if the provider is nil.

func NewLookupListHandler(provider engine.OverlayEngineProvider) *LookupListHandler {

	service, err := app.NewLookupListService(provider)

	if err != nil {

		panic(err)

	}

	return &LookupListHandler{service: service}

}
