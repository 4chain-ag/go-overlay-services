package gasp_test

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

func TestOverlayGASPStorage_AppendToGraph(t *testing.T) {
	t.Run("should append a new node to an empty graph", func(t *testing.T) {
		// given
		ctx := context.Background()
		mockEngine := &engine.Engine{
			Storage: &mockStorage{},
		}
		storage := engine.NewOverlayGASPStorage("test-topic", mockEngine, nil)

		// Create a minimal valid transaction
		tx := transaction.NewTransaction()
		tx.AddOutput(&transaction.TransactionOutput{
			Satoshis:      1000,
			LockingScript: &script.Script{},
		})
		
		graphID := &overlay.Outpoint{
			Txid:        *tx.TxID(),
			OutputIndex: 0,
		}
		
		gaspNode := &core.GASPNode{
			RawTx:       tx.Hex(),
			OutputIndex: 0,
			GraphID:     graphID,
		}

		// when
		err := storage.AppendToGraph(ctx, gaspNode, nil)

		// then
		require.NoError(t, err)
		// Verify node was added by trying to append a child
		childTx := transaction.NewTransaction()
		childTx.AddOutput(&transaction.TransactionOutput{
			Satoshis:      500,
			LockingScript: &script.Script{},
		})
		
		childNode := &core.GASPNode{
			RawTx:       childTx.Hex(),
			OutputIndex: 0,
			GraphID:     graphID,
		}
		
		err = storage.AppendToGraph(ctx, childNode, tx.TxID())
		require.NoError(t, err)
	})

	t.Run("should return error when max nodes exceeded", func(t *testing.T) {
		// given
		ctx := context.Background()
		maxNodes := 2
		mockEngine := &engine.Engine{
			Storage: &mockStorage{},
		}
		storage := engine.NewOverlayGASPStorage("test-topic", mockEngine, &maxNodes)

		// Add nodes up to the limit
		for i := 0; i < maxNodes; i++ {
			tx := transaction.NewTransaction()
			tx.AddOutput(&transaction.TransactionOutput{
				Satoshis:      1000,
				LockingScript: &script.Script{},
			})
			
			graphID := &overlay.Outpoint{
				Txid:        *tx.TxID(),
				OutputIndex: uint32(i),
			}
			
			gaspNode := &core.GASPNode{
				RawTx:       tx.Hex(),
				OutputIndex: uint32(i),
				GraphID:     graphID,
			}
			
			err := storage.AppendToGraph(ctx, gaspNode, nil)
			require.NoError(t, err)
		}

		// Try to add one more node
		tx := transaction.NewTransaction()
		tx.AddOutput(&transaction.TransactionOutput{
			Satoshis:      1000,
			LockingScript: &script.Script{},
		})
		
		graphID := &overlay.Outpoint{
			Txid:        *tx.TxID(),
			OutputIndex: 99,
		}
		
		gaspNode := &core.GASPNode{
			RawTx:       tx.Hex(),
			OutputIndex: 99,
			GraphID:     graphID,
		}

		// when
		err := storage.AppendToGraph(ctx, gaspNode, nil)

		// then
		require.Error(t, err)
		require.Equal(t, engine.ErrGraphFull, err)
	})

	t.Run("should return error for invalid transaction hex", func(t *testing.T) {
		// given
		ctx := context.Background()
		mockEngine := &engine.Engine{
			Storage: &mockStorage{},
		}
		storage := engine.NewOverlayGASPStorage("test-topic", mockEngine, nil)

		gaspNode := &core.GASPNode{
			RawTx:       "invalid-hex",
			OutputIndex: 0,
			GraphID: &overlay.Outpoint{
				Txid:        chainhash.Hash{},
				OutputIndex: 0,
			},
		}

		// when
		err := storage.AppendToGraph(ctx, gaspNode, nil)

		// then
		require.Error(t, err)
	})
}

func TestOverlayGASPStorage_FindKnownUTXOs(t *testing.T) {
	t.Run("should return known UTXOs since given timestamp", func(t *testing.T) {
		// given
		ctx := context.Background()
		since := uint32(1234567890)
		expectedUTXOs := []*engine.Output{
			{
				Outpoint: overlay.Outpoint{
					Txid:        chainhash.Hash{1},
					OutputIndex: 0,
				},
			},
			{
				Outpoint: overlay.Outpoint{
					Txid:        chainhash.Hash{2},
					OutputIndex: 1,
				},
			},
		}

		mockStorage := &mockStorage{
			findUTXOsForTopicFunc: func(ctx context.Context, topic string, since uint32, historical bool) ([]*engine.Output, error) {
				return expectedUTXOs, nil
			},
		}

		mockEngine := &engine.Engine{
			Storage: mockStorage,
		}
		storage := engine.NewOverlayGASPStorage("test-topic", mockEngine, nil)

		// when
		result, err := storage.FindKnownUTXOs(ctx, since)

		// then
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, &expectedUTXOs[0].Outpoint, result[0])
		require.Equal(t, &expectedUTXOs[1].Outpoint, result[1])
	})

	t.Run("should handle storage errors", func(t *testing.T) {
		// given
		ctx := context.Background()
		expectedErr := errors.New("database error")

		mockStorage := &mockStorage{
			findUTXOsForTopicFunc: func(ctx context.Context, topic string, since uint32, historical bool) ([]*engine.Output, error) {
				return nil, expectedErr
			},
		}

		mockEngine := &engine.Engine{
			Storage: mockStorage,
		}
		storage := engine.NewOverlayGASPStorage("test-topic", mockEngine, nil)

		// when
		result, err := storage.FindKnownUTXOs(ctx, 0)

		// then
		require.Error(t, err)
		require.Equal(t, expectedErr, err)
		require.Nil(t, result)
	})
}

func TestOverlayGASPStorage_DiscardGraph(t *testing.T) {
	t.Run("should discard graph and all its nodes", func(t *testing.T) {
		// given
		ctx := context.Background()
		mockEngine := &engine.Engine{
			Storage: &mockStorage{},
		}
		storage := engine.NewOverlayGASPStorage("test-topic", mockEngine, nil)

		// Create a graph with root and child nodes
		rootTx := transaction.NewTransaction()
		rootTx.AddOutput(&transaction.TransactionOutput{
			Satoshis:      1000,
			LockingScript: &script.Script{},
		})
		
		graphID := &overlay.Outpoint{
			Txid:        *rootTx.TxID(),
			OutputIndex: 0,
		}
		
		rootNode := &core.GASPNode{
			RawTx:       rootTx.Hex(),
			OutputIndex: 0,
			GraphID:     graphID,
		}

		// Add root node
		err := storage.AppendToGraph(ctx, rootNode, nil)
		require.NoError(t, err)

		// Add child node
		childTx := transaction.NewTransaction()
		childTx.AddOutput(&transaction.TransactionOutput{
			Satoshis:      500,
			LockingScript: &script.Script{},
		})
		
		childNode := &core.GASPNode{
			RawTx:       childTx.Hex(),
			OutputIndex: 0,
			GraphID:     graphID,
		}
		
		err = storage.AppendToGraph(ctx, childNode, rootTx.TxID())
		require.NoError(t, err)

		// when
		err = storage.DiscardGraph(ctx, graphID)

		// then
		require.NoError(t, err)
		
		// Verify graph is empty by trying to add to the discarded graph
		newNode := &core.GASPNode{
			RawTx:       rootTx.Hex(),
			OutputIndex: 1,
			GraphID:     graphID,
		}
		
		// This should fail because the parent node was discarded
		err = storage.AppendToGraph(ctx, newNode, rootTx.TxID())
		require.Error(t, err)
	})

	t.Run("should handle non-existent graphID gracefully", func(t *testing.T) {
		// given
		ctx := context.Background()
		mockEngine := &engine.Engine{
			Storage: &mockStorage{},
		}
		storage := engine.NewOverlayGASPStorage("test-topic", mockEngine, nil)

		nonExistentGraphID := &overlay.Outpoint{
			Txid:        chainhash.Hash{99, 99, 99},
			OutputIndex: 0,
		}

		// when
		err := storage.DiscardGraph(ctx, nonExistentGraphID)

		// then
		require.NoError(t, err)
	})
}

func TestOverlayGASPStorage_HydrateGASPNode(t *testing.T) {
	t.Run("should return error when no output found", func(t *testing.T) {
		// given
		ctx := context.Background()
		mockStorage := &mockStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, historical bool) (*engine.Output, error) {
				return nil, nil // No output found
			},
		}

		mockEngine := &engine.Engine{
			Storage: mockStorage,
		}
		storage := engine.NewOverlayGASPStorage("test-topic", mockEngine, nil)

		graphID := &overlay.Outpoint{
			Txid:        chainhash.Hash{1},
			OutputIndex: 0,
		}
		outpoint := &overlay.Outpoint{
			Txid:        chainhash.Hash{2},
			OutputIndex: 0,
		}

		// when
		result, err := storage.HydrateGASPNode(ctx, graphID, outpoint, false)

		// then
		require.Error(t, err)
		require.Equal(t, engine.ErrMissingInput, err)
		require.Nil(t, result)
	})

	t.Run("should hydrate node with valid BEEF", func(t *testing.T) {
		// given
		ctx := context.Background()
		
		// Create a transaction with merkle path
		tx := transaction.NewTransaction()
		tx.AddOutput(&transaction.TransactionOutput{
			Satoshis:      1000,
			LockingScript: &script.Script{},
		})
		
		// Create mock merkle path
		tx.MerklePath = &transaction.MerklePath{
			BlockHeight: 100,
			Path:        [][]*transaction.PathElement{},
		}
		
		beef, err := transaction.NewBeefFromTransaction(tx)
		require.NoError(t, err)
		beefBytes, err := beef.AtomicBytes(tx.TxID())
		require.NoError(t, err)

		mockStorage := &mockStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, historical bool) (*engine.Output, error) {
				return &engine.Output{
					Outpoint: *outpoint,
					Beef:     beefBytes,
				}, nil
			},
		}

		mockEngine := &engine.Engine{
			Storage: mockStorage,
		}
		storage := engine.NewOverlayGASPStorage("test-topic", mockEngine, nil)

		graphID := &overlay.Outpoint{
			Txid:        chainhash.Hash{1},
			OutputIndex: 0,
		}
		outpoint := &overlay.Outpoint{
			Txid:        *tx.TxID(),
			OutputIndex: 0,
		}

		// when
		result, err := storage.HydrateGASPNode(ctx, graphID, outpoint, false)

		// then
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, graphID, result.GraphID)
		require.Equal(t, uint32(0), result.OutputIndex)
		require.Equal(t, tx.Hex(), result.RawTx)
		require.NotNil(t, result.Proof)
	})
}

// Mock storage implementation
type mockStorage struct {
	findUTXOsForTopicFunc func(ctx context.Context, topic string, since uint32, historical bool) ([]*engine.Output, error)
	findOutputFunc        func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, historical bool) (*engine.Output, error)
	findOutputsFunc       func(ctx context.Context, outpoints []*overlay.Outpoint, topic string, spent *bool, historical bool) ([]*engine.Output, error)
}

func (m *mockStorage) FindUTXOsForTopic(ctx context.Context, topic string, since uint32, historical bool) ([]*engine.Output, error) {
	if m.findUTXOsForTopicFunc != nil {
		return m.findUTXOsForTopicFunc(ctx, topic, since, historical)
	}
	return nil, nil
}

func (m *mockStorage) FindOutput(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, historical bool) (*engine.Output, error) {
	if m.findOutputFunc != nil {
		return m.findOutputFunc(ctx, outpoint, topic, spent, historical)
	}
	return nil, nil
}

func (m *mockStorage) FindOutputs(ctx context.Context, outpoints []*overlay.Outpoint, topic string, spent *bool, historical bool) ([]*engine.Output, error) {
	if m.findOutputsFunc != nil {
		return m.findOutputsFunc(ctx, outpoints, topic, spent, historical)
	}
	return nil, nil
}

// Implement remaining Storage interface methods with empty implementations
func (m *mockStorage) SetIncoming(ctx context.Context, txs []*transaction.Transaction) error { return nil }
func (m *mockStorage) SetOutgoing(ctx context.Context, tx *transaction.Transaction, steak *overlay.Steak) error {
	return nil
}
func (m *mockStorage) UpdateConsumedBy(ctx context.Context, outpoint *overlay.Outpoint, consumedBy string, inputs []*overlay.Outpoint) error {
	return nil
}
func (m *mockStorage) DeleteOutput(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
	return nil
}
func (m *mockStorage) FindTransaction(ctx context.Context, txid chainhash.Hash, requireProof bool) (*transaction.Transaction, error) {
	return nil, nil
}
func (m *mockStorage) FindTransactionsCreatingUtxos(ctx context.Context) ([]*chainhash.Hash, error) {
	return nil, nil
}

func (m *mockStorage) DoesAppliedTransactionExist(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
	return false, nil
}

func (m *mockStorage) InsertAppliedTransaction(ctx context.Context, tx *overlay.AppliedTransaction) error {
	return nil
}

func (m *mockStorage) UpdateTransactionBEEF(ctx context.Context, txid *chainhash.Hash, beef []byte) error {
	return nil
}

func (m *mockStorage) MarkUTXOsAsSpent(ctx context.Context, utxos []*overlay.Outpoint, spentBy string, blockHash *chainhash.Hash) error {
	return nil
}

func (m *mockStorage) InsertOutput(ctx context.Context, output *engine.Output) error {
	return nil
}

func (m *mockStorage) FindOutputsForTransaction(ctx context.Context, txid *chainhash.Hash, includeBEEF bool) ([]*engine.Output, error) {
	return nil, nil
}

func (m *mockStorage) UpdateOutputBlockHeight(ctx context.Context, outpoint *overlay.Outpoint, topic string, blockHeight uint32, blockIndex uint64, ancillaryBeef []byte) error {
	return nil
}
