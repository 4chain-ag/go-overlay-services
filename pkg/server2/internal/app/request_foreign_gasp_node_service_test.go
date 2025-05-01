package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/assert"
)

type mockRequestForeignGASPNodeProvider struct {
	shouldFail bool
}

func (m *mockRequestForeignGASPNodeProvider) ProvideForeignGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error) {

	if m.shouldFail {

		return nil, errors.New("mock foreign GASP node error")

	}

	return &core.GASPNode{}, nil

}

func TestRequestForeignGASPNodeService_RequestForeignGASPNode_Success(t *testing.T) {

	// Given

	mockProvider := &mockRequestForeignGASPNodeProvider{shouldFail: false}

	service := app.NewRequestForeignGASPNodeService(mockProvider)

	// When

	node, err := service.RequestForeignGASPNode(context.Background(), &overlay.Outpoint{}, &overlay.Outpoint{}, "test-topic")

	// Then

	assert.NoError(t, err)

	assert.NotNil(t, node)

}

func TestRequestForeignGASPNodeService_RequestForeignGASPNode_Error(t *testing.T) {

	// Given

	mockProvider := &mockRequestForeignGASPNodeProvider{shouldFail: true}

	service := app.NewRequestForeignGASPNodeService(mockProvider)

	// When

	node, err := service.RequestForeignGASPNode(context.Background(), &overlay.Outpoint{}, &overlay.Outpoint{}, "test-topic")

	// Then

	assert.Error(t, err)

	assert.ErrorIs(t, err, app.ErrRequestForeignGASPNodeProvider)

	assert.Nil(t, node)

}

func TestNewRequestForeignGASPNodeService_NilProvider(t *testing.T) {

	// Given/When/Then

	assert.Panics(t, func() {

		app.NewRequestForeignGASPNodeService(nil)

	})

}
