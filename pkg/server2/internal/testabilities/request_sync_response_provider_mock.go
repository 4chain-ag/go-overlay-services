package testabilities

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/stretchr/testify/require"
)

// RequestSyncResponseProviderMockExpectations defines mock expectations.

type RequestSyncResponseProviderMockExpectations struct {
	Error error

	Response *core.GASPInitialResponse

	ProvideForeignSyncResponseCall bool
}

var DefaultRequestSyncResponseProviderMockExpectations = RequestSyncResponseProviderMockExpectations{

	Error: nil,

	Response: &core.GASPInitialResponse{},

	ProvideForeignSyncResponseCall: true,
}

// RequestSyncResponseProviderMock is a mock provider.

type RequestSyncResponseProviderMock struct {
	t *testing.T

	expectations RequestSyncResponseProviderMockExpectations

	requestSyncResponseCall bool
}

// ProvideForeignSyncResponse mocks the method.

func (m *RequestSyncResponseProviderMock) ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {

	m.t.Helper()

	m.requestSyncResponseCall = true

	if m.expectations.Error != nil {

		return nil, m.expectations.Error

	}

	return m.expectations.Response, nil

}

// AssertCalled verifies the method was called as expected.

func (m *RequestSyncResponseProviderMock) AssertCalled() {

	m.t.Helper()

	require.Equal(m.t, m.expectations.ProvideForeignSyncResponseCall, m.requestSyncResponseCall, "Discrepancy between expected and actual ProvideForeignSyncResponseCall")

}

// NewRequestSyncResponseProviderMock creates a new mock provider.

func NewRequestSyncResponseProviderMock(t *testing.T, expectations RequestSyncResponseProviderMockExpectations) *RequestSyncResponseProviderMock {

	return &RequestSyncResponseProviderMock{

		t: t,

		expectations: expectations,
	}

}
