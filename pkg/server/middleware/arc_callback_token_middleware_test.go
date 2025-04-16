package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/4chain-ag/go-overlay-services/pkg/server/middleware"
	"github.com/stretchr/testify/require"
)

func TestARCCallbackTokenMiddleware(t *testing.T) {
	tests := map[string]struct {
		setupRequest          func(r *http.Request)
		expectedStatus        int
		expectedCallbackToken string
		expectedResponse      middleware.FailureResponse
	}{
		"should succeed with 200 when ARC callback token matches the configured key": {
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer 234c13dd-db82-48a5-bb5d-69381aa5478a")
			},
			expectedStatus:        http.StatusOK,
			expectedCallbackToken: "234c13dd-db82-48a5-bb5d-69381aa5478a",
		},
		"should fail with 404 when ARC callback token is not configured": {
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer 7c3c81fa-f732-4e48-b088-7d29ec0bd3bf")
			},
			expectedStatus:        http.StatusNotFound,
			expectedCallbackToken: "",
			expectedResponse:      middleware.EndpointNotSupportedResponse,
		},
		"should fail with 401 when Authorization header is missing": {
			setupRequest:          func(r *http.Request) {},
			expectedStatus:        http.StatusUnauthorized,
			expectedCallbackToken: "valid-token",
			expectedResponse:      middleware.MissingAuthHeaderResponse,
		},
		"should fail with 401 when Authorization header doesn't have Bearer prefix": {
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Token 7c3c81fa-f732-4e48-b088-7d29ec0bd3bf")
			},
			expectedStatus:        http.StatusUnauthorized,
			expectedCallbackToken: "valid-token",
			expectedResponse:      middleware.MissingAuthHeaderValueResponse,
		},
		"should fail with 401 when Authorization header has Bearer prefix but no token": {
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer ")
			},
			expectedStatus:        http.StatusUnauthorized,
			expectedCallbackToken: "valid-token",
			expectedResponse:      middleware.MissingAuthHeaderValueResponse,
		},
		"should fail with 403 when token doesn't match expected token": {
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer wrong-token")
			},
			expectedStatus:        http.StatusForbidden,
			expectedCallbackToken: "valid-token",
			expectedResponse:      middleware.InvalidBearerTokenValueResponse,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			handler := middleware.ARCCallbackTokenMiddleware(tc.expectedCallbackToken)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			ts := httptest.NewServer(handler)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			require.NoError(t, err)

			if tc.setupRequest != nil {
				tc.setupRequest(req)
			}

			// when:
			resp, err := ts.Client().Do(req)

			// then:
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatus, resp.StatusCode)

			if resp.StatusCode != http.StatusOK {
				var actual middleware.FailureResponse
				require.NoError(t, jsonutil.DecodeResponseBody(resp, &actual))
				require.Equal(t, tc.expectedResponse, actual)
			}
		})
	}
}
