package ports

import (
	"fmt"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// LookupServicesListService defines the interface for a service responsible for retrieving

// and formatting lookup service provider metadata.

type LookupServicesListService interface {
	ListLookupServiceProviders() app.LookupServicesListResponse
}

// LookupServicesListHandler handles incoming requests for lookup service provider information.

// It delegates to the LookupServicesListService to retrieve the metadata and formats

// the response according to the API spec.

type LookupServicesListHandler struct {
	service LookupServicesListService
}

// Handle processes an HTTP request to list all lookup service providers.

// It returns an HTTP 200 OK with a LookupServicesListResponse.

func (h *LookupServicesListHandler) Handle(c *fiber.Ctx) error {

	response := h.service.ListLookupServiceProviders()

	return c.Status(http.StatusOK).JSON(response)

}

// NewLookupServicesListHandler creates a new LookupServicesListHandler with the given provider.

// It initializes the internal LookupServicesListService.

// Panics if the provider is nil.

func NewLookupServicesListHandler(provider app.LookupServicesListProvider) *LookupServicesListHandler {

	service, err := app.NewLookupServicesListService(provider)

	if err != nil {

		panic(fmt.Sprintf("failed to create lookup services list service: %v", err))

	}

	return &LookupServicesListHandler{

		service: service,
	}

}

// LookupServicesListServiceInternalError is the internal server error response for lookup services list.

// This error is returned when an internal issue occurs while retrieving the lookup services list.

var LookupServicesListServiceInternalError = openapi.InternalServerErrorResponse{

	Message: "Unable to retrieve lookup services list due to an error in the overlay engine.",
}
