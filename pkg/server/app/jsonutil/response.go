package jsonutil

import (
	"encoding/json"
	"net/http"
)

const (
	// ReasonBadRequest represents an HTTP 400 Bad Request error reason.
	ReasonBadRequest = "Bad Request"

	// ReasonInvalidRequest represents a request that is syntactically valid but semantically invalid.
	ReasonInvalidRequest = "Invalid Request"

	// ReasonInternalError represents a generic HTTP 500 Internal Server Error reason.
	ReasonInternalError = "Internal Server Error"

	// ReasonUnauthorized represents an HTTP 401 Unauthorized error reason.
	ReasonUnauthorized = "Unauthorized"

	// ReasonNotFound represents an HTTP 404 Not Found error reason.
	ReasonNotFound = "Resource Not Found"
)

// ResponseFailure defines the standardized JSON structure for HTTP failure responses.
type ResponseFailure struct {
	Reason string `json:"reason"`
	Hint   string `json:"hint"`
}

// SendHTTPFailureResponse sends a standardized JSON failure response to the client.
// It sets the Content-Type to application/json and writes the provided status code, reason, and hint.
func SendHTTPFailureResponse(w http.ResponseWriter, statusCode int, reason, hint string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(ResponseFailure{
		Reason: reason,
		Hint:   hint,
	})
}
