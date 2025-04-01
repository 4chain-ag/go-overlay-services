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

// TopicManagerListProviderAlwaysEmpty is an implementation that always returns an empty list
type TopicManagerListProviderAlwaysEmpty struct{}

func (*TopicManagerListProviderAlwaysEmpty) ListTopicManagers() map[string]*queries.MetaData {
	return map[string]*queries.MetaData{}
}

// TopicManagerListProviderAlwaysSuccess is an implementation that always returns a predefined set of topic managers
type TopicManagerListProviderAlwaysSuccess struct{}

func (*TopicManagerListProviderAlwaysSuccess) ListTopicManagers() map[string]*queries.MetaData {
	return map[string]*queries.MetaData{
		"manager1": {
			ShortDescription: "Description 1",
			IconURL:          "https://example.com/icon.png",
			Version:          "1.0.0",
			InformationURL:   "https://example.com/info",
		},
		"manager2": {
			ShortDescription: "Description 2",
			IconURL:          "https://example.com/icon2.png",
			Version:          "1.0.0",
			InformationURL:   "https://example.com/info",
		},
	}
}

func TestTopicManagerListHandler_Handle_EmptyList(t *testing.T) {
	// Given:
	handler := queries.NewTopicManagerListHandler(&TopicManagerListProviderAlwaysEmpty{})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var actual queries.TopicManagerListHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &actual))
	assert.Empty(t, actual)
}

func TestTopicManagerListHandler_Handle_WithManagers(t *testing.T) {
	// Given:
	handler := queries.NewTopicManagerListHandler(&TopicManagerListProviderAlwaysSuccess{})
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// When:
	res, err := ts.Client().Get(ts.URL)

	// Then:
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var actual queries.TopicManagerListHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &actual))

	expected := queries.TopicManagerListHandlerResponse{
		"manager1": queries.TopicManagerMetadata{
			Name:             "manager1",
			ShortDescription: "Description 1",
			IconURL:          toPtr("https://example.com/icon.png"),
			Version:          toPtr("1.0.0"),
			InformationURL:   toPtr("https://example.com/info"),
		},
		"manager2": queries.TopicManagerMetadata{
			Name:             "manager2",
			ShortDescription: "Description 2",
			IconURL:          toPtr("https://example.com/icon2.png"),
			Version:          toPtr("1.0.0"),
			InformationURL:   toPtr("https://example.com/info"),
		},
	}

	assert.EqualValues(t, expected, actual)
}

func TestNewTopicManagerListHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider queries.TopicManagerListProvider = nil

	// When & Then:
	assert.Panics(t, func() {
		queries.NewTopicManagerListHandler(provider)
	}, "Expected panic when provider is nil")
}

func toPtr[T any](x T) *T { return &x }
