package queries

import (
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
)

// TopicManagerDocumentationHandlerResponse defines the response body content that
// will be sent in JSON format after successfully processing the handler logic.
type TopicManagerDocumentationHandlerResponse struct {
	Message string `json:"message"`
}

// TopicManagerDocumentationProvider defines the contract that must be fulfilled
// to send a topic manager documentation request to the overlay engine for further processing.
// Note: The contract definition is still in development and will be updated after
// migrating the engine code.
type TopicManagerDocumentationProvider interface {
	GetDocumentationForTopicManager(provider string) (string, error)
}

// TopicManagerDocumentationHandler orchestrates the processing flow of a topic documentation
// request, including the request body validation, converting the request body
// into an overlay-engine-compatible format, and applying any other necessary
// logic before invoking the engine. It returns the requested topic manager
// documentation in the text/markdown format.
type TopicManagerDocumentationHandler struct {
	provider TopicManagerDocumentationProvider
}

// Handle orchestrates the processing flow of a topic manager documentation request.
// It prepares and sends a JSON response after invoking the engine and returns an HTTP response
// with the appropriate status code based on the engine's response.
func (t *TopicManagerDocumentationHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// TODO: Add custom validation logic.
	_, err := t.provider.GetDocumentationForTopicManager("")
	if err != nil {
		jsonutil.SendHTTPInternalServerErrorTextResponse(w)
	}

	jsonutil.SendHTTPResponse(w, http.StatusOK, TopicManagerDocumentationHandlerResponse{Message: "OK"})
}

// NewTopicManagerDocumentationHandler returns an instance of a TopicManagerDocumentationHandler, utilizing
// an implementation of TopicManagerDocumentationProvider. If the provided argument is nil, it triggers a panic.
func NewTopicManagerDocumentationHandler(provider TopicManagerDocumentationProvider) *TopicManagerDocumentationHandler {
	if provider == nil {
		panic("topic manager documentation provider is nil")
	}
	return &TopicManagerDocumentationHandler{
		provider: provider,
	}
}
