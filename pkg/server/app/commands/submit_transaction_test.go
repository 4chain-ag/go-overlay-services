package commands_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SubmitTransactionProviderAlwaysSuccess is an implementation that always succeeds
type SubmitTransactionProviderAlwaysSuccess struct{}

func (SubmitTransactionProviderAlwaysSuccess) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error) {
	// Call the onSteakReady callback to simulate async completion
	if onSteakReady != nil {
		steak := overlay.Steak{}
		onSteakReady(&steak)
	}
	return overlay.Steak{}, nil
}

// SubmitTransactionProviderAlwaysFailure is an implementation that always returns an error
type SubmitTransactionProviderAlwaysFailure struct{}

func (SubmitTransactionProviderAlwaysFailure) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error) {
	return overlay.Steak{}, errors.New("Submit transaction test error")
}

// SubmitTransactionProviderNeverCallback is an implementation that never calls the callback
type SubmitTransactionProviderNeverCallback struct{}

func (SubmitTransactionProviderNeverCallback) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error) {
	// Never call the callback which then should trigger the timeout
	return overlay.Steak{}, nil
}

// For testing purposes only - allows creating a handler with a custom body limit using middleware type approach
func createTestHandlerWithLimit(provider commands.SubmitTransactionProvider, limit int64) http.HandlerFunc {
	handler, err := commands.NewSubmitTransactionCommandHandler(provider)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > limit {
			http.Error(w, commands.ErrRequestBodyTooLarge.Error(), http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(io.LimitReader(r.Body, limit))
		handler.Handle(w, r)
	}
}

func TestSubmitTransactionHandler_Handle_SuccessfulSubmission(t *testing.T) {
	// Given:
	handler, err := commands.NewSubmitTransactionCommandHandler(&SubmitTransactionProviderAlwaysSuccess{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	requestBody := []byte("test transaction body")
	topics := []string{"topic1", "topic2"}
	topicsJSON, _ := json.Marshal(topics)

	req, _ := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(commands.XTopicsHeader, string(topicsJSON))

	// When:
	res, err := ts.Client().Do(req)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	var actual commands.SubmitTransactionHandlerResponse
	expected := commands.SubmitTransactionHandlerResponse{Steak: overlay.Steak{}}

	require.NoError(t, jsonutil.DecodeResponseBody(res, &actual))
	assert.Equal(t, expected, actual)
}

func TestSubmitTransactionHandler_Handle_InvalidMethod(t *testing.T) {
	// Given:
	handler, err := commands.NewSubmitTransactionCommandHandler(&SubmitTransactionProviderAlwaysSuccess{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	// When:
	res, err := ts.Client().Do(req)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, commands.ErrInvalidHTTPMethod.Error()+"\n", string(body))
}

func TestSubmitTransactionHandler_Handle_MissingTopicsHeader(t *testing.T) {
	// Given:
	handler, err := commands.NewSubmitTransactionCommandHandler(&SubmitTransactionProviderAlwaysSuccess{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	req, _ := http.NewRequest("POST", ts.URL, strings.NewReader("test body"))
	req.Header.Set("Content-Type", "application/json")

	// When:
	res, err := ts.Client().Do(req)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, commands.ErrMissingXTopicsHeader.Error()+"\n", string(body))
}

func TestSubmitTransactionHandler_Handle_InvalidTopicsFormat(t *testing.T) {
	// Given:
	handler, err := commands.NewSubmitTransactionCommandHandler(&SubmitTransactionProviderAlwaysSuccess{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	req, _ := http.NewRequest("POST", ts.URL, strings.NewReader("test body"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(commands.XTopicsHeader, "not-valid-json")

	// When:
	res, err := ts.Client().Do(req)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, commands.ErrInvalidXTopicsHeaderFormat.Error()+"\n", string(body))
}

func TestSubmitTransactionHandler_Handle_ProviderError(t *testing.T) {
	// Given:
	handler, err := commands.NewSubmitTransactionCommandHandler(&SubmitTransactionProviderAlwaysFailure{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	requestBody := []byte("test transaction body")
	topics := []string{"topic1", "topic2"}
	topicsJSON, _ := json.Marshal(topics)

	req, _ := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(commands.XTopicsHeader, string(topicsJSON))

	// When:
	res, err := ts.Client().Do(req)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestSubmitTransactionHandler_Handle_RequestTooLarge(t *testing.T) {
	// Given a handler with a small request body limit (10 bytes)
	testHandler := createTestHandlerWithLimit(&SubmitTransactionProviderAlwaysSuccess{}, 10)
	ts := httptest.NewServer(testHandler)
	defer ts.Close()

	requestBody := bytes.NewBufferString("this is more than 10 bytes of data")
	topics := []string{"topic1"}
	topicsJSON, _ := json.Marshal(topics)

	req, _ := http.NewRequest("POST", ts.URL, requestBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(commands.XTopicsHeader, string(topicsJSON))
	req.ContentLength = int64(requestBody.Len())

	// When:
	res, err := ts.Client().Do(req)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, commands.ErrRequestBodyTooLarge.Error()+"\n", string(body))
}

func TestSubmitTransactionHandler_Handle_Timeout(t *testing.T) {
	// Given:
	handler, err := commands.NewSubmitTransactionCommandHandler(&SubmitTransactionProviderNeverCallback{})
	require.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	requestBody := []byte("test transaction body")
	topics := []string{"topic1", "topic2"}
	topicsJSON, _ := json.Marshal(topics)

	req, _ := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(commands.XTopicsHeader, string(topicsJSON))

	// When:
	res, err := ts.Client().Do(req)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusRequestTimeout, res.StatusCode)
}

func TestNewSubmitTransactionCommandHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider commands.SubmitTransactionProvider = nil

	// When:
	handler, err := commands.NewSubmitTransactionCommandHandler(provider)

	// Then:
	assert.Nil(t, handler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "submit transaction provider is nil")
}
