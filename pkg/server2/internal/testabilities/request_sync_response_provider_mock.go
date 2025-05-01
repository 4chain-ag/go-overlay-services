package testabilities

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/assert"
)

// ErrMockRequestSyncResponse is the error returned by RequestSyncResponseProviderMock when shouldFail is true.

var ErrMockRequestSyncResponse = assert.AnError

// RequestSyncResponseProviderMock is a mock implementation of the RequestSyncResponseProvider interface.

// It's used for testing the RequestSyncResponseHandler.

type RequestSyncResponseProviderMock struct {
	t *testing.T

	shouldFail bool

	called bool
}

// ProvideForeignSyncResponse implements the RequestSyncResponseProvider interface.

// It records that the method was called and returns a mock response or error based on shouldFail.

func (m *RequestSyncResponseProviderMock) ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {

	m.called = true

	if m.shouldFail {

		return nil, ErrMockRequestSyncResponse

	}

	return &core.GASPInitialResponse{

		UTXOList: []*overlay.Outpoint{},

		Since: initialRequest.Since,
	}, nil

}

// AssertCalled verifies that ProvideForeignSyncResponse was called.

func (m *RequestSyncResponseProviderMock) AssertCalled(t *testing.T) {

	assert.True(t, m.called, "Expected ProvideForeignSyncResponse to be called")

}

// NewRequestSyncResponseProviderMock creates a new RequestSyncResponseProviderMock with the specified behavior.

func NewRequestSyncResponseProviderMock(t *testing.T, shouldFail bool) *RequestSyncResponseProviderMock {

	return &RequestSyncResponseProviderMock{

		t: t,

		shouldFail: shouldFail,

		called: false,
	}

}
