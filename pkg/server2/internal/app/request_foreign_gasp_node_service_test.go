package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestRequestForeignGASPNodeService_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		graphIDStr    string
		txIDStr       string
		outputIndex   uint32
		topic         string
		expectations  testabilities.RequestForeignGASPNodeProviderMockExpectations
		expectedError app.Error
	}{
		"Request foreign GASP node service fails to handle the request - missing topic": {
			graphIDStr:    testabilities.DefaultValidGraphID,
			txIDStr:       testabilities.DefaultValidTxID,
			outputIndex:   testabilities.DefaultValidOutputIndex,
			topic:         testabilities.DefaultEmptyTopic,
			expectations:  testabilities.RequestForeignGASPNodeProviderMockExpectations{},
			expectedError: app.NewRequestForeignGASPNodeMissingTopicError(),
		},
		"Request foreign GASP node service fails to handle the request - invalid txID format": {
			graphIDStr:    testabilities.DefaultValidGraphID,
			txIDStr:       testabilities.DefaultInvalidTxID,
			outputIndex:   testabilities.DefaultValidOutputIndex,
			topic:         testabilities.DefaultValidTopic,
			expectations:  testabilities.RequestForeignGASPNodeProviderMockExpectations{},
			expectedError: app.NewRequestForeignGASPNodeInvalidTxIDError(),
		},
		"Request foreign GASP node service fails to handle the request - invalid graphID format": {
			graphIDStr:    testabilities.DefaultInvalidGraphID,
			txIDStr:       testabilities.DefaultValidTxID,
			outputIndex:   testabilities.DefaultValidOutputIndex,
			topic:         testabilities.DefaultValidTopic,
			expectations:  testabilities.RequestForeignGASPNodeProviderMockExpectations{},
			expectedError: app.NewRequestForeignGASPNodeInvalidGraphIDError(),
		},
		"Request foreign GASP node service fails to handle the request - provider failure": {
			graphIDStr:  testabilities.DefaultValidGraphID,
			txIDStr:     testabilities.DefaultValidTxID,
			outputIndex: testabilities.DefaultValidOutputIndex,
			topic:       testabilities.DefaultValidTopic,
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: true,
				Error:                      errors.New("internal request foreign GASP node service test error"),
			},
			expectedError: app.NewRequestForeignGASPNodeProviderError(errors.New("internal request foreign GASP node service test error")),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewRequestForeignGASPNodeProviderMock(t, tc.expectations)
			service := app.NewRequestForeignGASPNodeService(mock)

			// when:
			node, err := service.RequestForeignGASPNode(context.Background(), tc.graphIDStr, tc.txIDStr, tc.outputIndex, tc.topic)

			// then:
			var appErr app.Error
			require.ErrorAs(t, err, &appErr)
			require.Equal(t, tc.expectedError, appErr)
			require.Nil(t, node)

			if tc.expectations.ProvideForeignGASPNodeCall {
				mock.AssertCalled()
			}
		})
	}
}

func TestRequestForeignGASPNodeService_ValidCase(t *testing.T) {
	// given:
	mock := testabilities.NewRequestForeignGASPNodeProviderMock(t, testabilities.DefaultRequestForeignGASPNodeProviderMockExpectations)
	service := app.NewRequestForeignGASPNodeService(mock)

	// when:
	node, err := service.RequestForeignGASPNode(context.Background(), testabilities.DefaultValidGraphID, testabilities.DefaultValidTxID, testabilities.DefaultValidOutputIndex, testabilities.DefaultValidTopic)

	// then:
	require.NoError(t, err)
	require.Equal(t, testabilities.DefaultRequestForeignGASPNodeProviderMockExpectations.Node, node)
	mock.AssertCalled()
}
