package commands_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/stretchr/testify/require"
)

// mockErrProvider forces StartGASPSync to fail
type mockErrProvider struct{}

func (m *mockErrProvider) StartGASPSync() error {
	return errors.New("sync error")
}

func TestStartGASPSyncHandler_Success(t *testing.T) {
	// Given:
	handler := commands.NewStartGASPSyncHandler(server.NewNoopEngineProvider())
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	req, err := http.NewRequest("POST", ts.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	// Then:
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestStartGASPSyncHandler_Failure(t *testing.T) {
	// Given:
	handler := commands.NewStartGASPSyncHandler(&mockErrProvider{})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	req, err := http.NewRequest("POST", ts.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	// Then:
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
