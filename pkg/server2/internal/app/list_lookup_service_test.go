package app_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestLookupListService_ValidCases(t *testing.T) {
	tests := map[string]struct {
		expectations testabilities.LookupListProviderMockExpectations
		expected     app.LookupServiceProviders
	}{
		"List lookup service success - empty list": {
			expectations: testabilities.LookupListProviderMockExpectations{
				MetadataList:          testabilities.EmptyMetadata,
				ListLookupServiceProvidersCall: true,
			},
			expected: testabilities.EmptyExpectedResponse,
		},
		"List lookup service success - default list": {
			expectations: testabilities.LookupListProviderMockExpectations{
				MetadataList:                   testabilities.DefaultMetadata,
				ListLookupServiceProvidersCall: true,
			},
			expected: testabilities.DefaultExpectedResponse,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewLookupListProviderMock(t, tc.expectations)
			service := app.NewLookupListService(mock)

			// when:
			response := service.ListLookupServiceProviders()

			// then:
			require.Equal(t, tc.expected, response)
			mock.AssertCalled()
		})
	}
}
