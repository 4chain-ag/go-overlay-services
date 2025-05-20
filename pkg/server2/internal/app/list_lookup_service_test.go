package app_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestLookupListService_EmptyList(t *testing.T) {

	// given:

	expectations := testabilities.LookupListProviderMockExpectations{

		MetadataList: testabilities.LookupListEmptyMetadata,

		ListLookupServiceProvidersCall: true,
	}

	mock := testabilities.NewLookupListProviderMock(t, expectations)

	service, err := app.NewLookupListService(mock)

	require.NoError(t, err)

	// when:

	response := service.ListLookup()

	// then:

	require.Equal(t, testabilities.EmptyLookupListExpectedResponse, response)

	mock.AssertCalled()

}

func TestLookupListService_WithProviders(t *testing.T) {

	// given:

	expectations := testabilities.LookupListProviderMockExpectations{

		MetadataList: testabilities.LookupDefaultMetadata,

		ListLookupServiceProvidersCall: true,
	}

	mock := testabilities.NewLookupListProviderMock(t, expectations)

	service, err := app.NewLookupListService(mock)

	require.NoError(t, err)

	// when:

	response := service.ListLookup()

	// then:

	require.Equal(t, testabilities.DefaultLookupListExpectedResponse, response)

	mock.AssertCalled()

}

func TestLookupListService_WithGenericSuccessProvider(t *testing.T) {

	// given:

	provider := &testabilities.LookupListProviderAlwaysSuccess{}

	service, err := app.NewLookupListService(provider)

	require.NoError(t, err)

	// when:

	response := service.ListLookup()

	// then:

	require.Equal(t, testabilities.DefaultLookupListExpectedResponse, response)

}

func TestLookupListService_WithGenericEmptyProvider(t *testing.T) {

	// given:

	provider := &testabilities.LookupListProviderAlwaysEmpty{}

	service, err := app.NewLookupListService(provider)

	require.NoError(t, err)

	// when:

	response := service.ListLookup()

	// then:

	require.Equal(t, testabilities.EmptyLookupListExpectedResponse, response)

}

func TestNewLookupListService_WithNilProvider(t *testing.T) {

	// given/when:

	service, err := app.NewLookupListService(nil)

	// then:

	require.Error(t, err)

	var appError app.Error

	require.ErrorAs(t, err, &appError)

	require.Equal(t, app.ErrorTypeIncorrectInput, appError.ErrorType())

	require.Nil(t, service)

}
