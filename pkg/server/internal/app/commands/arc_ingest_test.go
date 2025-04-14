package commands_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/commands"
	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockMerkleProofProvider bypasses actual merkle path validation for testing purposes.
type MockMerkleProofProvider struct {
	ExpectedError error
	WasCalled     bool
	CalledTxid    *chainhash.Hash
	CalledProof   *transaction.MerklePath
}

// HandleNewMerkleProof just records the call; returns ExpectedError if set.
func (m *MockMerkleProofProvider) HandleNewMerkleProof(ctx context.Context, txid *chainhash.Hash, proof *transaction.MerklePath) error {
	m.WasCalled = true
	m.CalledTxid = txid
	m.CalledProof = proof
	return m.ExpectedError
}

// Verify the type interfaces for the mock.
var _ commands.NewMerkleProofProvider = (*MockMerkleProofProvider)(nil)

// Test data for merkle proofs
// BEEF data containing transaction with merkle path information
const beefTransactionData = "0100beef01fef4f10c000902fd020100ab2d9b3bbfc2ecf5c834f7719bf10dcabef31cfa3b78fa714bd3a2b3958fc6b5fd0301029515629d935d81e704ca97b8dc02a698c39d78f06bbb3d8d46bcaa5178f3c827018000fc94b1da08b5d0850afcc4948a9e129a2c2c0bb629090ca84fef5d326ceadbc2014100fa78356192a22189dd3c2a03bfb2e043667a55a85386a8d0f0853bbaa54ba748012100e37888006acc37cbfc6a459117c3e957b00e8ff3a90d93bc9aaa3fc9d972453701110052564e444cbe4b1f24015d91f5cfd7eff63e5f403613a7b48e77efa32761f950010900acef65445ef232f1fcb527bab67fccc90ae9d0c610a4f752d84c5e10a3415d74010500e1fadd176551150808ae96ff2b4dff487c432788d23db25fd3fb5e8b3e16b96f0103003e530f146a992ed9378807fa43bc1e6033de0711d728ee3a0ea1507b5ffff2940100006c8eca8a9d680ab2ae4dc1f0e247f282de3f32b99175e6f9ddb4f8e76ae0bf1b010100000001792c8c563b17721744d8f101d79e873aae62816d0fd47dd9d90ada219c1b252d070000006b48304502210090368b5925795abb0415a3cd6265d7abae7d3541fb7f2fdb1fdffdd916b3fce3022068136da44df8ee367807c7c59c950ea8d428adfe142268cd6371cf4c1b4c3208412102858ea1804708d1e16e77b49ba1559b3742797f5f43efd3922fa3d21557521d7cffffffff03f401000000000000fde10121026afd108ddbc05e093207ebdf06557a3d66e15b96accbf1766c2e7d6efa3d1981ac22313852713463583334356b6a7461734467764b41455a435135564c5a72516f76786d2081d846bcae65ce281829b6d81eff3862a68b8eaf4e99ad97c0c4ad2e149631c020f936ace2d2b4034710dec819e343f69ddbb1a16d721b88413b697c2475f75f504c4e8470062e6966cb7ec80864e5e1217ecf2b3f113bb5307e1ec1d338678e54b426a11a267e990f7dcd85072c7b3870fabb54812d7f8a7aaa672ce01e4f8a2c6a1d532849877329bd23be0b3383fcb437ffd1ff05aad824343bfe5bae52d5775ed99ba9d5a72038dd50f81153e674a8b1c5d659aac361e660cc5657b291eda149de2505efb0bbac0a353331393633353238334c797f3040f01934bb25733d4700a0309194a7defd8ceea7a8a59ad279c6237e2917264a115773777c5746546288b15ec1cf6788b801dd3952fec804cadbba2ff4e05464735e2db83c3509b59119efc843716c34becd2884f011f2e7f61f382237891875115b6ea5bac1c3e9ca6c575b44547d589c38c1b5a8cadb463044022016e1fdb226b465f3151aaa1e033b122dd58f3b495d2846d532f784ec8f6fc9ae022031dc12155c02e7f588c9f6d975a79c9c26cf341c1c88578cd6a87a3c01f6a73d6d6d6d6dc8000000000000001976a914010b7bfa8ff585b6d95f8a18a998cbd87508bd4188aca9010000000000001976a914254c6f1ce67d27ae4dafb03d2ba08318df2883c588ac000000000100"

// Valid txid data containing transaction with merkle path information
const validTestTxid = "27c8f37851aabc468d3dbb6bf0789dc398a602dcb897ca04e7815d939d621595"

// Calculate merkle path from transaction
func getMerklePathFromTransaction() (string, error) {
	beefBytes, err := hex.DecodeString(beefTransactionData)
	if err != nil {
		return "", fmt.Errorf("failed to decode BEEF data: %w", err)
	}
	tx, err := transaction.NewTransactionFromBEEF(beefBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse transaction: %w", err)
	}
	return tx.MerklePath.Hex(), nil
}

// Get the valid merkle path for tests
func getValidTestMerklePath(t *testing.T) string {
	t.Helper()
	path, err := getMerklePathFromTransaction()
	if err != nil {
		t.Fatalf("Failed to get merkle path: %v", err)
	}
	return path
}

// newTestArcIngestHandler is a helper to create the ArcIngestHandler with a mock provider.
func newTestArcIngestHandler(provider commands.NewMerkleProofProvider, opts ...commands.ArcIngestHandlerOption) (*commands.ArcIngestHandler, error) {
	return commands.NewArcIngestHandler(provider, opts...)
}

// TestArcIngestHandler_Handle_SuccessfulRequest verifies that a valid request
// returns 200 OK and calls the provider with the correct parameters.
func TestArcIngestHandler_Handle_SuccessfulRequest(t *testing.T) {
	// Given:
	mock := &MockMerkleProofProvider{ExpectedError: nil}
	handler, err := newTestArcIngestHandler(mock)
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	requestBody := commands.ArcIngestRequest{
		Txid:        validTestTxid,
		MerklePath:  getValidTestMerklePath(t),
		BlockHeight: 848372,
	}

	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewBuffer(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// When:
	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	// Then:
	require.Equal(t, http.StatusOK, res.StatusCode)
	var response commands.ArcIngestHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &response))
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Transaction status updated", response.Message)
	assert.True(t, mock.WasCalled, "Provider should have been called on a valid request")
	require.NotNil(t, mock.CalledTxid)
	require.NotNil(t, mock.CalledProof)
	assert.Equal(t, requestBody.BlockHeight, mock.CalledProof.BlockHeight)
}

// TestArcIngestHandler_Handle_InvalidMethod ensures that non-POST methods are rejected.
func TestArcIngestHandler_Handle_InvalidMethod(t *testing.T) {
	// Given:
	mock := &MockMerkleProofProvider{ExpectedError: nil}
	handler, err := newTestArcIngestHandler(mock)
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()
	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(t, err)

	// When:
	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	// Then:
	require.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	var response commands.ArcIngestHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &response))
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, commands.ErrArcIngestInvalidHTTPMethod.Error(), response.Message)
	assert.False(t, mock.WasCalled, "Provider should not be called for invalid method")
}

// TestArcIngestHandler_Handle_InvalidRequestBody checks if a non-JSON body triggers an error.
func TestArcIngestHandler_Handle_InvalidRequestBody(t *testing.T) {
	// Given:
	mock := &MockMerkleProofProvider{ExpectedError: nil}
	handler, err := newTestArcIngestHandler(mock)
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()
	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewBuffer([]byte("invalid json")))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// When:
	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	// Then:
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	var response commands.ArcIngestHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &response))
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Message, "invalid request body")
	assert.False(t, mock.WasCalled, "Provider should not be called if JSON is invalid")
}

// TestArcIngestHandler_Handle_MissingRequiredFields ensures if Txid or MerklePath
// is empty, we return 400 and do not call the provider.
func TestArcIngestHandler_Handle_MissingRequiredFields(t *testing.T) {
	// Given:
	mock := &MockMerkleProofProvider{}
	handler, err := newTestArcIngestHandler(mock)
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// Missing MerklePath
	requestBody := commands.ArcIngestRequest{
		Txid:        validTestTxid,
		BlockHeight: 848372,
	}

	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewBuffer(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// When:
	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	// Then:
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	var response commands.ArcIngestHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &response))
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, commands.ErrMissingRequiredFields.Error(), response.Message)
	assert.False(t, mock.WasCalled, "Provider should not be called when required fields are missing")
}

// TestArcIngestHandler_Handle_EngineError ensures that if the provider returns an error,
// it surfaces as a 500 "Failed to process merkle proof."
func TestArcIngestHandler_Handle_EngineError(t *testing.T) {
	// Given:
	expectedErr := errors.New("merkle proof processing error")
	mock := &MockMerkleProofProvider{ExpectedError: expectedErr}
	handler, err := newTestArcIngestHandler(mock)
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	requestBody := commands.ArcIngestRequest{
		Txid:        validTestTxid,
		MerklePath:  getValidTestMerklePath(t),
		BlockHeight: 848372,
	}

	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewBuffer(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// When:
	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	// Then:
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	var response commands.ArcIngestHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &response))
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Message, "Failed to process merkle proof")
	assert.True(t, mock.WasCalled, "Provider should be called on valid request despite returning error")
}

// TestArcIngestHandler_Handle_InvalidTxidFormat verifies that a non-hex or wrong-length TxID returns 400.
func TestArcIngestHandler_Handle_InvalidTxidFormat(t *testing.T) {
	mock := &MockMerkleProofProvider{}
	handler, err := newTestArcIngestHandler(mock)
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	invalidTxid := "zzzzzzzzzzzzzzzz"
	requestBody := commands.ArcIngestRequest{
		Txid:        invalidTxid,
		MerklePath:  getValidTestMerklePath(t),
		BlockHeight: 848372,
	}

	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewBuffer(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// When:
	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	// Then:
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	var response commands.ArcIngestHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &response))
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Message, "invalid txid format")
	assert.False(t, mock.WasCalled, "Provider should not be called if TxID is invalid")
}

func TestArcIngestHandler_Handle_InvalidTxidLength(t *testing.T) {
	mock := &MockMerkleProofProvider{}
	handler, err := newTestArcIngestHandler(mock)
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	shortTxid := "0123456789"
	requestBody := commands.ArcIngestRequest{
		Txid:        shortTxid,
		MerklePath:  getValidTestMerklePath(t),
		BlockHeight: 848372,
	}

	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewBuffer(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// When:
	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	// Then:
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	var response commands.ArcIngestHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &response))
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Message, "invalid txid format")
	assert.False(t, mock.WasCalled, "Provider should not be called if TxID length is invalid")
}

// TestArcIngestHandler_Handle_InvalidMerklePathFormat checks that if the merkle path
// fails to decode from hex, we get a 400 error.
func TestArcIngestHandler_Handle_InvalidMerklePathFormat(t *testing.T) {
	mock := &MockMerkleProofProvider{}
	handler, err := newTestArcIngestHandler(mock)
	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	invalidPath := "ZZZZZZ"
	requestBody := commands.ArcIngestRequest{
		Txid:        validTestTxid,
		MerklePath:  invalidPath,
		BlockHeight: 848372,
	}

	bodyBytes, err := json.Marshal(requestBody)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewBuffer(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// When:
	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	// Then:
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	var response commands.ArcIngestHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &response))
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Message, "invalid merkle path format")
	assert.False(t, mock.WasCalled)
}
