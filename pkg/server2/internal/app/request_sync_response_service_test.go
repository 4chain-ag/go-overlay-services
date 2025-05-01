package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/stretchr/testify/assert"
)

type mockRequestSyncResponseProvider struct {
	shouldFail bool
}

func (m *mockRequestSyncResponseProvider) ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {
	if m.shouldFail {
		return nil, errors.New("mock sync response error")
	}
	return &core.GASPInitialResponse{
		UTXOList: nil,
		Since:    initialRequest.Since,
	}, nil
}

func TestRequestSyncResponseService_RequestSyncResponse_Success(t *testing.T) {
	// Given
	mockProvider := &mockRequestSyncResponseProvider{shouldFail: false}
	service := app.NewRequestSyncResponseService(mockProvider)
	initialRequest := &core.GASPInitialRequest{
		Version: 1,
		Since:   1000,
	}

	// When
	response, err := service.RequestSyncResponse(context.Background(), initialRequest, "test-topic")

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, uint32(1000), response.Since)
}

func TestRequestSyncResponseService_RequestSyncResponse_Error(t *testing.T) {
	// Given
	mockProvider := &mockRequestSyncResponseProvider{shouldFail: true}
	service := app.NewRequestSyncResponseService(mockProvider)
	initialRequest := &core.GASPInitialRequest{
		Version: 1,
		Since:   1000,
	}

	// When
	response, err := service.RequestSyncResponse(context.Background(), initialRequest, "test-topic")

	// Then
	assert.Error(t, err)
	assert.ErrorIs(t, err, app.ErrRequestSyncResponseProvider)
	assert.Nil(t, response)
}

func TestNewRequestSyncResponseService_NilProvider(t *testing.T) {
	// Given/When/Then
	assert.Panics(t, func() {
		app.NewRequestSyncResponseService(nil)
	})
}
