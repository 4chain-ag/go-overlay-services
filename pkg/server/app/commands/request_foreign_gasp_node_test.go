package commands_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp"
	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestRequestForeignGASPNodeHandler_ValidInput_ReturnsGASPNode(t *testing.T) {
	// Given
	app := fiber.New()
	handler := commands.NewRequestForeignGASPNodeHandler(server.NewNoopEngineProvider())
	app.Post("/", handler.Handle)

	payload := `{"graphID":"graph123", "txid":"tx789", "outputIndex":1}`

	// When
	req := httptestRequest("POST", "/", payload)
	resp, err := app.Test(req, -1)

	// Then
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var node gasp.GASPNode
	err = json.NewDecoder(resp.Body).Decode(&node)
	require.NoError(t, err)
	require.Equal(t, "", node.GraphID)
	require.Equal(t, "", node.RawTx)
	require.Equal(t, uint32(0), node.OutputIndex)
}

func TestRequestForeignGASPNodeHandler_InvalidJSON_Returns400(t *testing.T) {
	// Given
	app := fiber.New()
	handler := commands.NewRequestForeignGASPNodeHandler(server.NewNoopEngineProvider())
	app.Post("/", handler.Handle)

	// When
	req := httptestRequest("POST", "/", `INVALID_JSON`)
	resp, err := app.Test(req, -1)

	// Then
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestRequestForeignGASPNodeHandler_MissingFields_StillReturnsOK(t *testing.T) {
	// Given
	app := fiber.New()
	handler := commands.NewRequestForeignGASPNodeHandler(server.NewNoopEngineProvider())
	app.Post("/", handler.Handle)

	// When
	req := httptestRequest("POST", "/", `{}`)
	resp, err := app.Test(req, -1)

	// Then
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRequestForeignGASPNodeHandler_InvalidHTTPMethod_Returns405(t *testing.T) {
	// Given
	app := fiber.New()
	handler := commands.NewRequestForeignGASPNodeHandler(server.NewNoopEngineProvider())
	app.Post("/", handler.Handle)

	// When
	req := httptestRequest("GET", "/", ``)
	resp, err := app.Test(req, -1)

	// Then
	require.NoError(t, err)
	require.Equal(t, fiber.StatusMethodNotAllowed, resp.StatusCode) // 405
}

// helper
func httptestRequest(method, url, body string) *http.Request {
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}
