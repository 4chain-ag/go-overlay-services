package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/4chain-ag/go-overlay-services/pkg/server/middleware"
	"github.com/stretchr/testify/require"
)

// TODO: Implement missing test cases..
func TestARCCallbackTokenMiddleware(t *testing.T) {
	tests := map[string]struct {
		actutalToken          string
		expectedStatus        int
		expectedCallbackToken string
		expectedResponse      middleware.MiddlewareFailureResponse
	}{
		"should success with 200 when ARC callback token matches to the configured key": {
			expectedStatus:        http.StatusOK,
			actutalToken:          "234c13dd-db82-48a5-bb5d-69381aa5478a",
			expectedCallbackToken: "234c13dd-db82-48a5-bb5d-69381aa5478a",
		},
		"should fail with 404 when ARC callback token is not configured": {
			expectedStatus:        http.StatusNotFound,
			actutalToken:          "7c3c81fa-f732-4e48-b088-7d29ec0bd3bf",
			expectedCallbackToken: "",
			expectedResponse:      middleware.EndpointNotSupportedResponse,
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

			if tc.actutalToken != "" {
				req.Header.Set("Authorization", "Bearer "+tc.actutalToken)
			}

			// when:
			resp, err := ts.Client().Do(req)

			// then:
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatus, resp.StatusCode)

			if resp.StatusCode != http.StatusOK {
				var actual middleware.MiddlewareFailureResponse
				require.NoError(t, jsonutil.DecodeResponseBody(resp, &actual))
				require.Equal(t, tc.expectedResponse, actual)
			}
		})
	}
}
