package commands_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/commands"
	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/commands/testutil"
	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/stretchr/testify/require"
)

// Test_ArcIngestHandler_ShouldRespondsWith200AndCallsProvider tests the successful handling of a valid request
func Test_ArcIngestHandler_ShouldRespondsWith200AndCallsProvider(t *testing.T) {
	// given:
	payload := commands.ArcIngestRequest{
		TxID:        testutil.ValidTxId,
		MerklePath:  testutil.NewValidTestMerklePath(t),
		BlockHeight: 848372,
	}

	mock := testutil.NewMerkleProofProviderMock(nil, payload.BlockHeight)
	handler, err := commands.NewArcIngestHandler(mock)

	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL, testutil.RequestBody(t, payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// when:
	res, err := ts.Client().Do(req)

	// then:
	require.NoError(t, err)
	defer res.Body.Close()

	require.NotNil(t, res)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var actualResponse commands.ArcIngestHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &actualResponse))

	expectedResponse := commands.NewSuccessArcIngestHandlerResponse()
	require.Equal(t, expectedResponse, actualResponse)

	mock.AssertCalled(t)
}

// Test_ArcIngestHandler_ValidationTests tests error handling for various invalid inputs
func Test_ArcIngestHandler_ValidationTests(t *testing.T) {
	tests := map[string]struct {
		method         string
		payload        commands.ArcIngestRequest
		setupRequest   func(*http.Request)
		expectedStatus int
		expectedError  error
	}{
		"should fail with 405 when HTTP method is GET": {
			method: http.MethodGet,
			payload: commands.ArcIngestRequest{
				TxID:        testutil.ValidTxId,
				MerklePath:  testutil.NewValidTestMerklePath(t),
				BlockHeight: 848372,
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  commands.ErrInvalidHTTPMethod,
		},
		"should fail with 405 when HTTP method is PUT": {
			method: http.MethodPut,
			payload: commands.ArcIngestRequest{
				TxID:        testutil.ValidTxId,
				MerklePath:  testutil.NewValidTestMerklePath(t),
				BlockHeight: 848372,
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  commands.ErrInvalidHTTPMethod,
		},
		"should fail with 400 when all required fields are missing": {
			method: http.MethodPost,
			payload: commands.ArcIngestRequest{
				TxID:        "",
				MerklePath:  "",
				BlockHeight: 0,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  commands.ErrMissingRequiredRequestFieldsDefinition,
		},
		"should fail with 400 when TxID field is missing": {
			method: http.MethodPost,
			payload: commands.ArcIngestRequest{
				TxID:        "",
				MerklePath:  testutil.NewValidTestMerklePath(t),
				BlockHeight: 848372,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  commands.ErrMissingRequiredTxIDFieldDefinition,
		},
		"should fail with 400 when MerklePath field is missing": {
			method: http.MethodPost,
			payload: commands.ArcIngestRequest{
				TxID:        testutil.ValidTxId,
				MerklePath:  "",
				BlockHeight: 848372,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  commands.ErrMissingRequiredMerklePathFieldDefinition,
		},
		"should fail with 400 when TxID format is invalid": {
			method: http.MethodPost,
			payload: commands.ArcIngestRequest{
				TxID:        "invalid-hex-string",
				MerklePath:  testutil.NewValidTestMerklePath(t),
				BlockHeight: 848372,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  commands.ErrInvalidTxIDFormat,
		},
		"should fail with 400 when TxID length is invalid": {
			method: http.MethodPost,
			payload: commands.ArcIngestRequest{
				TxID:        "1234",
				MerklePath:  testutil.NewValidTestMerklePath(t),
				BlockHeight: 848372,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  commands.ErrInvalidTxIDLength,
		},
		"should fail with 400 when MerklePath format is invalid": {
			method: http.MethodPost,
			payload: commands.ArcIngestRequest{
				TxID:        testutil.ValidTxId,
				MerklePath:  "invalid-merkle-path",
				BlockHeight: 848372,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  commands.ErrInvalidMerklePathFormat,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testutil.NewMerkleProofProviderMock(nil, tc.payload.BlockHeight)
			handler, err := commands.NewArcIngestHandler(mock)
			require.NoError(t, err)

			ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
			defer ts.Close()

			var requestBody io.Reader
			if tc.setupRequest == nil {
				data, err := json.Marshal(tc.payload)
				require.NoError(t, err)
				requestBody = bytes.NewReader(data)
			} else {
				requestBody = bytes.NewReader([]byte{})
			}

			req, err := http.NewRequest(tc.method, ts.URL, requestBody)
			require.NoError(t, err)

			if tc.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/json")
			}

			if tc.setupRequest != nil {
				tc.setupRequest(req)
			}

			// when:
			res, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			// then:
			require.Equal(t, tc.expectedStatus, res.StatusCode)

			var response commands.ArcIngestHandlerResponse
			err = jsonutil.DecodeResponseBody(res, &response)
			require.NoError(t, err)

			require.Contains(t, response.Message, tc.expectedError.Error())
			require.Equal(t, "error", response.Status)
		})
	}
}
