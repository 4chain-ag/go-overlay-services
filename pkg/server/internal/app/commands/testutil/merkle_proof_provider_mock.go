package testutil

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/stretchr/testify/require"
)

// MerkleProofProviderMock is a test double for a Merkle proof handler.
// It records the inputs it was called with and returns a preconfigured error.
type MerkleProofProviderMock struct {
	// Expected values
	err            error
	expectedHeight uint32

	// Call tracking
	called      bool
	calledTxID  *chainhash.Hash
	calledProof *transaction.MerklePath
}

// NewMerkleProofProviderMock creates a new MerkleProofProviderMock
// with the expected error and block height configured.
func NewMerkleProofProviderMock(err error, expectedHeight uint32) *MerkleProofProviderMock {
	return &MerkleProofProviderMock{
		err:            err,
		expectedHeight: expectedHeight,
	}
}

// HandleNewMerkleProof records the input parameters and returns the configured error.
func (m *MerkleProofProviderMock) HandleNewMerkleProof(
	ctx context.Context,
	txID *chainhash.Hash,
	merklePath *transaction.MerklePath,
) error {
	m.called = true
	m.calledTxID = txID
	m.calledProof = merklePath
	return m.err
}

// AssertCalled verifies that HandleNewMerkleProof was called with non-nil values,
// and that the Merkle path had the expected block height.
func (m *MerkleProofProviderMock) AssertCalled(t *testing.T) {
	t.Helper()

	require.True(t, m.called, "HandleNewMerkleProof was not called")
	require.NotNil(t, m.calledTxID, "txID argument was nil")
	require.NotNil(t, m.calledProof, "merklePath argument was nil")
	require.Equal(t, m.expectedHeight, m.calledProof.BlockHeight, "unexpected block height")
}

// Error returns the expected error to be used by the mock.
func (m *MerkleProofProviderMock) Error() error {
	return m.err
}

// ExpectedBlockHeight returns the configured block height for assertion.
func (m *MerkleProofProviderMock) ExpectedBlockHeight() uint32 {
	return m.expectedHeight
}

// Test data for merkle proofs
// BEEF data containing transaction with merkle path information
const BEEFTransactionData = "0100beef01fef4f10c000902fd020100ab2d9b3bbfc2ecf5c834f7719bf10dcabef31cfa3b78fa714bd3a2b3958fc6b5fd0301029515629d935d81e704ca97b8dc02a698c39d78f06bbb3d8d46bcaa5178f3c827018000fc94b1da08b5d0850afcc4948a9e129a2c2c0bb629090ca84fef5d326ceadbc2014100fa78356192a22189dd3c2a03bfb2e043667a55a85386a8d0f0853bbaa54ba748012100e37888006acc37cbfc6a459117c3e957b00e8ff3a90d93bc9aaa3fc9d972453701110052564e444cbe4b1f24015d91f5cfd7eff63e5f403613a7b48e77efa32761f950010900acef65445ef232f1fcb527bab67fccc90ae9d0c610a4f752d84c5e10a3415d74010500e1fadd176551150808ae96ff2b4dff487c432788d23db25fd3fb5e8b3e16b96f0103003e530f146a992ed9378807fa43bc1e6033de0711d728ee3a0ea1507b5ffff2940100006c8eca8a9d680ab2ae4dc1f0e247f282de3f32b99175e6f9ddb4f8e76ae0bf1b010100000001792c8c563b17721744d8f101d79e873aae62816d0fd47dd9d90ada219c1b252d070000006b48304502210090368b5925795abb0415a3cd6265d7abae7d3541fb7f2fdb1fdffdd916b3fce3022068136da44df8ee367807c7c59c950ea8d428adfe142268cd6371cf4c1b4c3208412102858ea1804708d1e16e77b49ba1559b3742797f5f43efd3922fa3d21557521d7cffffffff03f401000000000000fde10121026afd108ddbc05e093207ebdf06557a3d66e15b96accbf1766c2e7d6efa3d1981ac22313852713463583334356b6a7461734467764b41455a435135564c5a72516f76786d2081d846bcae65ce281829b6d81eff3862a68b8eaf4e99ad97c0c4ad2e149631c020f936ace2d2b4034710dec819e343f69ddbb1a16d721b88413b697c2475f75f504c4e8470062e6966cb7ec80864e5e1217ecf2b3f113bb5307e1ec1d338678e54b426a11a267e990f7dcd85072c7b3870fabb54812d7f8a7aaa672ce01e4f8a2c6a1d532849877329bd23be0b3383fcb437ffd1ff05aad824343bfe5bae52d5775ed99ba9d5a72038dd50f81153e674a8b1c5d659aac361e660cc5657b291eda149de2505efb0bbac0a353331393633353238334c797f3040f01934bb25733d4700a0309194a7defd8ceea7a8a59ad279c6237e2917264a115773777c5746546288b15ec1cf6788b801dd3952fec804cadbba2ff4e05464735e2db83c3509b59119efc843716c34becd2884f011f2e7f61f382237891875115b6ea5bac1c3e9ca6c575b44547d589c38c1b5a8cadb463044022016e1fdb226b465f3151aaa1e033b122dd58f3b495d2846d532f784ec8f6fc9ae022031dc12155c02e7f588c9f6d975a79c9c26cf341c1c88578cd6a87a3c01f6a73d6d6d6d6dc8000000000000001976a914010b7bfa8ff585b6d95f8a18a998cbd87508bd4188aca9010000000000001976a914254c6f1ce67d27ae4dafb03d2ba08318df2883c588ac000000000100"

// Valid txid data containing transaction with merkle path information
const ValidTxId = "27c8f37851aabc468d3dbb6bf0789dc398a602dcb897ca04e7815d939d621595"

// NewValidTestMerklePath returns a valid Merkle path hex string for use in tests.
// It decodes a BEEF-encoded transaction and extracts the Merkle path.
func NewValidTestMerklePath(t *testing.T) string {
	t.Helper()

	data, err := hex.DecodeString(BEEFTransactionData)
	require.NoError(t, err, "failed to decode BEEFTransactionData")
	require.NotEmpty(t, data, "decoded transaction data is empty")

	tx, err := transaction.NewTransactionFromBEEF(data)
	require.NoError(t, err, "failed to parse transaction from BEEF format")
	require.NotNil(t, tx, "transaction is nil")

	return tx.MerklePath.Hex()
}
