package overlayhttp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type LookupDocumentationProviderAlwaysFailure struct{}

func (*LookupDocumentationProviderAlwaysFailure) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	return "", errors.New("documentation not found")
}

type LookupDocumentationProviderAlwaysSuccess struct{}

func (*LookupDocumentationProviderAlwaysSuccess) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	return "# Test Documentation\nThis is a test markdown document.", nil
}

func TestLookupServiceDocumentationHandler_Handle_SuccessfulRetrieval(t *testing.T) {
	// Given:
	handler := NewLookupServiceDocumentationHandler(&LookupDocumentationProviderAlwaysSuccess{})
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		params := openapi.LookupServiceDocumentationParams{
			LookupService: "example",
		}
		return handler.Handle(c, params)
	})

	// When:
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)

	// Then:
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var responseBody []byte
	responseBody, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	
	var result map[string]string
	err = json.Unmarshal(responseBody, &result)
	require.NoError(t, err)
	
	const expected = "# Test Documentation\nThis is a test markdown document."
	assert.Equal(t, expected, result["documentation"])
}

func TestLookupDocumentationHandler_Handle_ProviderError(t *testing.T) {
	// Given:
	handler := NewLookupServiceDocumentationHandler(&LookupDocumentationProviderAlwaysFailure{})
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		params := openapi.LookupServiceDocumentationParams{
			LookupService: "example",
		}
		return handler.Handle(c, params)
	})

	// When:
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)

	// Then:
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestLookupDocumentationHandler_Handle_EmptyLookupServiceParameter(t *testing.T) {
	// Given:
	handler := NewLookupServiceDocumentationHandler(&LookupDocumentationProviderAlwaysSuccess{})
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		params := openapi.LookupServiceDocumentationParams{
			LookupService: "",
		}
		return handler.Handle(c, params)
	})

	// When:
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)

	// Then:
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "lookupService query parameter is required", string(body))
}

func TestNewLookupServiceDocumentationHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider LookupServiceDocumentationProvider = nil

	// When/Then:
	assert.Panics(t, func() {
		NewLookupServiceDocumentationHandler(provider)
	})
}
