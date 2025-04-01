package queries

import (
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
)

// TopicManagerMetadata represents the metadata for a topic manager.
type TopicManagerMetadata struct {
	Name             string  `json:"name"`
	ShortDescription string  `json:"shortDescription"`
	IconURL          *string `json:"iconURL,omitempty"`
	Version          *string `json:"version,omitempty"`
	InformationURL   *string `json:"informationURL,omitempty"`
}

// TopicManagerListHandlerResponse defines the response body content that
// will be sent in JSON format after successfully processing the handler logic.
type TopicManagerListHandlerResponse map[string]TopicManagerMetadata

// MetaData represents the metadata information for topic managers coming from the engine.
type MetaData struct {
	ShortDescription string
	IconURL          string
	Version          string
	InformationURL   string
}

// TopicManagerListProvider defines the contract that must be fulfilled
// to retrieve a list of topic managers from the overlay engine.
type TopicManagerListProvider interface {
	ListTopicManagers() map[string]*MetaData
}

// TopicManagerListHandler orchestrates the processing flow of a topic manager list
// request, returning a map of topic manager metadata with appropriate HTTP status.
type TopicManagerListHandler struct {
	provider TopicManagerListProvider
}

// Handle processes the topic manager list request and sends a JSON response.
func (t *TopicManagerListHandler) Handle(w http.ResponseWriter, r *http.Request) {
	engineTopicManagers := t.provider.ListTopicManagers()
	result := make(TopicManagerListHandlerResponse, len(engineTopicManagers))

	for name, metadata := range engineTopicManagers {
		tmMetadata := TopicManagerMetadata{
			Name:             name,
			ShortDescription: "No description available",
		}
		if metadata != nil {
			if metadata.ShortDescription != "" {
				tmMetadata.ShortDescription = metadata.ShortDescription
			}
			if metadata.IconURL != "" {
				url := metadata.IconURL
				tmMetadata.IconURL = &url
			}
			if metadata.Version != "" {
				version := metadata.Version
				tmMetadata.Version = &version
			}
			if metadata.InformationURL != "" {
				info := metadata.InformationURL
				tmMetadata.InformationURL = &info
			}
		}
		result[name] = tmMetadata
	}

	jsonutil.SendHTTPResponse(w, http.StatusOK, result)
}

// NewTopicManagerListHandler returns an instance of TopicManagerListHandler.
// If the provided provider is nil, it panics.
func NewTopicManagerListHandler(provider TopicManagerListProvider) *TopicManagerListHandler {
	if provider == nil {
		panic("topic manager list provider is nil")
	}
	return &TopicManagerListHandler{provider: provider}
}
