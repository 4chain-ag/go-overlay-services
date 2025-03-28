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

// ErrorTopicManagerProvider extends NoopEngineProvider to simulate an error when retrieving documentation.
type ErrorTopicManagerProvider struct {
	*server.NoopEngineProvider
}

func (*ErrorTopicManagerProvider) GetDocumentationForTopicManager(provider string) (string, error) {
	return "", errors.New("documentation not found")
}

// CustomSuccessTopicManagerProvider extends NoopEngineProvider to return custom documentation.
type CustomSuccessTopicManagerProvider struct {
	*server.NoopEngineProvider
}

func (*CustomSuccessTopicManagerProvider) GetDocumentationForTopicManager(provider string) (string, error) {
	return "# Test Documentation\nThis is a test markdown document.", nil
}

func TestTopicManagerDocumentationHandler_Handle_SuccessfulRetrieval(t *testing.T) {
	// Given:
	handler := queries.NewTopicManagerDocumentationHandler(&CustomSuccessTopicManagerProvider{
		server.NewNoopEngineProvider().(*server.NoopEngineProvider),
	})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL + "?provider=example")

	// Then:
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))
	defer res.Body.Close()

	var response queries.TopicManagerDocumentationHandlerResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "# Test Documentation\nThis is a test markdown document.", response.Documentation)
}

func TestTopicManagerDocumentationHandler_Handle_ProviderError(t *testing.T) {
	// Given:
	handler := queries.NewTopicManagerDocumentationHandler(&ErrorTopicManagerProvider{
		server.NewNoopEngineProvider().(*server.NoopEngineProvider),
	})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL + "?provider=example")

	// Then:
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	defer res.Body.Close()
}

func TestTopicManagerDocumentationHandler_Handle_EmptyProviderParameter(t *testing.T) {
	// Given:
	handler := queries.NewTopicManagerDocumentationHandler(&CustomSuccessTopicManagerProvider{
		server.NewNoopEngineProvider().(*server.NoopEngineProvider),
	})
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
	assert.Equal(t, "provider query parameter is required", errorResponse["error"])
}

func TestNewTopicManagerDocumentationHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider queries.TopicManagerDocumentationProvider = nil

	// When & Then:
	assert.Panics(t, func() {
		queries.NewTopicManagerDocumentationHandler(provider)
	}, "Expected panic when provider is nil")
}
