package engine

import (
	"context"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

type OverlayEngineRepositoryAdapter struct {
}

func (*OverlayEngineRepositoryAdapter) InsertOutput(ctx context.Context, utxo *Output) error {
	return nil
}

func (*OverlayEngineRepositoryAdapter) FindOutput(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*Output, error) {
	return nil, nil
}

func (*OverlayEngineRepositoryAdapter) FindOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) ([]*Output, error) {
	return nil, nil
}

func (*OverlayEngineRepositoryAdapter) FindOutputsForTransaction(ctx context.Context, txid *chainhash.Hash, includeBEEF bool) ([]*Output, error) {
	return nil, nil
}

func (*OverlayEngineRepositoryAdapter) FindUTXOsForTopic(ctx context.Context, topic string, since float64, includeBEEF bool) ([]*Output, error) {
	return nil, nil
}

func (*OverlayEngineRepositoryAdapter) DeleteOutput(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	return nil
}

func (*OverlayEngineRepositoryAdapter) DeleteOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	return nil
}

func (*OverlayEngineRepositoryAdapter) MarkUTXOAsSpent(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	return nil
}

func (*OverlayEngineRepositoryAdapter) MarkUTXOsAsSpent(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	return nil
}

func (*OverlayEngineRepositoryAdapter) UpdateConsumedBy(ctx context.Context, outpoint *overlay.Outpoint, topic string, consumedBy []*overlay.Outpoint) error {
	return nil
}

func (*OverlayEngineRepositoryAdapter) UpdateTransactionBEEF(ctx context.Context, txid *chainhash.Hash, beef []byte) error {
	return nil
}

func (*OverlayEngineRepositoryAdapter) UpdateOutputBlockHeight(ctx context.Context, outpoint *overlay.Outpoint, topic string, blockHeight uint32, blockIndex uint64) error {
	return nil
}

func (*OverlayEngineRepositoryAdapter) InsertAppliedTransaction(ctx context.Context, tx *overlay.AppliedTransaction) error {
	return nil
}

func (*OverlayEngineRepositoryAdapter) DoesAppliedTransactionExist(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
	return false, nil
}
