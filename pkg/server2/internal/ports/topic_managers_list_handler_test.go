package ports_test

import (
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestTopicManagersListHandler_GetList_ShouldReturnEmptyList(t *testing.T) {

	// Given

	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithEmptyTopicManagersList())

	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))

	// When

	var actualResponse openapi.TopicManagersListResponse

	res, _ := fixture.Client().
		R().
		SetResult(&actualResponse).
		Get("/api/v1/listTopicManagers")

	// Then

	require.Equal(t, http.StatusOK, res.StatusCode())

	require.Empty(t, actualResponse)

}

func TestTopicManagersListHandler_GetList_ShouldReturnManagersList(t *testing.T) {

	// Given

	metadata := map[string]*testabilities.TopicManagerMetadataMock{

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
	}

	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithTopicManagersList(metadata))

	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))

	// When

	var actualResponse openapi.TopicManagersListResponse

	res, _ := fixture.Client().
		R().
		SetResult(&actualResponse).
		Get("/api/v1/listTopicManagers")

	// Then

	require.Equal(t, http.StatusOK, res.StatusCode())

	require.Len(t, actualResponse, 2)

	require.Contains(t, actualResponse, "manager1")

	manager1 := actualResponse["manager1"]

	require.Equal(t, "manager1", manager1.Name)

	require.Equal(t, "Description 1", manager1.ShortDescription)

	require.NotNil(t, manager1.IconURL)

	require.Equal(t, "https://example.com/icon.png", *manager1.IconURL)

	require.NotNil(t, manager1.Version)

	require.Equal(t, "1.0.0", *manager1.Version)

	require.NotNil(t, manager1.InformationURL)

	require.Equal(t, "https://example.com/info", *manager1.InformationURL)

	require.Contains(t, actualResponse, "manager2")

	manager2 := actualResponse["manager2"]

	require.Equal(t, "manager2", manager2.Name)

	require.Equal(t, "Description 2", manager2.ShortDescription)

	require.NotNil(t, manager2.IconURL)

	require.Equal(t, "https://example.com/icon2.png", *manager2.IconURL)

	require.NotNil(t, manager2.Version)

	require.Equal(t, "2.0.0", *manager2.Version)

	require.NotNil(t, manager2.InformationURL)

	require.Equal(t, "https://example.com/info2", *manager2.InformationURL)

}
