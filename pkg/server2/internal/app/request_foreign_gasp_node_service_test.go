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

func TestRequestForeignGASPNodeService_ProviderFailure(t *testing.T) {
	// given:
	expectations := testabilities.RequestForeignGASPNodeProviderMockExpectations{
		ProvideForeignGASPNodeCall: true,
		Error:                      errors.New("internal request foreign GASP node service test error"),
	}
	mock := testabilities.NewRequestForeignGASPNodeProviderMock(t, expectations)
	service, err := app.NewRequestForeignGASPNodeService(mock)
	require.NoError(t, err)
	graphID := &overlay.Outpoint{}
	outpoint := &overlay.Outpoint{}
	topic := "test-topic"

	// when:
	node, err := service.RequestForeignGASPNode(context.Background(), graphID, outpoint, topic)

	// then:
	var appErr app.Error
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, app.ErrorTypeProviderFailure, appErr.ErrorType())
	require.Nil(t, node)
	mock.AssertCalled()
}

func TestRequestForeignGASPNodeService_MissingTopic(t *testing.T) {
	// given:
	mock := testabilities.NewRequestForeignGASPNodeProviderMock(t, testabilities.RequestForeignGASPNodeProviderMockExpectations{})
	service, err := app.NewRequestForeignGASPNodeService(mock)
	require.NoError(t, err)
	graphID := &overlay.Outpoint{}
	outpoint := &overlay.Outpoint{}
	topic := ""

	// when:
	node, err := service.RequestForeignGASPNode(context.Background(), graphID, outpoint, topic)

	// then:
	var appErr app.Error
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, app.ErrorTypeIncorrectInput, appErr.ErrorType())
	require.Nil(t, node)
}

func TestRequestForeignGASPNodeService_MissingGraphID(t *testing.T) {
	// given:
	mock := testabilities.NewRequestForeignGASPNodeProviderMock(t, testabilities.RequestForeignGASPNodeProviderMockExpectations{})
	service, err := app.NewRequestForeignGASPNodeService(mock)
	require.NoError(t, err)
	outpoint := &overlay.Outpoint{}
	topic := "test-topic"

	// when:
	node, err := service.RequestForeignGASPNode(context.Background(), nil, outpoint, topic)

	// then:
	var appErr app.Error
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, app.ErrorTypeIncorrectInput, appErr.ErrorType())
	require.Nil(t, node)
}

func TestRequestForeignGASPNodeService_MissingOutpoint(t *testing.T) {
	// given:
	mock := testabilities.NewRequestForeignGASPNodeProviderMock(t, testabilities.RequestForeignGASPNodeProviderMockExpectations{})
	service, err := app.NewRequestForeignGASPNodeService(mock)
	require.NoError(t, err)
	graphID := &overlay.Outpoint{}
	topic := "test-topic"

	// when:
	node, err := service.RequestForeignGASPNode(context.Background(), graphID, nil, topic)

	// then:
	var appErr app.Error
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, app.ErrorTypeIncorrectInput, appErr.ErrorType())
	require.Nil(t, node)
}

func TestRequestForeignGASPNodeService_NilProvider(t *testing.T) {
	// given/when:
	service, err := app.NewRequestForeignGASPNodeService(nil)

	// then:
	require.Error(t, err)
	var appErr app.Error
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, app.ErrorTypeIncorrectInput, appErr.ErrorType())
	require.Nil(t, service)
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
	service, err := app.NewRequestForeignGASPNodeService(mock)
	require.NoError(t, err)
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
