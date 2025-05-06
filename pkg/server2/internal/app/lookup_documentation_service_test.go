package app_test

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

type mockLookupServicesListProvider struct {
	lookupServiceProviders map[string]*overlay.MetaData
}

func (m *mockLookupServicesListProvider) ListLookupServiceProviders() map[string]*overlay.MetaData {
	return m.lookupServiceProviders
}

func TestLookupServicesListService_GetList_EmptyList(t *testing.T) {
	// Given
	mockProvider := &mockLookupServicesListProvider{
		lookupServiceProviders: map[string]*overlay.MetaData{},
	}
	service := app.NewLookupServicesListService(mockProvider)

	// When
	result := service.GetList(context.Background())

	// Then
	require.Empty(t, result)
}

func TestLookupServicesListService_GetList_WithProviders(t *testing.T) {
	// Given
	mockProvider := &mockLookupServicesListProvider{
		lookupServiceProviders: map[string]*overlay.MetaData{
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
		},
	}
	service := app.NewLookupServicesListService(mockProvider)

	// When
	result := service.GetList(context.Background())

	// Then
	require.Len(t, result, 2)
	require.Contains(t, result, "provider1")
	provider1 := result["provider1"]
	require.Equal(t, "provider1", provider1.Name)
	require.Equal(t, "Description 1", provider1.ShortDescription)
	require.NotNil(t, provider1.IconURL)
	require.Equal(t, "https://example.com/icon.png", *provider1.IconURL)
	require.NotNil(t, provider1.Version)
	require.Equal(t, "1.0.0", *provider1.Version)
	require.NotNil(t, provider1.InformationURL)
	require.Equal(t, "https://example.com/info", *provider1.InformationURL)
	require.Contains(t, result, "provider2")
	provider2 := result["provider2"]
	require.Equal(t, "provider2", provider2.Name)
	require.Equal(t, "Description 2", provider2.ShortDescription)
	require.NotNil(t, provider2.IconURL)
	require.Equal(t, "https://example.com/icon2.png", *provider2.IconURL)
	require.NotNil(t, provider2.Version)
	require.Equal(t, "2.0.0", *provider2.Version)
	require.NotNil(t, provider2.InformationURL)
	require.Equal(t, "https://example.com/info2", *provider2.InformationURL)
}

func TestLookupServicesListService_GetList_WithNilMetadata(t *testing.T) {
	// Given
	mockProvider := &mockLookupServicesListProvider{
		lookupServiceProviders: map[string]*overlay.MetaData{
			"provider1": nil,
		},
	}
	service := app.NewLookupServicesListService(mockProvider)

	// When
	result := service.GetList(context.Background())

	// Then
	require.Len(t, result, 1)
	require.Contains(t, result, "provider1")
	provider1 := result["provider1"]
	require.Equal(t, "provider1", provider1.Name)
	require.Equal(t, "No description available", provider1.ShortDescription)
	require.Nil(t, provider1.IconURL)
	require.Nil(t, provider1.Version)
	require.Nil(t, provider1.InformationURL)
}

func TestNewLookupServicesListService_WithNilProvider(t *testing.T) {
	// Given, When, Then
	require.Panics(t, func() {
		app.NewLookupServicesListService(nil)
	})
}
