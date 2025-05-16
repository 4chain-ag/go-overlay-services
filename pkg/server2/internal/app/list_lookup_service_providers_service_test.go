package app_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

func TestLookupServicesListService_EmptyList(t *testing.T) {

	// given:

	expectations := testabilities.LookupServicesListProviderMockExpectations{

		MetadataList: map[string]*overlay.MetaData{},

		ListLookupServiceProvidersCall: true,
	}

	mock := testabilities.NewLookupServicesListProviderMock(t, expectations)

	service, err := app.NewLookupServicesListService(mock)

	require.NoError(t, err)

	// when:

	response := service.ListLookupServiceProviders()

	// then:

	require.Empty(t, response)

	mock.AssertCalled()

}

func TestLookupServicesListService_WithProviders(t *testing.T) {

	// given:

	expectations := testabilities.LookupServicesListProviderMockExpectations{

		MetadataList: map[string]*overlay.MetaData{

			"provider1": {

				Description: "Description 1",

				Icon: "https://example.com/icon.png",

				Version: "1.0.0",

				InfoUrl: "https://example.com/info",
			},

			"provider2": {

				Description: "Description 2",

				Icon: "https://example.com/icon2.png",

				Version: "2.0.0",

				InfoUrl: "https://example.com/info2",
			},
		},

		ListLookupServiceProvidersCall: true,
	}

	mock := testabilities.NewLookupServicesListProviderMock(t, expectations)

	service, err := app.NewLookupServicesListService(mock)

	require.NoError(t, err)

	// when:

	response := service.ListLookupServiceProviders()

	// then:

	expected := app.LookupServicesListResponse{

		"provider1": app.LookupMetadata{

			Name: "provider1",

			ShortDescription: "Description 1",

			IconURL: ptr.To("https://example.com/icon.png"),

			Version: ptr.To("1.0.0"),

			InformationURL: ptr.To("https://example.com/info"),
		},

		"provider2": app.LookupMetadata{

			Name: "provider2",

			ShortDescription: "Description 2",

			IconURL: ptr.To("https://example.com/icon2.png"),

			Version: ptr.To("2.0.0"),

			InformationURL: ptr.To("https://example.com/info2"),
		},
	}

	require.EqualValues(t, expected, response)

	mock.AssertCalled()

}

func TestLookupServicesListService_WithGenericSuccessProvider(t *testing.T) {

	// given:

	provider := &testabilities.LookupListProviderAlwaysSuccess{}

	service, err := app.NewLookupServicesListService(provider)

	require.NoError(t, err)

	// when:

	response := service.ListLookupServiceProviders()

	// then:

	expected := app.LookupServicesListResponse{

		"provider1": app.LookupMetadata{

			Name: "provider1",

			ShortDescription: "Description 1",

			IconURL: ptr.To("https://example.com/icon.png"),

			Version: ptr.To("1.0.0"),

			InformationURL: ptr.To("https://example.com/info"),
		},

		"provider2": app.LookupMetadata{

			Name: "provider2",

			ShortDescription: "Description 2",

			IconURL: ptr.To("https://example.com/icon2.png"),

			Version: ptr.To("2.0.0"),

			InformationURL: ptr.To("https://example.com/info2"),
		},
	}

	require.EqualValues(t, expected, response)

}

func TestLookupServicesListService_WithGenericEmptyProvider(t *testing.T) {

	// given:

	provider := &testabilities.LookupListProviderAlwaysEmpty{}

	service, err := app.NewLookupServicesListService(provider)

	require.NoError(t, err)

	// when:

	response := service.ListLookupServiceProviders()

	// then:

	require.Empty(t, response)

}

func TestNewLookupServicesListService_WithNilProvider(t *testing.T) {

	// given/when:

	service, err := app.NewLookupServicesListService(nil)

	// then:

	require.Error(t, err)

	var appError app.Error

	require.ErrorAs(t, err, &appError)

	require.Equal(t, app.ErrorTypeIncorrectInput, appError.ErrorType())

	require.Nil(t, service)

}
