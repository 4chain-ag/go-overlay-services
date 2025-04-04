package queries_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/queries"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

// TopicManagerListProviderAlwaysEmpty is an implementation that always returns an empty list
type TopicManagerListProviderAlwaysEmpty struct{}

func (*TopicManagerListProviderAlwaysEmpty) ListTopicManagers() map[string]*overlay.MetaData {
	return map[string]*overlay.MetaData{}
}

// TopicManagerListProviderAlwaysSuccess is an implementation that always returns a predefined set of topic managers
type TopicManagerListProviderAlwaysSuccess struct{}

func (*TopicManagerListProviderAlwaysSuccess) ListTopicManagers() map[string]*overlay.MetaData {
	return map[string]*overlay.MetaData{
		"manager1": {
			Description: "Description 1",
			Icon:        "https://example.com/icon.png",
			Version:     "1.0.0",
			InfoUrl:     "https://example.com/info",
			Name:        "Manager 1",
		},
		"manager2": {
			Description: "Description 2",
			Icon:        "https://example.com/icon2.png",
			Version:     "1.0.0",
			InfoUrl:     "https://example.com/info",
			Name:        "Manager 2",
		},
	}
}

func TestTopicManagerListHandler_Handle_EmptyList(t *testing.T) {
	// Given:
	handler, err := queries.NewTopicManagerListHandler(&TopicManagerListProviderAlwaysEmpty{})
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

	var actual queries.TopicManagerListHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &actual))
	require.Empty(t, actual)
}

func TestTopicManagerListHandler_Handle_WithManagers(t *testing.T) {
	// Given:
	handler, err := queries.NewTopicManagerListHandler(&TopicManagerListProviderAlwaysSuccess{})
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

	var actual queries.TopicManagerListHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &actual))

	expected := queries.TopicManagerListHandlerResponse{
		"manager1": queries.TopicManagerMetadata{
			Name:           "manager1",
			Description:    "Description 1",
			IconURL:        ptr.To("https://example.com/icon.png"),
			Version:        ptr.To("1.0.0"),
			InformationURL: ptr.To("https://example.com/info"),
		},
		"manager2": queries.TopicManagerMetadata{
			Name:           "manager2",
			Description:    "Description 2",
			IconURL:        ptr.To("https://example.com/icon2.png"),
			Version:        ptr.To("1.0.0"),
			InformationURL: ptr.To("https://example.com/info"),
		},
	}

	require.EqualValues(t, expected, actual)
}

func TestNewTopicManagerListHandler_WithNilProvider(t *testing.T) {
	// Given:
	var provider queries.TopicManagerListProvider = nil

	// When:
	handler, err := queries.NewTopicManagerListHandler(provider)
	require.Error(t, err)

	// Then:
	require.Nil(t, handler)
}
