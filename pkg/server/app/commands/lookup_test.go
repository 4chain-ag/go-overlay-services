package commands_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// alwaysSucceedsLookup implements the LookupQuestionProvider interface for successful test cases
type alwaysSucceedsLookup struct{}

func (s *alwaysSucceedsLookup) Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
	return &lookup.LookupAnswer{
		Type: lookup.AnswerTypeFreeform,
		Result: map[string]interface{}{
			"data": "test data",
		},
	}, nil
}

// alwaysFailsLookup implements the LookupQuestionProvider interface for failure test cases
type alwaysFailsLookup struct{}

func (s *alwaysFailsLookup) Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
	return nil, errors.New("lookup failed")
}

func TestLookupHandler_ValidInput_ReturnsAnswer(t *testing.T) {
	// Given:
	handler, err := commands.NewLookupHandler(&alwaysSucceedsLookup{})
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

func TestLookupHandler_InvalidJSON_Returns400(t *testing.T) {
	// Given:
	handler, err := commands.NewLookupHandler(&alwaysSucceedsLookup{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	resp, err := http.Post(ts.URL, "application/json", bytes.NewBufferString(`INVALID_JSON`))

	// Then:
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLookupHandler_MissingFields_Returns400(t *testing.T) {
	// Given:
	handler, err := commands.NewLookupHandler(&alwaysSucceedsLookup{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	resp, err := ts.Client().Post(ts.URL, "application/json", bytes.NewBufferString(`{}`))

	// Then:
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLookupHandler_InvalidHTTPMethod_Returns405(t *testing.T) {
	// Given:
	handler, err := commands.NewLookupHandler(&alwaysSucceedsLookup{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	resp, err := ts.Client().Do(req)

	// Then:
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestLookupHandler_EngineError_Returns400(t *testing.T) {
	// Given:
	handler, err := commands.NewLookupHandler(&alwaysFailsLookup{})
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
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
} 
