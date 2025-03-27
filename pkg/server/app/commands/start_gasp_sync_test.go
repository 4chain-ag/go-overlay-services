package commands_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// mockErrProvider forces StartGASPSync to fail
type mockErrProvider struct{}

func (m *mockErrProvider) StartGASPSync() error {
	return errors.New("sync error")
}

func TestStartGASPSyncHandler_Success(t *testing.T) {
	// Given
	app := fiber.New()
	handler := commands.NewStartGASPSyncHandler(server.NewNoopEngineProvider())
	app.Post("/", handler.Handle)

	// When
	req, _ := http.NewRequest("POST", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Then
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestStartGASPSyncHandler_Failure(t *testing.T) {
	// Given
	app := fiber.New()
	handler := commands.NewStartGASPSyncHandler(&mockErrProvider{})
	app.Post("/", handler.Handle)

	// When
	req, _ := http.NewRequest("POST", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	// Then
	require.NoError(t, err)
	require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
