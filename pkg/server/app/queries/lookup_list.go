package queries

import (
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
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

// MetaDataLookup represents the metadata information for lookup service providers coming from the engine.
type MetaDataLookup struct {
	ShortDescription string
	IconURL          string
	Version          string
	InformationURL   string
}

// LookupListProvider defines the contract that must be fulfilled
// to retrieve a list of lookup service providers from the overlay engine.
type LookupListProvider interface {
	ListLookupServiceProviders() map[string]*MetaDataLookup
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

	for name, metadata := range engineLookupProviders {
		lookupMetadata := LookupMetadata{
			Name:             name,
			ShortDescription: "No description available",
		}
		if metadata != nil {
			if metadata.ShortDescription != "" {
				lookupMetadata.ShortDescription = metadata.ShortDescription
			}
			if metadata.IconURL != "" {
				url := metadata.IconURL
				lookupMetadata.IconURL = &url
			}
			if metadata.Version != "" {
				version := metadata.Version
				lookupMetadata.Version = &version
			}
			if metadata.InformationURL != "" {
				info := metadata.InformationURL
				lookupMetadata.InformationURL = &info
			}
		}
		result[name] = lookupMetadata
	}

	jsonutil.SendHTTPResponse(w, http.StatusOK, result)
}

// NewLookupListHandler returns an instance of LookupListHandler.
// If the provided provider is nil, it panics.
func NewLookupListHandler(provider LookupListProvider) *LookupListHandler {
	if provider == nil {
		panic("lookup list provider is nil")
	}
	return &LookupListHandler{provider: provider}
}
