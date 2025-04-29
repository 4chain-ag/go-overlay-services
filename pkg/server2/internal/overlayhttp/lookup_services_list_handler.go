package overlayhttp

import (
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
	"k8s.io/utils/ptr"
)

// LookupServicesListProvider defines the contract that must be fulfilled
// to retrieve a list of lookup service providers from the overlay engine.
type LookupServicesListProvider interface {
	ListLookupServiceProviders() map[string]*overlay.MetaData
}

// LookupServicesListHandler orchestrates the retrieval of available lookup service providers
// and returns their metadata.
type LookupServicesListHandler struct {
	provider LookupServicesListProvider
}

// Handle processes a request for listing lookup service providers.
// It retrieves the list of available lookup service providers and returns their metadata as JSON.
func (l *LookupServicesListHandler) Handle(c *fiber.Ctx) error {
	providers := l.provider.ListLookupServiceProviders()
	if providers == nil {
		return c.Status(fiber.StatusOK).JSON(make(map[string]openapi.LookupMetadata))
	}

	response := make(map[string]openapi.LookupMetadata, len(providers))

	setIfNotEmpty := func(s string) *string {
		if s == "" {
			return nil
		}
		return ptr.To(s)
	}

	coalesce := func(primary, fallback string) string {
		if primary != "" {
			return primary
		}
		return fallback
	}

	for name, metadata := range providers {
		lookupMetadata := openapi.LookupMetadata{
			Name:             name,
			ShortDescription: "No description available",
		}

		if metadata != nil {
			lookupMetadata.ShortDescription = coalesce(metadata.Description, "No description available")
			lookupMetadata.IconURL = setIfNotEmpty(metadata.Icon)
			lookupMetadata.Version = setIfNotEmpty(metadata.Version)
			lookupMetadata.InformationURL = setIfNotEmpty(metadata.InfoUrl)
		}

		response[name] = lookupMetadata
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// NewLookupServicesListHandler returns a new instance of a LookupServicesListHandler.
// If the provider argument is nil, it triggers a panic.
func NewLookupServicesListHandler(provider LookupServicesListProvider) *LookupServicesListHandler {
	if provider == nil {
		panic("lookup services list provider is nil")
	}

	return &LookupServicesListHandler{
		provider: provider,
	}
} 
