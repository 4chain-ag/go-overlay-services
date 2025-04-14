package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

func TestGASP_RequestNode_ShouldReturnNode(t *testing.T) {
	// given
	ctx := context.Background()
	expectedNode := &core.GASPNode{}
	gasp := &core.GASP{
		Storage: fakeStorage{
			hydrateGASPNodeFunc: func(ctx context.Context, graphID, outpoint *overlay.Outpoint, metadata bool) (*core.GASPNode, error) {
				return expectedNode, nil
			},
		},
	}

	// when
	node, err := gasp.RequestNode(ctx, &overlay.Outpoint{}, &overlay.Outpoint{}, true)

	// then
	require.NoError(t, err)
	require.Equal(t, expectedNode, node)
}

func TestGASP_RequestNode_ShouldReturnError(t *testing.T) {
	// given
	ctx := context.Background()
	gasp := &core.GASP{
		Storage: fakeStorage{
			hydrateGASPNodeFunc: func(ctx context.Context, graphID, outpoint *overlay.Outpoint, metadata bool) (*core.GASPNode, error) {
				return nil, errors.New("hydrate-fail")
			},
		},
	}

	// when
	node, err := gasp.RequestNode(ctx, &overlay.Outpoint{}, &overlay.Outpoint{}, true)

	// then
	require.Error(t, err)
	require.Nil(t, node)
	require.EqualError(t, err, "hydrate-fail")
}
