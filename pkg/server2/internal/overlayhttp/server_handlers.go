package overlayhttp

import (
	"strings"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
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
	return BearearTokenAuthorizationMiddleware(s.token, s.advertisementsSyncHandler.Handle)(c)
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

// BearearTokenAuthorizationMiddleware is a middleware function that checks if the request
// contains a valid Bearer token in the Authorization header. If the token is invalid or
// missing, it responds with an appropriate error.
func BearearTokenAuthorizationMiddleware(expectedToken string, next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the Authorization header from the request
		auth := c.Get("Authorization")

		// Check if the Authorization header is missing
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(openapi.BadRequestResponse{Message: "Missing Authorization header in the request"})
		}

		// Check if the Authorization header does not start with 'Bearer '
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(openapi.BadRequestResponse{Message: "Missing Authorization header Bearer token value"})
		}

		// Extract the token from the Authorization header
		token := strings.TrimPrefix(auth, "Bearer ")

		// Check if the token does not match the expected token
		if token != expectedToken {
			return c.Status(fiber.StatusForbidden).JSON(openapi.BadRequestResponse{Message: "Invalid Bearer token value"})
		}

		// Proceed with the next handler if the token is valid
		return next(c)
	}
}
