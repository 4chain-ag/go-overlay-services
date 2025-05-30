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
		dto           app.RequestForeignGASPNodeDTO
		expectations  testabilities.RequestForeignGASPNodeProviderMockExpectations
		expectedError app.Error
	}{
		"Request foreign GASP node service fails to handle the request with missing topic": {
			dto: app.RequestForeignGASPNodeDTO{
				GraphID:     testabilities.DefaultValidGraphID,
				TxID:        testabilities.DefaultValidTxID,
				OutputIndex: testabilities.DefaultValidOutputIndex,
				Topic:       "",
			},
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: false,
			},
			expectedError: app.NewRequestForeignGASPNodeMissingTopicError(),
		},
		"Request foreign GASP node service fails to handle the request with invalid txID format": {
			dto: app.RequestForeignGASPNodeDTO{
				GraphID:     testabilities.DefaultValidGraphID,
				TxID:        "invalid-txid",
				OutputIndex: testabilities.DefaultValidOutputIndex,
				Topic:       testabilities.DefaultValidTopic,
			},
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: false,
			},
			expectedError: app.NewRequestForeignGASPNodeInvalidTxIDError(),
		},
		"Request foreign GASP node service fails to handle the request with invalid graphID format": {
			dto: app.RequestForeignGASPNodeDTO{
				GraphID:     "invalid-graphid",
				TxID:        testabilities.DefaultValidTxID,
				OutputIndex: testabilities.DefaultValidOutputIndex,
				Topic:       testabilities.DefaultValidTopic,
			},
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: false,
			},
			expectedError: app.NewRequestForeignGASPNodeInvalidGraphIDError(),
		},
		"Request foreign GASP node service fails to handle the request with provider failure": {
			dto: app.RequestForeignGASPNodeDTO{
				GraphID:     testabilities.DefaultValidGraphID,
				TxID:        testabilities.DefaultValidTxID,
				OutputIndex: testabilities.DefaultValidOutputIndex,
				Topic:       testabilities.DefaultValidTopic,
			},
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
			node, err := service.RequestForeignGASPNode(context.Background(), tc.dto)

			// then:
			var appErr app.Error
			require.ErrorAs(t, err, &appErr)
			require.Equal(t, tc.expectedError, appErr)
			require.Nil(t, node)
			mock.AssertCalled()
		})
	}
}

func TestRequestForeignGASPNodeService_ValidCase(t *testing.T) {
	// given:
	mock := testabilities.NewRequestForeignGASPNodeProviderMock(t, testabilities.DefaultRequestForeignGASPNodeProviderMockExpectations)
	service := app.NewRequestForeignGASPNodeService(mock)
	dto := app.RequestForeignGASPNodeDTO{
		GraphID:     testabilities.DefaultValidGraphID,
		TxID:        testabilities.DefaultValidTxID,
		OutputIndex: testabilities.DefaultValidOutputIndex,
		Topic:       testabilities.DefaultValidTopic,
	}

	// when:
	node, err := service.RequestForeignGASPNode(context.Background(), dto)

	// then:
	require.NoError(t, err)
	require.Equal(t, testabilities.DefaultRequestForeignGASPNodeProviderMockExpectations.Node, node)
	mock.AssertCalled()
}
