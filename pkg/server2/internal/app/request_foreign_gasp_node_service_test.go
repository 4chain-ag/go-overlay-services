package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

func TestRequestForeignGASPNodeService_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		expectations      testabilities.RequestForeignGASPNodeProviderMockExpectations
		expectedErrorType app.ErrorType
	}{
		"Request foreign GASP node service fails - internal error": {
			expectations: testabilities.RequestForeignGASPNodeProviderMockExpectations{
				ProvideForeignGASPNodeCall: true,
				Error:                      errors.New("internal request foreign GASP node service test error"),
			},
			expectedErrorType: app.ErrorTypeProviderFailure,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewRequestForeignGASPNodeProviderMock(t, tc.expectations)
			service := app.NewRequestForeignGASPNodeService(mock)
			graphID := &overlay.Outpoint{}
			outpoint := &overlay.Outpoint{}
			topic := "test-topic"

			// when:
			node, err := service.RequestForeignGASPNode(context.Background(), graphID, outpoint, topic)

			// then:
			var appErr app.Error
			require.ErrorAs(t, err, &appErr)
			require.Equal(t, tc.expectedErrorType, appErr.ErrorType())
			require.Nil(t, node)
			mock.AssertCalled()
		})
	}
}

func TestRequestForeignGASPNodeService_ValidCase(t *testing.T) {
	// given:
	expectedNode := &core.GASPNode{}
	expectations := testabilities.RequestForeignGASPNodeProviderMockExpectations{
		ProvideForeignGASPNodeCall: true,
		Node:                       expectedNode,
		Error:                      nil,
	}
	mock := testabilities.NewRequestForeignGASPNodeProviderMock(t, expectations)
	service := app.NewRequestForeignGASPNodeService(mock)
	graphID := &overlay.Outpoint{}
	outpoint := &overlay.Outpoint{}
	topic := "test-topic"

	// when:
	node, err := service.RequestForeignGASPNode(context.Background(), graphID, outpoint, topic)

	// then:
	require.NoError(t, err)
	require.Equal(t, expectedNode, node)
	mock.AssertCalled()
}
