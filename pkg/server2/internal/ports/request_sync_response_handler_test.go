package ports_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestRequestSyncResponseHandler_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		payload            interface{}
		headers            map[string]string
		expectations       testabilities.RequestSyncResponseProviderMockExpectations
		expectedStatusCode int
	}{
		"Request sync response handler fails due to missing topic header": {
			payload:            testabilities.DefaultMockRequestPayload,
			headers:            testabilities.MissingTopicHeaders,
			expectedStatusCode: fiber.StatusBadRequest,
		},
		"Request sync response handler fails due to invalid JSON": {
			payload:            "INVALID_JSON",
			headers:            testabilities.DefaultMockHeaders,
			expectedStatusCode: fiber.StatusBadRequest,
		},
		"Request sync response handler fails due to provider error": {
			payload:            testabilities.DefaultMockRequestPayload,
			headers:            testabilities.DefaultMockHeaders,
			expectations:       testabilities.NewErrorResponseExpectations(errors.New("provider error")),
			expectedStatusCode: fiber.StatusInternalServerError,
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
			res, _ := fixture.Client().
				R().
				SetHeaders(tc.headers).
				SetBody(tc.payload).
				Post("/api/v1/requestSyncResponse")

			// then:
			require.Equal(t, tc.expectedStatusCode, res.StatusCode())
			if tc.expectations.ProvideForeignSyncResponseCall {
				stub.AssertProvidersState()
			}
		})
	}
}

func TestRequestSyncResponseHandler_ValidCases(t *testing.T) {
	tests := map[string]struct {
		payload            interface{}
		headers            map[string]string
		expectations       testabilities.RequestSyncResponseProviderMockExpectations
		expectedStatusCode int
		expectedResponse   openapi.RequestSyncResResponse
	}{
		"Request sync response handler succeeds with empty UTXO list": {
			payload:            testabilities.DefaultMockRequestPayload,
			headers:            testabilities.DefaultMockHeaders,
			expectations:       testabilities.NewEmptyResponseExpectations(),
			expectedStatusCode: http.StatusOK,
			expectedResponse: openapi.RequestSyncResResponse{
				UTXOList: []openapi.UTXOItem{},
				Since:    0,
			},
		},
		"Request sync response handler succeeds with single UTXO": {
			payload:            testabilities.DefaultMockRequestPayload,
			headers:            testabilities.DefaultMockHeaders,
			expectations:       testabilities.NewSingleUTXOResponseExpectations(),
			expectedStatusCode: http.StatusOK,
			expectedResponse: openapi.RequestSyncResResponse{
				UTXOList: []openapi.UTXOItem{
					{
						Txid: "03895fb984362a4196bc9931629318fcbb2aeba7c6293638119ea653fa31d119",
						Vout: 0,
					},
				},
				Since: 1000000,
			},
		},
		"Request sync response handler succeeds with multiple UTXOs": {
			payload:            testabilities.DefaultMockRequestPayload,
			headers:            testabilities.DefaultMockHeaders,
			expectations:       testabilities.DefaultRequestSyncResponseProviderMockExpectations,
			expectedStatusCode: http.StatusOK,
			expectedResponse: openapi.RequestSyncResResponse{
				UTXOList: []openapi.UTXOItem{
					{
						Txid: "03895fb984362a4196bc9931629318fcbb2aeba7c6293638119ea653fa31d119",
						Vout: 0,
					},
					{
						Txid: "27c8f37851aabc468d3dbb6bf0789dc398a602dcb897ca04e7815d939d621595",
						Vout: 1,
					},
					{
						Txid: "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
						Vout: 2,
					},
				},
				Since: 1234567890,
			},
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
			var actualResponse openapi.RequestSyncResResponse
			res, _ := fixture.Client().
				R().
				SetHeaders(tc.headers).
				SetBody(tc.payload).
				SetResult(&actualResponse).
				Post("/api/v1/requestSyncResponse")

			// then:
			require.Equal(t, tc.expectedStatusCode, res.StatusCode())
			require.Equal(t, tc.expectedResponse, actualResponse)
			stub.AssertProvidersState()
		})
	}
}
