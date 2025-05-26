package ports_test

import (
	"errors"
	"net/http"
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
		expectedStatusCode int
		headers            map[string]string
		payload            interface{}
		expectedResponse   openapi.Error
		expectations       testabilities.RequestForeignGASPNodeProviderMockExpectations
	}{
		"Request foreign GASP node service fails to handle the request - missing topic header": {
			expectedStatusCode: fiber.StatusBadRequest,
			payload: map[string]interface{}{
				"graphID":     testabilities.DefaultValidGraphID,
				"txID":        testabilities.DefaultValidTxID,
				"outputIndex": testabilities.DefaultValidOutputIndex,
			},
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
			},
			expectedResponse: openapi.Error{
				Message: "The submitted request does not include required header: X-BSV-Topic.",
			},
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: false,
			},
		},
		"Request foreign GASP node service fails to handle the request - invalid JSON body": {
			expectedStatusCode: fiber.StatusBadRequest,
			payload:            "INVALID_JSON",
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
				"X-BSV-Topic":           testabilities.DefaultValidTopic,
			},
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, ports.NewInvalidRequestBodyError()),
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: false,
			},
		},
		"Request foreign GASP node service fails to handle the request - missing topic": {
			expectedStatusCode: fiber.StatusBadRequest,
			payload: map[string]interface{}{
				"graphID":     testabilities.DefaultValidGraphID,
				"txID":        testabilities.DefaultValidTxID,
				"outputIndex": testabilities.DefaultValidOutputIndex,
			},
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
				"X-BSV-Topic":           testabilities.DefaultEmptyTopic,
			},
			expectedResponse: openapi.Error{
				Message: "One or more topics are in an invalid format. Empty string values are not allowed.",
			},
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: false,
			},
		},
		"Request foreign GASP node service fails to handle the request - invalid txID format": {
			expectedStatusCode: fiber.StatusBadRequest,
			payload: map[string]interface{}{
				"graphID":     testabilities.DefaultValidGraphID,
				"txID":        testabilities.DefaultInvalidTxID,
				"outputIndex": testabilities.DefaultValidOutputIndex,
			},
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
				"X-BSV-Topic":           testabilities.DefaultValidTopic,
			},
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, app.NewRequestForeignGASPNodeInvalidTxIDError()),
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: false,
			},
		},
		"Request foreign GASP node service fails to handle the request - invalid graphID format": {
			expectedStatusCode: fiber.StatusBadRequest,
			payload: map[string]interface{}{
				"graphID":     testabilities.DefaultInvalidGraphID,
				"txID":        testabilities.DefaultValidTxID,
				"outputIndex": testabilities.DefaultValidOutputIndex,
			},
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
				"X-BSV-Topic":           testabilities.DefaultValidTopic,
			},
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t, app.NewRequestForeignGASPNodeInvalidGraphIDError()),
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: false,
			},
		},
		"Request foreign GASP node service fails to handle the request - provider failure": {
			expectedStatusCode: fiber.StatusInternalServerError,
			payload: map[string]interface{}{
				"graphID":     testabilities.DefaultValidGraphID,
				"txID":        testabilities.DefaultValidTxID,
				"outputIndex": testabilities.DefaultValidOutputIndex,
			},
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEApplicationJSON,
				"X-BSV-Topic":           testabilities.DefaultValidTopic,
			},
			expectedResponse: testabilities.NewTestOpenapiErrorResponse(t,
				app.NewRequestForeignGASPNodeProviderError(
					errors.New("internal request foreign GASP node provider error during handler unit test"),
				),
			),
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: true,
				Error:                      errors.New("internal request foreign GASP node provider error during handler unit test"),
			},
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

			if tc.expectations.ProvideForeignGASPNodeCall {
				stub.AssertProvidersState()
			}
		})
	}
}

func TestRequestForeignGASPNodeHandler_ValidCase(t *testing.T) {
	// given:
	expectedNode := &core.GASPNode{}
	expectations := testabilities.RequestForeignGASPNodeProviderMockExpectations{
		ProvideForeignGASPNodeCall: true,
		Node:                       expectedNode,
	}

	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestForeignGASPNodeProvider(
		testabilities.NewRequestForeignGASPNodeProviderMock(t, expectations),
	))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	headers := map[string]string{
		fiber.HeaderContentType: fiber.MIMEApplicationJSON,
		"X-BSV-Topic":           testabilities.DefaultValidTopic,
	}

	payload := map[string]interface{}{
		"graphID":     testabilities.DefaultValidGraphID,
		"txID":        testabilities.DefaultValidTxID,
		"outputIndex": testabilities.DefaultValidOutputIndex,
	}

	// when:
	var actualNode core.GASPNode
	res, _ := fixture.Client().
		R().
		SetHeaders(headers).
		SetBody(payload).
		SetResult(&actualNode).
		Post("/api/v1/requestForeignGASPNode")

	// then:
	expectedResponse := ports.NewRequestForeignGASPNodeSuccessResponse(expectedNode)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.Equal(t, expectedResponse, &actualNode)
	stub.AssertProvidersState()
}
