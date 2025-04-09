package core__test

import (
	"context"
	"errors"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

type fakeStorage struct {
	hydrateGASPNodeFunc func(ctx context.Context, graphID, outpoint *overlay.Outpoint, metadata bool) (*core.GASPNode, error)
}

func (f fakeStorage) HydrateGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, metadata bool) (*core.GASPNode, error) {
	if f.hydrateGASPNodeFunc != nil {
		return f.hydrateGASPNodeFunc(ctx, graphID, outpoint, metadata)
	}
	return nil, errors.New("not implemented")
}

func (f fakeStorage) FindKnownUTXOs(ctx context.Context, since uint32) ([]*overlay.Outpoint, error) {
	return nil, nil
}
func (f fakeStorage) FindNeededInputs(ctx context.Context, tx *core.GASPNode) (*core.GASPNodeResponse, error) {
	return nil, nil
}
func (f fakeStorage) AppendToGraph(ctx context.Context, tx *core.GASPNode, spentBy *chainhash.Hash) error {
	return nil
}
func (f fakeStorage) ValidateGraphAnchor(ctx context.Context, graphID *overlay.Outpoint) error {
	return nil
}
func (f fakeStorage) DiscardGraph(ctx context.Context, graphID *overlay.Outpoint) error {
	return nil
}
func (f fakeStorage) FinalizeGraph(ctx context.Context, graphID *overlay.Outpoint) error {
	return nil
}
