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

func TestSubmitTransactionHandler_Handle_SuccessfulSubmission(t *testing.T) {
	// Given:
	handler := commands.NewSubmitTransactionCommandHandler(&SubmitTransactionProviderAlwaysSuccess{})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	requestBody := []byte("test transaction body")
	topics := []string{"topic1", "topic2"}
	topicsJSON, _ := json.Marshal(topics)

	req, _ := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-topics", string(topicsJSON))

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
	handler := commands.NewSubmitTransactionCommandHandler(&SubmitTransactionProviderAlwaysSuccess{})
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
	assert.Equal(t, "Method not allowed\n", string(body))
}

func TestSubmitTransactionHandler_Handle_MissingTopicsHeader(t *testing.T) {
	// Given:
	handler := commands.NewSubmitTransactionCommandHandler(&SubmitTransactionProviderAlwaysSuccess{})
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
	assert.Equal(t, "Missing x-topics header\n", string(body))
}

func TestSubmitTransactionHandler_Handle_InvalidTopicsFormat(t *testing.T) {
	// Given:
	handler := commands.NewSubmitTransactionCommandHandler(&SubmitTransactionProviderAlwaysSuccess{})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	req, _ := http.NewRequest("POST", ts.URL, strings.NewReader("test body"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-topics", "not-valid-json")

	// When:
	res, err := ts.Client().Do(req)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "Invalid x-topics header format\n", string(body))
}

func TestSubmitTransactionHandler_Handle_ProviderError(t *testing.T) {
	// Given:
	handler := commands.NewSubmitTransactionCommandHandler(&SubmitTransactionProviderAlwaysFailure{})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	requestBody := []byte("test transaction body")
	topics := []string{"topic1", "topic2"}
	topicsJSON, _ := json.Marshal(topics)

	req, _ := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-topics", string(topicsJSON))

	// When:
	res, err := ts.Client().Do(req)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "Submit transaction test error\n", string(body))
}

func TestNewSubmitTransactionCommandHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider commands.SubmitTransactionProvider = nil

	// When & Then:
	assert.Panics(t, func() {
		commands.NewSubmitTransactionCommandHandler(provider)
	}, "Expected panic when provider is nil")
}
