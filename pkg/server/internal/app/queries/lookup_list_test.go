package queries_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/queries"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

// LookupListProviderAlwaysEmpty is an implementation that always returns an empty list
type LookupListProviderAlwaysEmpty struct{}

func (*LookupListProviderAlwaysEmpty) ListLookupServiceProviders() map[string]*overlay.MetaData {
	return map[string]*overlay.MetaData{}
}

// LookupListProviderAlwaysSuccess is an implementation that always returns a predefined set of lookup providers
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

func TestLookupListHandler_Handle_EmptyList(t *testing.T) {
	// Given:
	handler, err := queries.NewLookupListHandler(&LookupListProviderAlwaysEmpty{})
	require.NoError(t, err)
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
	require.Empty(t, actual)
}

func TestLookupListHandler_Handle_WithProviders(t *testing.T) {
	// Given:
	handler, err := queries.NewLookupListHandler(&LookupListProviderAlwaysSuccess{})
	require.NoError(t, err)
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
			IconURL:          ptr.To("https://example.com/icon.png"),
			Version:          ptr.To("1.0.0"),
			InformationURL:   ptr.To("https://example.com/info"),
		},
		"provider2": queries.LookupMetadata{
			Name:             "provider2",
			ShortDescription: "Description 2",
			IconURL:          ptr.To("https://example.com/icon2.png"),
			Version:          ptr.To("2.0.0"),
			InformationURL:   ptr.To("https://example.com/info2"),
		},
	}

	require.EqualValues(t, expected, actual)
}

func TestNewLookupListHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider queries.LookupListProvider = nil

	// When:
	handler, err := queries.NewLookupListHandler(provider)
	require.Error(t, err)

	// Then:
	require.Nil(t, handler)
}
