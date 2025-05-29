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

var (
	// DefaultMockRequestPayload provides a standard request payload using OpenAPI model.
	DefaultMockRequestPayload = openapi.RequestSyncResponseBody{
		Version: 1,
		Since:   100000,
	}

	// DefaultMockHeaders provides standard headers for testing.
	DefaultMockHeaders = map[string]string{
		"Content-Type": "application/json",
		"X-BSV-Topic":  "test-topic",
	}

	// MissingTopicHeaders simulates missing topic header scenario.
	MissingTopicHeaders = map[string]string{
		"Content-Type": "application/json",
	}
)

// NewRequestSyncResponseProviderMock creates a new mock provider.
func NewRequestSyncResponseProviderMock(t *testing.T, expectations RequestSyncResponseProviderMockExpectations) *RequestSyncResponseProviderMock {
	return &RequestSyncResponseProviderMock{
		t:            t,
		expectations: expectations,
	}
}

// NewMockRequestPayload creates a custom request payload using OpenAPI model.
func NewMockRequestPayload(version, since int) openapi.RequestSyncResponseBody {
	return openapi.RequestSyncResponseBody{
		Version: version,
		Since:   since,
	}
}

// NewMockHeaders creates custom headers for testing.
func NewMockHeaders(contentType, topic string) map[string]string {
	headers := map[string]string{
		"Content-Type": contentType,
	}
	if topic != "" {
		headers["X-BSV-Topic"] = topic
	}
	return headers
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

// NewEmptyResponseExpectations creates expectations for an empty UTXO list response.
func NewEmptyResponseExpectations() RequestSyncResponseProviderMockExpectations {
	return RequestSyncResponseProviderMockExpectations{
		Error: nil,
		Response: &core.GASPInitialResponse{
			UTXOList: []*overlay.Outpoint{},
			Since:    0,
		},
		ProvideForeignSyncResponseCall: true,
	}
}

// NewErrorResponseExpectations creates expectations that return an error.
func NewErrorResponseExpectations(err error) RequestSyncResponseProviderMockExpectations {
	return RequestSyncResponseProviderMockExpectations{
		Error:                          err,
		Response:                       nil,
		ProvideForeignSyncResponseCall: true,
	}
}

// NewCustomUTXOListExpectations creates expectations with a custom list of UTXOs.
func NewCustomUTXOListExpectations(utxos []*overlay.Outpoint, since uint32) RequestSyncResponseProviderMockExpectations {
	return RequestSyncResponseProviderMockExpectations{
		Error: nil,
		Response: &core.GASPInitialResponse{
			UTXOList: utxos,
			Since:    since,
		},
		ProvideForeignSyncResponseCall: true,
	}
}
