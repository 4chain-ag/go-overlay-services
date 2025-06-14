package ports_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestRequestSyncResponseHandler_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		payload            any
		headers            map[string]string
		expectations       testabilities.RequestSyncResponseProviderMockExpectations
		expectedStatusCode int
		expectedResponse   openapi.Error
	}{
		"Request sync response handler fails due to missing topic header": {
			payload: testabilities.NewDefaultRequestSyncResponseBody(),
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				ProvideForeignSyncResponseCall: false,
			},
			expectedResponse: openapi.Error{Message: "The submitted request does not include required header: X-BSV-Topic."},
		},
		"Request sync response handler fails due to invalid JSON": {
			payload: "INVALID_JSON",
			headers: map[string]string{
				"Content-Type": "application/json",
				"X-BSV-Topic":  testabilities.DefaultTopic,
			},
			expectedStatusCode: fiber.StatusInternalServerError,
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				ProvideForeignSyncResponseCall: false,
			},
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, ports.NewRequestBodyParserError(testabilities.ErrTestNoopOpFailure)),
		},
		"Request sync response handler fails due to provider error": {
			payload: testabilities.NewDefaultRequestSyncResponseBody(),
			headers: map[string]string{
				"Content-Type": "application/json",
				"X-BSV-Topic":  testabilities.DefaultTopic,
			},
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				Error:                          testabilities.ErrTestNoopOpFailure,
				ProvideForeignSyncResponseCall: true,
				InitialRequest: &core.GASPInitialRequest{
					Version: testabilities.DefaultVersion,
					Since:   uint32(testabilities.DefaultSince),
				},
				Topic: testabilities.DefaultTopic,
			},
			expectedStatusCode: fiber.StatusInternalServerError,
			expectedResponse:   testabilities.NewTestOpenapiErrorResponse(t, app.NewRequestSyncResponseProviderError(testabilities.ErrTestNoopOpFailure)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestSyncResponseProvider(
				testabilities.NewRequestSyncResponseProviderMock(t, tc.expectations),
			))
			fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

			// when:
			var actualResponse openapi.Error

			res, _ := fixture.Client().
				R().
				SetHeaders(tc.headers).
				SetBody(tc.payload).
				SetError(&actualResponse).
				Post("/api/v1/requestSyncResponse")

			// then:
			require.Equal(t, tc.expectedStatusCode, res.StatusCode())
			require.Equal(t, &tc.expectedResponse, &actualResponse)
			stub.AssertProvidersState()
		})
	}
}

func TestRequestSyncResponseHandler_ValidCase(t *testing.T) {
	// given:
	expectations := testabilities.RequestSyncResponseProviderMockExpectations{
		ProvideForeignSyncResponseCall: true,
		InitialRequest: &core.GASPInitialRequest{
			Version: testabilities.DefaultVersion,
			Since:   testabilities.DefaultSince,
		},
		Topic: testabilities.DefaultTopic,
		Response: &core.GASPInitialResponse{
			Since: testabilities.DefaultSince,
			UTXOList: []*transaction.Outpoint{
				{
					Txid:  *testabilities.DummyTxHash(t, "03895fb984362a4196bc9931629318fcbb2aeba7c6293638119ea653fa31d119"),
					Index: 0,
				},
				{
					Txid:  *testabilities.DummyTxHash(t, "27c8f37851aabc468d3dbb6bf0789dc398a602dcb897ca04e7815d939d621595"),
					Index: 1,
				},
				{
					Txid:  *testabilities.DummyTxHash(t, "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"),
					Index: 2,
				},
			},
		},
	}

	expectedDTO := app.NewRequestSyncResponseDTO(expectations.Response)
	expectedResponse := ports.NewRequestSyncResponseSuccessResponse(expectedDTO)
	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestSyncResponseProvider(testabilities.NewRequestSyncResponseProviderMock(t, expectations)))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	headers := map[string]string{
		"X-BSV-Topic":  testabilities.DefaultTopic,
		"Content-Type": fiber.MIMEApplicationJSON,
	}

	// when:
	var actualResponse openapi.RequestSyncResResponse

	res, _ := fixture.Client().
		R().
		SetHeaders(headers).
		SetBody(testabilities.NewDefaultRequestSyncResponseBody()).
		SetResult(&actualResponse).
		Post("/api/v1/requestSyncResponse")

	// then:
	require.Equal(t, fiber.StatusOK, res.StatusCode())
	require.Equal(t, expectedResponse, &actualResponse)
	stub.AssertProvidersState()
}
