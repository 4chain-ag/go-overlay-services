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
		expectedDTO  app.LookupServicesMetadataDTO
	}{
		"List lookup service returns a default lookup service providers list.": {
			expectations: testabilities.LookupListProviderMockExpectations{
				MetadataList:                   testabilities.DefaultOverlayMetadata,
				ListLookupServiceProvidersCall: true,
			},
			expectedDTO: app.NewLookupServicesMetadataDTO(testabilities.DefaultOverlayMetadata),
		},
		"List lookup service returns an empty lookup service providers list.": {
			expectations: testabilities.LookupListProviderMockExpectations{
				MetadataList:                   map[string]*overlay.MetaData{},
				ListLookupServiceProvidersCall: true,
			},
			expectedDTO: app.LookupServicesMetadataDTO{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewLookupListProviderMock(t, tc.expectations)
			service := app.NewLookupListService(mock)

			// when:
			actualDTO := service.ListLookupServiceProviders()

			// then:
			require.Equal(t, tc.expectedDTO, actualDTO)
			mock.AssertCalled()
		})
	}
}
