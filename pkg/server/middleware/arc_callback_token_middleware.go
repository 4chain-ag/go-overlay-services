package middleware

import (
	"net/http"
	"strings"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
)

// TODO: Add missing docs.
// TODO: Maybe it's better to define functions that return
// specific types of middleware failure response :-)

type MiddlewareFailureResponse struct {
	Message string `json:"message"`
}

var MissingAuthHeaderResponse = MiddlewareFailureResponse{
	Message: "Authorization header is missing from the request.",
}

var MissingAuthHeaderValueResponse = MiddlewareFailureResponse{
	Message: "Authorization header is present, but the Bearer token is missing.",
}

var InvalidBearerTokenValueResponse = MiddlewareFailureResponse{
	Message: "The Bearer token provided is invalid.",
}

var EndpointNotSupportedResponse = MiddlewareFailureResponse{
	Message: "This endpoint is not supported by the current service configuration.",
}

func ARCCallbackTokenMiddleware(expectedToken string) func(http.Handler) http.Handler {
	const schema = "Bearer "
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if expectedToken == "" {
				jsonutil.SendHTTPResponse(w, http.StatusNotFound, EndpointNotSupportedResponse)
				return
			}

			auth := r.Header.Get("Authorization")
			if auth == "" {
				jsonutil.SendHTTPResponse(w, http.StatusUnauthorized, MissingAuthHeaderResponse)
				return
			}

			if !strings.HasPrefix(auth, schema) || len(auth) <= len(schema) {
				jsonutil.SendHTTPResponse(w, http.StatusUnauthorized, MissingAuthHeaderValueResponse)
				return
			}

			token := strings.TrimPrefix(auth, schema)
			if token != expectedToken {
				jsonutil.SendHTTPResponse(w, http.StatusForbidden, InvalidBearerTokenValueResponse)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
