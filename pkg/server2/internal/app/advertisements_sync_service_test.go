package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestAdvertisementsSyncService_ValidCase(t *testing.T) {
	// given:
	mock := testabilities.NewSyncAdvertisementsProviderMock(t, testabilities.SyncAdvertisementsProviderMockExpectations{
		SyncAdvertisementsCall: true,
		Err:                    nil,
	})
	service := app.NewAdvertisementsSyncService(mock)

	// when:
	err := service.SyncAdvertisements(context.Background())

	// then:
	require.NoError(t, err)
	mock.AssertCalled()
}

func TestAdvertisementsSyncService_InvalidCase(t *testing.T) {
	// given:
	mock := testabilities.NewSyncAdvertisementsProviderMock(t, testabilities.SyncAdvertisementsProviderMockExpectations{
		SyncAdvertisementsCall: true,
		Err:                    errors.New("internal test error"),
	})
	service := app.NewAdvertisementsSyncService(mock)

	// when:
	err := service.SyncAdvertisements(context.Background())

	// then:
	var as app.Error
	require.ErrorAs(t, err, &as)
	require.Equal(t, app.ErrorTypeProviderFailure, as.ErrorType())

	mock.AssertCalled()
}
