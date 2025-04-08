package queries

import (
	"fmt"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
)

// TopicManagerDocumentationHandlerResponse defines the response body content that
// will be sent in JSON format after successfully processing the handler logic.
type TopicManagerDocumentationHandlerResponse struct {
	Documentation string `json:"documentation"`
}

// TopicManagerDocumentationProvider defines the contract that must be fulfilled
// to send a topic manager documentation request to the overlay engine for further processing.
// Note: The contract definition is still in development and will be updated after
// migrating the engine code.
type TopicManagerDocumentationProvider interface {
	GetDocumentationForTopicManager(topicManager string) (string, error)
}

// TopicManagerDocumentationHandler orchestrates the processing flow of a topic manager documentation
// request, including the request parameter validation, converting the request
// into an overlay-engine-compatible format, and applying any other necessary
// logic before invoking the engine. It returns the requested topic manager
// documentation in the text/markdown format.
type TopicManagerDocumentationHandler struct {
	provider TopicManagerDocumentationProvider
}

// Handle orchestrates the processing flow of a topic manager documentation request.
// It extracts the topicManager query parameter, invokes the engine provider,
// and returns a Markdown-formatted documentation string as JSON with the appropriate status code.
func (t *TopicManagerDocumentationHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonutil.SendHTTPFailureResponse(w, http.StatusMethodNotAllowed, jsonutil.ReasonBadRequest, "only GET method is allowed")
		return
	}

	topicManager := r.URL.Query().Get("topicManager")
	if topicManager == "" {
		jsonutil.SendHTTPFailureResponse(w, http.StatusBadRequest, jsonutil.ReasonInvalidRequest, "topicManager query parameter is required")
		return
	}

	documentation, err := t.provider.GetDocumentationForTopicManager(topicManager)
	if err != nil {
		jsonutil.SendHTTPFailureResponse(w, http.StatusInternalServerError, jsonutil.ReasonInternalError, "failed to fetch topic manager documentation")
		return
	}

	jsonutil.SendHTTPResponse(w, http.StatusOK, TopicManagerDocumentationHandlerResponse{
		Documentation: documentation,
	})
}

// NewTopicManagerDocumentationHandler returns an instance of a TopicManagerDocumentationHandler, utilizing
// an implementation of TopicManagerDocumentationProvider. If the provided argument is nil, it panics.
func NewTopicManagerDocumentationHandler(provider TopicManagerDocumentationProvider) (*TopicManagerDocumentationHandler, error) {
	if provider == nil {
		return nil, fmt.Errorf("topic manager documentation provider cannot be nil")
	}
	return &TopicManagerDocumentationHandler{provider: provider}, nil
}
