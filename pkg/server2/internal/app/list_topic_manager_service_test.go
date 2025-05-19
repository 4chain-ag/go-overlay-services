package app_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestTopicManagersListService_EmptyList(t *testing.T) {
	// given:
	expectations := testabilities.TopicManagersListProviderMockExpectations{
		MetadataList:          testabilities.EmptyMetadata,
		ListTopicManagersCall: true,
	}
	mock := testabilities.NewTopicManagersListProviderMock(t, expectations)
	service, err := app.NewTopicManagersListService(mock)
	require.NoError(t, err)

	// when:
	response := service.ListTopicManagers()

	// then:
	require.Equal(t, testabilities.EmptyExpectedResponse, response)
	mock.AssertCalled()
}

func TestTopicManagersListService_WithProviders(t *testing.T) {
	// given:
	expectations := testabilities.TopicManagersListProviderMockExpectations{
		MetadataList:          testabilities.DefaultMetadata,
		ListTopicManagersCall: true,
	}
	mock := testabilities.NewTopicManagersListProviderMock(t, expectations)
	service, err := app.NewTopicManagersListService(mock)
	require.NoError(t, err)

	// when:
	response := service.ListTopicManagers()

	// then:
	require.Equal(t, testabilities.DefaultExpectedResponse, response)
	mock.AssertCalled()
}

func TestTopicManagersListService_WithGenericSuccessProvider(t *testing.T) {
	// given:
	provider := &testabilities.TopicManagersListProviderAlwaysSuccess{}
	service, err := app.NewTopicManagersListService(provider)
	require.NoError(t, err)

	// when:
	response := service.ListTopicManagers()

	// then:
	require.Equal(t, testabilities.DefaultExpectedResponse, response)
}

func TestTopicManagersListService_WithGenericEmptyProvider(t *testing.T) {
	// given:
	provider := &testabilities.TopicManagersListProviderAlwaysEmpty{}
	service, err := app.NewTopicManagersListService(provider)
	require.NoError(t, err)

	// when:
	response := service.ListTopicManagers()

	// then:
	require.Equal(t, testabilities.EmptyExpectedResponse, response)
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
