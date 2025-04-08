package queries_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/queries"
	"github.com/stretchr/testify/require"
)

// TopicManagerDocumentationProviderAlwaysFailure is an implementation that always returns an error
type TopicManagerDocumentationProviderAlwaysFailure struct{}

func (*TopicManagerDocumentationProviderAlwaysFailure) GetDocumentationForTopicManager(topicManager string) (string, error) {
	return "", errors.New("documentation not found")
}

// TopicManagerDocumentationProviderAlwaysSuccess extends NoopEngineProvider to return custom documentation
type TopicManagerDocumentationProviderAlwaysSuccess struct{}

func (*TopicManagerDocumentationProviderAlwaysSuccess) GetDocumentationForTopicManager(topicManager string) (string, error) {
	return "# Test Documentation\nThis is a test markdown document.", nil
}

func TestTopicManagerDocumentationHandler_Handle_SuccessfulRetrieval(t *testing.T) {
	// Given:
	handler, err := queries.NewTopicManagerDocumentationHandler(&TopicManagerDocumentationProviderAlwaysSuccess{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL + "?topicManager=example")

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var actual queries.TopicManagerDocumentationHandlerResponse
	expected := "# Test Documentation\nThis is a test markdown document."

	require.NoError(t, jsonutil.DecodeResponseBody(res, &actual))
	require.Equal(t, expected, actual.Documentation)
}

func TestTopicManagerDocumentationHandler_Handle_ProviderError(t *testing.T) {
	// Given:
	handler, err := queries.NewTopicManagerDocumentationHandler(&TopicManagerDocumentationProviderAlwaysFailure{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL + "?topicManager=example")

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestTopicManagerDocumentationHandler_Handle_EmptyTopicManagerParameter(t *testing.T) {
	// Given:
	handler, err := queries.NewTopicManagerDocumentationHandler(&TopicManagerDocumentationProviderAlwaysSuccess{})
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var failureResp jsonutil.ResponseFailure
	err = json.NewDecoder(res.Body).Decode(&failureResp)
	require.NoError(t, err)

	require.Equal(t, jsonutil.ReasonInvalidRequest, failureResp.Reason)
	require.Equal(t, "topicManager query parameter is required", failureResp.Hint)
}

func TestNewTopicManagerDocumentationHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider queries.TopicManagerDocumentationProvider = nil

	// When:
	handler, err := queries.NewTopicManagerDocumentationHandler(provider)
	require.Error(t, err)

	// Then:
	require.Nil(t, handler)
}
