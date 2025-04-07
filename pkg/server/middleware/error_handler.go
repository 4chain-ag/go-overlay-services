package middleware

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents the standardized JSON error format
type ErrorResponse struct {
	Error string `json:"error"`
}

// ErrorHandlerMiddleware is a middleware that catches panics and standardizes error responses
func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				json.NewEncoder(w).Encode(ErrorResponse{
					Error: http.StatusText(http.StatusInternalServerError),
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
