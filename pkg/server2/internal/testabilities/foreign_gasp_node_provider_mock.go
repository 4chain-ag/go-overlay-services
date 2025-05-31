package testabilities

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

// Default test values for RequestForeignGASPNode operations.
const (
	DefaultValidGraphID     = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef.0"
	DefaultValidTxID        = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	DefaultValidOutputIndex = uint32(0)
	DefaultValidTopic       = "test-topic"
	DefaultInvalidTxID      = "invalid-txid"
	DefaultInvalidGraphID   = "invalid-graphid"
	DefaultEmptyTopic       = ""
)

// ForeignGASPNodeDefaultDTO provides a default DTO for RequestForeignGASPNode tests.
var ForeignGASPNodeDefaultDTO = app.RequestForeignGASPNodeDTO{
	GraphID:     DefaultValidGraphID,
	TxID:        DefaultValidTxID,
	OutputIndex: DefaultValidOutputIndex,
	Topic:       DefaultValidTopic,
}

// ForeignGASPNodeProviderMockExpectations defines expected behavior for the mock provider,
// including whether a call is expected, what node to return, and any error to simulate.
type ForeignGASPNodeProviderMockExpectations struct {
	Error                      error
	Node                       *core.GASPNode
	ProvideForeignGASPNodeCall bool
}

// ForeignGASPNodeProviderMock is a mock implementation of the RequestForeignGASPNodeProvider interface
// used for testing. It checks whether the expected methods were called and returns predefined results.
type ForeignGASPNodeProviderMock struct {
	t            *testing.T
	expectations ForeignGASPNodeProviderMockExpectations
	called       bool
}

// ProvideForeignGASPNode simulates the provider method. It marks the method as called,
// returns the expected node or error based on the test configuration.
func (m *ForeignGASPNodeProviderMock) ProvideForeignGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error) {
	m.t.Helper()
	m.called = true

	if m.expectations.Error != nil {
		return nil, m.expectations.Error
	}

	return m.expectations.Node, nil
}

// AssertCalled verifies whether ProvideForeignGASPNode was called as expected during the test.
// It fails the test if there's a mismatch.
func (m *ForeignGASPNodeProviderMock) AssertCalled() {
	m.t.Helper()
	require.Equal(m.t, m.expectations.ProvideForeignGASPNodeCall, m.called, "Discrepancy between expected and actual ProvideForeignGASPNode call")
}

// NewForeignGASPNodeProviderMock initializes and returns a new mock provider with the given expectations.
func NewForeignGASPNodeProviderMock(t *testing.T, expectations ForeignGASPNodeProviderMockExpectations) *ForeignGASPNodeProviderMock {
	return &ForeignGASPNodeProviderMock{
		t:            t,
		expectations: expectations,
	}
}
