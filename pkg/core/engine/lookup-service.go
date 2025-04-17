package engine

import (
	"context"

	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
)

type LookupService interface {
	OutputAdded(ctx context.Context, outpoint *overlay.Outpoint, topic string, beef []byte) error
	OutputSpent(ctx context.Context, outpoint *overlay.Outpoint, topic string, beef []byte) error
	OutputDeleted(ctx context.Context, outpoint *overlay.Outpoint, topic string) error
	OutputBlockHeightUpdated(ctx context.Context, outpoint *overlay.Outpoint, blockHeight uint32, blockIndex uint64) error
	Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error)
	GetDocumentation() string
	GetMetaData() *overlay.MetaData
}
