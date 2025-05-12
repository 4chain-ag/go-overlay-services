package ports_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// TODO: check draft if all test cases are covered for http and fiber
const (
	testArcCallbackToken = "test-arc-callback-token"
	testArcApiKey        = "test-arc-api-key"
)

func Test_Debug_DirectArcIngestHandling(t *testing.T) {
	tests := []struct {
		name           string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Context Deadline Exceeded",
			mockError:      context.DeadlineExceeded,
			expectedStatus: fiber.StatusGatewayTimeout,
		},
		{
			name:           "Context Canceled",
			mockError:      context.Canceled,
			expectedStatus: fiber.StatusRequestTimeout,
		},
		{
			name:           "General Error",
			mockError:      errors.New("general error"),
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
		Status:  "success",
		Message: "Transaction status updated",
	}

	mockProvider := testabilities.NewMerkleProofProviderMockWithBlockHeight(t, nil, testabilities.DefaultBlockHeight)
	server := fiber.New()
	handler := ports.NewArcIngestHandler(mockProvider)
	server.Post("/api/v1/arc-ingest", handler.HandleArcIngest)

	// when:
	reqBody := fmt.Sprintf(`{"txid":"%s","merklePath":"%s","blockHeight":%d}`,
		testabilities.ValidTxId, testabilities.NewValidTestMerklePath(t), testabilities.DefaultBlockHeight)
	req, err := http.NewRequest("POST", "/api/v1/arc-ingest", strings.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testArcCallbackToken)

	resp, err := server.Test(req)
	require.NoError(t, err)

	// then:
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Read and verify response body
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var actualResponse ports.ArcIngestResponse
	err = json.Unmarshal(body, &actualResponse)
	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}

func Test_ArcIngestHandler_ValidationAndErrorTests(t *testing.T) {
	tests := map[string]struct {
		requestBody        map[string]interface{}
		mockError          error
		expectedStatusCode int
		expectedErrorMsg   string
	}{
		"should fail with 400 when the request body is missing txid": {
			requestBody: map[string]interface{}{
				"merklePath":  testabilities.NewValidTestMerklePath(t),
				"blockHeight": testabilities.DefaultBlockHeight,
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedErrorMsg:   "missing required field: txid",
		},
		"should fail with 400 when the request body is missing merklePath": {
			requestBody: map[string]interface{}{
				"txid":        testabilities.ValidTxId,
				"blockHeight": testabilities.DefaultBlockHeight,
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedErrorMsg:   "missing required field: merkle path",
		},
		"should fail with 400 when txid is invalid format": {
			requestBody: map[string]interface{}{
				"txid":        "not-a-hex-string",
				"merklePath":  testabilities.NewValidTestMerklePath(t),
				"blockHeight": testabilities.DefaultBlockHeight,
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedErrorMsg:   "Invalid transaction ID format",
		},
		"should fail with 400 when txid is invalid length": {
			requestBody: map[string]interface{}{
				"txid":        "1234abcd",
				"merklePath":  testabilities.NewValidTestMerklePath(t),
				"blockHeight": testabilities.DefaultBlockHeight,
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedErrorMsg:   "transaction ID does not match the expected length",
		},
		"should fail with 400 when merklePath is invalid format": {
			requestBody: map[string]interface{}{
				"txid":        testabilities.ValidTxId,
				"merklePath":  "not-a-hex-string",
				"blockHeight": testabilities.DefaultBlockHeight,
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedErrorMsg:   "Merkle path format is invalid",
		},
		"should fail with 504 when context deadline exceeded": {
			requestBody: map[string]interface{}{
				"txid":        testabilities.ValidTxId,
				"merklePath":  testabilities.NewValidTestMerklePath(t),
				"blockHeight": testabilities.DefaultBlockHeight,
			},
			mockError:          context.DeadlineExceeded,
			expectedStatusCode: fiber.StatusGatewayTimeout,
			expectedErrorMsg:   "timeout limit",
		},
		"should fail with 408 when context canceled": {
			requestBody: map[string]interface{}{
				"txid":        testabilities.ValidTxId,
				"merklePath":  testabilities.NewValidTestMerklePath(t),
				"blockHeight": testabilities.DefaultBlockHeight,
			},
			mockError:          context.Canceled,
			expectedStatusCode: fiber.StatusRequestTimeout,
			expectedErrorMsg:   "canceled",
		},
		"should fail with 500 when general error occurs": {
			requestBody: map[string]interface{}{
				"txid":        testabilities.ValidTxId,
				"merklePath":  testabilities.NewValidTestMerklePath(t),
				"blockHeight": testabilities.DefaultBlockHeight,
			},
			mockError:          errors.New("general error"),
			expectedStatusCode: fiber.StatusInternalServerError,
			expectedErrorMsg:   "Internal server error",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mockProvider := testabilities.NewMerkleProofProviderMockWithBlockHeight(t, tc.mockError, testabilities.DefaultBlockHeight)
			server := fiber.New()
			handler := ports.NewArcIngestHandler(mockProvider)
			server.Post("/api/v1/arc-ingest", handler.HandleArcIngest)

			// when:
			reqBodyJSON, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/arc-ingest", bytes.NewReader(reqBodyJSON))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testArcCallbackToken)

			resp, err := server.Test(req)
			require.NoError(t, err)

			// then:
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.expectedErrorMsg != "" {
				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)
				require.Contains(t, string(body), tc.expectedErrorMsg)
			}
		})
	}
}
