package ports_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestArcAuthorizationDecorator_InvalidCases(t *testing.T) {
	const arcCallbackToken = "valid-token"
	const arcApiKey = "test-api-key"

	validTxID := testabilities.NewValidTestTxID(t).String()
	validMerklePath := testabilities.NewValidTestMerklePath(t)

	payload := map[string]interface{}{
		"txid":        validTxID,
		"merklePath":  validMerklePath,
		"blockHeight": 0,
	}

	tests := map[string]struct {
		expectedStatus   int
		expectedResponse openapi.Error
		headers          map[string]string
		expectations     testabilities.ServiceTestMerkleProofProviderExpectations
	}{
		"Missing Authorization header": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, ports.NewArcMissingAuthHeaderError()),
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
			},
			expectations: testabilities.ServiceTestMerkleProofProviderExpectations{
				ArcIngestCall: false,
			},
		},
		"Authorization header without Bearer prefix": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, ports.NewArcMissingBearerTokenError()),
			headers: map[string]string{
				fiber.HeaderContentType:   fiber.MIMEApplicationJSON,
				fiber.HeaderAuthorization: "Basic sometoken",
			},
			expectations: testabilities.ServiceTestMerkleProofProviderExpectations{
				ArcIngestCall: false,
			},
		},
		"Authorization header with Bearer but no token": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, ports.NewArcMissingBearerTokenError()),
			headers: map[string]string{
				fiber.HeaderContentType:   fiber.MIMEApplicationJSON,
				fiber.HeaderAuthorization: "Bearer ",
			},
			expectations: testabilities.ServiceTestMerkleProofProviderExpectations{
				ArcIngestCall: false,
			},
		},
		"Authorization header with Bearer prefix only": {
			expectedStatus:   fiber.StatusUnauthorized,
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, ports.NewArcMissingBearerTokenError()),
			headers: map[string]string{
				fiber.HeaderContentType:   fiber.MIMEApplicationJSON,
				fiber.HeaderAuthorization: "Bearer",
			},
			expectations: testabilities.ServiceTestMerkleProofProviderExpectations{
				ArcIngestCall: false,
			},
		},
		"Authorization header with invalid Bearer token": {
			expectedStatus:   fiber.StatusForbidden,
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, ports.NewArcInvalidBearerTokenError()),
			headers: map[string]string{
				fiber.HeaderContentType:   fiber.MIMEApplicationJSON,
				fiber.HeaderAuthorization: "Bearer invalidtoken",
			},
			expectations: testabilities.ServiceTestMerkleProofProviderExpectations{
				ArcIngestCall: false,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewServiceTestMerkleProofProviderMock(t, tc.expectations)
			stub := testabilities.NewTestOverlayEngineStub(t,
				testabilities.WithArcIngestProvider(mock),
			)

			fixture := server2.NewServerTestFixture(t,
				server2.WithArcApiKey(arcApiKey),
				server2.WithArcCallbackToken(arcCallbackToken),
				server2.WithEngine(stub),
			)

			// when:
			var actualResponse openapi.Error
			res, _ := fixture.Client().
				R().
				SetHeaders(tc.headers).
				SetBody(payload).
				SetError(&actualResponse).
				Post("/api/v1/arc-ingest")

			// then:
			require.Equal(t, tc.expectedStatus, res.StatusCode())
			require.Equal(t, tc.expectedResponse, actualResponse)
			mock.AssertCalled()
		})
	}
}

func TestArcAuthorizationDecorator_ValidCase(t *testing.T) {
	// given:
	const arcCallbackToken = "valid-token"
	const arcApiKey = "test-api-key"

	validTxID := testabilities.NewValidTestTxID(t).String()
	validMerklePath := testabilities.NewValidTestMerklePath(t)

	expectations := testabilities.ServiceTestMerkleProofProviderExpectations{
		ArcIngestCall:      true,
		ExpectedTxID:       validTxID,
		ExpectedMerklePath: validMerklePath,
		Error:              nil,
	}

	mock := testabilities.NewServiceTestMerkleProofProviderMock(t, expectations)
	stub := testabilities.NewTestOverlayEngineStub(t,
		testabilities.WithArcIngestProvider(mock),
	)

	fixture := server2.NewServerTestFixture(t,
		server2.WithArcApiKey(arcApiKey),
		server2.WithArcCallbackToken(arcCallbackToken),
		server2.WithEngine(stub),
	)

	headers := map[string]string{
		fiber.HeaderContentType:   fiber.MIMEApplicationJSON,
		fiber.HeaderAuthorization: "Bearer " + arcCallbackToken,
	}

	body := map[string]interface{}{
		"txid":        validTxID,
		"merklePath":  validMerklePath,
		"blockHeight": 0,
	}

	expectedResponse := ports.NewArcIngestSuccessResponse()

	// when:
	var actualResponse openapi.ArcIngest
	res, _ := fixture.Client().
		R().
		SetHeaders(headers).
		SetBody(body).
		SetResult(&actualResponse).
		Post("/api/v1/arc-ingest")

	// then:
	require.Equal(t, fiber.StatusOK, res.StatusCode())
	require.Equal(t, expectedResponse, &actualResponse)
	mock.AssertCalled()
}
