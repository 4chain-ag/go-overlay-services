package overlayhttp

import (
	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/overlayhttp/middleware"
	"github.com/gofiber/fiber/v2"
)

// ServerHandlers is a struct that holds handlers for different server routes.
type ServerHandlers struct {
	token                     string                     // token used for Bearer token authentication
	submitTransactionHandler  *SubmitTransactionHandler  // handler for submitting transactions
	advertisementsSyncHandler *AdvertisementsSyncHandler // handler for synchronizing advertisements
}

// AdvertisementsSync handles the synchronization of advertisements by verifying
// the Bearer token and delegating the request to the advertisementsSyncHandler.
func (s *ServerHandlers) AdvertisementsSync(c *fiber.Ctx) error {
	// Middleware to authorize requests using Bearer token authentication
	return middleware.BearerTokenAuthorizationMiddleware(s.token, s.advertisementsSyncHandler.Handle)(c)
}

// SubmitTransaction handles the submission of a transaction. It delegates the
// request to the submitTransactionHandler, passing along the provided params.
func (s *ServerHandlers) SubmitTransaction(c *fiber.Ctx, params openapi.SubmitTransactionParams) error {
	// Delegates the request to the SubmitTransactionHandler
	return s.submitTransactionHandler.Handle(c, params)
}

// NewServerHandlers creates a new instance of ServerHandlers with the specified token
// and overlay engine provider. It initializes both non-admin and admin handlers.
func NewServerHandlers(token string, provider engine.OverlayEngineProvider) openapi.ServerInterface {
	if provider == nil {
		panic("overlay engine provider is nil")
	}

	// Return a new instance of ServerHandlers
	return &ServerHandlers{
		token: token,
		// Non-admin handlers:
		submitTransactionHandler: NewSubmitTransactionHandler(provider),
		// Admin handlers:
		advertisementsSyncHandler: NewAdvertisementsSyncHandler(provider),
	}
}
