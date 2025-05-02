package middleware_test

import (
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
	fixture := server2.NewTestFixture(
		t,
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
		header           map[string]string
	}{
		"Authorization header with a valid HTTP server token": {
			expectedStatus: fiber.StatusOK,
			header: map[string]string{
				fiber.HeaderAuthorization: "Bearer " + bearerToken,
			},
		},
		"Authorization header with ivalid HTTP server token": {
			expectedStatus:   fiber.StatusForbidden,
			expectedResponse: middleware.InvalidBearerTokenValueResponse,
			header: map[string]string{
				fiber.HeaderAuthorization: "Bearer " + "1234",
			},
		},
		"Missing Authorization header in the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: middleware.MissingAuthorizationHeaderResponse,
			header: map[string]string{
				"RandomHeader": "Bearer " + bearerToken,
			},
		},
		"Missing Authorization header value in the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: middleware.MissingAuthorizationHeaderResponse,
			header: map[string]string{
				fiber.HeaderAuthorization: "",
			},
		},
		"Invalid Bearer scheme in the Authorization header appended to the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: middleware.MissingAuthorizationHeaderBearerTokenValueResponse,
			header: map[string]string{
				fiber.HeaderAuthorization: "InvalidScheme " + bearerToken,
			},
		},
	}

	for _, path := range testPaths {
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				// when:
				var actual openapi.BadRequestResponse

				res, _ := fixture.Client().
					R().
					SetHeaders(tc.header).
					SetError(&actual).
					Execute(path.method, path.endpoint)

				// then:
				require.Equal(t, tc.expectedStatus, res.StatusCode())
				require.Equal(t, tc.expectedResponse, actual, "unexpected error response from Bearer token authorization middleware")
			})
		}
	}
}
