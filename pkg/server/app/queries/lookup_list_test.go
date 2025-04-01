package queries_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
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
			IconURL:          "https://example.com/icon.png",
			Version:          "1.0.0",
			InformationURL:   "https://example.com/info",
		},
		"provider2": {
			ShortDescription: "Description 2",
			IconURL:          "https://example.com/icon2.png",
			Version:          "2.0.0",
			InformationURL:   "https://example.com/info2",
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

	var actual queries.LookupListHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &actual))
	assert.Empty(t, actual)
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

	var actual queries.LookupListHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &actual))

	expected := queries.LookupListHandlerResponse{
		"provider1": queries.LookupMetadata{
			Name:             "provider1",
			ShortDescription: "Description 1",
			IconURL:          toPtr("https://example.com/icon.png"),
			Version:          toPtr("1.0.0"),
			InformationURL:   toPtr("https://example.com/info"),
		},
		"provider2": queries.LookupMetadata{
			Name:             "provider2",
			ShortDescription: "Description 2",
			IconURL:          toPtr("https://example.com/icon2.png"),
			Version:          toPtr("2.0.0"),
			InformationURL:   toPtr("https://example.com/info2"),
		},
	}

	assert.EqualValues(t, expected, actual)
}

func TestNewLookupListHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider queries.LookupListProvider = nil

	// When & Then:
	assert.Panics(t, func() {
		queries.NewLookupListHandler(provider)
	}, "Expected panic when provider is nil")
}
