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
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestRequestForeignGASPNodeHandler_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		payload            interface{}
		headers            map[string]string
		expectations       testabilities.RequestForeignGASPNodeProviderMockExpectations
		expectedStatusCode int
		expectedResponse   openapi.Error
	}{
		"Request foreign GASP node service fails to handle the request with missing topic header": {
			payload: map[string]interface{}{
				"graphID":     testabilities.DefaultValidGraphID,
				"txID":        testabilities.DefaultValidTxID,
				"outputIndex": testabilities.DefaultValidOutputIndex,
			},
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
			},
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: false,
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse: openapi.Error{
				Message: "The submitted request does not include required header: X-BSV-Topic.",
			},
		},
		"Request foreign GASP node service fails to handle the request with invalid JSON body": {
			payload: "INVALID_JSON",
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
				"X-BSV-Topic":           testabilities.DefaultValidTopic,
			},
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: false,
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse:   testabilities.NewTestOpenapiErrorResponse(t, app.NewRequestForeignGASPNodeInvalidJSONError()),
		},
		"Request foreign GASP node service fails to handle the request with empty topic": {
			payload: map[string]interface{}{
				"graphID":     testabilities.DefaultValidGraphID,
				"txID":        testabilities.DefaultValidTxID,
				"outputIndex": testabilities.DefaultValidOutputIndex,
			},
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
				"X-BSV-Topic":           testabilities.DefaultEmptyTopic,
			},
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: false,
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse: openapi.Error{
				Message: "One or more topics are in an invalid format. Empty string values are not allowed.",
			},
		},
		"Request foreign GASP node service fails to handle the request with invalid txID format": {
			payload: map[string]interface{}{
				"graphID":     testabilities.DefaultValidGraphID,
				"txID":        testabilities.DefaultInvalidTxID,
				"outputIndex": testabilities.DefaultValidOutputIndex,
			},
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
				"X-BSV-Topic":           testabilities.DefaultValidTopic,
			},
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: false,
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse:   testabilities.NewTestOpenapiErrorResponse(t, app.NewRequestForeignGASPNodeInvalidTxIDError()),
		},
		"Request foreign GASP node service fails to handle the request with invalid graphID format": {
			payload: map[string]interface{}{
				"graphID":     testabilities.DefaultInvalidGraphID,
				"txID":        testabilities.DefaultValidTxID,
				"outputIndex": testabilities.DefaultValidOutputIndex,
			},
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
				"X-BSV-Topic":           testabilities.DefaultValidTopic,
			},
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: false,
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedResponse:   testabilities.NewTestOpenapiErrorResponse(t, app.NewRequestForeignGASPNodeInvalidGraphIDError()),
		},
		"Request foreign GASP node service fails to handle the request with provider failure": {
			payload: map[string]interface{}{
				"graphID":     testabilities.DefaultValidGraphID,
				"txID":        testabilities.DefaultValidTxID,
				"outputIndex": testabilities.DefaultValidOutputIndex,
			},
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
				"X-BSV-Topic":           testabilities.DefaultValidTopic,
			},
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: true,
				Error:                      errors.New("internal request foreign GASP node provider error during handler unit test"),
			},
			expectedStatusCode: fiber.StatusInternalServerError,
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t,
				app.NewRequestForeignGASPNodeProviderError(
					errors.New("internal request foreign GASP node provider error during handler unit test"),
				),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestForeignGASPNodeProvider(
				testabilities.NewRequestForeignGASPNodeProviderMock(t, tc.expectations),
			))
			fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

			// when:
			var actualResponse openapi.BadRequestResponse
			res, _ := fixture.Client().
				R().
				SetHeaders(tc.headers).
				SetBody(tc.payload).
				SetError(&actualResponse).
				Post("/api/v1/requestForeignGASPNode")

			// then:
			require.Equal(t, tc.expectedStatusCode, res.StatusCode())
			require.Equal(t, &tc.expectedResponse, &actualResponse)
			stub.AssertProvidersState()
		})
	}
}

func TestRequestForeignGASPNodeHandler_ValidCase(t *testing.T) {
	// given:
	expectations := testabilities.RequestForeignGASPNodeProviderMockExpectations{
		ProvideForeignGASPNodeCall: true,
		Node:                       &core.GASPNode{},
	}

	expectedResponse := ports.NewRequestForeignGASPNodeSuccessResponse(expectations.Node)
	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestForeignGASPNodeProvider(
		testabilities.NewRequestForeignGASPNodeProviderMock(t, expectations),
	))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	headers := map[string]string{
		"X-BSV-Topic":           testabilities.DefaultValidTopic,
		fiber.HeaderContentType: fiber.MIMEApplicationJSON,
	}

	payload := map[string]interface{}{
		"graphID":     testabilities.DefaultValidGraphID,
		"txID":        testabilities.DefaultValidTxID,
		"outputIndex": testabilities.DefaultValidOutputIndex,
	}

	// when:
	var actualResponse openapi.GASPNode
	res, _ := fixture.Client().
		R().
		SetHeaders(headers).
		SetBody(payload).
		SetResult(actualResponse).
		Post("/api/v1/requestForeignGASPNode")

	// then:
	require.Equal(t, fiber.StatusOK, res.StatusCode())
	require.Equal(t, expectedResponse, actualResponse)
	stub.AssertProvidersState()
}
