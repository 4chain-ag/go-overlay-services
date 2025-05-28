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

func TestRequestSyncResponseService_ValidCase(t *testing.T) {
	// given:
	expectedResponse := &core.GASPInitialResponse{UTXOList: []*overlay.Outpoint{{}}, Since: 1000}
	provider := testabilities.NewRequestSyncResponseProviderMock(t, testabilities.RequestSyncResponseProviderMockExpectations{ProvideForeignSyncResponseCall: true, Response: expectedResponse})
	service := app.NewRequestSyncResponseService(provider)

	// when:
	response, err := service.RequestSyncResponse(context.Background(), &app.RequestSyncResponseDTO{Version: 1, Since: 500}, "test-topic")

	// then:
	require.NoError(t, err)
	require.Equal(t, expectedResponse, response)
	provider.AssertCalled()
}

func TestRequestSyncResponseService_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		expectedError app.Error
		expectations  testabilities.RequestSyncResponseProviderMockExpectations
		topic         string
		dto           app.RequestSyncResponseDTO
	}{
		" Request sync response service fails due to invalid input ": {
			topic: "",
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				ProvideForeignSyncResponseCall: false,
			},
			dto: app.RequestSyncResponseDTO{
				Version: 1,
				Since:   500,
			},
			expectedError: app.NewRequestSyncResponseInvalidInputError(),
		},
		"Request sync response service fails due to provider error": {
			topic: "test-topic",
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				ProvideForeignSyncResponseCall: true,
				Error:                          errors.New("provider error"),
			},
			dto: app.RequestSyncResponseDTO{
				Version: 1,
				Since:   500,
			},
			expectedError: app.NewRequestSyncResponseProviderError(errors.New("provider error")),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewRequestSyncResponseProviderMock(t, tc.expectations)
			service := app.NewRequestSyncResponseService(mock)

			// when:
			document, err := service.RequestSyncResponse(context.Background(), &tc.dto, tc.topic)

			// then:
			var actualErr app.Error
			require.ErrorAs(t, err, &actualErr)
			require.Equal(t, tc.expectedError, actualErr)
			require.Empty(t, document)
			mock.AssertCalled()
		})
	}
}
