package testabilities

import (
	"context"
	"testing"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/stretchr/testify/require"
)

// TODO: consolidate mocks for arc ingest
// MerkleProofProviderMock is a mock implementation of the NewMerkleProofProvider interface,
// used for testing components that interact with Merkle proof processing.
type MerkleProofProviderMock struct {
	t *testing.T

	// error is the error that will be returned from HandleNewMerkleProof
	error error

	// called indicates whether HandleNewMerkleProof was called
	called bool

	// calledWithTxID stores the txID passed to HandleNewMerkleProof
	calledWithTxID *chainhash.Hash

	// calledWithProof stores the MerklePath passed to HandleNewMerkleProof
	calledWithProof *transaction.MerklePath

	// expectedBlockHeight is the block height that's expected to be provided in the MerklePath
	expectedBlockHeight uint32
}

// HandleNewMerkleProof simulates the handling of a new Merkle proof. It records the call
// and returns the predefined error if set.
func (m *MerkleProofProviderMock) HandleNewMerkleProof(ctx context.Context, txid *chainhash.Hash, proof *transaction.MerklePath) error {
	m.t.Helper()
	m.called = true
	m.calledWithTxID = txid
	m.calledWithProof = proof

	// Verify block height if expected
	if m.expectedBlockHeight > 0 {
		require.Equal(m.t, m.expectedBlockHeight, proof.BlockHeight, "Block height mismatch")
	}

	return m.error
}

// AssertCalled verifies that HandleNewMerkleProof was called as expected.
func (m *MerkleProofProviderMock) AssertCalled() {
	m.t.Helper()

	// A non-nil error indicates we expect HandleNewMerkleProof to be called
	requireCall := m.error != nil || m.called

	if requireCall && !m.called {
		m.t.Error("Expected HandleNewMerkleProof to be called, but it wasn't")
	}
}

// NewMerkleProofProviderMock creates a new instance of MerkleProofProviderMock with the given error response.
func NewMerkleProofProviderMock(t *testing.T, err error) *MerkleProofProviderMock {
	return &MerkleProofProviderMock{
		t:     t,
		error: err,
	}
}

// NewMerkleProofProviderMockWithBlockHeight creates a new instance of MerkleProofProviderMock with the
// given error response and expected block height for validation.
func NewMerkleProofProviderMockWithBlockHeight(t *testing.T, err error, blockHeight uint32) *MerkleProofProviderMock {
	return &MerkleProofProviderMock{
		t:                   t,
		error:               err,
		expectedBlockHeight: blockHeight,
	}
}
