package app_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

func TestLookupListService_ValidCases(t *testing.T) {
	tests := map[string]struct {
		expectations testabilities.LookupListProviderMockExpectations
		expected     map[string]*overlay.MetaData
	}{
		"List lookup service returns an empty lookup service providers list.": {
			expectations: testabilities.LookupListProviderMockExpectations{
				MetadataList:                   testabilities.LookupListDefaultMetadata,
				ListLookupServiceProvidersCall: true,
			},
			expected: testabilities.LookupListDefaultMetadata,
		},
		"List lookup service returns a default lookup service providers list.": {
			expectations: testabilities.LookupListProviderMockExpectations{
				MetadataList:                   map[string]*overlay.MetaData{},
				ListLookupServiceProvidersCall: true,
			},
			expected: map[string]*overlay.MetaData{},
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
