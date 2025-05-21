package ports_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
	"github.com/gofiber/fiber/v2"

)


func TestTopicManagersListHandler_ValidCases(t *testing.T) {
	tests := map[string]struct {
		expectations testabilities.TopicManagersListProviderMockExpectations
		expected     openapi.TopicManagersListResponse
		expectedStatusCode int
	}{
		"empty list": {
			expectations: testabilities.TopicManagersListProviderMockExpectations{
				MetadataList:          testabilities.EmptyMetadata,
				ListTopicManagersCall: true,
			},
			expected: ports.NewTopicManagersListSuccessResponse(testabilities.EmptyExpectedResponse),
			expectedStatusCode: fiber.StatusOK,
		},
		"default list": {
			expectations: testabilities.TopicManagersListProviderMockExpectations{
				MetadataList:          testabilities.DefaultMetadata,
				ListTopicManagersCall: true,
			},
			expected: ports.NewTopicManagersListSuccessResponse(testabilities.DefaultExpectedResponse),
			expectedStatusCode: fiber.StatusOK,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithTopicManagersListProvider(testabilities.NewTopicManagersListProviderMock(t, tc.expectations))) 
			fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))


			// when:
			var actualResponse openapi.TopicManagersListResponse

			res, _ := fixture.Client().
				R().
				SetResult(&actualResponse).
				Get("/api/v1/listTopicManagers")

			// then:
			require.Equal(t, tc.expectedStatusCode, res.StatusCode())
			require.Equal(t, tc.expected, actualResponse)
			stub.AssertProvidersState()
		})
	}
}

// func TestTopicManagersListHandler_EmptyList(t *testing.T) {
// 	// given:
// 	expectations := testabilities.TopicManagersListProviderMockExpectations{
// 		MetadataList:          testabilities.EmptyMetadata,
// 		ListTopicManagersCall: true,
// 	}
// 	mock := testabilities.NewTopicManagersListProviderMock(t, expectations)
// 	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithTopicManagersListProvider(mock))
// 	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

// 	// when:
// 	var actualResponse openapi.TopicManagersListResponse
// 	res, err := fixture.Client().
// 		R().
// 		SetResult(&actualResponse).
// 		Get("/api/v1/listTopicManagers")

// 	// then:
// 	require.NoError(t, err)
// 	require.Equal(t, fiber.StatusOK, res.StatusCode())
// 	require.Equal(t, testabilities.EmptyExpectedResponse, actualResponse)
// 	stub.AssertProvidersState()
// }

// func TestTopicManagersListHandler_ValidCase(t *testing.T) {
// 	// given:
// 	expectations := testabilities.TopicManagersListProviderMockExpectations{
// 		MetadataList:          testabilities.DefaultMetadata,
// 		ListTopicManagersCall: true,
// 	}
// 	mock := testabilities.NewTopicManagersListProviderMock(t, expectations)
// 	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithTopicManagersListProvider(mock))
// 	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

// 	// when:
// 	var actualResponse openapi.TopicManagersListResponse
// 	res, err := fixture.Client().
// 		R().
// 		SetResult(&actualResponse).
// 		Get("/api/v1/listTopicManagers")

// 	// then:
// 	require.NoError(t, err)
// 	require.Equal(t, fiber.StatusOK, res.StatusCode())
// 	require.Equal(t, , actualResponse)
// 	stub.AssertProvidersState()
// }
