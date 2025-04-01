package queries_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/queries"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LookupListProviderAlwaysEmpty is an implementation that always returns an empty list
type LookupListProviderAlwaysEmpty struct{}

func (*LookupListProviderAlwaysEmpty) ListLookupServiceProviders() map[string]*queries.MetaDataLookup {
	return map[string]*queries.MetaDataLookup{}
}

// LookupListProviderAlwaysSuccess is an implementation that always returns a predefined set of lookup providers
type LookupListProviderAlwaysSuccess struct{}

func (*LookupListProviderAlwaysSuccess) ListLookupServiceProviders() map[string]*queries.MetaDataLookup {
	return map[string]*queries.MetaDataLookup{
		"provider1": {
			ShortDescription: "Description 1",
			IconURL: "https://example.com/icon.png",
			Version: "1.0.0",
			InformationURL: "https://example.com/info",
		},
		"provider2": {
			ShortDescription: "Description 2",
			IconURL: "",
			Version: "",
			InformationURL: "",
		},
	}
}

func TestLookupListHandler_Handle_EmptyList(t *testing.T) {
	// Given:
	handler := queries.NewLookupListHandler(&LookupListProviderAlwaysEmpty{})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var result map[string]queries.LookupMetadata
	require.NoError(t, json.NewDecoder(res.Body).Decode(&result))
	assert.Empty(t, result)
}

func TestLookupListHandler_Handle_WithProviders(t *testing.T) {
	// Given:
	handler := queries.NewLookupListHandler(&LookupListProviderAlwaysSuccess{})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var result map[string]queries.LookupMetadata
	require.NoError(t, json.NewDecoder(res.Body).Decode(&result))

	// Check provider1 data
	require.Contains(t, result, "provider1")
	assert.Equal(t, "provider1", result["provider1"].Name)
	assert.Equal(t, "Description 1", result["provider1"].ShortDescription)
	
	iconURL := "https://example.com/icon.png"
	version := "1.0.0"
	infoURL := "https://example.com/info"
	
	assert.Equal(t, &iconURL, result["provider1"].IconURL)
	assert.Equal(t, &version, result["provider1"].Version)
	assert.Equal(t, &infoURL, result["provider1"].InformationURL)

	// Check provider2 data
	require.Contains(t, result, "provider2")
	assert.Equal(t, "provider2", result["provider2"].Name)
	assert.Equal(t, "Description 2", result["provider2"].ShortDescription)
	assert.Nil(t, result["provider2"].IconURL)
	assert.Nil(t, result["provider2"].Version)
	assert.Nil(t, result["provider2"].InformationURL)
}

func TestNewLookupListHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider queries.LookupListProvider = nil

	// When & Then:
	assert.Panics(t, func() {
		queries.NewLookupListHandler(provider)
	}, "Expected panic when provider is nil")
}
