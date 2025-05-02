package ports_test

import (
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestTopicManagerDocumentationHandler_GetDocumentation_ShouldReturnBadRequestResponse(t *testing.T) {

	// Given

	engine := testabilities.NewTestOverlayEngineStub(t)

	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))

	// When

	var actualResponse openapi.Error

	res, _ := fixture.Client().
		R().
		SetError(&actualResponse).
		Get("/api/v1/getDocumentationForTopicManager")

	// Then

	expectedResponse := ports.NewMissingTopicManagerParameterResponse()

	require.Equal(t, fiber.StatusBadRequest, res.StatusCode())

	require.Equal(t, expectedResponse.Message, actualResponse.Message)

}

func TestTopicManagerDocumentationHandler_GetDocumentation_ShouldReturnInternalServerErrorResponse(t *testing.T) {

	// Given

	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithTopicManagerDocumentationError())

	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))

	// When

	var actualResponse openapi.Error

	res, _ := fixture.Client().
		R().
		SetError(&actualResponse).
		Get("/api/v1/getDocumentationForTopicManager?topicManager=testTopicManager")

	// Then

	expectedResponse := ports.NewTopicManagerProviderErrorResponse()

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode())

	require.Equal(t, expectedResponse.Message, actualResponse.Message)

}

func TestTopicManagerDocumentationHandler_GetDocumentation_ShouldReturnSuccessResponse(t *testing.T) {

	// Given

	expectedDocumentation := "# Test Documentation\nThis is a test markdown document."

	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithTopicManagerDocumentation(expectedDocumentation))

	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))

	// When

	var actualResponse openapi.TopicManagerDocumentationResponse

	res, _ := fixture.Client().
		R().
		SetResult(&actualResponse).
		Get("/api/v1/getDocumentationForTopicManager?topicManager=testTopicManager")

	// Then

	require.Equal(t, http.StatusOK, res.StatusCode())

	require.Equal(t, expectedDocumentation, actualResponse.Documentation)

}
