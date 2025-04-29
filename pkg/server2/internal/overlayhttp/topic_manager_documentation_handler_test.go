package overlayhttp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TopicManagerDocumentationProviderAlwaysFailure struct{}

func (*TopicManagerDocumentationProviderAlwaysFailure) GetDocumentationForTopicManager(topicManager string) (string, error) {

	return "", errors.New("documentation not found")

}

type TopicManagerDocumentationProviderAlwaysSuccess struct{}

func (*TopicManagerDocumentationProviderAlwaysSuccess) GetDocumentationForTopicManager(topicManager string) (string, error) {

	return "# Test Documentation\nThis is a test markdown document.", nil

}

func TestTopicManagerDocumentationHandler_Handle_SuccessfulRetrieval(t *testing.T) {

	// Given:

	handler := NewTopicManagerDocumentationHandler(&TopicManagerDocumentationProviderAlwaysSuccess{})

	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {

		params := openapi.TopicManagerDocumentationParams{

			TopicManager: "example",
		}

		return handler.Handle(c, params)

	})

	// When:

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	resp, err := app.Test(req)

	// Then:

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var responseBody []byte

	responseBody, err = io.ReadAll(resp.Body)

	require.NoError(t, err)

	var result map[string]string

	err = json.Unmarshal(responseBody, &result)

	require.NoError(t, err)

	const expected = "# Test Documentation\nThis is a test markdown document."

	assert.Equal(t, expected, result["documentation"])

}

func TestTopicManagerDocumentationHandler_Handle_ProviderError(t *testing.T) {

	// Given:

	handler := NewTopicManagerDocumentationHandler(&TopicManagerDocumentationProviderAlwaysFailure{})

	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {

		params := openapi.TopicManagerDocumentationParams{

			TopicManager: "example",
		}

		return handler.Handle(c, params)

	})

	// When:

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	resp, err := app.Test(req)

	// Then:

	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func TestTopicManagerDocumentationHandler_Handle_EmptyTopicManagerParameter(t *testing.T) {

	// Given:

	handler := NewTopicManagerDocumentationHandler(&TopicManagerDocumentationProviderAlwaysSuccess{})

	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {

		params := openapi.TopicManagerDocumentationParams{

			TopicManager: "",
		}

		return handler.Handle(c, params)

	})

	// When:

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	resp, err := app.Test(req)

	// Then:

	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)

	require.NoError(t, err)

	assert.Equal(t, "topicManager query parameter is required", string(body))

}

func TestNewTopicManagerDocumentationHandler_WithNilProvider(t *testing.T) {

	// Given:

	var provider TopicManagerDocumentationProvider = nil

	// When/Then:

	assert.Panics(t, func() {

		NewTopicManagerDocumentationHandler(provider)

	})

}
