package ports_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock implementation that satisfies the RequestForeignGASPNodeProvider interface
type mockRequestForeignGASPNodeProvider struct {
	shouldFail bool
	called     bool
}

func (m *mockRequestForeignGASPNodeProvider) ProvideForeignGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error) {
	m.called = true
	if m.shouldFail {
		return nil, errors.New("mock service error")
	}
	return &core.GASPNode{}, nil
}

func (m *mockRequestForeignGASPNodeProvider) AssertCalled(t *testing.T) {
	assert.True(t, m.called, "Expected ProvideForeignGASPNode to be called")
}

func TestRequestForeignGASPNodeHandler_Handle_Success(t *testing.T) {
	// Given
	mock := &mockRequestForeignGASPNodeProvider{shouldFail: false}
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestForeignGASPNodeProvider(mock))
	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))

	payload := ports.RequestForeignGASPNodePayload{
		GraphID:     "0000000000000000000000000000000000000000000000000000000000000000.1",
		TxID:        "0000000000000000000000000000000000000000000000000000000000000000",
		OutputIndex: 1,
	}

	// When
	var response core.GASPNode
	res, _ := fixture.Client().
		R().
		SetHeader("X-BSV-Topic", "test-topic").
		SetBody(payload).
		SetResult(&response).
		Post("/api/v1/requestForeignGASPNode")

	// Then
	require.Equal(t, http.StatusOK, res.StatusCode())
	mock.AssertCalled(t)
}

func TestRequestForeignGASPNodeHandler_Handle_MissingTopic(t *testing.T) {
	// Given
	mock := &mockRequestForeignGASPNodeProvider{shouldFail: false}
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestForeignGASPNodeProvider(mock))
	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))

	payload := ports.RequestForeignGASPNodePayload{
		GraphID:     "0000000000000000000000000000000000000000000000000000000000000000.1",
		TxID:        "0000000000000000000000000000000000000000000000000000000000000000",
		OutputIndex: 1,
	}

	// When - Missing X-BSV-Topic header
	var errorResponse openapi.Error
	res, _ := fixture.Client().
		R().
		SetBody(payload).
		SetError(&errorResponse).
		Post("/api/v1/requestForeignGASPNode")

	// Then
	require.Equal(t, http.StatusBadRequest, res.StatusCode())
	assert.Contains(t, errorResponse.Message, "Missing 'X-BSV-Topic' header")
}

func TestRequestForeignGASPNodeHandler_Handle_InvalidJSON(t *testing.T) {
	// Given
	mock := &mockRequestForeignGASPNodeProvider{shouldFail: false}
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestForeignGASPNodeProvider(mock))
	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))

	// When
	var errorResponse openapi.Error
	res, _ := fixture.Client().
		R().
		SetHeader("X-BSV-Topic", "test-topic").
		SetBody("invalid json").
		SetError(&errorResponse).
		Post("/api/v1/requestForeignGASPNode")

	// Then
	require.Equal(t, http.StatusBadRequest, res.StatusCode())
	assert.Contains(t, errorResponse.Message, "Invalid request body")
}

func TestRequestForeignGASPNodeHandler_Handle_InvalidTxID(t *testing.T) {
	// Given
	mock := &mockRequestForeignGASPNodeProvider{shouldFail: false}
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestForeignGASPNodeProvider(mock))
	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))

	payload := ports.RequestForeignGASPNodePayload{
		GraphID:     "0000000000000000000000000000000000000000000000000000000000000000.1",
		TxID:        "invalid-txid",
		OutputIndex: 1,
	}

	// When
	var errorResponse openapi.Error
	res, _ := fixture.Client().
		R().
		SetHeader("X-BSV-Topic", "test-topic").
		SetBody(payload).
		SetError(&errorResponse).
		Post("/api/v1/requestForeignGASPNode")

	// Then
	require.Equal(t, http.StatusBadRequest, res.StatusCode())
	assert.Contains(t, errorResponse.Message, "Invalid txID format")
}

func TestRequestForeignGASPNodeHandler_Handle_ServiceError(t *testing.T) {
	// Given
	mock := &mockRequestForeignGASPNodeProvider{shouldFail: true}
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestForeignGASPNodeProvider(mock))
	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))

	payload := ports.RequestForeignGASPNodePayload{
		GraphID:     "0000000000000000000000000000000000000000000000000000000000000000.1",
		TxID:        "0000000000000000000000000000000000000000000000000000000000000000",
		OutputIndex: 1,
	}

	// When
	var errorResponse openapi.Error
	res, _ := fixture.Client().
		R().
		SetHeader("X-BSV-Topic", "test-topic").
		SetBody(payload).
		SetError(&errorResponse).
		Post("/api/v1/requestForeignGASPNode")

	// Then
	require.Equal(t, http.StatusInternalServerError, res.StatusCode())
	assert.Contains(t, errorResponse.Message, "Unable to process foreign GASP node request")
	mock.AssertCalled(t)
}
