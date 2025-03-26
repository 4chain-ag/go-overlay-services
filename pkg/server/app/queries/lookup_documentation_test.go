package queries_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/dto"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/queries"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ErrorEngineProvider is an implementation that always returns an error
type ErrorEngineProvider struct{}

func (*ErrorEngineProvider) GetDocumentationForLookupServiceProvider(provider string) (string, error) {
	return "", errors.New("documentation not found")
}

// setupTest sets up a new Fiber app and handler for testing
func setupTest(provider queries.LookupDocumentationProvider) (*fiber.App, *queries.LookupDocumentationHandler) {
	app := fiber.New()
	handler := queries.NewLookupDocumentationHandler(provider)
	return app, handler
}

func TestLookupDocumentationHandler_Handle_SuccessfulRetrieval(t *testing.T) {
	// Given:
	noopProvider := server.NewNoopEngineProvider()
	app, handler := setupTest(noopProvider)
	app.Get("/docs", handler.Handle)

	// When:
	req := httptest.NewRequest(http.MethodGet, "/docs?lookupService=example", nil)
	resp, err := app.Test(req)

	// Then:
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	contentType := resp.Header.Get("Content-Type")
	require.NotEmpty(t, contentType, "Content-Type header should not be empty")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "", string(body))
}

func TestLookupDocumentationHandler_Handle_ProviderError(t *testing.T) {
	// Given:
	errorProvider := &ErrorEngineProvider{}
	app, handler := setupTest(errorProvider)
	app.Get("/docs", handler.Handle)

	// When:
	req := httptest.NewRequest(http.MethodGet, "/docs?lookupService=example", nil)
	resp, err := app.Test(req)

	// Then:
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response dto.HandlerResponse
	err = readJSONResponse(resp, &response)
	require.NoError(t, err)
	assert.Equal(t, dto.HandlerResponseNonOK.Message, response.Message)
}

func TestLookupDocumentationHandler_Handle_EmptyLookupServiceParameter(t *testing.T) {
	// Given:
	noopProvider := server.NewNoopEngineProvider()
	app, handler := setupTest(noopProvider)
	app.Get("/docs", handler.Handle)

	// When:
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	resp, err := app.Test(req)

	// Then:
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var errorResponse map[string]string
	err = readJSONResponse(resp, &errorResponse)
	require.NoError(t, err)
	assert.Equal(t, "lookupService query parameter is required", errorResponse["error"])
}

func TestNewLookupDocumentationHandler_WithNilProvider(t *testing.T) {
	// Given:
	// When:
	handler := queries.NewLookupDocumentationHandler(nil)

	// Then:
	assert.Nil(t, handler, "Expected nil when provider is nil")
}

// Helper function to read JSON response
func readJSONResponse(resp *http.Response, v interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
