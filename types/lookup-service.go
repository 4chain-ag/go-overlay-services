package types

import (
	"context"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
)

type LookupService interface {
	OutputAdded(ctx context.Context, output *Output) error
	OutputSpent(ctx context.Context, txid *chainhash.Hash, outputIndex uint32, topic string) error
	OutputDeleted(ctx context.Context, outpoint *overlay.Outpoint, topic string) error
	Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error)
	GetDocumentation() string
	GetMetaData() *overlay.MetaData
}
