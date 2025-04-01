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
			IconURL:          "",
			Version:          "",
			InformationURL:   "",
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

	var result map[string]queries.TopicManagerMetadata
	require.NoError(t, json.NewDecoder(res.Body).Decode(&result))
	assert.Empty(t, result)
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

	var result map[string]queries.TopicManagerMetadata
	require.NoError(t, json.NewDecoder(res.Body).Decode(&result))

	// Check manager1 data
	require.Contains(t, result, "manager1")
	assert.Equal(t, "manager1", result["manager1"].Name)
	assert.Equal(t, "Description 1", result["manager1"].ShortDescription)

	iconURL := "https://example.com/icon.png"
	version := "1.0.0"
	infoURL := "https://example.com/info"

	assert.Equal(t, &iconURL, result["manager1"].IconURL)
	assert.Equal(t, &version, result["manager1"].Version)
	assert.Equal(t, &infoURL, result["manager1"].InformationURL)

	// Check manager2 data
	require.Contains(t, result, "manager2")
	assert.Equal(t, "manager2", result["manager2"].Name)
	assert.Equal(t, "Description 2", result["manager2"].ShortDescription)
	assert.Nil(t, result["manager2"].IconURL)
	assert.Nil(t, result["manager2"].Version)
	assert.Nil(t, result["manager2"].InformationURL)
}

func TestNewTopicManagerListHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider queries.TopicManagerListProvider = nil

	// When & Then:
	assert.Panics(t, func() {
		queries.NewTopicManagerListHandler(provider)
	}, "Expected panic when provider is nil")
}
