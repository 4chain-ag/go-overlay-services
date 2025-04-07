package queries

import (
	"fmt"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// LookupMetadata represents the metadata for a lookup service provider.
type LookupMetadata struct {
	Name             string  `json:"name"`
	ShortDescription string  `json:"shortDescription"`
	IconURL          *string `json:"iconURL,omitempty"`
	Version          *string `json:"version,omitempty"`
	InformationURL   *string `json:"informationURL,omitempty"`
}

// LookupListHandlerResponse defines the response body content that
// will be sent in JSON format after successfully processing the handler logic.
type LookupListHandlerResponse map[string]LookupMetadata

// LookupListProvider defines the contract that must be fulfilled
// to retrieve a list of lookup service providers from the overlay engine.
type LookupListProvider interface {
	ListLookupServiceProviders() map[string]*overlay.MetaData
}

// LookupListHandler orchestrates the processing flow of a lookup service provider list
// request, returning a map of lookup service provider metadata with appropriate HTTP status.
type LookupListHandler struct {
	provider LookupListProvider
}

// Handle processes the lookup service provider list request and sends a JSON response.
func (l *LookupListHandler) Handle(w http.ResponseWriter, r *http.Request) {
	engineLookupProviders := l.provider.ListLookupServiceProviders()
	result := make(LookupListHandlerResponse, len(engineLookupProviders))

	setIfNotEmpty := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	coalesce := func(primary, fallback string) string {
		if primary != "" {
			return primary
		}
		return fallback
	}

	for name, metadata := range engineLookupProviders {
		lookupMetadata := LookupMetadata{
			Name:             name,
			ShortDescription: "No description available",
		}

		if metadata != nil {
			lookupMetadata.ShortDescription = coalesce(metadata.Description, "No description available")
			lookupMetadata.IconURL = setIfNotEmpty(metadata.Icon)
			lookupMetadata.Version = setIfNotEmpty(metadata.Version)
			lookupMetadata.InformationURL = setIfNotEmpty(metadata.InfoUrl)
		}

		result[name] = lookupMetadata
	}

	jsonutil.SendHTTPResponse(w, http.StatusOK, result)
}

// NewLookupListHandler returns an instance of LookupListHandler.
// If the provided provider is nil, it panics.
func NewLookupListHandler(provider LookupListProvider) (*LookupListHandler, error) {
	if provider == nil {
		return nil, fmt.Errorf("lookup list provider cannot be nil")
	}
	return &LookupListHandler{provider: provider}, nil
}
