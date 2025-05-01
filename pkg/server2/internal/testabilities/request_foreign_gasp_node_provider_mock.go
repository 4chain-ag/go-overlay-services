package testabilities

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/assert"
)

// ErrMockRequestForeignGASPNode is the error returned by RequestForeignGASPNodeProviderMock when shouldFail is true.

var ErrMockRequestForeignGASPNode = assert.AnError

// RequestForeignGASPNodeProviderMock is a mock implementation of the RequestForeignGASPNodeProvider interface.

// It's used for testing the RequestForeignGASPNodeHandler.

type RequestForeignGASPNodeProviderMock struct {
	t *testing.T

	shouldFail bool

	called bool
}

// ProvideForeignGASPNode implements the RequestForeignGASPNodeProvider interface.

// It records that the method was called and returns a mock response or error based on shouldFail.

func (m *RequestForeignGASPNodeProviderMock) ProvideForeignGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error) {

	m.called = true

	if m.shouldFail {

		return nil, ErrMockRequestForeignGASPNode

	}

	return &core.GASPNode{}, nil

}

// AssertCalled verifies that ProvideForeignGASPNode was called.

func (m *RequestForeignGASPNodeProviderMock) AssertCalled(t *testing.T) {

	assert.True(t, m.called, "Expected ProvideForeignGASPNode to be called")

}

// NewRequestForeignGASPNodeProviderMock creates a new RequestForeignGASPNodeProviderMock with the specified behavior.

func NewRequestForeignGASPNodeProviderMock(t *testing.T, shouldFail bool) *RequestForeignGASPNodeProviderMock {

	return &RequestForeignGASPNodeProviderMock{

		t: t,

		shouldFail: shouldFail,

		called: false,
	}

}
