package app_test

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

type mockTopicManagersListProvider struct {
	topicManagers map[string]*overlay.MetaData
}

func (m *mockTopicManagersListProvider) ListTopicManagers() map[string]*overlay.MetaData {
	return m.topicManagers
}

func TestTopicManagersListService_GetList_EmptyList(t *testing.T) {
	// Given
	mockProvider := &mockTopicManagersListProvider{
		topicManagers: map[string]*overlay.MetaData{},
	}
	service := app.NewTopicManagersListService(mockProvider)

	// When
	result := service.GetList(context.Background())

	// Then
	require.Empty(t, result)
}

func TestTopicManagersListService_GetList_WithProviders(t *testing.T) {
	// Given
	mockProvider := &mockTopicManagersListProvider{
		topicManagers: map[string]*overlay.MetaData{
			"manager1": {
				Description: "Description 1",
				Icon:        "https://example.com/icon.png",
				Version:     "1.0.0",
				InfoUrl:     "https://example.com/info",
			},
			"manager2": {
				Description: "Description 2",
				Icon:        "https://example.com/icon2.png",
				Version:     "2.0.0",
				InfoUrl:     "https://example.com/info2",
			},
		},
	}
	service := app.NewTopicManagersListService(mockProvider)

	// When
	result := service.GetList(context.Background())

	// Then
	require.Len(t, result, 2)
	require.Contains(t, result, "manager1")
	manager1 := result["manager1"]
	require.Equal(t, "manager1", manager1.Name)
	require.Equal(t, "Description 1", manager1.ShortDescription)
	require.NotNil(t, manager1.IconURL)
	require.Equal(t, "https://example.com/icon.png", *manager1.IconURL)
	require.NotNil(t, manager1.Version)
	require.Equal(t, "1.0.0", *manager1.Version)
	require.NotNil(t, manager1.InformationURL)
	require.Equal(t, "https://example.com/info", *manager1.InformationURL)
	require.Contains(t, result, "manager2")
	manager2 := result["manager2"]
	require.Equal(t, "manager2", manager2.Name)
	require.Equal(t, "Description 2", manager2.ShortDescription)
	require.NotNil(t, manager2.IconURL)
	require.Equal(t, "https://example.com/icon2.png", *manager2.IconURL)
	require.NotNil(t, manager2.Version)
	require.Equal(t, "2.0.0", *manager2.Version)
	require.NotNil(t, manager2.InformationURL)
	require.Equal(t, "https://example.com/info2", *manager2.InformationURL)
}

func TestTopicManagersListService_GetList_WithNilMetadata(t *testing.T) {
	// Given
	mockProvider := &mockTopicManagersListProvider{
		topicManagers: map[string]*overlay.MetaData{
			"manager1": nil,
		},
	}
	service := app.NewTopicManagersListService(mockProvider)

	// When
	result := service.GetList(context.Background())

	// Then
	require.Len(t, result, 1)
	require.Contains(t, result, "manager1")
	manager1 := result["manager1"]
	require.Equal(t, "manager1", manager1.Name)
	require.Equal(t, "No description available", manager1.ShortDescription)
	require.Nil(t, manager1.IconURL)
	require.Nil(t, manager1.Version)
	require.Nil(t, manager1.InformationURL)
}

func TestNewTopicManagersListService_WithNilProvider(t *testing.T) {
	// Given, When, Then
	require.Panics(t, func() {
		app.NewTopicManagersListService(nil)
	})
}
