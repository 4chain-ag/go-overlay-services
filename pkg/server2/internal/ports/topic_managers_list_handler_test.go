package overlayhttp

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

type TopicManagerListProviderAlwaysEmpty struct{}

func (*TopicManagerListProviderAlwaysEmpty) ListTopicManagers() map[string]*overlay.MetaData {

	return map[string]*overlay.MetaData{}

}

type TopicManagerListProviderAlwaysSuccess struct{}

func (*TopicManagerListProviderAlwaysSuccess) ListTopicManagers() map[string]*overlay.MetaData {

	return map[string]*overlay.MetaData{

		"manager1": {

			Description: "Description 1",

			Icon: "https://example.com/icon.png",

			Version: "1.0.0",

			InfoUrl: "https://example.com/info",

			Name: "Manager 1",
		},

		"manager2": {

			Description: "Description 2",

			Icon: "https://example.com/icon2.png",

			Version: "1.0.0",

			InfoUrl: "https://example.com/info",

			Name: "Manager 2",
		},
	}

}

func TestTopicManagersListHandler_Handle_EmptyList(t *testing.T) {

	// Given:

	handler := NewTopicManagersListHandler(&TopicManagerListProviderAlwaysEmpty{})

	app := fiber.New()

	app.Get("/test", handler.Handle)

	// When:

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	resp, err := app.Test(req)

	// Then:

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)

	require.NoError(t, err)

	assert.Equal(t, "{}", string(body))

}

func TestTopicManagersListHandler_Handle_WithManagers(t *testing.T) {

	// Given:

	handler := NewTopicManagersListHandler(&TopicManagerListProviderAlwaysSuccess{})

	app := fiber.New()

	app.Get("/test", handler.Handle)

	// When:

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	resp, err := app.Test(req)

	// Then:

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var actual map[string]openapi.TopicManagerMetadata

	body, err := io.ReadAll(resp.Body)

	require.NoError(t, err)

	require.NoError(t, json.Unmarshal(body, &actual))

	expected := map[string]openapi.TopicManagerMetadata{

		"manager1": {

			Name: "manager1",

			Description: "Description 1",

			IconURL: ptr.To("https://example.com/icon.png"),

			Version: ptr.To("1.0.0"),

			InformationURL: ptr.To("https://example.com/info"),
		},

		"manager2": {

			Name: "manager2",

			Description: "Description 2",

			IconURL: ptr.To("https://example.com/icon2.png"),

			Version: ptr.To("1.0.0"),

			InformationURL: ptr.To("https://example.com/info"),
		},
	}

	require.Equal(t, len(expected), len(actual))

	for manager, expectedMetadata := range expected {

		actualMetadata, exists := actual[manager]

		require.True(t, exists, "Manager %s missing from response", manager)

		assert.Equal(t, expectedMetadata.Name, actualMetadata.Name)

		assert.Equal(t, expectedMetadata.Description, actualMetadata.Description)

		assert.Equal(t, expectedMetadata.IconURL, actualMetadata.IconURL)

		assert.Equal(t, expectedMetadata.Version, actualMetadata.Version)

		assert.Equal(t, expectedMetadata.InformationURL, actualMetadata.InformationURL)

	}

}

func TestNewTopicManagersListHandler_WithNilProvider(t *testing.T) {

	// Given:

	var provider TopicManagersListProvider = nil

	// When/Then:

	assert.Panics(t, func() {

		NewTopicManagersListHandler(provider)

	})

}
