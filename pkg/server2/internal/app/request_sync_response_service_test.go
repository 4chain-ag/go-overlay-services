package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

func TestNewRequestSyncResponseService_ShouldPanicWithNilProvider(t *testing.T) {
	require.Panics(t, func() {
		app.NewRequestSyncResponseService(nil)
	})
}

func TestRequestSyncResponseService_ShouldCallProviderAndReturnSuccessfully(t *testing.T) {
	// given:
	expectedResponse := &core.GASPInitialResponse{
		UTXOList: []*overlay.Outpoint{
			{},
		},
		Since: 1000,
	}

	provider := testabilities.NewRequestSyncResponseProviderMock(t, testabilities.RequestSyncResponseProviderMockExpectations{
		ProvideForeignSyncResponseCall: true,
		Response:                       expectedResponse,
	})

	service := app.NewRequestSyncResponseService(provider)
	initialRequest := &core.GASPInitialRequest{
		Version: 1,
		Since:   500,
	}
	topic := "test-topic"

	// when:
	response, err := service.RequestSyncResponse(context.Background(), initialRequest, topic)

	// then:
	require.NoError(t, err)
	require.Equal(t, expectedResponse, response)
	provider.AssertCalled()
}

func TestRequestSyncResponseService_ShouldReturnErrorOnProviderFailure(t *testing.T) {
	// given:
	providerError := errors.New("provider error")
	provider := testabilities.NewRequestSyncResponseProviderMock(t, testabilities.RequestSyncResponseProviderMockExpectations{
		ProvideForeignSyncResponseCall: true,
		Error:                          providerError,
	})

	service := app.NewRequestSyncResponseService(provider)
	initialRequest := &core.GASPInitialRequest{
		Version: 1,
		Since:   500,
	}
	topic := "test-topic"

	// when:
	response, err := service.RequestSyncResponse(context.Background(), initialRequest, topic)

	// then:
	require.Error(t, err)
	require.Nil(t, response)

	var appErr app.Error
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, app.ErrorTypeProviderFailure, appErr.ErrorType())

	provider.AssertCalled()
}
