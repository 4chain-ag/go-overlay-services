package ports_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

const (
	testArcCallbackToken = "test-arc-callback-token"

	testArcApiKey = "test-arc-api-key"
)

func Test_Debug_DirectArcIngestHandling(t *testing.T) {

	tests := []struct {
		name string

		mockError error

		expectedStatus int
	}{

		{

			name: "Context Deadline Exceeded",

			mockError: context.DeadlineExceeded,

			expectedStatus: fiber.StatusGatewayTimeout,
		},

		{

			name: "Context Canceled",

			mockError: context.Canceled,

			expectedStatus: fiber.StatusRequestTimeout,
		},

		{

			name: "General Error",

			mockError: errors.New("general error"),

			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			mockProvider := testabilities.NewMerkleProofProviderMockWithBlockHeight(t, tc.mockError, testabilities.DefaultBlockHeight)

			handler := ports.NewArcIngestHandler(mockProvider)

			app := fiber.New()

			app.Post("/test", handler.HandleArcIngest)

			reqBody := fmt.Sprintf(`{"txid":"%s","merklePath":"%s","blockHeight":%d}`,

				testabilities.ValidTxId, testabilities.NewValidTestMerklePath(t), testabilities.DefaultBlockHeight)

			req, err := http.NewRequest("POST", "/test", strings.NewReader(reqBody))

			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, resp.StatusCode, "Expected status %d but got %d", tc.expectedStatus, resp.StatusCode)

		})

	}

}

func Test_ArcIngestHandler_ShouldRespondWith200AndCallsProvider(t *testing.T) {

	// given:

	expectedResponse := ports.ArcIngestResponse{

		Status: "success",

		Message: "Transaction status updated",
	}

	mockEngineStub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithArcIngestProvider(

		testabilities.NewMerkleProofProviderMockWithBlockHeight(t, nil, testabilities.DefaultBlockHeight),
	))

	fixture := server2.NewServerTestFixture(t,

		server2.WithEngine(mockEngineStub),

		server2.WithArcConfiguration(testArcApiKey, testArcCallbackToken),
	)

	requestBody := map[string]interface{}{

		"txid": testabilities.ValidTxId,

		"merklePath": testabilities.NewValidTestMerklePath(t),

		"blockHeight": testabilities.DefaultBlockHeight,
	}

	// when:

	var actualResponse ports.ArcIngestResponse

	res, _ := fixture.Client().
		R().
		SetBody(requestBody).
		SetResult(&actualResponse).
		SetHeader("Authorization", "Bearer "+testArcCallbackToken).
		Post("/api/v1/arc-ingest")

	// then:

	require.Equal(t, http.StatusOK, res.StatusCode())

	require.Equal(t, expectedResponse, actualResponse)

	mockEngineStub.AssertProvidersState()

}

func Test_ArcIngestHandler_ValidationAndErrorTests(t *testing.T) {

	tests := map[string]struct {
		requestBody map[string]interface{}

		mockError error

		expectedStatusCode int

		expectedErrorMsg string
	}{

		"should fail with 400 when the request body is missing txid": {

			requestBody: map[string]interface{}{

				"merklePath": testabilities.NewValidTestMerklePath(t),

				"blockHeight": testabilities.DefaultBlockHeight,
			},

			expectedStatusCode: fiber.StatusBadRequest,

			expectedErrorMsg: "missing required field: txid",
		},

		"should fail with 400 when the request body is missing merklePath": {

			requestBody: map[string]interface{}{

				"txid": testabilities.ValidTxId,

				"blockHeight": testabilities.DefaultBlockHeight,
			},

			expectedStatusCode: fiber.StatusBadRequest,

			expectedErrorMsg: "missing required field: merkle path",
		},

		"should fail with 400 when txid is invalid format": {

			requestBody: map[string]interface{}{

				"txid": "not-a-hex-string",

				"merklePath": testabilities.NewValidTestMerklePath(t),

				"blockHeight": testabilities.DefaultBlockHeight,
			},

			expectedStatusCode: fiber.StatusBadRequest,

			expectedErrorMsg: app.ErrInvalidTxIDFormat.Error(),
		},

		"should fail with 400 when txid is invalid length": {

			requestBody: map[string]interface{}{

				"txid": "1234abcd",

				"merklePath": testabilities.NewValidTestMerklePath(t),

				"blockHeight": testabilities.DefaultBlockHeight,
			},

			expectedStatusCode: fiber.StatusBadRequest,

			expectedErrorMsg: app.ErrInvalidTxIDLength.Error(),
		},

		"should fail with 400 when merklePath is invalid format": {

			requestBody: map[string]interface{}{

				"txid": testabilities.ValidTxId,

				"merklePath": "not-a-hex-string",

				"blockHeight": testabilities.DefaultBlockHeight,
			},

			expectedStatusCode: fiber.StatusBadRequest,

			expectedErrorMsg: app.ErrInvalidMerklePathFormat.Error(),
		},

		"should fail with 504 when context deadline exceeded": {

			requestBody: map[string]interface{}{

				"txid": testabilities.ValidTxId,

				"merklePath": testabilities.NewValidTestMerklePath(t),

				"blockHeight": testabilities.DefaultBlockHeight,
			},

			mockError: context.DeadlineExceeded,

			expectedStatusCode: fiber.StatusGatewayTimeout,

			expectedErrorMsg: app.ErrMerkleProofProcessingTimeout.Error(),
		},

		"should fail with 408 when context canceled": {

			requestBody: map[string]interface{}{

				"txid": testabilities.ValidTxId,

				"merklePath": testabilities.NewValidTestMerklePath(t),

				"blockHeight": testabilities.DefaultBlockHeight,
			},

			mockError: context.Canceled,

			expectedStatusCode: fiber.StatusRequestTimeout,

			expectedErrorMsg: app.ErrMerkleProofProcessingCanceled.Error(),
		},

		"should fail with 500 when general error occurs": {

			requestBody: map[string]interface{}{

				"txid": testabilities.ValidTxId,

				"merklePath": testabilities.NewValidTestMerklePath(t),

				"blockHeight": testabilities.DefaultBlockHeight,
			},

			mockError: errors.New("general error"),

			expectedStatusCode: fiber.StatusInternalServerError,

			expectedErrorMsg: app.ErrMerkleProofProcessingFailed.Error(),
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			// given:

			mockEngineStub := testabilities.NewTestOverlayEngineStub(t,

				testabilities.WithArcIngestProvider(

					testabilities.NewMerkleProofProviderMockWithBlockHeight(t, tc.mockError, testabilities.DefaultBlockHeight),
				),
			)

			fixture := server2.NewServerTestFixture(t,

				server2.WithEngine(mockEngineStub),

				server2.WithArcConfiguration(testArcApiKey, testArcCallbackToken),
			)

			// when:

			var errorResponse openapi.Error

			res, _ := fixture.Client().
				R().
				SetBody(tc.requestBody).
				SetError(&errorResponse).
				SetHeader("Authorization", "Bearer "+testArcCallbackToken).
				Post("/api/v1/arc-ingest")

			// then:

			require.Equal(t, tc.expectedStatusCode, res.StatusCode())

			if tc.expectedErrorMsg != "" {

				require.Contains(t, errorResponse.Message, tc.expectedErrorMsg)

			}

			mockEngineStub.AssertProvidersState()

		})

	}

}
