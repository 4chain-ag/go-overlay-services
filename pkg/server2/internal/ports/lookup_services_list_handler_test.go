package ports_test

import (
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestLookupServicesListHandler_GetList_ShouldReturnEmptyList(t *testing.T) {
	// Given
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithEmptyLookupServicesList())
	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))

	// When
	var actualResponse openapi.LookupServicesListResponse

	res, _ := fixture.Client().
		R().
		SetResult(&actualResponse).
		Get("/api/v1/listLookupServiceProviders")

	// Then
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.Empty(t, actualResponse)
}

func TestLookupServicesListHandler_GetList_ShouldReturnServicesList(t *testing.T) {
	// Given
	metadata := map[string]*testabilities.LookupMetadataMock{
		"provider1": {
			Description: "Description 1",
			Icon:        "https://example.com/icon.png",
			Version:     "1.0.0",
			InfoUrl:     "https://example.com/info",
		},
		"provider2": {
			Description: "Description 2",
			Icon:        "https://example.com/icon2.png",
			Version:     "2.0.0",
			InfoUrl:     "https://example.com/info2",
		},
	}

	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupServicesList(metadata))
	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))

	// When
	var actualResponse openapi.LookupServicesListResponse

	res, _ := fixture.Client().
		R().
		SetResult(&actualResponse).
		Get("/api/v1/listLookupServiceProviders")

	// Then
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.Len(t, actualResponse, 2)

	require.Contains(t, actualResponse, "provider1")
	provider1 := actualResponse["provider1"]
	require.Equal(t, "provider1", provider1.Name)
	require.Equal(t, "Description 1", provider1.ShortDescription)
	require.NotNil(t, provider1.IconURL)
	require.Equal(t, "https://example.com/icon.png", *provider1.IconURL)
	require.NotNil(t, provider1.Version)
	require.Equal(t, "1.0.0", *provider1.Version)
	require.NotNil(t, provider1.InformationURL)
	require.Equal(t, "https://example.com/info", *provider1.InformationURL)

	require.Contains(t, actualResponse, "provider2")
	provider2 := actualResponse["provider2"]
	require.Equal(t, "provider2", provider2.Name)
	require.Equal(t, "Description 2", provider2.ShortDescription)
	require.NotNil(t, provider2.IconURL)
	require.Equal(t, "https://example.com/icon2.png", *provider2.IconURL)
	require.NotNil(t, provider2.Version)
	require.Equal(t, "2.0.0", *provider2.Version)
	require.NotNil(t, provider2.InformationURL)
	require.Equal(t, "https://example.com/info2", *provider2.InformationURL)
}
