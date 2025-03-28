package commands_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp"
	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestForeignGASPNodeHandler_ValidInput_ReturnsGASPNode(t *testing.T) {
	// Given:
	handler := commands.NewRequestForeignGASPNodeHandler(server.NewNoopEngineProvider())
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()
	payload := `{"graphID":"graph123", "txid":"tx789", "outputIndex":1}`

	// When:
	resp, err := http.Post(ts.URL, "application/json", bytes.NewBufferString(payload))

	// Then:
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var actual gasp.GASPNode
	expected := gasp.GASPNode{}
	err = json.NewDecoder(resp.Body).Decode(&actual)
	require.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestRequestForeignGASPNodeHandler_InvalidJSON_Returns400(t *testing.T) {
	// Given:
	handler := commands.NewRequestForeignGASPNodeHandler(server.NewNoopEngineProvider())
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	resp, err := http.Post(ts.URL, "application/json", bytes.NewBufferString(`INVALID_JSON`))

	// Then:
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRequestForeignGASPNodeHandler_MissingFields_StillReturnsOK(t *testing.T) {
	// Given:
	handler := commands.NewRequestForeignGASPNodeHandler(server.NewNoopEngineProvider())
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	resp, err := http.Post(ts.URL, "application/json", bytes.NewBufferString(`{}`))

	// Then:
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRequestForeignGASPNodeHandler_InvalidHTTPMethod_Returns405(t *testing.T) {
	// Given:
	handler := commands.NewRequestForeignGASPNodeHandler(server.NewNoopEngineProvider())
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	req, _ := http.NewRequest("GET", ts.URL, nil)
	resp, err := http.DefaultClient.Do(req)

	// Then:
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}
