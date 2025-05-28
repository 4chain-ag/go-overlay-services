package testabilities

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

// RequestSyncResponseProviderMockExpectations defines mock expectations.
type RequestSyncResponseProviderMockExpectations struct {
	Error                          error
	Response                       *core.GASPInitialResponse
	ProvideForeignSyncResponseCall bool
}

// RequestSyncResponseProviderMock is a mock provider.
type RequestSyncResponseProviderMock struct {
	t            *testing.T
	expectations RequestSyncResponseProviderMockExpectations
	called       bool
}

// MockRequestPayload represents a typical request payload for testing.
type MockRequestPayload struct {
	Version int `json:"version"`
	Since   int `json:"since"`
}

// MockRequestHeaders represents common headers for testing.
type MockRequestHeaders map[string]string

// DefaultRequestSyncResponseProviderMockExpectations provides realistic mock values
// with sample UTXOs that demonstrate typical sync response data.
var DefaultRequestSyncResponseProviderMockExpectations = RequestSyncResponseProviderMockExpectations{
	Error: nil,
	Response: &core.GASPInitialResponse{
		UTXOList: []*overlay.Outpoint{
			{
				Txid:        createMockTxHash("03895fb984362a4196bc9931629318fcbb2aeba7c6293638119ea653fa31d119"),
				OutputIndex: 0,
			},
			{
				Txid:        createMockTxHash("27c8f37851aabc468d3dbb6bf0789dc398a602dcb897ca04e7815d939d621595"),
				OutputIndex: 1,
			},
			{
				Txid:        createMockTxHash("4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"),
				OutputIndex: 2,
			},
		},
		Since: 1234567890,
	},
	ProvideForeignSyncResponseCall: true,
}

// Common mock request scenarios
var (
	// DefaultMockRequestPayload provides a standard request payload for testing.
	DefaultMockRequestPayload = MockRequestPayload{
		Version: 1,
		Since:   100000,
	}

	// DefaultMockHeaders provides standard headers for testing.
	DefaultMockHeaders = MockRequestHeaders{
		"Content-Type": "application/json",
		"X-BSV-Topic":  "test-topic",
	}

	// MissingTopicHeaders simulates missing topic header scenario.
	MissingTopicHeaders = MockRequestHeaders{
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

// NewMockRequestPayload creates a custom request payload for testing.
func NewMockRequestPayload(version, since int) MockRequestPayload {
	return MockRequestPayload{
		Version: version,
		Since:   since,
	}
}

// NewMockHeaders creates custom headers for testing.
func NewMockHeaders(contentType, topic string) MockRequestHeaders {
	headers := MockRequestHeaders{
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

	if m.expectations.Error != nil {
		return nil, m.expectations.Error
	}

	return m.expectations.Response, nil
}

// AssertCalled verifies the method was called as expected.
func (m *RequestSyncResponseProviderMock) AssertCalled() {
	m.t.Helper()
	require.Equal(m.t, m.expectations.ProvideForeignSyncResponseCall, m.called, "Discrepancy between expected and actual ProvideForeignSyncResponseCall")
}

// createMockTxHash creates a chainhash.Hash from a hex string for testing.
func createMockTxHash(hexStr string) chainhash.Hash {
	hash, err := chainhash.NewHashFromHex(hexStr)
	if err != nil {
		panic("invalid hex string for mock transaction hash: " + err.Error())
	}
	return *hash
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

// NewSingleUTXOResponseExpectations creates expectations with a single UTXO for basic testing.
func NewSingleUTXOResponseExpectations() RequestSyncResponseProviderMockExpectations {
	return RequestSyncResponseProviderMockExpectations{
		Error: nil,
		Response: &core.GASPInitialResponse{
			UTXOList: []*overlay.Outpoint{
				{
					Txid:        createMockTxHash("03895fb984362a4196bc9931629318fcbb2aeba7c6293638119ea653fa31d119"),
					OutputIndex: 0,
				},
			},
			Since: 1000000,
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
