package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// TestAdminRoutesProtection verifies that the admin routes are protected by the AdminAuthMiddlewareFiber.
func TestAdminRoutesProtection(t *testing.T) {
	// Given
	adminToken := "valid_admin_token"
	app := fiber.New()

	adminGroup := app.Group("/admin", server.AdminAuth(adminToken))
	adminGroup.Post("/advertisements-sync", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	tests := []struct {
		name           string
		method         string
		url            string
		token          string
		expectedStatus int
	}{
		{
			name:           "Access admin route without token",
			method:         http.MethodPost,
			url:            "/admin/advertisements-sync",
			token:          "",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "Access admin route with invalid token",
			method:         http.MethodPost,
			url:            "/admin/advertisements-sync",
			token:          "invalid_token",
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:           "Access admin route with valid token",
			method:         http.MethodPost,
			url:            "/admin/advertisements-sync",
			token:          "valid_admin_token",
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			req := httptest.NewRequest(tt.method, tt.url, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			resp, err := app.Test(req)

			// Then
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
