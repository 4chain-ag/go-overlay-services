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

func TestTopicManagersListHandler_EmptyList(t *testing.T) {

	// given:

	expectations := testabilities.TopicManagersListProviderMockExpectations{

		MetadataList: map[string]*overlay.MetaData{},

		ListTopicManagersCall: true,
	}

	mockProvider := testabilities.NewTopicManagersListProviderMock(t, expectations)

	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithTopicManagersListProvider(mockProvider))

	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	// when:

	var response app.TopicManagersListResponse

	res, err := fixture.Client().
		R().
		SetResult(&response).
		Get("/api/v1/listTopicManagers")

	// then:

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode())

	require.Empty(t, response)

	stub.AssertProvidersState()

}

func TestTopicManagersListHandler_WithManagers(t *testing.T) {

	// given:

	expectations := testabilities.TopicManagersListProviderMockExpectations{

		MetadataList: map[string]*overlay.MetaData{

			"manager1": {

				Description: "Description 1",

				Icon: "https://example.com/icon.png",

				Version: "1.0.0",

				InfoUrl: "https://example.com/info",
			},

			"manager2": {

				Description: "Description 2",

				Icon: "https://example.com/icon2.png",

				Version: "2.0.0",

				InfoUrl: "https://example.com/info2",
			},
		},

		ListTopicManagersCall: true,
	}

	mockProvider := testabilities.NewTopicManagersListProviderMock(t, expectations)

	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithTopicManagersListProvider(mockProvider))

	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	// when:

	var response app.TopicManagersListResponse

	res, err := fixture.Client().
		R().
		SetResult(&response).
		Get("/api/v1/listTopicManagers")

	// then:

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode())

	expected := app.TopicManagersListResponse{

		"manager1": app.TopicManagerMetadata{

			Name: "manager1",

			ShortDescription: "Description 1",

			IconURL: ptr.To("https://example.com/icon.png"),

			Version: ptr.To("1.0.0"),

			InformationURL: ptr.To("https://example.com/info"),
		},

		"manager2": app.TopicManagerMetadata{

			Name: "manager2",

			ShortDescription: "Description 2",

			IconURL: ptr.To("https://example.com/icon2.png"),

			Version: ptr.To("2.0.0"),

			InformationURL: ptr.To("https://example.com/info2"),
		},
	}

	require.Equal(t, expected, response)

	stub.AssertProvidersState()

}
