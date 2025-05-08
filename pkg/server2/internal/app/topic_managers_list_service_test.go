package app_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

func TestTopicManagersListService_EmptyList(t *testing.T) {
	// given:
	service, err := app.NewTopicManagersListService(&testabilities.TopicManagerListProviderAlwaysEmpty{})
	require.NoError(t, err)

	// when:
	response := service.ListTopicManagers()

	// then:
	require.Empty(t, response)
}

func TestTopicManagersListService_WithManagers(t *testing.T) {
	// given:
	service, err := app.NewTopicManagersListService(&testabilities.TopicManagerListProviderAlwaysSuccess{})
	require.NoError(t, err)

	// when:
	response := service.ListTopicManagers()

	// then:
	expected := app.TopicManagersListResponse{
		"manager1": app.TopicManagerMetadata{
			Name:             "manager1",
			ShortDescription: "Description 1",
			IconURL:          ptr.To("https://example.com/icon.png"),
			Version:          ptr.To("1.0.0"),
			InformationURL:   ptr.To("https://example.com/info"),
		},
		"manager2": app.TopicManagerMetadata{
			Name:             "manager2",
			ShortDescription: "Description 2",
			IconURL:          ptr.To("https://example.com/icon2.png"),
			Version:          ptr.To("2.0.0"),
			InformationURL:   ptr.To("https://example.com/info2"),
		},
	}
	require.EqualValues(t, expected, response)
}

func TestNewTopicManagersListService_WithNilProvider(t *testing.T) {
	// given/when:
	service, err := app.NewTopicManagersListService(nil)

	// then:
	require.Error(t, err)
	var appError app.Error
	require.ErrorAs(t, err, &appError)
	require.Equal(t, app.ErrorTypeIncorrectInput, appError.ErrorType())
	require.Nil(t, service)
}
