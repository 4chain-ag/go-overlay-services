package testabilities

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

// RequestSyncResponseProviderMockExpectations defines mock expectations.
type RequestSyncResponseProviderMockExpectations struct {
	Error                          error
	Response                       *core.GASPInitialResponse
	ProvideForeignSyncResponseCall bool
	InitialRequest                 *core.GASPInitialRequest
	Topic                          string
}

// RequestSyncResponseProviderMock is a mock provider.
type RequestSyncResponseProviderMock struct {
	t              *testing.T
	expectations   RequestSyncResponseProviderMockExpectations
	called         bool
	topic          string
	initialRequest *core.GASPInitialRequest
}

const (
	DefaultVersion = 1
	DefaultSince   = 100000
	DefaultTopic   = "test-topic"
)

func NewDefaultGASPInitialResponseTestHelper(t *testing.T) *core.GASPInitialResponse {
	t.Helper()

	return &core.GASPInitialResponse{
		UTXOList: []*overlay.Outpoint{
			{
				Txid:        *DummyTxHash(t, "03895fb984362a4196bc9931629318fcbb2aeba7c6293638119ea653fa31d119"),
				OutputIndex: 0,
			},
		},
		Since: 1000000,
	}
}

// NewDefaultRequestSyncResponseBody creates a new request sync response body.
func NewDefaultRequestSyncResponseBody() openapi.RequestSyncResponseBody {
	return openapi.RequestSyncResponseBody{
		Version: DefaultVersion,
		Since:   DefaultSince,
	}
}

// NewRequestSyncResponseProviderMock creates a new mock provider.
func NewRequestSyncResponseProviderMock(t *testing.T, expectations RequestSyncResponseProviderMockExpectations) *RequestSyncResponseProviderMock {
	return &RequestSyncResponseProviderMock{
		t:            t,
		expectations: expectations,
	}
}

// ProvideForeignSyncResponse mocks the method.
func (m *RequestSyncResponseProviderMock) ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {
	m.t.Helper()
	m.called = true
	m.topic = topic
	m.initialRequest = initialRequest

	if m.expectations.Error != nil {
		return nil, m.expectations.Error
	}

	return m.expectations.Response, nil
}

// AssertCalled verifies the method was called as expected.
func (m *RequestSyncResponseProviderMock) AssertCalled() {
	m.t.Helper()
	require.Equal(m.t, m.expectations.ProvideForeignSyncResponseCall, m.called, "Discrepancy between expected and actual ProvideForeignSyncResponseCall")
	require.Equal(m.t, m.expectations.InitialRequest, m.initialRequest, "Discrepancy between expected and actual InitialRequest")
	require.Equal(m.t, m.expectations.Topic, m.topic, "Discrepancy between expected and actual Topic")
}
