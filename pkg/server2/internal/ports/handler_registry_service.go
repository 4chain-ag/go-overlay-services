package ports

import (
	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Handle(c *fiber.Ctx) error
}

// HandlerRegistryService defines the main point for registering HTTP handler dependencies.
// It acts as a central registry for mapping API endpoints to their handler implementations.
type HandlerRegistryService struct {
	lookupList                *LookupListHandler
	topicManagersList         *TopicManagersListHandler
	lookupDocumentation       *LookupProviderDocumentationHandler
	startGASPSync             *StartGASPSyncHandler
	topicManagerDocumentation *TopicManagerDocumentationHandler
	submitTransaction         *SubmitTransactionHandler
	syncAdvertisements        *SyncAdvertisementsHandler
	arcIngest                 Handler
	requestSyncResponse       *RequestSyncResponseHandler
}

// ListLookupServiceProviders method delegates the request to the configured lookup list handler.
func (h *HandlerRegistryService) ListLookupServiceProviders(c *fiber.Ctx) error {
	return h.lookupList.Handle(c)
}

// AdvertisementsSync method delegates the request to the configured sync advertisements handler.
func (h *HandlerRegistryService) AdvertisementsSync(c *fiber.Ctx) error {
	return h.syncAdvertisements.Handle(c)
}

// GetLookupServiceProviderDocumentation method delegates the request to the configured lookup service provider documentation handler.
func (h *HandlerRegistryService) GetLookupServiceProviderDocumentation(c *fiber.Ctx, params openapi.GetLookupServiceProviderDocumentationParams) error {
	return h.lookupDocumentation.Handle(c, params)
}

// GetTopicManagerDocumentation method delegates the request to the configured topic manager documentation handler.
func (h *HandlerRegistryService) GetTopicManagerDocumentation(c *fiber.Ctx, params openapi.GetTopicManagerDocumentationParams) error {
	return h.topicManagerDocumentation.Handle(c, params)
}

// SubmitTransaction method delegates the request to the configured submit transaction handler.
func (h *HandlerRegistryService) SubmitTransaction(c *fiber.Ctx, params openapi.SubmitTransactionParams) error {
	return h.submitTransaction.Handle(c, params)
}

// ListTopicManagers method delegates the request to the configured topic managers list handler.
func (h *HandlerRegistryService) ListTopicManagers(c *fiber.Ctx) error {
	return h.topicManagersList.Handle(c)
}

// StartGASPSync method delegates the request to the configured start GASP sync handler.
func (h *HandlerRegistryService) StartGASPSync(c *fiber.Ctx) error {
	return h.startGASPSync.Handle(c)
}

// ArcIngest method delegates the request to the configured arc ingest handler.
func (h *HandlerRegistryService) ArcIngest(c *fiber.Ctx) error {
	return h.arcIngest.Handle(c)
}

// RequestSyncResponse method delegates the request to the configured request sync response handler.
func (h *HandlerRegistryService) RequestSyncResponse(c *fiber.Ctx, params openapi.RequestSyncResponseParams) error {
	return h.requestSyncResponse.Handle(c, params)
}

// NewHandlerRegistryService creates and returns a new HandlerRegistryService instance.
// It initializes all handler implementations with their required dependencies.
func NewHandlerRegistryService(provider engine.OverlayEngineProvider, arcIngestHandlerConfig ArcIngestHandlerConfig) *HandlerRegistryService {
	return &HandlerRegistryService{
		lookupList:                NewLookupListHandler(provider),
		topicManagersList:         NewTopicManagersListHandler(provider),
		lookupDocumentation:       NewLookupProviderDocumentationHandler(provider),
		startGASPSync:             NewStartGASPSyncHandler(provider),
		topicManagerDocumentation: NewTopicManagerDocumentationHandler(provider),
		submitTransaction:         NewSubmitTransactionHandler(provider),
		syncAdvertisements:        NewSyncAdvertisementsHandler(provider),
		arcIngest:                 NewArcAuthorizationDecorator(NewArcIngestHandler(provider), arcIngestHandlerConfig),
		requestSyncResponse:       NewRequestSyncResponseHandler(provider),
	}
}
