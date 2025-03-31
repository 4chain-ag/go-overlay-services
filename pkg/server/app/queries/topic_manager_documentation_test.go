package queries_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/queries"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TopicManagerDocumentationProviderAlwaysFailure is an implementation that always returns an error
type TopicManagerDocumentationProviderAlwaysFailure struct{}

func (*TopicManagerDocumentationProviderAlwaysFailure) GetDocumentationForTopicManager(topicManager string) (string, error) {
	return "", errors.New("documentation not found")
}

// TopicManagerDocumentationProviderAlwaysSuccess extends NoopEngineProvider to return custom documentation
type TopicManagerDocumentationProviderAlwaysSuccess struct {
	*server.NoopEngineProvider
}

func (*TopicManagerDocumentationProviderAlwaysSuccess) GetDocumentationForTopicManager(topicManager string) (string, error) {
	return "# Test Documentation\nThis is a test markdown document.", nil
}

func TestTopicManagerDocumentationHandler_Handle_SuccessfulRetrieval(t *testing.T) {
	// Given:
	handler := queries.NewTopicManagerDocumentationHandler(&TopicManagerDocumentationProviderAlwaysSuccess{server.NewNoopEngineProvider().(*server.NoopEngineProvider)})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL + "?topicManager=example")

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
	handler := queries.NewTopicManagerDocumentationHandler(&TopicManagerDocumentationProviderAlwaysFailure{})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL + "?topicManager=example")

	// Then:
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	defer res.Body.Close()
}

func TestTopicManagerDocumentationHandler_Handle_EmptyTopicManagerParameter(t *testing.T) {
	// Given:
	handler := queries.NewTopicManagerDocumentationHandler(server.NewNoopEngineProvider())
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL)

	// Then:
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "topicManager query parameter is required\n", string(body))
}

func TestNewTopicManagerDocumentationHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider queries.TopicManagerDocumentationProvider = nil

	// When & Then:
	assert.Panics(t, func() {
		queries.NewTopicManagerDocumentationHandler(provider)
	}, "Expected panic when provider is nil")
}
