package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/stretchr/testify/assert"
)

type mockStartGASPSyncProvider struct {
	shouldFail bool
}

func (m *mockStartGASPSyncProvider) StartGASPSync(ctx context.Context) error {
	if m.shouldFail {
		return errors.New("mock GASP sync error")
	}
	return nil
}

func TestStartGASPSyncService_StartGASPSync_Success(t *testing.T) {
	// Given
	mockProvider := &mockStartGASPSyncProvider{shouldFail: false}
	service := app.NewStartGASPSyncService(mockProvider)

	// When
	err := service.StartGASPSync(context.Background())

	// Then
	assert.NoError(t, err)
}

func TestStartGASPSyncService_StartGASPSync_Error(t *testing.T) {
	// Given
	mockProvider := &mockStartGASPSyncProvider{shouldFail: true}
	service := app.NewStartGASPSyncService(mockProvider)

	// When
	err := service.StartGASPSync(context.Background())

	// Then
	assert.Error(t, err)
	assert.ErrorIs(t, err, app.ErrStartGASPSyncProvider)
}

func TestNewStartGASPSyncService_NilProvider(t *testing.T) {
	// Given/When/Then
	assert.Panics(t, func() {
		app.NewStartGASPSyncService(nil)
	})
}
