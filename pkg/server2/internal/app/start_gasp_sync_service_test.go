package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestStartGASPSyncService_InvalidCase(t *testing.T) {
	// given:
	providerError := errors.New("internal start GASP sync service test error")
	expectations := testabilities.StartGASPSyncProviderMockExpectations{
		StartGASPSyncCall: true,
		Error:             providerError,
	}
	expectedErr := app.NewStartGASPSyncProviderError(providerError)
	mock := testabilities.NewStartGASPSyncProviderMock(t, expectations)
	service := app.NewStartGASPSyncService(mock)

	// when:
	err := service.StartGASPSync(context.Background())

	// then:
	var actualErr app.Error
	require.ErrorAs(t, err, &actualErr)
	require.Equal(t, expectedErr, actualErr)
	mock.AssertCalled()
}

func TestStartGASPSyncService_ValidCase(t *testing.T) {
	// given:
	expectations := testabilities.StartGASPSyncProviderMockExpectations{
		StartGASPSyncCall: true,
		Error:             nil,
	}

	mock := testabilities.NewStartGASPSyncProviderMock(t, expectations)
	service := app.NewStartGASPSyncService(mock)

	// when:
	err := service.StartGASPSync(context.Background())

	// then:
	require.NoError(t, err)
	mock.AssertCalled()
}
