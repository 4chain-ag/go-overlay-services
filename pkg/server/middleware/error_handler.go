package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// ErrorResponse represents the standardized JSON error format
type ErrorResponse struct {
	Error string `json:"error"`
}

// ErrorHandlerMiddleware catches panics and standardizes errors into JSON.
func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal Server Error"}); err != nil {
					slog.Error("Failed to write error response", "error", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
