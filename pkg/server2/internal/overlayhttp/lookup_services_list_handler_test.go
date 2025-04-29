package overlayhttp

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

type LookupListProviderAlwaysEmpty struct{}

func (*LookupListProviderAlwaysEmpty) ListLookupServiceProviders() map[string]*overlay.MetaData {
	return map[string]*overlay.MetaData{}
}

type LookupListProviderAlwaysSuccess struct{}

func (*LookupListProviderAlwaysSuccess) ListLookupServiceProviders() map[string]*overlay.MetaData {
	return map[string]*overlay.MetaData{
		"provider1": {
			Description: "Description 1",
			Icon:        "https://example.com/icon.png",
			Version:     "1.0.0",
			InfoUrl:     "https://example.com/info",
		},
		"provider2": {
			Description: "Description 2",
			Icon:        "https://example.com/icon2.png",
			Version:     "2.0.0",
			InfoUrl:     "https://example.com/info2",
		},
	}
}

func TestLookupServicesListHandler_Handle_EmptyList(t *testing.T) {
	// Given:
	handler := NewLookupServicesListHandler(&LookupListProviderAlwaysEmpty{})
	app := fiber.New()
	
	app.Get("/test", handler.Handle)

	// When:
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)

	// Then:
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "{}", string(body))
}

func TestLookupServicesListHandler_Handle_WithProviders(t *testing.T) {
	// Given:
	handler := NewLookupServicesListHandler(&LookupListProviderAlwaysSuccess{})
	app := fiber.New()
	
	app.Get("/test", handler.Handle)

	// When:
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)

	// Then:
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var actual map[string]openapi.LookupMetadata
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(body, &actual))
	
	expected := map[string]openapi.LookupMetadata{
		"provider1": {
			Name:             "provider1",
			ShortDescription: "Description 1",
			IconURL:          ptr.To("https://example.com/icon.png"),
			Version:          ptr.To("1.0.0"),
			InformationURL:   ptr.To("https://example.com/info"),
		},
		"provider2": {
			Name:             "provider2",
			ShortDescription: "Description 2",
			IconURL:          ptr.To("https://example.com/icon2.png"),
			Version:          ptr.To("2.0.0"),
			InformationURL:   ptr.To("https://example.com/info2"),
		},
	}
	
	require.Equal(t, len(expected), len(actual))
	for provider, expectedMetadata := range expected {
		actualMetadata, exists := actual[provider]
		require.True(t, exists, "Provider %s missing from response", provider)
		
		assert.Equal(t, expectedMetadata.Name, actualMetadata.Name)
		assert.Equal(t, expectedMetadata.ShortDescription, actualMetadata.ShortDescription)
		assert.Equal(t, expectedMetadata.IconURL, actualMetadata.IconURL)
		assert.Equal(t, expectedMetadata.Version, actualMetadata.Version)
		assert.Equal(t, expectedMetadata.InformationURL, actualMetadata.InformationURL)
	}
}

func TestNewLookupServicesListHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider LookupServicesListProvider = nil

	// When/Then:
	assert.Panics(t, func() {
		NewLookupServicesListHandler(provider)
	})
} 
