package testabilities

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

// RequestForeignGASPNodeProviderMockExpectations defines the expected behavior of the mock provider.
type RequestForeignGASPNodeProviderMockExpectations struct {
	// Error is the error to return from ProvideForeignGASPNode.
	Error error

	// Node is the GASP node to return from ProvideForeignGASPNode.
	Node *core.GASPNode

	// ProvideForeignGASPNodeCall indicates whether the method should be called.
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
