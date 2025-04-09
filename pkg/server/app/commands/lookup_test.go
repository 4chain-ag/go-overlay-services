package commands_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands/testutil"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookupHandler_ValidInput_ReturnsAnswer(t *testing.T) {
	// Given:
	handler, err := commands.NewLookupHandler(&testutil.AlwaysSucceedsLookup{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	payload := lookup.LookupQuestion{
		Service: "test-service",
		Query:   json.RawMessage(`{"test":"query"}`),
	}
	jsonData, err := json.Marshal(payload)
	require.NoError(t, err)

	// When:
	resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(jsonData))

	// Then:
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	type lookupResponse struct {
		Answer *lookup.LookupAnswer `json:"answer"`
	}
	var response lookupResponse
	require.NoError(t, jsonutil.DecodeResponseBody(resp, &response))
	assert.Equal(t, lookup.AnswerTypeFreeform, response.Answer.Type)
	assert.Contains(t, response.Answer.Result, "data")

	resultMap, ok := response.Answer.Result.(map[string]interface{})
	require.True(t, ok, "Result should be a map[string]interface{}")
	assert.Equal(t, "test data", resultMap["data"])
}

func TestLookupHandler_ErrorCases(t *testing.T) {
	tests := []struct {
		name             string
		mockProvider     commands.LookupQuestionProvider
		requestMethod    string
		requestBody      string
		expectedStatus   int
		validateResponse func(t *testing.T, body []byte)
	}{
		{
			name:           "Invalid JSON",
			mockProvider:   &testutil.AlwaysSucceedsLookup{},
			requestMethod:  http.MethodPost,
			requestBody:    `INVALID_JSON`,
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "invalid")
			},
		},
		{
			name:           "Missing Fields",
			mockProvider:   &testutil.AlwaysSucceedsLookup{},
			requestMethod:  http.MethodPost,
			requestBody:    `{}`,
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "missing")
			},
		},
		{
			name:           "Invalid HTTP Method",
			mockProvider:   &testutil.AlwaysSucceedsLookup{},
			requestMethod:  http.MethodGet,
			requestBody:    `{}`,
			expectedStatus: http.StatusMethodNotAllowed,
			validateResponse: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "method")
			},
		},
		{
			name:           "Engine Error",
			mockProvider:   &testutil.AlwaysFailsLookup{},
			requestMethod:  http.MethodPost,
			requestBody:    `{"service":"test-service","query":{"test":"query"}}`,
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "lookup failed")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given:
			handler, err := commands.NewLookupHandler(tc.mockProvider)
			require.NoError(t, err)
			ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
			defer ts.Close()

			// When:
			var resp *http.Response
			if tc.requestMethod == http.MethodPost {
				resp, err = http.Post(ts.URL, "application/json", bytes.NewBufferString(tc.requestBody))
			} else {
				req, _ := http.NewRequest(tc.requestMethod, ts.URL, nil)
				resp, err = ts.Client().Do(req)
			}

			// Then:
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, tc.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if tc.validateResponse != nil {
				tc.validateResponse(t, body)
			}
		})
	}
}
