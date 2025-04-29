package overlayhttp

import (
	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/overlayhttp/middleware"
	"github.com/gofiber/fiber/v2"
)

// ServerHandlers is a struct that holds handlers for different server routes.
type ServerHandlers struct {
	token                           string                           // token used for Bearer token authentication
	submitTransactionHandler        *SubmitTransactionHandler        // handler for submitting transactions
	advertisementsSyncHandler       *AdvertisementsSyncHandler       // handler for synchronizing advertisements
	lookupServiceDocumentationHandler *LookupServiceDocumentationHandler // handler for lookup service documentation
	topicManagerDocumentationHandler  *TopicManagerDocumentationHandler  // handler for topic manager documentation
	topicManagersListHandler          *TopicManagersListHandler          // handler for listing topic managers
	lookupServicesListHandler         *LookupServicesListHandler         // handler for listing lookup service providers
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

// LookupServiceDocumentation handles a request for documentation of a lookup service.
// It delegates the request to the lookupServiceDocumentationHandler.
func (s *ServerHandlers) LookupServiceDocumentation(c *fiber.Ctx, params openapi.LookupServiceDocumentationParams) error {
	// Delegates the request to the LookupServiceDocumentationHandler
	return s.lookupServiceDocumentationHandler.Handle(c, params)
}

// TopicManagerDocumentation handles a request for documentation of a topic manager.
// It delegates the request to the topicManagerDocumentationHandler.
func (s *ServerHandlers) TopicManagerDocumentation(c *fiber.Ctx, params openapi.TopicManagerDocumentationParams) error {
	// Delegates the request to the TopicManagerDocumentationHandler
	return s.topicManagerDocumentationHandler.Handle(c, params)
}

// TopicManagersList handles a request for listing available topic managers.
// It delegates the request to the topicManagersListHandler.
func (s *ServerHandlers) TopicManagersList(c *fiber.Ctx) error {
	// Delegates the request to the TopicManagersListHandler
	return s.topicManagersListHandler.Handle(c)
}

// LookupServicesList handles a request for listing available lookup service providers.
// It delegates the request to the lookupServicesListHandler.
func (s *ServerHandlers) LookupServicesList(c *fiber.Ctx) error {
	// Delegates the request to the LookupServicesListHandler
	return s.lookupServicesListHandler.Handle(c)
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
		lookupServiceDocumentationHandler: NewLookupServiceDocumentationHandler(provider),
		topicManagerDocumentationHandler: NewTopicManagerDocumentationHandler(provider),
		topicManagersListHandler: NewTopicManagersListHandler(provider),
		lookupServicesListHandler: NewLookupServicesListHandler(provider),
		// Admin handlers:
		advertisementsSyncHandler: NewAdvertisementsSyncHandler(provider),
	}
}
