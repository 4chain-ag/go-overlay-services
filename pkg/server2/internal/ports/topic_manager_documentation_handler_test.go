package ports_test

import (
	"testing"

	server2 "github.com/4chain-ag/go-overlay-services/pkg/server2/internal"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestTopicManagerDocumentationHandler_GetDocumentation_ShouldReturnBadRequestResponse(t *testing.T) {
	// Given
	mock := testabilities.NewTopicManagerDocumentationProviderMock(t, testabilities.TopicManagerDocumentationProviderMockExpectations{DocumentationCall: false})
	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithTopicManagerDocumentationProvider(mock))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))
	expectatedResponse := testabilities.NewTestOpenapiErrorResponse(t, app.NewEmptyTopicManagerNameError())

	// When
	var actualResponse openapi.Error
	res, _ := fixture.Client().
		R().
		SetError(&actualResponse).
		Get("/api/v1/getDocumentationForTopicManager")

	// Then
	require.Equal(t, fiber.StatusBadRequest, res.StatusCode())
	require.Equal(t, expectatedResponse, actualResponse)
	mock.AssertCalled()
}

func TestTopicManagerDocumentationHandler_GetDocumentation_ShouldReturnInternalServerErrorResponse(t *testing.T) {
	// Given
	providerError := app.NewTopicManagerDocumentationError(nil)
	mock := testabilities.NewTopicManagerDocumentationProviderMock(t, testabilities.TopicManagerDocumentationProviderMockExpectations{
		DocumentationCall: true,
		Error:             providerError,
	})
	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithTopicManagerDocumentationProvider(mock))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))
	expectedResponse := testabilities.NewTestOpenapiErrorResponse(t, providerError)

	// When
	var actualResponse openapi.Error
	res, _ := fixture.Client().
		R().
		SetError(&actualResponse).
		Get("/api/v1/getDocumentationForTopicManager?topicManager=testProvider")

	// Then
	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode())
	require.Equal(t, expectedResponse, actualResponse)
	mock.AssertCalled()
}

func TestTopicManagerDocumentationHandler_GetDocumentation_ShouldReturnSuccessResponse(t *testing.T) {
	// Given
	mock := testabilities.NewTopicManagerDocumentationProviderMock(t, testabilities.TopicManagerDocumentationProviderMockExpectations{
		DocumentationCall: true,
		Documentation:     testabilities.DefaultTopicManagerDocumentationProviderMockExpectations.Documentation,
	})
	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithTopicManagerDocumentationProvider(mock))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	// When
	var actualResponse openapi.TopicManagerDocumentationResponse
	res, _ := fixture.Client().
		R().
		SetResult(&actualResponse).
		Get("/api/v1/getDocumentationForTopicManager?topicManager=testProvider")

	// Then
	require.Equal(t, fiber.StatusOK, res.StatusCode())
	require.Equal(t, testabilities.DefaultTopicManagerDocumentationProviderMockExpectations.Documentation, actualResponse.Documentation)
	mock.AssertCalled()
}
