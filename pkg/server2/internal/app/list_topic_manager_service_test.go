package app_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestTopicManagersListService_ValidCases(t *testing.T) {
	tests := map[string]struct {
		expectations testabilities.TopicManagersListProviderMockExpectations
		expected     app.TopicManagers
	}{
		"List topic manager service success - empty list": {
			expectations: testabilities.TopicManagersListProviderMockExpectations{
				MetadataList:          testabilities.EmptyMetadata,
				ListTopicManagersCall: true,
			},
			expected: testabilities.EmptyExpectedResponse,
		},
		"List topic manager service success - default list": {
			expectations: testabilities.TopicManagersListProviderMockExpectations{
				MetadataList:          testabilities.DefaultMetadata,
				ListTopicManagersCall: true,
			},
			expected: testabilities.DefaultExpectedResponse,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewTopicManagersListProviderMock(t, tc.expectations)
			service := app.NewTopicManagersListService(mock)

			// when:
			response := service.ListTopicManagers()

			// then:
			require.Equal(t, tc.expected, response)
			mock.AssertCalled()
		})
	}
}
