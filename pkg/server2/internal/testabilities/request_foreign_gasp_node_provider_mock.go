package testabilities

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

// Default test values for RequestForeignGASPNode operations
const (
	DefaultValidGraphID     = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef.0"
	DefaultValidTxID        = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	DefaultValidOutputIndex = uint32(0)
	DefaultValidTopic       = "test-topic"
	DefaultInvalidTxID      = "invalid-txid"
	DefaultInvalidGraphID   = "invalid-graphid"
	DefaultEmptyTopic       = ""
)

// Default expectations for successful RequestForeignGASPNode operations
var DefaultRequestForeignGASPNodeProviderMockExpectations = RequestForeignGASPNodeProviderMockExpectations{
	ProvideForeignGASPNodeCall: true,
	Error:                      nil,
	Node:                       &core.GASPNode{},
}

// Default DTO for RequestForeignGASPNode operations
var RequestForeignGASPNodeDefaultDTO = app.RequestForeignGASPNodeDTO{
	GraphID:     DefaultValidGraphID,
	TxID:        DefaultValidTxID,
	OutputIndex: DefaultValidOutputIndex,
	Topic:       DefaultValidTopic,
}

// RequestForeignGASPNodeProviderMockExpectations defines the expected behavior of the mock provider.
type RequestForeignGASPNodeProviderMockExpectations struct {
	Error                      error
	Node                       *core.GASPNode
	ProvideForeignGASPNodeCall bool
}

// RequestForeignGASPNodeProviderMock is a mock implementation for testing.
type RequestForeignGASPNodeProviderMock struct {
	t            *testing.T
	expectations RequestForeignGASPNodeProviderMockExpectations
	called       bool
}

// ProvideForeignGASPNode mocks the ProvideForeignGASPNode method.
func (m *RequestForeignGASPNodeProviderMock) ProvideForeignGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error) {
	m.t.Helper()
	m.called = true

	if m.expectations.Error != nil {
		return nil, m.expectations.Error
	}

	return m.expectations.Node, nil
}

// AssertCalled verifies the method was called as expected.
func (m *RequestForeignGASPNodeProviderMock) AssertCalled() {
	m.t.Helper()
	require.Equal(m.t, m.expectations.ProvideForeignGASPNodeCall, m.called, "Discrepancy between expected and actual ProvideForeignGASPNode call")
}

// NewRequestForeignGASPNodeProviderMock creates a new mock provider.
func NewRequestForeignGASPNodeProviderMock(t *testing.T, expectations RequestForeignGASPNodeProviderMockExpectations) *RequestForeignGASPNodeProviderMock {
	return &RequestForeignGASPNodeProviderMock{
		t:            t,
		expectations: expectations,
	}
}
