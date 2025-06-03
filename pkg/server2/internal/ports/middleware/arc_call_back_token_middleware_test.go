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

func TestArcCallbackTokenMiddleware_ValidCases(t *testing.T) {
	testPaths := []struct {
		endpoint               string
		method                 string
		expectedProviderToCall testabilities.TestOverlayEngineStubOption
	}{
		{
			endpoint:               "/api/v1/arc-ingest",
			method:                 fiber.MethodPost,
			expectedProviderToCall: testabilities.WithArcIngestProvider(testabilities.NewServiceTestMerkleProofProviderMock(t, testabilities.ServiceTestMerkleProofProviderExpectations{ArcIngestCall: true})),
		},
	}

	const arcCallbackToken = "valid_arc_callback_token"
	const arcApiKey = "valid_arc_api_key"

	validTxID := testabilities.NewValidTestTxID(t).String()
	validMerklePath := testabilities.NewValidTestMerklePath(t)

	tests := map[string]struct {
		expectedStatus int
		headers        map[string]string
		body           interface{}
	}{
		"Authorization header with a valid ARC callback token": {
			expectedStatus: fiber.StatusOK,
			headers: map[string]string{
				fiber.HeaderAuthorization: "Bearer " + arcCallbackToken,
				fiber.HeaderContentType:   fiber.MIMEApplicationJSON,
			},
			body: map[string]interface{}{
				"txid":        validTxID,
				"merklePath":  validMerklePath,
				"blockHeight": 0,
			},
		},
		"No Authorization header when ARC callback token is empty": {
			expectedStatus: fiber.StatusOK,
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
			},
			body: map[string]interface{}{
				"txid":        validTxID,
				"merklePath":  validMerklePath,
				"blockHeight": 0,
			},
		},
	}

	for _, path := range testPaths {
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				// given:
				expectations := testabilities.ServiceTestMerkleProofProviderExpectations{
					ArcIngestCall:      true,
					ExpectedTxID:       validTxID,
					ExpectedMerklePath: validMerklePath,
				}
				stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithArcIngestProvider(testabilities.NewServiceTestMerkleProofProviderMock(t, expectations)))
				var fixture *server2.ServerTestFixture

				if name == "No Authorization header when ARC callback token is empty" {
					fixture = server2.NewServerTestFixture(
						t,
						server2.WithArcCallbackToken(""),
						server2.WithArcApiKey(arcApiKey),
						server2.WithEngine(stub),
					)
				} else {
					fixture = server2.NewServerTestFixture(
						t,
						server2.WithArcCallbackToken(arcCallbackToken),
						server2.WithArcApiKey(arcApiKey),
						server2.WithEngine(stub),
					)
				}

				// when:
				res, _ := fixture.Client().
					R().
					SetHeaders(tc.headers).
					SetBody(tc.body).
					Execute(path.method, path.endpoint)

				// then:
				require.Equal(t, tc.expectedStatus, res.StatusCode())
				stub.AssertProvidersState()
			})
		}
	}
}

func TestArcCallbackTokenMiddleware_InvalidCases(t *testing.T) {
	testPaths := []struct {
		endpoint string
		method   string
	}{
		{
			endpoint: "/api/v1/arc-ingest",
			method:   fiber.MethodPost,
		},
	}

	const arcCallbackToken = "valid_arc_callback_token"
	const arcApiKey = "valid_arc_api_key"

	validTxID := testabilities.NewValidTestTxID(t).String()
	validMerklePath := testabilities.NewValidTestMerklePath(t)

	tests := map[string]struct {
		expectedResponse openapi.BadRequestResponse
		expectedStatus   int
		headers          map[string]string
		body             interface{}
	}{
		"Authorization header with invalid ARC callback token": {
			expectedStatus:   fiber.StatusForbidden,
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, middleware.NewArcInvalidBearerTokenError()),
			headers: map[string]string{
				fiber.HeaderAuthorization: "Bearer " + "invalid_token",
				fiber.HeaderContentType:   fiber.MIMEApplicationJSON,
			},
			body: map[string]interface{}{
				"txid":        validTxID,
				"merklePath":  validMerklePath,
				"blockHeight": 0,
			},
		},
		"Missing Authorization header in the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, middleware.NewArcMissingAuthHeaderError()),
			headers: map[string]string{
				"RandomHeader":          "Bearer " + arcCallbackToken,
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
			},
			body: map[string]interface{}{
				"txid":        validTxID,
				"merklePath":  validMerklePath,
				"blockHeight": 0,
			},
		},
		"Missing Authorization header value in the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, middleware.NewArcMissingAuthHeaderError()),
			headers: map[string]string{
				fiber.HeaderAuthorization: "",
				fiber.HeaderContentType:   fiber.MIMEApplicationJSON,
			},
			body: map[string]interface{}{
				"txid":        validTxID,
				"merklePath":  validMerklePath,
				"blockHeight": 0,
			},
		},
		"Invalid Bearer scheme in the Authorization header appended to the HTTP request": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, middleware.NewArcMissingBearerTokenError()),
			headers: map[string]string{
				fiber.HeaderAuthorization: "InvalidScheme " + arcCallbackToken,
				fiber.HeaderContentType:   fiber.MIMEApplicationJSON,
			},
			body: map[string]interface{}{
				"txid":        validTxID,
				"merklePath":  validMerklePath,
				"blockHeight": 0,
			},
		},
		"ARC endpoint not supported when API key is empty": {
			expectedStatus:   fiber.StatusBadRequest,
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, middleware.NewArcEndpointNotSupportedError()),
			headers: map[string]string{
				fiber.HeaderAuthorization: "Bearer " + arcCallbackToken,
				fiber.HeaderContentType:   fiber.MIMEApplicationJSON,
			},
			body: map[string]interface{}{
				"txid":        validTxID,
				"merklePath":  validMerklePath,
				"blockHeight": 0,
			},
		},
	}

	for _, path := range testPaths {
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				// given:
				stub := testabilities.NewTestOverlayEngineStub(t)
				var fixture *server2.ServerTestFixture

				if name == "ARC endpoint not supported when API key is empty" {
					fixture = server2.NewServerTestFixture(
						t,
						server2.WithArcCallbackToken(arcCallbackToken),
						server2.WithArcApiKey(""),
						server2.WithEngine(stub),
					)
				} else {
					fixture = server2.NewServerTestFixture(
						t,
						server2.WithArcCallbackToken(arcCallbackToken),
						server2.WithArcApiKey(arcApiKey),
						server2.WithEngine(stub),
					)
				}

				// when:
				var actual openapi.BadRequestResponse

				res, _ := fixture.Client().
					R().
					SetHeaders(tc.headers).
					SetBody(tc.body).
					SetError(&actual).
					Execute(path.method, path.endpoint)

				// then:
				require.Equal(t, tc.expectedStatus, res.StatusCode())
				require.Equal(t, tc.expectedResponse, actual)
				stub.AssertProvidersState()
			})
		}
	}
}
