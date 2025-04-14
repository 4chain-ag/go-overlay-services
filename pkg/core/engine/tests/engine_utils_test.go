package engine_test

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/advertiser"
	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/bsv-blockchain/go-sdk/script"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/stretchr/testify/require"
)

var errFakeStorage = errors.New("fakeStorage: method not implemented")

type fakeStorage struct {
	findOutputFunc                  func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error)
	doesAppliedTransactionExistFunc func(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error)
	insertOutputFunc                func(ctx context.Context, utxo *engine.Output) error
	markUTXOAsSpentFunc             func(ctx context.Context, outpoint *overlay.Outpoint, topic string) error
	insertAppliedTransactionFunc    func(ctx context.Context, tx *overlay.AppliedTransaction) error
	updateConsumedByFunc            func(ctx context.Context, outpoint *overlay.Outpoint, topic string, consumedBy []*overlay.Outpoint) error
	deleteOutputFunc                func(ctx context.Context, outpoint *overlay.Outpoint, topic string) error
	findUTXOsForTopicFunc           func(ctx context.Context, topic string, since uint32, includeBEEF bool) ([]*engine.Output, error)
}

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
func (f fakeStorage) UpdateConsumedBy(ctx context.Context, outpoint *overlay.Outpoint, topic string, consumedBy []*overlay.Outpoint) error {
	if f.updateConsumedByFunc != nil {
		return f.updateConsumedByFunc(ctx, outpoint, topic, consumedBy)
	}
	return nil
}
func (f fakeStorage) DeleteOutput(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	if f.deleteOutputFunc != nil {
		return f.deleteOutputFunc(ctx, outpoint, topic)
	}
	return nil
}
func (f fakeStorage) FindOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) ([]*engine.Output, error) {
	return nil, errFakeStorage
}
func (f fakeStorage) FindOutputsForTransaction(ctx context.Context, txid *chainhash.Hash, includeBEEF bool) ([]*engine.Output, error) {
	return nil, errFakeStorage
}
func (f fakeStorage) FindUTXOsForTopic(ctx context.Context, topic string, since uint32, includeBEEF bool) ([]*engine.Output, error) {
	if f.findUTXOsForTopicFunc != nil {
		return f.findUTXOsForTopicFunc(ctx, topic, since, includeBEEF)
	}
	return nil, errFakeStorage
}
func (f fakeStorage) DeleteOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	return errFakeStorage
}
func (f fakeStorage) MarkUTXOsAsSpent(ctx context.Context, outpoints []*overlay.Outpoint, topic string) error {
	return errFakeStorage
}
func (f fakeStorage) UpdateTransactionBEEF(ctx context.Context, txid *chainhash.Hash, beef []byte) error {
	return errFakeStorage
}
func (f fakeStorage) UpdateOutputBlockHeight(ctx context.Context, outpoint *overlay.Outpoint, topic string, blockHeight uint32, blockIndex uint64, ancillaryBeef []byte) error {
	return errFakeStorage
}

type fakeManager struct{}

func (f fakeManager) IdentifyAdmissableOutputs(ctx context.Context, beef []byte, previousCoins []uint32) (overlay.AdmittanceInstructions, error) {
	return overlay.AdmittanceInstructions{OutputsToAdmit: []uint32{0}}, nil
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
	return nil, &transaction.BroadcastFailure{Code: "broadcast-failed", Description: "forced failure for testing"}
}
func (f fakeBroadcasterFail) BroadcastCtx(ctx context.Context, tx *transaction.Transaction) (*transaction.BroadcastSuccess, *transaction.BroadcastFailure) {
	return nil, &transaction.BroadcastFailure{Code: "broadcast-failed", Description: "forced failure for testing"}
}

var errFakeLookup = errors.New("lookup error")

type fakeLookupService struct {
	lookupFunc func(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error)
}

func (f fakeLookupService) Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
	if f.lookupFunc != nil {
		return f.lookupFunc(ctx, question)
	}
	return nil, errors.New("lookup not implemented")
}

func (f fakeLookupService) OutputAdded(ctx context.Context, outpoint *overlay.Outpoint, outputScript *script.Script, topic string, blockHeight uint32, blockIndex uint64) error {
	return nil
}

func (f fakeLookupService) OutputSpent(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	return nil
}

func (f fakeLookupService) OutputDeleted(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	return nil
}

func (f fakeLookupService) OutputBlockHeightUpdated(ctx context.Context, outpoint *overlay.Outpoint, blockHeight uint32, blockIndex uint64) error {
	return nil
}

func (f fakeLookupService) GetDocumentation() string {
	return ""
}

func (f fakeLookupService) GetMetaData() *overlay.MetaData {
	return &overlay.MetaData{}
}

type fakeAdvertiser struct {
	findAllAdvertisements     func(protocol overlay.Protocol) ([]*advertiser.Advertisement, error)
	createAdvertisements      func(data []*advertiser.AdvertisementData) (overlay.TaggedBEEF, error)
	revokeAdvertisements      func(data []*advertiser.Advertisement) (overlay.TaggedBEEF, error)
	parseAdvertisement        func(script *script.Script) (*advertiser.Advertisement, error)
	findAllAdvertisementsFunc func(protocol overlay.Protocol) ([]*advertiser.Advertisement, error)
	createAdvertisementsFunc  func(data []*advertiser.AdvertisementData) (overlay.TaggedBEEF, error)
	revokeAdvertisementsFunc  func(data []*advertiser.Advertisement) (overlay.TaggedBEEF, error)
	parseAdvertisementFunc    func(script *script.Script) (*advertiser.Advertisement, error)
}

func (f fakeAdvertiser) FindAllAdvertisements(protocol overlay.Protocol) ([]*advertiser.Advertisement, error) {
	if f.findAllAdvertisements != nil {
		return f.findAllAdvertisements(protocol)
	}
	return nil, nil
}
func (f fakeAdvertiser) CreateAdvertisements(data []*advertiser.AdvertisementData) (overlay.TaggedBEEF, error) {
	if f.createAdvertisements != nil {
		return f.createAdvertisements(data)
	}
	return overlay.TaggedBEEF{}, nil
}
func (f fakeAdvertiser) RevokeAdvertisements(data []*advertiser.Advertisement) (overlay.TaggedBEEF, error) {
	if f.revokeAdvertisements != nil {
		return f.revokeAdvertisements(data)
	}
	return overlay.TaggedBEEF{}, nil
}
func (f fakeAdvertiser) ParseAdvertisement(script *script.Script) (*advertiser.Advertisement, error) {
	if f.parseAdvertisement != nil {
		return f.parseAdvertisement(script)
	}
	return nil, nil
}

type fakeTopicManager struct{}

func (fakeTopicManager) IdentifyAdmissableOutputs(ctx context.Context, beef []byte, previousCoins []uint32) (overlay.AdmittanceInstructions, error) {
	return overlay.AdmittanceInstructions{}, nil
}
func (fakeTopicManager) IdentifyNeededInputs(ctx context.Context, beef []byte) ([]*overlay.Outpoint, error) {
	return nil, nil
}
func (fakeTopicManager) GetMetaData() *overlay.MetaData {
	return &overlay.MetaData{}
}
func (fakeTopicManager) GetDocumentation() string {
	return ""
}

// helper function to create a dummy BEEF transaction
// This function creates a dummy BEEF transaction with a single output and no inputs.
// It returns the serialized bytes of the BEEF transaction.
// The transaction is created with a dummy locking script that contains an OP_RETURN opcode.
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

func fakeTxID() chainhash.Hash {
	b, _ := hex.DecodeString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	var h chainhash.Hash
	copy(h[:], b)
	return h
}
