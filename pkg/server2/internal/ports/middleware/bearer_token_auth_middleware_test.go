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

func TestBearerTokenAuthMiddleware_ValidCases(t *testing.T) {
	testPaths := []struct {
		endpoint               string
		method                 string
		expectedProviderToCall testabilities.TestOverlayEngineStubOption
	}{
		{
			endpoint:               "/api/v1/admin/syncAdvertisements",
			method:                 fiber.MethodPost,
			expectedProviderToCall: testabilities.WithSyncAdvertisementsProvider(testabilities.NewSyncAdvertisementsProviderMock(t, testabilities.SyncAdvertisementsProviderMockExpectations{SyncAdvertisementsCall: true})),
		},
	}

	const bearerToken = "valid_admin_token"

	tests := map[string]struct {
		expectedStatus int
		headers        map[string]string
	}{
		"Authorization header with a valid HTTP server token": {
			expectedStatus: fiber.StatusOK,
			headers: map[string]string{
				fiber.HeaderAuthorization: "Bearer " + bearerToken,
			},
		},
	}

	for _, path := range testPaths {
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				// given:
				stub := testabilities.NewTestOverlayEngineStub(t,
					path.expectedProviderToCall,
					// .. TODO: Update the providers list that should not be called.
					testabilities.WithSubmitTransactionProvider(testabilities.NewSubmitTransactionProviderMock(t, testabilities.SubmitTransactionProviderMockExpectations{SubmitCall: false})),
				)

				fixture := server2.NewServerTestFixture(
					t,
					server2.WithAdminBearerToken(bearerToken),
					server2.WithEngine(stub),
				)

				// when:
				res, _ := fixture.Client().
					R().
					SetHeaders(tc.headers).
					Execute(path.method, path.endpoint)

				// then:
				require.Equal(t, tc.expectedStatus, res.StatusCode(), "mismatch between the expected and actual response status codes")
				stub.AssertProvidersState()
			})
		}
	}
}

func TestBearerTokenAuthMiddleware_InvalidCases(t *testing.T) {
	testPaths := []struct {
		endpoint string
		method   string
	}{
		{
			endpoint: "/api/v1/admin/syncAdvertisements",
			method:   fiber.MethodPost,
		},
	}

	const bearerToken = "valid_admin_token"

	tests := map[string]struct {
		expectedResponse openapi.BadRequestResponse
		expectedStatus   int
		headers          map[string]string
	}{
		"Authorization header with ivalid HTTP server token": {
			expectedStatus:   fiber.StatusForbidden,
			expectedResponse: middleware.InvalidBearerTokenValueResponse,
			headers: map[string]string{
				fiber.HeaderAuthorization: "Bearer " + "1234",
			},
		},
		"Missing Authorization header in the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: middleware.MissingAuthorizationHeaderResponse,
			headers: map[string]string{
				"RandomHeader": "Bearer " + bearerToken,
			},
		},
		"Missing Authorization header value in the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: middleware.MissingAuthorizationHeaderResponse,
			headers: map[string]string{
				fiber.HeaderAuthorization: "",
			},
		},
		"Invalid Bearer scheme in the Authorization header appended to the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: middleware.MissingAuthorizationHeaderBearerTokenValueResponse,
			headers: map[string]string{
				fiber.HeaderAuthorization: "InvalidScheme " + bearerToken,
			},
		},
	}

	for _, path := range testPaths {
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				// given:
				stub := testabilities.NewTestOverlayEngineStub(t,
					// .. TODO: Update the providers list that should not be called.
					testabilities.WithSubmitTransactionProvider(testabilities.NewSubmitTransactionProviderMock(t, testabilities.SubmitTransactionProviderMockExpectations{SubmitCall: false})),
					testabilities.WithSyncAdvertisementsProvider(testabilities.NewSyncAdvertisementsProviderMock(t, testabilities.SyncAdvertisementsProviderMockExpectations{SyncAdvertisementsCall: false})),
				)

				fixture := server2.NewServerTestFixture(
					t,
					server2.WithAdminBearerToken(bearerToken),
					server2.WithEngine(stub),
				)

				// when:
				var actual openapi.BadRequestResponse

				res, _ := fixture.Client().
					R().
					SetHeaders(tc.headers).
					SetError(&actual).
					Execute(path.method, path.endpoint)

				// then:
				require.Equal(t, tc.expectedStatus, res.StatusCode())
				require.Equal(t, tc.expectedResponse, actual, "unexpected error response from Bearer token authorization middleware")
				stub.AssertProvidersState()
			})
		}
	}
}
