package ports_test

import (
	"net/http"
	"testing"

	server2 "github.com/4chain-ag/go-overlay-services/pkg/server2/internal"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestTopicManagersListHandler_EmptyList(t *testing.T) {

	// given:

	expectations := testabilities.TopicManagersListProviderMockExpectations{

		MetadataList: testabilities.EmptyMetadata,

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

	require.Equal(t, testabilities.EmptyExpectedResponse, response)

	stub.AssertProvidersState()

}

func TestTopicManagersListHandler_WithDefaultManagers(t *testing.T) {

	// given:

	expectations := testabilities.TopicManagersListProviderMockExpectations{

		MetadataList: testabilities.DefaultMetadata,

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

	require.Equal(t, testabilities.DefaultExpectedResponse, response)

	stub.AssertProvidersState()

}
