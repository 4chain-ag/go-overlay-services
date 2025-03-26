package queries_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/dto"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/queries"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockLookupDocumentationProvider is a mock implementation of the LookupDocumentationProvider interface
type mockLookupDocumentationProvider struct {
	services map[string]string
}

// NewMockLookupDocumentationProvider creates a new mock provider with the given services
func NewMockLookupDocumentationProvider(services map[string]string) *mockLookupDocumentationProvider {
	if services == nil {
		services = make(map[string]string)
	}
	return &mockLookupDocumentationProvider{
		services: services,
	}
}

func (m *mockLookupDocumentationProvider) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	if doc, ok := m.services[lookupService]; ok {
		return doc, nil
	}
	return "", errors.New("no documentation found")
}

// setupTest sets up a new Fiber app and handler for testing
func setupTest(provider queries.LookupDocumentationProvider) (*fiber.App, *queries.LookupDocumentationHandler) {
	app := fiber.New()
	handler := queries.NewLookupDocumentationHandler(provider)
	return app, handler
}

func TestLookupDocumentationHandler_Handle_SuccessfulRetrieval(t *testing.T) {
	// Given:
	expectedDocumentation := "# Lookup Service Documentation\n\nThis is the documentation for the lookup service."
	services := map[string]string{
		"example": expectedDocumentation,
	}
	mockProvider := NewMockLookupDocumentationProvider(services)
	app, handler := setupTest(mockProvider)
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
	assert.Equal(t, expectedDocumentation, string(body))
}

func TestLookupDocumentationHandler_Handle_ProviderNotFound(t *testing.T) {
	// Given:
	services := map[string]string{
		"existing": "Some documentation",
	}
	mockProvider := NewMockLookupDocumentationProvider(services)
	app, handler := setupTest(mockProvider)
	app.Get("/docs", handler.Handle)

	// When:
	req := httptest.NewRequest(http.MethodGet, "/docs?lookupService=nonexistent", nil)
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
	mockProvider := NewMockLookupDocumentationProvider(nil)
	app, handler := setupTest(mockProvider)
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
