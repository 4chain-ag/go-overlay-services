package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/4chain-ag/go-overlay-services/pkg/server/middleware"
)

func TestErrorHandlerMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		handlerFunc    http.HandlerFunc
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Catch panic and return 500",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				panic("something went wrong")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal Server Error",
		},
		{
			name: "Pass through without panic (200 OK)",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Everything fine"))
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Everything fine",
		},
		{
			name: "Explicitly return 400 Bad Request",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Bad Request", http.StatusBadRequest)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad Request",
		},
		{
			name: "Explicitly return 401 Unauthorized",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Unauthorized",
		},
		{
			name: "Explicitly return 403 Forbidden",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Forbidden", http.StatusForbidden)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Forbidden",
		},
		{
			name: "Explicitly return 404 Not Found",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", http.StatusNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Not Found",
		},
		{
			name: "Explicitly return 405 Method Not Allowed",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method Not Allowed",
		},
		{
			name: "Explicitly return 429 Too Many Requests",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			},
			expectedStatus: http.StatusTooManyRequests,
			expectedBody:   "Too Many Requests",
		},
		{
			name: "Explicitly return 500 Internal Server Error",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal Server Error",
		},
		{
			name: "Explicitly return 503 Service Unavailable",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   "Service Unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			handler := middleware.ErrorHandlerMiddleware(tt.handlerFunc)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			// When
			handler.ServeHTTP(rec, req)

			// Then
			require.Equal(t, tt.expectedStatus, rec.Code)
			require.Contains(t, rec.Body.String(), tt.expectedBody)
		})
	}
}
