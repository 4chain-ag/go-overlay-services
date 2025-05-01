package ports_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRequestSyncResponseProvider struct {
	shouldFail bool
	called     bool
}

func (m *mockRequestSyncResponseProvider) ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {
	m.called = true
	if m.shouldFail {
		return nil, errors.New("mock service error")
	}
	return &core.GASPInitialResponse{
		UTXOList: []*overlay.Outpoint{},
		Since:    initialRequest.Since,
	}, nil
}

func (m *mockRequestSyncResponseProvider) AssertCalled(t *testing.T) {
	assert.True(t, m.called, "Expected ProvideForeignSyncResponse to be called")
}

func TestRequestSyncResponseHandler_Handle_Success(t *testing.T) {
	// Given
	mock := &mockRequestSyncResponseProvider{shouldFail: false}
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestSyncResponseProvider(mock))
	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))
	
	payload := core.GASPInitialRequest{
		Version: 1,
		Since:   1000,
	}
	
	// When
	var response core.GASPInitialResponse
	res, _ := fixture.Client().
		R().
		SetHeader(ports.XBSVTopicHeader, "test-topic").
		SetBody(payload).
		SetResult(&response).
		Post("/api/v1/requestSyncResponse")
	
	// Then
	require.Equal(t, http.StatusOK, res.StatusCode())
	assert.Equal(t, uint32(1000), response.Since)
	mock.AssertCalled(t)
}

func TestRequestSyncResponseHandler_Handle_MissingTopic(t *testing.T) {
	// Given
	mock := &mockRequestSyncResponseProvider{shouldFail: false}
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestSyncResponseProvider(mock))
	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))
	
	payload := core.GASPInitialRequest{
		Version: 1,
		Since:   1000,
	}
	
	// When - Missing X-BSV-Topic header
	var errorResponse map[string]interface{}
	res, _ := fixture.Client().
		R().
		SetBody(payload).
		SetResult(&errorResponse).
		Post("/api/v1/requestSyncResponse")
	
	// Then
	require.Equal(t, http.StatusBadRequest, res.StatusCode())
}

func TestRequestSyncResponseHandler_Handle_InvalidJSON(t *testing.T) {
	// Given
	mock := &mockRequestSyncResponseProvider{shouldFail: false}
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestSyncResponseProvider(mock))
	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))
	
	// When
	var errorResponse map[string]interface{}
	res, _ := fixture.Client().
		R().
		SetHeader(ports.XBSVTopicHeader, "test-topic").
		SetBody("invalid json").
		SetResult(&errorResponse).
		Post("/api/v1/requestSyncResponse")
	
	// Then
	require.Equal(t, http.StatusBadRequest, res.StatusCode())
}

func TestRequestSyncResponseHandler_Handle_ServiceError(t *testing.T) {
	// Given
	mock := &mockRequestSyncResponseProvider{shouldFail: true}
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithRequestSyncResponseProvider(mock))
	fixture := server2.NewTestFixture(t, server2.WithEngine(engine))
	
	payload := core.GASPInitialRequest{
		Version: 1,
		Since:   1000,
	}
	
	// When
	var errorResponse map[string]interface{}
	res, _ := fixture.Client().
		R().
		SetHeader(ports.XBSVTopicHeader, "test-topic").
		SetBody(payload).
		SetResult(&errorResponse).
		Post("/api/v1/requestSyncResponse")
	
	// Then
	require.Equal(t, http.StatusInternalServerError, res.StatusCode())
	mock.AssertCalled(t)
} 
