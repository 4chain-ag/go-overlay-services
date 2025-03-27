package engine

import (
	"context"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/transaction"
)

type TopicManager interface {
	IdentifyAdmissableOutputs(ctx context.Context, beef *transaction.Beef, txid *chainhash.Hash, previousCoins []uint32) (overlay.AdmittanceInstructions, error)
	IdentifyNeededInputs(ctx context.Context, beef *transaction.Beef, txid *chainhash.Hash) ([]*overlay.Outpoint, error)
	GetDocumentation() string
	GetMetaData() *overlay.MetaData
}
