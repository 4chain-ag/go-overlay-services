package core

import (
	"context"

	"github.com/bsv-blockchain/go-sdk/overlay"
)

type GASPStorage interface {
	FindKnownUTXOs(ctx context.Context, since uint32) ([]*overlay.Outpoint, error)
	HydrateGASPNode(ctx context.Context, graphID *overlay.Outpoint, outpoint *overlay.Outpoint, metadata bool) (*GASPNode, error)
	FindNeededInputs(ctx context.Context, tx *GASPNode) (*GASPNodeResponse, error)
	AppendToGraph(ctx context.Context, tx *GASPNode, spentBy *overlay.Outpoint) error
	ValidateGraphAnchor(ctx context.Context, graphID *overlay.Outpoint) error
	DiscardGraph(ctx context.Context, graphID *overlay.Outpoint) error
	FinalizeGraph(ctx context.Context, graphID *overlay.Outpoint) error
}
