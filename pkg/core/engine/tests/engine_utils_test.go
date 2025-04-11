package engine_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/script"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/stretchr/testify/require"
)

var errFakeStorage = errors.New("fakeStorage: method not implemented")

type fakeStorage struct {
	findOutputFunc                  func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error)
	findOutputsFunc                 func(ctx context.Context, outpoints []*overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) ([]*engine.Output, error)
	doesAppliedTransactionExistFunc func(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error)
	insertOutputFunc                func(ctx context.Context, utxo *engine.Output) error
	markUTXOAsSpentFunc             func(ctx context.Context, outpoint *overlay.Outpoint, topic string) error
	insertAppliedTransactionFunc    func(ctx context.Context, tx *overlay.AppliedTransaction) error
	updateConsumedByFunc            func(ctx context.Context, outpoint *overlay.Outpoint, topic string, consumedBy []*overlay.Outpoint) error
	deleteOutputFunc                func(ctx context.Context, outpoint *overlay.Outpoint, topic string) error
}

func (f fakeStorage) FindOutput(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
	if f.findOutputFunc != nil {
		return f.findOutputFunc(ctx, outpoint, topic, spent, includeBEEF)
	}
	panic("func not defined")
}
func (f fakeStorage) DoesAppliedTransactionExist(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
	if f.doesAppliedTransactionExistFunc != nil {
		return f.doesAppliedTransactionExistFunc(ctx, tx)
	}
	panic("func not defined")
}
func (f fakeStorage) InsertOutput(ctx context.Context, utxo *engine.Output) error {
	if f.insertOutputFunc != nil {
		return f.insertOutputFunc(ctx, utxo)
	}
	panic("func not defined")
}
func (f fakeStorage) MarkUTXOAsSpent(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	if f.markUTXOAsSpentFunc != nil {
		return f.markUTXOAsSpentFunc(ctx, outpoint, topic)
	}
	panic("func not defined")
}
func (f fakeStorage) InsertAppliedTransaction(ctx context.Context, tx *overlay.AppliedTransaction) error {
	if f.insertAppliedTransactionFunc != nil {
		return f.insertAppliedTransactionFunc(ctx, tx)
	}
	panic("func not defined")
}
func (f fakeStorage) UpdateConsumedBy(ctx context.Context, outpoint *overlay.Outpoint, topic string, consumedBy []*overlay.Outpoint) error {
	if f.updateConsumedByFunc != nil {
		return f.updateConsumedByFunc(ctx, outpoint, topic, consumedBy)
	}
	panic("func not defined")
}
func (f fakeStorage) DeleteOutput(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	if f.deleteOutputFunc != nil {
		return f.deleteOutputFunc(ctx, outpoint, topic)
	}
	panic("DeleteOutput not defined")
}
func (f fakeStorage) FindOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) ([]*engine.Output, error) {
	if f.findOutputsFunc != nil {
		return f.findOutputsFunc(ctx, outpoints, topic, spent, includeBEEF)
	}
	panic("FindOutputs not defined")
}

func (f fakeStorage) FindOutputsForTransaction(ctx context.Context, txid *chainhash.Hash, includeBEEF bool) ([]*engine.Output, error) {
	panic("func not defined")
}

func (f fakeStorage) FindUTXOsForTopic(ctx context.Context, topic string, since uint32, includeBEEF bool) ([]*engine.Output, error) {
	panic("func not defined")
}

func (f fakeStorage) DeleteOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	panic("func not defined")
}

func (f fakeStorage) MarkUTXOsAsSpent(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	panic("func not defined")
}

func (f fakeStorage) UpdateTransactionBEEF(ctx context.Context, txid *chainhash.Hash, beef []byte) error {
	panic("func not defined")
}

func (f fakeStorage) UpdateOutputBlockHeight(ctx context.Context, outpoint *overlay.Outpoint, topic string, blockHeight uint32, blockIndex uint64, ancillaryBeef []byte) error {
	panic("func not defined")
}

type fakeManager struct {
	identifyAdmissableOutputsFunc func(ctx context.Context, beef []byte, previousCoins []uint32) (overlay.AdmittanceInstructions, error)
	identifyNeededInputsFunc      func(ctx context.Context, beef []byte) ([]*overlay.Outpoint, error)
}

func (f fakeManager) IdentifyAdmissableOutputs(ctx context.Context, beef []byte, previousCoins []uint32) (overlay.AdmittanceInstructions, error) {
	if f.identifyAdmissableOutputsFunc != nil {
		return f.identifyAdmissableOutputsFunc(ctx, beef, previousCoins)
	}
	panic("IdentifyAdmissableOutputs not defined")
}

func (f fakeManager) IdentifyNeededInputs(ctx context.Context, beef []byte) ([]*overlay.Outpoint, error) {
	if f.identifyNeededInputsFunc != nil {
		return f.identifyNeededInputsFunc(ctx, beef)
	}
	panic("IdentifyNeededInputs not defined")
}

func (f fakeManager) GetMetaData() *overlay.MetaData {
	panic("GetMetaData not defined")
}

func (f fakeManager) GetDocumentation() string {
	panic("GetDocumentation not defined")
}

type fakeChainTracker struct {
	verifyFunc             func(tx *transaction.Transaction, options ...any) (bool, error)
	isValidRootForHeight   func(root *chainhash.Hash, height uint32) (bool, error)
	findHeaderFunc         func(height uint32) ([]byte, error)
	findPreviousHeaderFunc func(tx *transaction.Transaction) ([]byte, error)
}

func (f fakeChainTracker) Verify(tx *transaction.Transaction, options ...any) (bool, error) {
	if f.verifyFunc != nil {
		return f.verifyFunc(tx, options...)
	}
	panic("Verify not defined")
}

func (f fakeChainTracker) IsValidRootForHeight(root *chainhash.Hash, height uint32) (bool, error) {
	if f.isValidRootForHeight != nil {
		return f.isValidRootForHeight(root, height)
	}
	panic("IsValidRootForHeight not defined")
}

func (f fakeChainTracker) FindHeader(height uint32) ([]byte, error) {
	if f.findHeaderFunc != nil {
		return f.findHeaderFunc(height)
	}
	panic("FindHeader not defined")
}

func (f fakeChainTracker) FindPreviousHeader(tx *transaction.Transaction) ([]byte, error) {
	if f.findPreviousHeaderFunc != nil {
		return f.findPreviousHeaderFunc(tx)
	}
	panic("FindPreviousHeader not defined")
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

type fakeBroadcasterFail struct {
	broadcastFunc    func(tx *transaction.Transaction) (*transaction.BroadcastSuccess, *transaction.BroadcastFailure)
	broadcastCtxFunc func(ctx context.Context, tx *transaction.Transaction) (*transaction.BroadcastSuccess, *transaction.BroadcastFailure)
}

func (f fakeBroadcasterFail) Broadcast(tx *transaction.Transaction) (*transaction.BroadcastSuccess, *transaction.BroadcastFailure) {
	if f.broadcastFunc != nil {
		return f.broadcastFunc(tx)
	}
	panic("Broadcast not defined")
}

func (f fakeBroadcasterFail) BroadcastCtx(ctx context.Context, tx *transaction.Transaction) (*transaction.BroadcastSuccess, *transaction.BroadcastFailure) {
	if f.broadcastCtxFunc != nil {
		return f.broadcastCtxFunc(ctx, tx)
	}
	panic("BroadcastCtx not defined")
}

// helper function to create a dummy BEEF transaction
// This function creates a dummy BEEF transaction with a single output and no inputs.
// It returns the serialized bytes of the BEEF transaction.
// The transaction is created with a dummy locking script that contains an OP_RETURN opcode.
func createDummyBeef(t *testing.T) []byte {
	t.Helper()
	dummyLockingScript := script.Script{script.OpRETURN}
	dummyTx := transaction.Transaction{
		Inputs:  []*transaction.TransactionInput{},
		Outputs: []*transaction.TransactionOutput{{Satoshis: 1000, LockingScript: &dummyLockingScript}},
	}
	beef, err := transaction.NewBeefFromTransaction(&dummyTx)
	require.NoError(t, err)
	serializedBytes, err := beef.AtomicBytes(dummyTx.TxID())
	require.NoError(t, err)
	return serializedBytes
}

// createDummyValidTaggedBEEF creates a dummy valid tagged BEEF transaction for testing.
// It creates a previous transaction and a current transaction, both with dummy locking scripts.
// The previous transaction is used as an input for the current transaction.
// It returns the tagged BEEF and the transaction ID of the previous transaction.
// The tagged BEEF contains a list of topics and the serialized bytes of the BEEF transaction.
func createDummyValidTaggedBEEF(t *testing.T) (overlay.TaggedBEEF, *chainhash.Hash) {
	t.Helper()
	prevTx := &transaction.Transaction{
		Inputs:  []*transaction.TransactionInput{},
		Outputs: []*transaction.TransactionOutput{{Satoshis: 1000, LockingScript: &script.Script{script.OpTRUE}}},
	}
	prevTxID := prevTx.TxID()

	currentTx := &transaction.Transaction{
		Inputs:  []*transaction.TransactionInput{{SourceTXID: prevTxID, SourceTxOutIndex: 0}},
		Outputs: []*transaction.TransactionOutput{{Satoshis: 900, LockingScript: &script.Script{script.OpTRUE}}},
	}
	currentTxID := currentTx.TxID()

	beef := &transaction.Beef{
		Version: transaction.BEEF_V2,
		Transactions: map[string]*transaction.BeefTx{
			prevTxID.String():    {Transaction: prevTx},
			currentTxID.String(): {Transaction: currentTx},
		},
	}
	beefBytes, err := beef.AtomicBytes(currentTxID)
	require.NoError(t, err)

	return overlay.TaggedBEEF{Topics: []string{"test-topic"}, Beef: beefBytes}, prevTxID
}

// createDummyBeefWithInputs creates a dummy BEEF transaction with inputs for testing.
// It creates a previous transaction with a dummy locking script and a current transaction
// that uses the previous transaction as an input. The current transaction also has a dummy locking script.
// It returns the serialized bytes of the BEEF transaction.
func createDummyBeefWithInputs(t *testing.T) []byte {
	t.Helper()

	prevTxID := chainhash.DoubleHashH([]byte("dummy prev tx"))

	dummyLockingScript := script.Script{script.OpTRUE}

	prevTx := &transaction.Transaction{
		Inputs:  []*transaction.TransactionInput{},
		Outputs: []*transaction.TransactionOutput{{Satoshis: 1000, LockingScript: &dummyLockingScript}},
	}

	currentTx := &transaction.Transaction{
		Inputs: []*transaction.TransactionInput{
			{SourceTXID: &prevTxID, SourceTxOutIndex: 0},
		},
		Outputs: []*transaction.TransactionOutput{
			{Satoshis: 900, LockingScript: &dummyLockingScript},
		},
	}

	beef := &transaction.Beef{
		Version: transaction.BEEF_V2,
		Transactions: map[string]*transaction.BeefTx{
			prevTx.TxID().String():    {Transaction: prevTx},
			currentTx.TxID().String(): {Transaction: currentTx},
		},
	}

	beefBytes, err := beef.AtomicBytes(currentTx.TxID())
	require.NoError(t, err)

	return beefBytes
}
