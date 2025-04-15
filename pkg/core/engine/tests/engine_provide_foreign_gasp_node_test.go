package engine_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/script"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/stretchr/testify/require"
)

// fakeStorage provides minimal stub for Storage interface
type fakeStorage struct {
	findOutputFunc func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error)
}

func (f fakeStorage) InsertOutput(ctx context.Context, utxo *engine.Output) error {
	panic("method not implemented")
}

func (f fakeStorage) FindOutput(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
	return f.findOutputFunc(ctx, outpoint, topic, spent, includeBEEF)
}

func (f fakeStorage) FindOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) ([]*engine.Output, error) {
	panic("method not implemented")
}

func (f fakeStorage) FindOutputsForTransaction(ctx context.Context, txid *chainhash.Hash, includeBEEF bool) ([]*engine.Output, error) {
	panic("method not implemented")
}

func (f fakeStorage) FindUTXOsForTopic(ctx context.Context, topic string, since uint32, includeBEEF bool) ([]*engine.Output, error) {
	panic("method not implemented")
}

func (f fakeStorage) DeleteOutput(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	panic("method not implemented")
}

func (f fakeStorage) DeleteOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	panic("method not implemented")
}

func (f fakeStorage) MarkUTXOAsSpent(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	panic("method not implemented")
}

func (f fakeStorage) MarkUTXOsAsSpent(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	panic("method not implemented")
}

func (f fakeStorage) UpdateConsumedBy(ctx context.Context, outpoint *overlay.Outpoint, topic string, consumedBy []*overlay.Outpoint) error {
	panic("method not implemented")
}

func (f fakeStorage) UpdateTransactionBEEF(ctx context.Context, txid *chainhash.Hash, beef []byte) error {
	panic("method not implemented")
}

func (f fakeStorage) UpdateOutputBlockHeight(ctx context.Context, outpoint *overlay.Outpoint, topic string, blockHeight uint32, blockIndex uint64, ancillaryBeef []byte) error {
	panic("method not implemented")
}

func (f fakeStorage) InsertAppliedTransaction(ctx context.Context, tx *overlay.AppliedTransaction) error {
	panic("method not implemented")
}

func (f fakeStorage) DoesAppliedTransactionExist(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
	panic("method not implemented")
}

func createDummyBEEF(t *testing.T) []byte {
	t.Helper()

	dummyTx := transaction.Transaction{
		Inputs: []*transaction.TransactionInput{},
		Outputs: []*transaction.TransactionOutput{
			{
				Satoshis:      1000,
				LockingScript: &script.Script{script.OpRETURN},
			},
		},
	}

	BEEF, err := transaction.NewBeefFromTransaction(&dummyTx)
	require.NoError(t, err)

	bytes, err := BEEF.AtomicBytes(dummyTx.TxID())
	require.NoError(t, err)
	return bytes
}

func parseBEEFToTx(t *testing.T, bytes []byte) *transaction.Transaction {
	t.Helper()

	_, tx, _, err := transaction.ParseBeef(bytes)
	require.NoError(t, err)
	return tx
}

func TestEngine_ProvideForeignGASPNode_Success(t *testing.T) {
	// given:
	ctx := context.Background()
	graphID := &overlay.Outpoint{}
	outpoint := &overlay.Outpoint{OutputIndex: 1}
	BEEF := createDummyBEEF(t)

	expectedNode := &core.GASPNode{
		GraphID:     graphID,
		RawTx:       parseBEEFToTx(t, BEEF).Hex(),
		OutputIndex: outpoint.OutputIndex,
	}

	sut := &engine.Engine{
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				return &engine.Output{Beef: BEEF}, nil
			},
		},
	}

	// when:
	node, err := sut.ProvideForeignGASPNode(ctx, graphID, outpoint, "test-topic")

	// then:
	require.NoError(t, err)
	require.Equal(t, expectedNode, node)
}

func TestEngine_ProvideForeignGASPNode_MissingBeef_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	graphID := &overlay.Outpoint{}
	outpoint := &overlay.Outpoint{}

	sut := &engine.Engine{
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				return &engine.Output{}, nil // Missing Beef
			},
		},
	}

	// when:
	node, err := sut.ProvideForeignGASPNode(ctx, graphID, outpoint, "test-topic")

	// then:
	require.ErrorIs(t, err, engine.ErrMissingInput)
	require.Nil(t, node)
}

func TestEngine_ProvideForeignGASPNode_CannotFindOutput_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	graphID := &overlay.Outpoint{}
	outpoint := &overlay.Outpoint{}
	expectedErr := errors.New("forced error")

	sut := &engine.Engine{
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				return nil, expectedErr
			},
		},
	}

	// when:
	node, err := sut.ProvideForeignGASPNode(ctx, graphID, outpoint, "test-topic")

	// then:
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, node)
}

func TestEngine_ProvideForeignGASPNode_TransactionNotFound_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	graphID := &overlay.Outpoint{}
	outpoint := &overlay.Outpoint{}

	sut := &engine.Engine{
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				return &engine.Output{Beef: []byte{0x00}}, nil
			},
		},
	}

	// when:
	node, err := sut.ProvideForeignGASPNode(ctx, graphID, outpoint, "test-topic")

	// then:
	require.ErrorContains(t, err, "invalid-atomic-beef") // temp solution
	require.Nil(t, node)
}
