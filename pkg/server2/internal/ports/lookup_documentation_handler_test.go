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

func TestLookupProviderDocumentationHandler_GetDocumentation_ShouldReturnBadRequestResponse(t *testing.T) {
	// Given
	mock := testabilities.NewLookupServiceDocumentationProviderMock(t, testabilities.LookupServiceDocumentationProviderMockExpectations{DocumentationCall: false})
	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupDocumentationProvider(mock))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))
	expectatedResponse := testabilities.NewTestOpenapiErrorResponse(t, app.NewEmptyLookupServiceNameError())

	// When
	var actualResponse openapi.Error
	res, _ := fixture.Client().
		R().
		SetError(&actualResponse).
		Get("/api/v1/getDocumentationForLookupServiceProvider")

	// Then
	require.Equal(t, fiber.StatusBadRequest, res.StatusCode())
	require.Equal(t, expectatedResponse, actualResponse)
	mock.AssertCalled()
}

// TODO: Implement these tests

// func TestLookupProviderDocumentationHandler_GetDocumentation_ShouldReturnInternalServerErrorResponse(t *testing.T) {
// 	// Given
// 	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupProviderDocumentationError())
// 	fixture := server2.NewServerTestFixture(t, server2.WithEngine(engine))

// 	// When
// 	var actualResponse openapi.Error
// 	res, _ := fixture.Client().
// 		R().
// 		SetError(&actualResponse).
// 		Get("/api/v1/getDocumentationForLookupServiceProvider?lookupService=testProvider")

// 	// Then
// 	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode())
// 	require.Equal(t, ports.LookupProviderError.Message, actualResponse.Message)
// }

// func TestLookupProviderDocumentationHandler_GetDocumentation_ShouldReturnSuccessResponse(t *testing.T) {
// 	// Given
// 	expectedDocumentation := "# Test Documentation\nThis is a test markdown document."
// 	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupProviderDocumentation(expectedDocumentation))
// 	fixture := server2.NewServerTestFixture(t, server2.WithEngine(engine))

// 	// When
// 	var actualResponse openapi.LookupServiceDocumentationResponse
// 	res, _ := fixture.Client().
// 		R().
// 		SetResult(&actualResponse).
// 		Get("/api/v1/getDocumentationForLookupServiceProvider?lookupService=testProvider")

// 	// Then
// 	require.Equal(t, http.StatusOK, res.StatusCode())
// 	require.Equal(t, expectedDocumentation, actualResponse.Documentation)
// }
