package ports_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestTopicManagersListHandler_ValidCases(t *testing.T) {
	tests := map[string]struct {
		expectations       testabilities.TopicManagersListProviderMockExpectations
		expected           openapi.TopicManagersListResponse
		expectedStatusCode int
	}{
		"List topic manager service returns an empty topic mangers list.": {
			expectations: testabilities.TopicManagersListProviderMockExpectations{
				MetadataList:          testabilities.TopicManagerEmptyMetadata,
				ListTopicManagersCall: true,
			},
			expected:           ports.NewTopicManagersListSuccessResponse(testabilities.TopicManagerEmptyExpectedResponse),
			expectedStatusCode: fiber.StatusOK,
		},
		"List topic manager service returns default topic managers list.": {
			expectations: testabilities.TopicManagersListProviderMockExpectations{
				MetadataList:          testabilities.TopicManagerDefaultMetadata,
				ListTopicManagersCall: true,
			},
			expected:           ports.NewTopicManagersListSuccessResponse(testabilities.TopicManagerDefaultExpectedResponse),
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
