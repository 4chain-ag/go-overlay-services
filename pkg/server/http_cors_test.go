package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/4chain-ag/go-overlay-services/pkg/server/config"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

const disableTimout = -1

func createTestServerWithCORS() *fiber.App {
	cfg := config.DefaultConfig()
	return server.New(
		server.WithConfig(cfg),
		server.WithCORS(),
	).App()
}

func assertCORS(t *testing.T, resp *http.Response, expectMethods bool, expectedMethods, expectedHeaders []string) {
	t.Helper()

	require.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"), "missing or invalid Access-Control-Allow-Origin")

	if expectMethods {
		methods := resp.Header.Get("Access-Control-Allow-Methods")
		for _, m := range expectedMethods {
			require.Contains(t, methods, m, "missing method in Access-Control-Allow-Methods")
		}
	}

	headers := resp.Header.Get("Access-Control-Allow-Headers")
	for _, h := range expectedHeaders {
		require.Contains(t, headers, h, "missing header in Access-Control-Allow-Headers")
	}
}

func TestCORSOptionsRequest_ShouldReturnCORSHeaders(t *testing.T) {
	// Given:
	app := createTestServerWithCORS()

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/submit", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Origin, Content-Type")

	// When:
	resp, err := app.Test(req, disableTimout)

	// Then:
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	assertCORS(t, resp, true,
		[]string{"POST"},
		[]string{"Origin", "Content-Type"},
	)
}

func TestCORSPostRequest_ShouldReturnCreatedWithCORS(t *testing.T) {
	// Given:
	app := createTestServerWithCORS()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/submit", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Content-Type", "application/json")

	// When:
	resp, err := app.Test(req, disableTimout)

	// Then:
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	assertCORS(t, resp, false, nil, nil)
}
