package queries_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/queries"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ErrorEngineProvider is an implementation that always returns an error
type ErrorEngineProvider struct{}

func (*ErrorEngineProvider) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	return "", errors.New("documentation not found")
}

// CustomSuccessEngineProvider extends NoopEngineProvider to return custom documentation
type CustomSuccessEngineProvider struct {
	*server.NoopEngineProvider
}

func (*CustomSuccessEngineProvider) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	return "# Test Documentation\nThis is a test markdown document.", nil
}

func TestLookupDocumentationHandler_Handle_SuccessfulRetrieval(t *testing.T) {
	// Given:
	handler := queries.NewLookupDocumentationHandler(&CustomSuccessEngineProvider{server.NewNoopEngineProvider().(*server.NoopEngineProvider)})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL + "?lookupService=example")

	// Then:
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))
	defer res.Body.Close()

	var response queries.LookupDocumentationHandlerResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "# Test Documentation\nThis is a test markdown document.", response.Documentation)
}

func TestLookupDocumentationHandler_Handle_ProviderError(t *testing.T) {
	// Given:
	handler := queries.NewLookupDocumentationHandler(&ErrorEngineProvider{})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL + "?lookupService=example")

	// Then:
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	defer res.Body.Close()
}

func TestLookupDocumentationHandler_Handle_EmptyLookupServiceParameter(t *testing.T) {
	// Given:
	handler := queries.NewLookupDocumentationHandler(server.NewNoopEngineProvider())
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL)

	// Then:
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))
	defer res.Body.Close()

	var errorResponse map[string]string
	err = json.NewDecoder(res.Body).Decode(&errorResponse)
	require.NoError(t, err)
	assert.Equal(t, "lookupService query parameter is required", errorResponse["error"])
}

func TestNewLookupDocumentationHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider queries.LookupDocumentationProvider = nil

	// When:
	handler := queries.NewLookupDocumentationHandler(provider)

	// Then:
	assert.Nil(t, handler, "Expected nil when provider is nil")
}
