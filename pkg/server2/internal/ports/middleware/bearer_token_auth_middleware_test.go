package middleware_test

import (
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/middleware"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func Test_BearearTokenAuthorizationMiddleware(t *testing.T) {
	// given:
	const bearerToken = "valid_admin_token"

	mock := testabilities.NewSubmitTransactionProviderMock(t)
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithSubmitTransactionProvider(mock))
	serverAPI := server2.NewServerTestAdapter(
		server2.WithAdminBearerToken(bearerToken),
		server2.WithEngine(engine),
	)

	testPaths := []struct {
		endpoint string
		method   string
	}{
		{
			endpoint: "/api/v1/admin/syncAdvertisements",
			method:   fiber.MethodPost,
		},
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
				headerValue: "Bearer " + bearerToken,
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
				headerValue: "Bearer " + bearerToken,
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
				headerValue: "InvalidScheme " + bearerToken,
			},
		},
	}

	for _, path := range testPaths {
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				// given:
				req, err := http.NewRequest(path.method, path.endpoint, nil)
				require.NoError(t, err, "failed to create a new HTTP request")
				req.Header.Set(tc.header.headerKey, tc.header.headerValue)

				// when:
				res, err := serverAPI.TestRequest(req, -1)

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
