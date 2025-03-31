package queries_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/queries"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LookupDocumentationProviderAlwaysFailure is an implementation that always returns an error
type LookupDocumentationProviderAlwaysFailure struct{}

func (*LookupDocumentationProviderAlwaysFailure) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	return "", errors.New("documentation not found")
}

// LookupDocumentationProviderAlwaysSuccess is an implementation that always returns an success
type LookupDocumentationProviderAlwaysSuccess struct{}

func (*LookupDocumentationProviderAlwaysSuccess) GetDocumentationForLookupServiceProvider(lookupService string) (string, error) {
	return "# Test Documentation\nThis is a test markdown document.", nil
}

func TestLookupDocumentationHandler_Handle_SuccessfulRetrieval(t *testing.T) {
	// Given:
	handler := queries.NewLookupDocumentationHandler(&LookupDocumentationProviderAlwaysSuccess{})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL + "?lookupService=example")

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var actual queries.LookupDocumentationHandlerResponse
	expected := "# Test Documentation\nThis is a test markdown document."

	require.NoError(t, jsonutil.DecodeResponseBody(res, &actual))
	assert.Equal(t, expected, actual.Documentation)
}

func TestLookupDocumentationHandler_Handle_ProviderError(t *testing.T) {
	// Given:
	handler := queries.NewLookupDocumentationHandler(&LookupDocumentationProviderAlwaysFailure{})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL + "?lookupService=example")

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestLookupDocumentationHandler_Handle_EmptyLookupServiceParameter(t *testing.T) {
	// Given:
	handler := queries.NewLookupDocumentationHandler(&LookupDocumentationProviderAlwaysSuccess{})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "lookupService query parameter is required\n", string(body))
}

func TestNewLookupDocumentationHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider queries.LookupDocumentationProvider = nil

	// When & Then:
	assert.Panics(t, func() {
		queries.NewLookupDocumentationHandler(provider)
	}, "Expected panic when provider is nil")
}
