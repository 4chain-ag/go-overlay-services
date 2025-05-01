package ports_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStartGASPSyncProvider struct {
	shouldFail bool
	called     bool
}

func (m *mockStartGASPSyncProvider) StartGASPSync(ctx context.Context) error {
	m.called = true
	if m.shouldFail {
		return errors.New("mock service error")
	}
	return nil
}

func (m *mockStartGASPSyncProvider) AssertCalled(t *testing.T) {
	assert.True(t, m.called, "Expected StartGASPSync to be called")
}

func TestStartGASPSyncHandler_Handle_Success(t *testing.T) {
	// Given
	adminToken := uuid.NewString()
	mock := &mockStartGASPSyncProvider{shouldFail: false}
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithStartGASPSyncProvider(mock))
	fixture := server2.NewTestFixture(t, 
		server2.WithEngine(engine),
		server2.WithAdminBearerToken(adminToken),
	)

	// When
	var successResponse openapi.StartGASPSyncResponse
	res, _ := fixture.Client().
		R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", adminToken)).
		SetResult(&successResponse).
		Post("/api/v1/admin/startGASPSync")

	// Then
	require.Equal(t, http.StatusOK, res.StatusCode())
	assert.Contains(t, successResponse.Message, "successfully")
	mock.AssertCalled(t)
}

func TestStartGASPSyncHandler_Handle_Error(t *testing.T) {
	// Given
	adminToken := uuid.NewString()
	mock := &mockStartGASPSyncProvider{shouldFail: true}
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithStartGASPSyncProvider(mock))
	fixture := server2.NewTestFixture(t, 
		server2.WithEngine(engine),
		server2.WithAdminBearerToken(adminToken),
	)

	// When
	var errorResponse openapi.Error
	res, _ := fixture.Client().
		R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", adminToken)).
		SetError(&errorResponse).
		Post("/api/v1/admin/startGASPSync")

	// Then
	require.Equal(t, http.StatusInternalServerError, res.StatusCode())
	assert.Contains(t, errorResponse.Message, "Unable to process GASP sync request")
	mock.AssertCalled(t)
}
