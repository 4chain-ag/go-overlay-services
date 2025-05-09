package ports_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
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
	}{
		"missing topic header": {
			payload: map[string]interface{}{
				"version": 1,
				"since":   1000,
			},
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
			},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		"invalid JSON": {
			payload: "INVALID_JSON",
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
				ports.XBSVTopicHeader:   "test-topic",
			},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		"provider error": {
			payload: map[string]interface{}{
				"version": 1,
				"since":   1000,
			},
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
				ports.XBSVTopicHeader:   "test-topic",
			},
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				ProvideForeignSyncResponseCall: true,
				Error:                          errors.New("provider error"),
			},
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

func TestRequestSyncResponseHandler_ValidCase(t *testing.T) {
	// given:
	expectedResponse := &core.GASPInitialResponse{
		UTXOList: []*overlay.Outpoint{},
		Since:    1000,
	}

	expectations := testabilities.RequestSyncResponseProviderMockExpectations{
		ProvideForeignSyncResponseCall: true,
		Response:                       expectedResponse,
	}

	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestSyncResponseProvider(
		testabilities.NewRequestSyncResponseProviderMock(t, expectations),
	))

	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	// when:
	var actualResponse core.GASPInitialResponse

	res, _ := fixture.Client().
		R().
		SetHeader(fiber.HeaderContentType, fiber.MIMEApplicationJSON).
		SetHeader(ports.XBSVTopicHeader, "test-topic").
		SetBody(map[string]interface{}{
			"version": 1,
			"since":   1000,
		}).
		SetResult(&actualResponse).
		Post("/api/v1/requestSyncResponse")

	// then:
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.Equal(t, *expectedResponse, actualResponse)
	stub.AssertProvidersState()
}
