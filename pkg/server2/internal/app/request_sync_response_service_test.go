package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestRequestSyncResponseService_ValidCases(t *testing.T) {
	tests := map[string]struct {
		dto               app.RequestSyncResponseDTO
		topic             string
		expectations      testabilities.RequestSyncResponseProviderMockExpectations
		expectedUTXOCount int
		expectedSince     uint32
	}{
		"Request sync response service succeeds with empty UTXO list": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultMockRequestPayload.Version,
				Since:   testabilities.DefaultMockRequestPayload.Since,
			},
			topic:             "test-topic",
			expectations:      testabilities.NewEmptyResponseExpectations(),
			expectedUTXOCount: 0,
			expectedSince:     0,
		},
		"Request sync response service succeeds with single UTXO": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultMockRequestPayload.Version,
				Since:   testabilities.DefaultMockRequestPayload.Since,
			},
			topic:             "test-topic",
			expectations:      testabilities.NewSingleUTXOResponseExpectations(),
			expectedUTXOCount: 1,
			expectedSince:     1000000,
		},
		"Request sync response service succeeds with multiple UTXOs": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultMockRequestPayload.Version,
				Since:   testabilities.DefaultMockRequestPayload.Since,
			},
			topic:             "test-topic",
			expectations:      testabilities.DefaultRequestSyncResponseProviderMockExpectations,
			expectedUTXOCount: 3,
			expectedSince:     1234567890,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			provider := testabilities.NewRequestSyncResponseProviderMock(t, tc.expectations)
			service := app.NewRequestSyncResponseService(provider)

			// when:
			response, err := service.RequestSyncResponse(context.Background(), &tc.dto, tc.topic)

			// then:
			require.NoError(t, err)
			require.NotNil(t, response)
			require.Len(t, response.UTXOList, tc.expectedUTXOCount)
			require.Equal(t, tc.expectedSince, response.Since)

			if tc.expectedUTXOCount > 0 {
				require.NotNil(t, response.UTXOList[0])
				require.NotEmpty(t, response.UTXOList[0].Txid.String())
			}

			provider.AssertCalled()
		})
	}
}

func TestRequestSyncResponseService_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		dto           app.RequestSyncResponseDTO
		topic         string
		expectations  testabilities.RequestSyncResponseProviderMockExpectations
		expectedError app.Error
	}{
		"Request sync response service fails due to empty topic": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultMockRequestPayload.Version,
				Since:   testabilities.DefaultMockRequestPayload.Since,
			},
			topic: "",
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				ProvideForeignSyncResponseCall: false,
			},
			expectedError: app.NewRequestSyncResponseInvalidInputError(),
		},
		"Request sync response service fails due to invalid version": {
			dto: app.RequestSyncResponseDTO{
				Version: 0,
				Since:   testabilities.DefaultMockRequestPayload.Since,
			},
			topic: "test-topic",
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				ProvideForeignSyncResponseCall: false,
			},
			expectedError: app.NewRequestSyncResponseInvalidVersionError(),
		},
		"Request sync response service fails due to invalid since value": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultMockRequestPayload.Version,
				Since:   0,
			},
			topic: "test-topic",
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				ProvideForeignSyncResponseCall: false,
			},
			expectedError: app.NewRequestSyncResponseInvalidSinceError(),
		},
		"Request sync response service fails due to provider error": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultMockRequestPayload.Version,
				Since:   testabilities.DefaultMockRequestPayload.Since,
			},
			topic:         "test-topic",
			expectations:  testabilities.NewErrorResponseExpectations(errors.New("provider error")),
			expectedError: app.NewRequestSyncResponseProviderError(errors.New("provider error")),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewRequestSyncResponseProviderMock(t, tc.expectations)
			service := app.NewRequestSyncResponseService(mock)

			// when:
			response, err := service.RequestSyncResponse(context.Background(), &tc.dto, tc.topic)

			// then:
			var actualErr app.Error
			require.ErrorAs(t, err, &actualErr)
			require.Equal(t, tc.expectedError, actualErr)
			require.Nil(t, response)
			mock.AssertCalled()
		})
	}
}
