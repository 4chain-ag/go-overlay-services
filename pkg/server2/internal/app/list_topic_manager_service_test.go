package app_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

func TestTopicManagerListService_EmptyList(t *testing.T) {

	// given:

	expectations := testabilities.TopicManagerListProviderMockExpectations{

		MetadataList: map[string]*overlay.MetaData{},

		ListTopicManagersCall: true,
	}

	mock := testabilities.NewTopicManagerListProviderMock(t, expectations)

	service, err := app.NewTopicManagerListService(mock)

	require.NoError(t, err)

	// when:

	response := service.ListTopicManagers()

	// then:

	require.Empty(t, response)

	mock.AssertCalled()

}

func TestTopicManagerListService_WithProviders(t *testing.T) {

	// given:

	expectations := testabilities.TopicManagerListProviderMockExpectations{

		MetadataList: map[string]*overlay.MetaData{

			"topic_manager1": {

				Description: "Description 1",

				Icon: "https://example.com/icon.png",

				Version: "1.0.0",

				InfoUrl: "https://example.com/info",
			},

			"topic_manager2": {

				Description: "Description 2",

				Icon: "https://example.com/icon2.png",

				Version: "2.0.0",

				InfoUrl: "https://example.com/info2",
			},
		},

		ListTopicManagersCall: true,
	}

	mock := testabilities.NewTopicManagerListProviderMock(t, expectations)

	service, err := app.NewTopicManagerListService(mock)

	require.NoError(t, err)

	// when:

	response := service.ListTopicManagers()

	// then:

	expected := app.TopicManagerListResponse{

		"topic_manager1": app.TopicManagerMetadata{

			Name: "topic_manager1",

			ShortDescription: "Description 1",

			IconURL: ptr.To("https://example.com/icon.png"),

			Version: ptr.To("1.0.0"),

			InformationURL: ptr.To("https://example.com/info"),
		},

		"topic_manager2": app.TopicManagerMetadata{

			Name: "topic_manager2",

			ShortDescription: "Description 2",

			IconURL: ptr.To("https://example.com/icon2.png"),

			Version: ptr.To("2.0.0"),

			InformationURL: ptr.To("https://example.com/info2"),
		},
	}

	require.EqualValues(t, expected, response)

	mock.AssertCalled()

}

func TestTopicManagerListService_WithGenericSuccessProvider(t *testing.T) {

	// given:

	provider := &testabilities.TopicManagerListProviderAlwaysSuccess{}

	service, err := app.NewTopicManagerListService(provider)

	require.NoError(t, err)

	// when:

	response := service.ListTopicManagers()

	// then:

	expected := app.TopicManagerListResponse{

		"topic_manager1": app.TopicManagerMetadata{

			Name: "topic_manager1",

			ShortDescription: "Description 1",

			IconURL: ptr.To("https://example.com/icon.png"),

			Version: ptr.To("1.0.0"),

			InformationURL: ptr.To("https://example.com/info"),
		},

		"topic_manager2": app.TopicManagerMetadata{

			Name: "topic_manager2",

			ShortDescription: "Description 2",

			IconURL: ptr.To("https://example.com/icon2.png"),

			Version: ptr.To("2.0.0"),

			InformationURL: ptr.To("https://example.com/info2"),
		},
	}

	require.EqualValues(t, expected, response)

}

func TestTopicManagerListService_WithGenericEmptyProvider(t *testing.T) {

	// given:

	provider := &testabilities.TopicManagerListProviderAlwaysEmpty{}

	service, err := app.NewTopicManagerListService(provider)

	require.NoError(t, err)

	// when:

	response := service.ListTopicManagers()

	// then:

	require.Empty(t, response)

}

func TestNewTopicManagerListService_WithNilProvider(t *testing.T) {

	// given/when:

	service, err := app.NewTopicManagerListService(nil)

	// then:

	require.Error(t, err)

	var appError app.Error

	require.ErrorAs(t, err, &appError)

	require.Equal(t, app.ErrorTypeIncorrectInput, appError.ErrorType())

	require.Nil(t, service)

}
