package testabilities

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/stretchr/testify/require"
)

// RequestSyncResponseProviderMockExpectations defines mock expectations.

type RequestSyncResponseProviderMockExpectations struct {

	// Error is the error to return.

	Error error

	// Response is the response to return.

	Response *core.GASPInitialResponse

	// ProvideForeignSyncResponseCall indicates if method should be called.

	ProvideForeignSyncResponseCall bool
}

// RequestSyncResponseProviderMock is a mock provider.

type RequestSyncResponseProviderMock struct {
	t *testing.T

	expectations RequestSyncResponseProviderMockExpectations

	called bool
}

// ProvideForeignSyncResponse mocks the method.

func (m *RequestSyncResponseProviderMock) ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {

	m.t.Helper()

	m.called = true

	if m.expectations.Error != nil {

		return nil, m.expectations.Error

	}

	return m.expectations.Response, nil

}

// AssertCalled verifies the method was called as expected.

func (m *RequestSyncResponseProviderMock) AssertCalled() {

	m.t.Helper()

	require.Equal(m.t, m.expectations.ProvideForeignSyncResponseCall, m.called)

}

// NewRequestSyncResponseProviderMock creates a new mock provider.

func NewRequestSyncResponseProviderMock(t *testing.T, expectations RequestSyncResponseProviderMockExpectations) *RequestSyncResponseProviderMock {

	return &RequestSyncResponseProviderMock{

		t: t,

		expectations: expectations,
	}

}
