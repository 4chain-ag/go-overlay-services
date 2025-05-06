package testabilities

import (
	"testing"

	testvectors "github.com/bsv-blockchain/universal-test-vectors/pkg/testabilities"
	"github.com/stretchr/testify/require"
)

// DummyTxBEEF returns a valid transaction serialized in BEEF format for use in tests.
// It creates a dummy transaction with predefined input and output values.
// The test fails immediately if the transaction cannot be serialized or results in an empty byte slice.
func DummyTxBEEF(t *testing.T) []byte {
	t.Helper()

	dummyTx := testvectors.GivenTX().
		WithInput(1000).
		WithP2PKHOutput(999).
		TX()

	bb, err := dummyTx.BEEF()
	require.NoError(t, err)
	require.NotEmpty(t, bb)
	return bb
}
