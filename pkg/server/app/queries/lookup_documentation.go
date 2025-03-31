package queries

import (
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
)

// LookupDocumentationHandlerResponse defines the response body content that
// will be sent in JSON format after successfully processing the handler logic.
type LookupDocumentationHandlerResponse struct {
	Documentation string `json:"documentation"`
}

// ErrorResponse defines a simple error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

// LookupDocumentationProvider defines the contract that must be fulfilled
// to send a lookup service documentation request to the overlay engine for further processing.
// Note: The contract definition is still in development and will be updated after
// migrating the engine code.
type LookupDocumentationProvider interface {
	GetDocumentationForLookupServiceProvider(lookupService string) (string, error)
}

// LookupDocumentationHandler orchestrates the processing flow of a lookup documentation
// request, including the request parameter validation, converting the request
// into an overlay-engine-compatible format, and applying any other necessary
// logic before invoking the engine. It returns the requested lookup service
// documentation in the text/markdown format.
type LookupDocumentationHandler struct {
	provider LookupDocumentationProvider
}

// Handle orchestrates the processing flow of a lookup documentation request.
// It extracts the lookupService query parameter, invokes the engine provider,
// and returns the a Markdown-formatted documentation string as JSON with the appropriate status code.
func (l *LookupDocumentationHandler) Handle(w http.ResponseWriter, r *http.Request) {
	lookupService := r.URL.Query().Get("lookupService")
	if lookupService == "" {
		http.Error(w, "lookupService query parameter is required", http.StatusBadRequest)
		return
	}

	documentation, err := l.provider.GetDocumentationForLookupServiceProvider(lookupService)
	if err != nil {
		jsonutil.SendHTTPInternalServerErrorTextResponse(w)
		return
	}

	jsonutil.SendHTTPResponse(w, http.StatusOK, LookupDocumentationHandlerResponse{
		Documentation: documentation,
	})
}

// NewLookupDocumentationHandler returns an instance of a LookupDocumentationHandler, utilizing
// an implementation of LookupDocumentationProvider. If the provided argument is nil, it panics.
func NewLookupDocumentationHandler(provider LookupDocumentationProvider) *LookupDocumentationHandler {
	if provider == nil {
		panic("provider cannot be nil")
	}
	return &LookupDocumentationHandler{
		provider: provider,
	}
}
