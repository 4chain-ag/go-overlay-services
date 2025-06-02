package ports_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestLookupListHandler_ValidCases(t *testing.T) {
	tests := map[string]struct {
		expectations       testabilities.LookupListProviderMockExpectations
		expectedResponse   openapi.LookupServiceProvidersListResponse
		expectedStatusCode int
	}{
		"List lookup service returns a default lookup service providers list.": {
			expectations: testabilities.LookupListProviderMockExpectations{
				MetadataList:                   testabilities.DefaultOverlayMetadata,
				ListLookupServiceProvidersCall: true,
			},
			expectedResponse: ports.NewLookupServicesMetadataSuccessResponse(app.LookupServicesMetadataDTO{
				"lookup_service1": {
					Description: "Description 1",
					IconURL:     "https://example.com/icon.png",
					Version:     "1.0.0",
					InfoURL:     "https://example.com/info",
				},
				"lookup_service2": {
					Description: "Description 2",
					IconURL:     "https://example.com/icon2.png",
					Version:     "2.0.0",
					InfoURL:     "https://example.com/info2",
				},
			}),
			expectedStatusCode: fiber.StatusOK,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupListProvider(testabilities.NewLookupListProviderMock(t, tc.expectations)))
			fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

			// when:
			var actualResponse openapi.LookupServiceProvidersListResponse
			res, _ := fixture.Client().
				R().
				SetResult(&actualResponse).
				Get("/api/v1/listLookupServiceProviders")

			// then:
			require.Equal(t, tc.expectedStatusCode, res.StatusCode())
			require.Equal(t, tc.expectedResponse, actualResponse)
			stub.AssertProvidersState()
		})
	}
}
