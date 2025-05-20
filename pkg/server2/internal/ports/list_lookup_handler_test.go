package ports_test

import (
	"net/http"
	"testing"

	server2 "github.com/4chain-ag/go-overlay-services/pkg/server2/internal"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestLookupListHandler_EmptyList(t *testing.T) {

	// given:

	expectations := testabilities.LookupListProviderMockExpectations{

		MetadataList: testabilities.LookupListEmptyMetadata,

		ListLookupServiceProvidersCall: true,
	}

	mockProvider := testabilities.NewLookupListProviderMock(t, expectations)

	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupListProvider(mockProvider))

	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	// when:

	var response app.LookupListResponse

	res, err := fixture.Client().
		R().
		SetResult(&response).
		Get("/api/v1/listLookupServiceProviders")

	// then:

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode())

	require.Equal(t, testabilities.EmptyLookupListExpectedResponse, response)

	stub.AssertProvidersState()

}

func TestLookupListHandler_WithDefaultProviders(t *testing.T) {

	// given:

	expectations := testabilities.LookupListProviderMockExpectations{

		MetadataList: testabilities.LookupDefaultMetadata,

		ListLookupServiceProvidersCall: true,
	}

	mockProvider := testabilities.NewLookupListProviderMock(t, expectations)

	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupListProvider(mockProvider))

	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	// when:

	var response app.LookupListResponse

	res, err := fixture.Client().
		R().
		SetResult(&response).
		Get("/api/v1/listLookupServiceProviders")

	// then:

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode())

	require.Equal(t, testabilities.DefaultLookupListExpectedResponse, response)

	stub.AssertProvidersState()

}
