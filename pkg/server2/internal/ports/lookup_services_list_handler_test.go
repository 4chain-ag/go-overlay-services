package ports_test

import (
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

func TestLookupServicesListHandler_EmptyList(t *testing.T) {

	// given:

	expectations := testabilities.LookupServicesListProviderMockExpectations{

		MetadataList: map[string]*overlay.MetaData{},

		ListLookupServiceProvidersCall: true,
	}

	mockProvider := testabilities.NewLookupServicesListProviderMock(t, expectations)

	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupServicesListProvider(mockProvider))

	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	// when:

	var response app.LookupServicesListResponse

	res, err := fixture.Client().
		R().
		SetResult(&response).
		Get("/api/v1/listLookupServiceProviders")

	// then:

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode())

	require.Empty(t, response)

	stub.AssertProvidersState()

}

func TestLookupServicesListHandler_WithProviders(t *testing.T) {

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

	mockProvider := testabilities.NewLookupServicesListProviderMock(t, expectations)

	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupServicesListProvider(mockProvider))

	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	// when:

	var response app.LookupServicesListResponse

	res, err := fixture.Client().
		R().
		SetResult(&response).
		Get("/api/v1/listLookupServiceProviders")

	// then:

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode())

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

	require.Equal(t, expected, response)

	stub.AssertProvidersState()

}
