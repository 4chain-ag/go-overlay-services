package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/stretchr/testify/require"
)

func Test_AuthorizationBearerTokenMiddleware(t *testing.T) {
	// Given
	adminToken := "valid_admin_token"
	ts := httptest.NewServer(server.AuthorizationBearerTokenMiddleware(adminToken))
	defer ts.Close()

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "Route access without a token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Route access with an invalid token",
			token:          "invalid_token",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Route access with a valid token",
			token:          "valid_admin_token",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			require.NoError(t, err)

			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			resp, err := ts.Client().Do(req)

			// Then
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
