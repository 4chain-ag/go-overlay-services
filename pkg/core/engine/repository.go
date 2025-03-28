package engine

import (
	"context"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

type Repository struct{}

func (*Repository) InsertOutput(ctx context.Context, utxo *Output) error {
	return nil
}

func (*Repository) FindOutput(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*Output, error) {
	return nil, nil
}

func (*Repository) FindOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) ([]*Output, error) {
	return nil, nil
}

func (*Repository) FindOutputsForTransaction(ctx context.Context, txid *chainhash.Hash, includeBEEF bool) ([]*Output, error) {
	return nil, nil
}

func (*Repository) FindUTXOsForTopic(ctx context.Context, topic string, since float64, includeBEEF bool) ([]*Output, error) {
	return nil, nil
}

func (*Repository) DeleteOutput(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	return nil
}

func (*Repository) DeleteOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	return nil
}

func (*Repository) MarkUTXOAsSpent(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	return nil
}

func (*Repository) MarkUTXOsAsSpent(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	return nil
}

func (*Repository) UpdateConsumedBy(ctx context.Context, outpoint *overlay.Outpoint, topic string, consumedBy []*overlay.Outpoint) error {
	return nil
}

func (*Repository) UpdateTransactionBEEF(ctx context.Context, txid *chainhash.Hash, beef []byte) error {
	return nil
}

func (*Repository) UpdateOutputBlockHeight(ctx context.Context, outpoint *overlay.Outpoint, topic string, blockHeight uint32, blockIndex uint64) error {
	return nil
}

func (*Repository) InsertAppliedTransaction(ctx context.Context, tx *overlay.AppliedTransaction) error {
	return nil
}

func (*Repository) DoesAppliedTransactionExist(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
	return false, nil
}
