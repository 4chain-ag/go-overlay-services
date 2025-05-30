package ports_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestLookupListHandler_ValidCases(t *testing.T) {
	tests := map[string]struct {
		expectations       testabilities.LookupListProviderMockExpectations
		expected           openapi.LookupServiceProvidersListResponse
		expectedStatusCode int
	}{
		"List lookup service returns an empty lookup service providers list.": {
			expectations: testabilities.LookupListProviderMockExpectations{
				MetadataList:                   testabilities.LookupListEmptyMetadata,
				ListLookupServiceProvidersCall: true,
			},
			expected:           ports.NewLookupListSuccessResponse(testabilities.LookupListEmptyExpectedResponse),
			expectedStatusCode: fiber.StatusOK,
		},
		"List lookup service returns a default lookup service providers list.": {
			expectations: testabilities.LookupListProviderMockExpectations{
				MetadataList:                   testabilities.LookupListDefaultMetadata,
				ListLookupServiceProvidersCall: true,
			},
			expected:           ports.NewLookupListSuccessResponse(testabilities.LookupListDefaultExpectedResponse),
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
			require.Equal(t, tc.expected, actualResponse)
			stub.AssertProvidersState()
		})
	}
}
