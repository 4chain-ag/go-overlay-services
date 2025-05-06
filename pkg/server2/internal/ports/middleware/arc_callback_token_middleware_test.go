package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestArcCallbackTokenMiddleware(t *testing.T) {
	tests := map[string]struct {
		authHeader           string
		expectedStatus       int
		arcCallbackToken     string
		arcApiKey            string
		expectedResponseBody string
	}{
		"should succeed with 200 when Arc api key is provided and Arc callback token matches the configured key": {
			authHeader:           "Bearer valid-callback-token",
			expectedStatus:       fiber.StatusOK,
			arcApiKey:            "valid-arc-api-key",
			arcCallbackToken:     "valid-callback-token",
			expectedResponseBody: "success",
		},
		"should succeed with 200 when Arc api key is provided and Arc callback token is empty": {
			authHeader:           "Bearer anything",
			expectedStatus:       fiber.StatusOK,
			arcApiKey:            "valid-arc-api-key",
			arcCallbackToken:     "",
			expectedResponseBody: "success",
		},
		"should fail with 404 when Arc api key token is not configured and Arc callback token matches the configured key": {
			authHeader:       "Bearer valid-callback-token",
			expectedStatus:   fiber.StatusNotFound,
			arcCallbackToken: "valid-callback-token",
			arcApiKey:        "",
		},
		"should fail with 404 when Arc api key token is not configured and Arc callback token is empty": {
			authHeader:       "",
			expectedStatus:   fiber.StatusNotFound,
			arcCallbackToken: "",
			arcApiKey:        "",
		},
		"should fail with 401 when Authorization header is missing": {
			authHeader:       "",
			expectedStatus:   fiber.StatusUnauthorized,
			arcCallbackToken: "valid-callback-token",
			arcApiKey:        "valid-arc-api-key",
		},
		"should fail with 401 when Authorization header doesn't have Bearer prefix": {
			authHeader:       "IncorrectPrefix valid-callback-token",
			expectedStatus:   fiber.StatusUnauthorized,
			arcCallbackToken: "valid-callback-token",
			arcApiKey:        "valid-arc-api-key",
		},
		"should fail with 401 when Authorization header has Bearer prefix but no token": {
			authHeader:       "Bearer ",
			expectedStatus:   fiber.StatusUnauthorized,
			arcCallbackToken: "valid-callback-token",
			arcApiKey:        "valid-arc-api-key",
		},
		"should fail with 403 when call back token doesn't match expected token": {
			authHeader:       "Bearer wrong-callback-token",
			expectedStatus:   fiber.StatusForbidden,
			arcCallbackToken: "valid-callback-token",
			arcApiKey:        "valid-arc-api-key",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Set up a Fiber app with the middleware
			app := fiber.New()
			app.Use(middleware.ArcCallbackTokenMiddleware(tc.arcCallbackToken, tc.arcApiKey))

			// Set up a success handler
			app.Get("/test", func(c *fiber.Ctx) error {
				return c.SendString("success")
			})

			// Create a new HTTP request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			// Perform the request
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify status code
			require.Equal(t, tc.expectedStatus, resp.StatusCode)

			// For success cases, verify the response body
			if tc.expectedStatus == fiber.StatusOK && tc.expectedResponseBody != "" {
				bodyBytes, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				require.Contains(t, string(bodyBytes), tc.expectedResponseBody)
			}
		})
	}
}
