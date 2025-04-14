package core_test

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

type fakeStorage struct {
	hydrateGASPNodeFunc func(ctx context.Context, graphID, outpoint *overlay.Outpoint, metadata bool) (*core.GASPNode, error)
	findKnownUTXOsFunc  func(ctx context.Context, since uint32) ([]*overlay.Outpoint, error)
}

func (f fakeStorage) HydrateGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, metadata bool) (*core.GASPNode, error) {
	if f.hydrateGASPNodeFunc != nil {
		return f.hydrateGASPNodeFunc(ctx, graphID, outpoint, metadata)
	}
	panic("hydrateGASPNodeFunc not set")
}

func (f fakeStorage) FindKnownUTXOs(ctx context.Context, since uint32) ([]*overlay.Outpoint, error) {
	if f.findKnownUTXOsFunc != nil {
		return f.findKnownUTXOsFunc(ctx, since)
	}
	panic("findKnownUTXOsFunc not set")
}
func (f fakeStorage) FindNeededInputs(ctx context.Context, tx *core.GASPNode) (*core.GASPNodeResponse, error) {
	panic("findNeededInputsFunc not set")
}
func (f fakeStorage) AppendToGraph(ctx context.Context, tx *core.GASPNode, spentBy *chainhash.Hash) error {
	panic("appendToGraphFunc not set")
}
func (f fakeStorage) ValidateGraphAnchor(ctx context.Context, graphID *overlay.Outpoint) error {
	panic("validateGraphAnchorFunc not set")
}
func (f fakeStorage) DiscardGraph(ctx context.Context, graphID *overlay.Outpoint) error {
	panic("discardGraphFunc not set")
}
func (f fakeStorage) FinalizeGraph(ctx context.Context, graphID *overlay.Outpoint) error {
	panic("finalizeGraphFunc not set")
}
