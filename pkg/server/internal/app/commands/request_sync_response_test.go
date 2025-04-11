package commands_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/commands"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

// Mock provider that always succeeds.
type foreignSyncProviderAlwaysSuccess struct{}

func (foreignSyncProviderAlwaysSuccess) ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {
	return &core.GASPInitialResponse{
		UTXOList: []*overlay.Outpoint{},
		Since:    initialRequest.Since,
	}, nil
}

// Mock provider that always fails.
type foreignSyncProviderAlwaysFailure struct{}

func (foreignSyncProviderAlwaysFailure) ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {
	return nil, fmt.Errorf("simulated sync failure")
}

func TestRequestSyncResponseHandler_Success(t *testing.T) {
	// Given:
	handler, err := commands.NewRequestSyncResponseHandler(&foreignSyncProviderAlwaysSuccess{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	payload := core.GASPInitialRequest{
		Version: 1,
		Since:   1000,
	}
	body, err := json.Marshal(payload)
	require.NotEmpty(t, body)
	require.NoError(t, err)

	// When:
	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-bsv-topic", "example-topic")
	
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Then:
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRequestSyncResponseHandler_MissingTopic(t *testing.T) {
	// Given:
	handler, err := commands.NewRequestSyncResponseHandler(&foreignSyncProviderAlwaysSuccess{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	payload := core.GASPInitialRequest{
		Version: 1,
		Since:   1000,
	}
	body, err := json.Marshal(payload)
	require.NotEmpty(t, body)
	require.NoError(t, err)

	// When:
	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Then:
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(respBody), "missing 'x-bsv-topic'")
}

func TestRequestSyncResponseHandler_InvalidJSON(t *testing.T) {
	// Given:
	handler, err := commands.NewRequestSyncResponseHandler(&foreignSyncProviderAlwaysSuccess{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewReader([]byte(`{invalid-json}`)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-bsv-topic", "example-topic")
	
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Then:
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRequestSyncResponseHandler_InternalServerError(t *testing.T) {
	// Given:
	handler, err := commands.NewRequestSyncResponseHandler(&foreignSyncProviderAlwaysFailure{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	payload := core.GASPInitialRequest{
		Version: 1,
		Since:   1000,
	}
	body, err := json.Marshal(payload)
	require.NotEmpty(t, body)
	require.NoError(t, err)

	// When:
	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-bsv-topic", "example-topic")
	
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Then:
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
