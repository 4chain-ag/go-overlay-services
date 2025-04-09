package engine_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/script"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/stretchr/testify/require"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
)

// errFakeStorage is returned for unimplemented methods.
var errFakeStorage = errors.New("fakeStorage: method not implemented")

// fakeStorage is a test double for engine.Storage interface.
type fakeStorage struct {
	findOutputFunc               func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error)
	doesAppliedTransactionExistFunc func(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error)
	insertOutputFunc              func(ctx context.Context, utxo *engine.Output) error
	markUTXOAsSpentFunc           func(ctx context.Context, outpoint *overlay.Outpoint, topic string) error
	insertAppliedTransactionFunc  func(ctx context.Context, tx *overlay.AppliedTransaction) error
	updateConsumedByFunc             func(ctx context.Context, outpoint *overlay.Outpoint, topic string, consumedBy []*overlay.Outpoint) error
}

// Provide dynamic behavior for key methods
func (f fakeStorage) FindOutput(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
	if f.findOutputFunc != nil {
		return f.findOutputFunc(ctx, outpoint, topic, spent, includeBEEF)
	}
	return nil, errFakeStorage
}

func (f fakeStorage) DoesAppliedTransactionExist(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
	if f.doesAppliedTransactionExistFunc != nil {
		return f.doesAppliedTransactionExistFunc(ctx, tx)
	}
	return false, errFakeStorage
}

func (f fakeStorage) InsertOutput(ctx context.Context, utxo *engine.Output) error {
	if f.insertOutputFunc != nil {
		return f.insertOutputFunc(ctx, utxo)
	}
	return errFakeStorage
}

func (f fakeStorage) MarkUTXOAsSpent(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	if f.markUTXOAsSpentFunc != nil {
		return f.markUTXOAsSpentFunc(ctx, outpoint, topic)
	}
	return errFakeStorage
}

func (f fakeStorage) InsertAppliedTransaction(ctx context.Context, tx *overlay.AppliedTransaction) error {
	if f.insertAppliedTransactionFunc != nil {
		return f.insertAppliedTransactionFunc(ctx, tx)
	}
	return errFakeStorage
}

func (f fakeStorage) FindOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) ([]*engine.Output, error) {
	return nil, errFakeStorage
}

func (f fakeStorage) FindOutputsForTransaction(ctx context.Context, txid *chainhash.Hash, includeBEEF bool) ([]*engine.Output, error) {
	return nil, errFakeStorage
}

func (f fakeStorage) FindUTXOsForTopic(ctx context.Context, topic string, since uint32, includeBEEF bool) ([]*engine.Output, error) {
	return nil, errFakeStorage
}

func (f fakeStorage) DeleteOutput(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	return errFakeStorage
}

func (f fakeStorage) DeleteOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	return errFakeStorage
}

func (f fakeStorage) MarkUTXOsAsSpent(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	return errFakeStorage
}

func (f fakeStorage) UpdateConsumedBy(ctx context.Context, outpoint *overlay.Outpoint, topic string, consumedBy []*overlay.Outpoint) error {
	if f.updateConsumedByFunc != nil {
		return f.updateConsumedByFunc(ctx, outpoint, topic, consumedBy)
	}
	return nil
}

func (f fakeStorage) UpdateTransactionBEEF(ctx context.Context, txid *chainhash.Hash, beef []byte) error {
	return errFakeStorage
}

func (f fakeStorage) UpdateOutputBlockHeight(ctx context.Context, outpoint *overlay.Outpoint, topic string, blockHeight uint32, blockIndex uint64, ancillaryBeef []byte) error {
	return errFakeStorage
}

type fakeManager struct{}

func (f fakeManager) IdentifyAdmissableOutputs(ctx context.Context, beef []byte, previousCoins []uint32) (overlay.AdmittanceInstructions, error) {
	return overlay.AdmittanceInstructions{
		OutputsToAdmit: []uint32{0},
		CoinsToRetain:  nil,             
		CoinsRemoved:   nil,              
		AncillaryTxids: nil,             
	}, nil
}

func (f fakeManager) IdentifyNeededInputs(ctx context.Context, beef []byte) ([]*overlay.Outpoint, error) {
	return nil, nil
}

func (f fakeManager) GetMetaData() *overlay.MetaData {
	return nil
}

func (f fakeManager) GetDocumentation() string {
	return ""
}


type fakeChainTracker struct{}

func (f fakeChainTracker) Verify(tx *transaction.Transaction, options ...any) (bool, error) {
	return true, nil
}

func (f fakeChainTracker) IsValidRootForHeight(root *chainhash.Hash, height uint32) (bool, error) {
	return true, nil
}

func (f fakeChainTracker) FindHeader(height uint32) ([]byte, error) {
	return nil, nil
}

func (f fakeChainTracker) FindPreviousHeader(tx *transaction.Transaction) ([]byte, error) {
	return nil, nil
}

type fakeChainTrackerSPVFail struct{}

func (f fakeChainTrackerSPVFail) Verify(tx *transaction.Transaction, options ...any) (bool, error) {
	return false, nil
}

func (f fakeChainTrackerSPVFail) IsValidRootForHeight(root *chainhash.Hash, height uint32) (bool, error) {
	return true, nil
}

func (f fakeChainTrackerSPVFail) FindHeader(height uint32) ([]byte, error) {
	return nil, nil
}

func (f fakeChainTrackerSPVFail) FindPreviousHeader(tx *transaction.Transaction) ([]byte, error) {
	return nil, nil
}

type fakeBroadcasterFail struct{}

func (f fakeBroadcasterFail) Broadcast(tx *transaction.Transaction) (*transaction.BroadcastSuccess, *transaction.BroadcastFailure) {
	return nil, &transaction.BroadcastFailure{
		Code:        "broadcast-failed",
		Description: "forced failure for testing",
	}
}

func (f fakeBroadcasterFail) BroadcastCtx(ctx context.Context, tx *transaction.Transaction) (*transaction.BroadcastSuccess, *transaction.BroadcastFailure) {
	return nil, &transaction.BroadcastFailure{
		Code:        "broadcast-failed",
		Description: "forced failure for testing",
	}
}

func createDummyBeef(t *testing.T) []byte {
	t.Helper()

	dummyLockingScript := script.Script{script.OpRETURN}

	dummyTx := transaction.Transaction{
		Inputs: []*transaction.TransactionInput{},
		Outputs: []*transaction.TransactionOutput{
			{
				Satoshis:      1000,
				LockingScript: &dummyLockingScript,
			},
		},
	}

	beef, err := transaction.NewBeefFromTransaction(&dummyTx)
	require.NoError(t, err)

	serializedBytes, err := beef.AtomicBytes(dummyTx.TxID())
	require.NoError(t, err)

	return serializedBytes
}

func createDummyValidTaggedBEEF(t *testing.T) (overlay.TaggedBEEF, *chainhash.Hash) {
	t.Helper()

	prevTx := &transaction.Transaction{
		Inputs: []*transaction.TransactionInput{},
		Outputs: []*transaction.TransactionOutput{
			{
				Satoshis:      1000,
				LockingScript: &script.Script{script.OpTRUE},
			},
		},
	}
	prevTxID := prevTx.TxID()

	currentTx := &transaction.Transaction{
		Inputs: []*transaction.TransactionInput{
			{
				SourceTXID:       prevTxID,
				SourceTxOutIndex: 0,
			},
		},
		Outputs: []*transaction.TransactionOutput{
			{
				Satoshis:      900,
				LockingScript: &script.Script{script.OpTRUE},
			},
		},
	}
	currentTxID := currentTx.TxID()

	beef := &transaction.Beef{
		Version:      transaction.BEEF_V2,
		Transactions: make(map[string]*transaction.BeefTx),
	}

	beef.Transactions[prevTxID.String()] = &transaction.BeefTx{Transaction: prevTx}
	beef.Transactions[currentTxID.String()] = &transaction.BeefTx{Transaction: currentTx}

	beefBytes, err := beef.AtomicBytes(currentTxID)
	require.NoError(t, err)

	return overlay.TaggedBEEF{
		Topics: []string{"test-topic"},
		Beef:   beefBytes,
	}, prevTxID
}


func mustParseHash(t *testing.T, hexStr string) *chainhash.Hash {
	t.Helper()

	h, err := chainhash.NewHashFromHex(hexStr)
	if err != nil {
		t.Fail()
	}
	return h
}

func mustBytes(t *testing.T, beef *transaction.Beef) []byte {
	t.Helper()
	b, err := beef.Bytes()
	require.NoError(t, err)
	return b
}
