package ports_test

import (
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestRequestSyncResponseHandler_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		payload            interface{}
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
			expectedStatusCode: fiber.StatusBadRequest,
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				ProvideForeignSyncResponseCall: false,
			},
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, app.NewRequestSyncResponseInvalidJSONError()),
		},
		"Request sync response handler fails due to provider error": {
			payload: testabilities.NewDefaultRequestSyncResponseBody(),
			headers: map[string]string{
				"Content-Type": "application/json",
				"X-BSV-Topic":  testabilities.DefaultTopic,
			},
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				Error:                          errors.New("internal submit transaction provider error during submit transaction handler unit test"),
				ProvideForeignSyncResponseCall: true,
				InitialRequest: &core.GASPInitialRequest{
					Version: testabilities.DefaultVersion,
					Since:   uint32(testabilities.DefaultSince),
				},
				Topic: testabilities.DefaultTopic,
			},
			expectedStatusCode: fiber.StatusInternalServerError,
			expectedResponse:   testabilities.NewTestOpenapiErrorResponse(t, app.NewRequestSyncResponseProviderError(errors.New("internal submit transaction provider error during submit transaction handler unit test"))),
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
			Since:   uint32(testabilities.DefaultSince),
		},
		Topic: testabilities.DefaultTopic,
		Response: &core.GASPInitialResponse{
			UTXOList: []*overlay.Outpoint{},
			Since:    0,
		},
	}

	expectedResponse := ports.NewRequestSyncResponseSuccessResponse(expectations.Response)
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
