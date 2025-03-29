package core

import (
	"context"

	"github.com/bsv-blockchain/go-sdk/overlay"
)

type GASPRemote interface {
	GetInitialResponse(ctx context.Context, request *GASPInitialRequest) (*GASPInitialResponse, error)
	GetInitialReplay(ctx context.Context, response *GASPInitialResponse) (*GASPInitialReply, error)
	RequestNode(ctx context.Context, graphID *overlay.Outpoint, outpoint *overlay.Outpoint, meatdata bool) (*GASPNode, error)
	SubmitNode(ctx context.Context, node *GASPNode) (*GASPNodeResponse, error)
}
