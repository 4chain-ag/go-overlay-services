package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/overlayhttp/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/stretchr/testify/require"
)

func Test_BearearTokenAuthorizationMiddleware(t *testing.T) {
	// given:
	expectedToken := "valid_admin_token"
	next := adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := adaptor.FiberHandler(middleware.BearerTokenAuthorizationMiddleware(expectedToken, next))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	tests := map[string]struct {
		expectedResponse openapi.BadRequestResponse
		expectedStatus   int
		request          *http.Request
	}{
		"Authorization header with a valid HTTP server token": {
			expectedStatus: fiber.StatusOK,
			request: NewBearerTokenAuthorizationMiddlewareRequest(t, BearerTokenAuthorizationMiddlewareRequestParams{
				url:         ts.URL,
				headerKey:   fiber.HeaderAuthorization,
				headerValue: "Bearer " + expectedToken,
			}),
		},
		"Authorization header with ivalid HTTP server token": {
			expectedStatus:   fiber.StatusForbidden,
			expectedResponse: middleware.InvalidBearerTokenValueResponse,
			request: NewBearerTokenAuthorizationMiddlewareRequest(t, BearerTokenAuthorizationMiddlewareRequestParams{
				url:         ts.URL,
				headerKey:   fiber.HeaderAuthorization,
				headerValue: "Bearer " + "1234",
			}),
		},
		"Missing Authorization header in the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: middleware.MissingAuthorizationHeaderResponse,
			request: NewBearerTokenAuthorizationMiddlewareRequest(t, BearerTokenAuthorizationMiddlewareRequestParams{
				url:         ts.URL,
				headerKey:   "RandomHeader",
				headerValue: "Bearer " + expectedToken,
			}),
		},
		"Missing Authorization header value in the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: middleware.MissingAuthorizationHeaderResponse,
			request: NewBearerTokenAuthorizationMiddlewareRequest(t, BearerTokenAuthorizationMiddlewareRequestParams{
				url:         ts.URL,
				headerKey:   fiber.HeaderAuthorization,
				headerValue: "",
			}),
		},
		"Invalid Bearer scheme in the Authorization header appended to the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: middleware.MissingAuthorizationHeaderBearerTokenValueResponse,
			request: NewBearerTokenAuthorizationMiddlewareRequest(t, BearerTokenAuthorizationMiddlewareRequestParams{
				url:         ts.URL,
				headerKey:   fiber.HeaderAuthorization,
				headerValue: "InvalidScheme " + expectedToken,
			}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			res, err := ts.Client().Do(tc.request)

			// then:
			require.NoError(t, err)
			require.Equal(t, tc.expectedStatus, res.StatusCode)
			AssertBearerTokenAuthorizationMiddlewareResponse(t, res, tc.expectedResponse)
		})
	}
}

// BearerTokenAuthorizationMiddlewareRequestParams holds parameters for creating a request
// used to test Bearer token authorization middleware.
type BearerTokenAuthorizationMiddlewareRequestParams struct {
	url         string
	headerKey   string
	headerValue string
}

// NewBearerTokenAuthorizationMiddlewareRequest creates a new HTTP GET request using the given parameters,
// intended for testing Bearer token authorization middleware.
// It fails the test immediately if request creation encounters an error.
func NewBearerTokenAuthorizationMiddlewareRequest(t *testing.T, params BearerTokenAuthorizationMiddlewareRequestParams) *http.Request {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, params.url, nil)
	require.NoError(t, err, "failed to create a new HTTP request")
	req.Header.Set(params.headerKey, params.headerValue)

	return req
}

// DecodeResponseBody attempts to decode the HTTP response body into given destination
// argument. It returns an error if the internal decoding operation fails; otherwise,
// it returns nil, indicating successful processing.
func DecodeResponseBody(t *testing.T, res *http.Response, dst any) {
	t.Helper()

	dec := json.NewDecoder(res.Body)
	err := dec.Decode(dst)
	require.NoError(t, err, "decoding http response body op failure")
}

// AssertBearerTokenAuthorizationMiddlewareResponse verifies the HTTP response from the Bearer token authorization middleware.
// If the response has a 200 OK status, it passes silently. Otherwise, it decodes the response body
// into an openapi.BadRequestResponse and asserts that it matches the expected response.
func AssertBearerTokenAuthorizationMiddlewareResponse(t *testing.T, res *http.Response, expectedResponse any) {
	t.Helper()

	if res.StatusCode == fiber.StatusOK {
		return
	}

	var actual openapi.BadRequestResponse
	DecodeResponseBody(t, res, &actual)
	require.Equal(t, expectedResponse, actual, "unexpected error response from Bearer token authorization middleware")
}
