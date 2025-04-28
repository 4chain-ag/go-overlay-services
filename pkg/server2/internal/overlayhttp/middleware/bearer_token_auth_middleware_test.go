package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/overlayhttp"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/overlayhttp/middleware"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/overlayhttp/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/stretchr/testify/require"
)

func Test_BearearTokenAuthorizationMiddleware(t *testing.T) {
	// given:
	expectedToken := "valid_admin_token"
	serverAPI := &openapi.ServerInterfaceWrapper{Handler: overlayhttp.NewServerHandlers(expectedToken, testabilities.NewTestOverlayEngineStub(t))}

	httpHandlers := []http.Handler{
		adaptor.FiberHandler(serverAPI.AdvertisementsSync), // TODO: Add the missing handlers that require auth check during access.
	}

	tests := map[string]struct {
		expectedResponse openapi.BadRequestResponse
		expectedStatus   int
		header           bearerTokenAuthorizationMiddlewareHeader
	}{
		"Authorization header with a valid HTTP server token": {
			expectedStatus: fiber.StatusOK,
			header: bearerTokenAuthorizationMiddlewareHeader{
				headerKey:   fiber.HeaderAuthorization,
				headerValue: "Bearer " + expectedToken,
			},
		},
		"Authorization header with ivalid HTTP server token": {
			expectedStatus:   fiber.StatusForbidden,
			expectedResponse: middleware.InvalidBearerTokenValueResponse,
			header: bearerTokenAuthorizationMiddlewareHeader{
				headerKey:   fiber.HeaderAuthorization,
				headerValue: "Bearer " + "1234",
			},
		},
		"Missing Authorization header in the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: middleware.MissingAuthorizationHeaderResponse,
			header: bearerTokenAuthorizationMiddlewareHeader{
				headerKey:   "RandomHeader",
				headerValue: "Bearer " + expectedToken,
			},
		},
		"Missing Authorization header value in the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: middleware.MissingAuthorizationHeaderResponse,
			header: bearerTokenAuthorizationMiddlewareHeader{
				headerKey:   fiber.HeaderAuthorization,
				headerValue: "",
			},
		},
		"Invalid Bearer scheme in the Authorization header appended to the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: middleware.MissingAuthorizationHeaderBearerTokenValueResponse,
			header: bearerTokenAuthorizationMiddlewareHeader{
				headerKey:   fiber.HeaderAuthorization,
				headerValue: "InvalidScheme " + expectedToken,
			},
		},
	}

	for _, handler := range httpHandlers {
		ts := httptest.NewServer(handler)
		defer ts.Close()

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				// when:
				req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
				require.NoError(t, err, "failed to create a new HTTP request")
				req.Header.Set(tc.header.headerKey, tc.header.headerValue)

				res, err := ts.Client().Do(req)

				// then:
				require.NoError(t, err)
				require.Equal(t, tc.expectedStatus, res.StatusCode)
				assertBearerTokenAuthorizationMiddlewareResponse(t, res, tc.expectedResponse)
			})
		}
	}
}

// bearerTokenAuthorizationMiddlewareHeader holds parameters for creating a request
// header used to test Bearer token authorization middleware.
type bearerTokenAuthorizationMiddlewareHeader struct {
	headerKey   string
	headerValue string
}

// assertBearerTokenAuthorizationMiddlewareResponse verifies the HTTP response from the Bearer token authorization middleware.
// If the response has a 200 OK status, it passes silently. Otherwise, it decodes the response body
// into an openapi.BadRequestResponse and asserts that it matches the expected response.
func assertBearerTokenAuthorizationMiddlewareResponse(t *testing.T, res *http.Response, expectedResponse any) {
	t.Helper()

	if res.StatusCode == fiber.StatusOK {
		return
	}

	var actual openapi.BadRequestResponse

	testabilities.DecodeResponseBody(t, res, &actual)
	require.Equal(t, expectedResponse, actual, "unexpected error response from Bearer token authorization middleware")
}
