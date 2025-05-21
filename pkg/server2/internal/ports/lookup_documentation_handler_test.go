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
	// given:
	mock := testabilities.NewLookupServiceDocumentationProviderMock(t, testabilities.LookupServiceDocumentationProviderMockExpectations{DocumentationCall: false})
	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupDocumentationProvider(mock))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))
	expectatedResponse := testabilities.NewTestOpenapiErrorResponse(t, app.NewEmptyLookupServiceNameError())

	// when:
	var actualResponse openapi.Error
	res, _ := fixture.Client().
		R().
		SetError(&actualResponse).
		Get("/api/v1/getDocumentationForLookupServiceProvider")

	// then:
	require.Equal(t, fiber.StatusBadRequest, res.StatusCode())
	require.Equal(t, expectatedResponse, actualResponse)
	mock.AssertCalled()
}

func TestLookupProviderDocumentationHandler_GetDocumentation_ShouldReturnInternalServerErrorResponse(t *testing.T) {
	// given:
	providerError := app.NewLookupServiceProviderDocumentationError(nil)
	mock := testabilities.NewLookupServiceDocumentationProviderMock(t, testabilities.LookupServiceDocumentationProviderMockExpectations{
		DocumentationCall: true,
		Error:             providerError,
	})
	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupDocumentationProvider(mock))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))
	expectedResponse := testabilities.NewTestOpenapiErrorResponse(t, providerError)

	// when:
	var actualResponse openapi.Error
	res, _ := fixture.Client().
		R().
		SetError(&actualResponse).
		Get("/api/v1/getDocumentationForLookupServiceProvider?lookupService=testProvider")

	// then:
	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode())
	require.Equal(t, expectedResponse, actualResponse)
	mock.AssertCalled()
}

func TestLookupProviderDocumentationHandler_GetDocumentation_ShouldReturnSuccessResponse(t *testing.T) {
	// given:
	expectedDocumentation := "# Test Documentation\nThis is a test markdown document."
	mock := testabilities.NewLookupServiceDocumentationProviderMock(t, testabilities.LookupServiceDocumentationProviderMockExpectations{
		DocumentationCall: true,
		Documentation:     expectedDocumentation,
	})
	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupDocumentationProvider(mock))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	// when:
	var actualResponse openapi.LookupServiceDocumentationResponse
	res, _ := fixture.Client().
		R().
		SetResult(&actualResponse).
		Get("/api/v1/getDocumentationForLookupServiceProvider?lookupService=testProvider")

	// then:
	require.Equal(t, fiber.StatusOK, res.StatusCode())
	require.Equal(t, expectedDocumentation, actualResponse.Documentation)
	mock.AssertCalled()
}
