package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestStartGASPSyncService_ProviderFailure(t *testing.T) {
	// given:
	expectations := testabilities.StartGASPSyncProviderMockExpectations{
		StartGASPSyncCall: true,
		Error:             errors.New("internal start GASP sync service test error"),
	}
	mock := testabilities.NewStartGASPSyncProviderMock(t, expectations)
	service, err := app.NewStartGASPSyncService(mock)
	require.NoError(t, err)

	// when:
	syncErr := service.StartGASPSync(context.Background())

	// then:
	var appErr app.Error
	require.ErrorAs(t, syncErr, &appErr)
	require.Equal(t, app.ErrorTypeProviderFailure, appErr.ErrorType())
	mock.AssertCalled()
}

func TestStartGASPSyncService_NilProvider(t *testing.T) {
	// given/when:
	service, err := app.NewStartGASPSyncService(nil)

	// then:
	require.Error(t, err)
	var appErr app.Error
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, app.ErrorTypeIncorrectInput, appErr.ErrorType())
	require.Nil(t, service)
}

func TestStartGASPSyncService_ValidCase(t *testing.T) {
	// given:
	expectations := testabilities.StartGASPSyncProviderMockExpectations{
		StartGASPSyncCall: true,
		Error:             nil,
	}
	mock := testabilities.NewStartGASPSyncProviderMock(t, expectations)
	service, err := app.NewStartGASPSyncService(mock)
	require.NoError(t, err)

	// when:
	err = service.StartGASPSync(context.Background())

	// then:
	require.NoError(t, err)
	mock.AssertCalled()
}
