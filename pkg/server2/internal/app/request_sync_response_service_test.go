package app_test

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

func TestRequestSyncResponseService_ValidCases(t *testing.T) {
	tests := map[string]struct {
		dto          app.RequestSyncResponseDTO
		expectations testabilities.RequestSyncResponseProviderMockExpectations
	}{

		"Request sync response service succeeds with empty UTXO list": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultVersion,
				Since:   testabilities.DefaultSince,
				Topic:   testabilities.DefaultTopic,
			},
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				InitialRequest: &core.GASPInitialRequest{
					Version: testabilities.DefaultVersion,
					Since:   uint32(testabilities.DefaultSince),
				},
				Topic: testabilities.DefaultTopic,
				Response: &core.GASPInitialResponse{
					UTXOList: []*overlay.Outpoint{},
					Since:    0,
				},
				ProvideForeignSyncResponseCall: true,
			},
		},

		"Request sync response service succeeds with minimum since value": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultVersion,
				Since:   0,
				Topic:   testabilities.DefaultTopic,
			},
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				InitialRequest: &core.GASPInitialRequest{
					Version: testabilities.DefaultVersion,
					Since:   uint32(0),
				},
				Topic: testabilities.DefaultTopic,
				Response: &core.GASPInitialResponse{
					UTXOList: []*overlay.Outpoint{},
					Since:    0,
				},
				ProvideForeignSyncResponseCall: true,
			},
		},

		"Request sync response service succeeds with maximum since value": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultVersion,
				Since:   math.MaxUint32,
				Topic:   testabilities.DefaultTopic,
			},
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				InitialRequest: &core.GASPInitialRequest{
					Version: testabilities.DefaultVersion,
					Since:   uint32(math.MaxUint32),
				},
				Topic: testabilities.DefaultTopic,
				Response: &core.GASPInitialResponse{
					UTXOList: []*overlay.Outpoint{},
					Since:    0,
				},
				ProvideForeignSyncResponseCall: true,
			},
		},

		"Request sync response service succeeds with single UTXO": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultVersion,
				Since:   testabilities.DefaultSince,
				Topic:   testabilities.DefaultTopic,
			},
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				InitialRequest: &core.GASPInitialRequest{
					Version: testabilities.DefaultVersion,
					Since:   uint32(testabilities.DefaultSince),
				},
				Topic: testabilities.DefaultTopic,
				Response: &core.GASPInitialResponse{
					UTXOList: []*overlay.Outpoint{
						{
							Txid:        *testabilities.DummyTxHash(t, "03895fb984362a4196bc9931629318fcbb2aeba7c6293638119ea653fa31d119"),
							OutputIndex: 0,
						},
					},
					Since: 1000000,
				},
				ProvideForeignSyncResponseCall: true,
			},
		},

		"Request sync response service succeeds with multiple UTXOs": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultVersion,
				Since:   testabilities.DefaultSince,
				Topic:   testabilities.DefaultTopic,
			},
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				InitialRequest: &core.GASPInitialRequest{
					Version: testabilities.DefaultVersion,
					Since:   uint32(testabilities.DefaultSince),
				},
				Topic: testabilities.DefaultTopic,
				Response: &core.GASPInitialResponse{
					UTXOList: []*overlay.Outpoint{
						{
							Txid:        *testabilities.DummyTxHash(t, "03895fb984362a4196bc9931629318fcbb2aeba7c6293638119ea653fa31d119"),
							OutputIndex: 0,
						},
						{
							Txid:        *testabilities.DummyTxHash(t, "27c8f37851aabc468d3dbb6bf0789dc398a602dcb897ca04e7815d939d621595"),
							OutputIndex: 1,
						},
						{
							Txid:        *testabilities.DummyTxHash(t, "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"),
							OutputIndex: 2,
						},
					},
					Since: 1234567890,
				},
				ProvideForeignSyncResponseCall: true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			provider := testabilities.NewRequestSyncResponseProviderMock(t, tc.expectations)
			service := app.NewRequestSyncResponseService(provider)

			// when:
			response, err := service.RequestSyncResponse(context.Background(), &tc.dto)

			// then:
			require.NoError(t, err)
			require.Equal(t, tc.expectations.Response, response)

			provider.AssertCalled()
		})
	}
}

func TestRequestSyncResponseService_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		dto           app.RequestSyncResponseDTO
		expectations  testabilities.RequestSyncResponseProviderMockExpectations
		expectedError app.Error
	}{
		"Request sync response service fails due to empty topic": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultVersion,
				Since:   testabilities.DefaultSince,
				Topic:   "",
			},
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				InitialRequest:                 nil,
				Topic:                          "",
				ProvideForeignSyncResponseCall: false,
			},
			expectedError: app.NewRequestSyncResponseInvalidInputError(),
		},

		"Request sync response service fails due to invalid version": {
			dto: app.RequestSyncResponseDTO{
				Version: -1,
				Since:   testabilities.DefaultSince,
				Topic:   testabilities.DefaultTopic,
			},
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				InitialRequest:                 nil,
				Topic:                          "",
				ProvideForeignSyncResponseCall: false,
			},
			expectedError: app.NewRequestSyncResponseInvalidVersionError(),
		},

		"Request sync response service fails due to invalid negative since value": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultVersion,
				Since:   -1,
				Topic:   testabilities.DefaultTopic,
			},
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				InitialRequest:                 nil,
				Topic:                          "",
				ProvideForeignSyncResponseCall: false,
			},
			expectedError: app.NewRequestSyncResponseInvalidSinceError(),
		},

		"Request sync response service fails due to maximum since value exceeded": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultVersion,
				Since:   math.MaxUint32 + 1,
				Topic:   testabilities.DefaultTopic,
			},
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				InitialRequest:                 nil,
				Topic:                          "",
				ProvideForeignSyncResponseCall: false,
			},
			expectedError: app.NewRequestSyncResponseInvalidSinceError(),
		},

		"Request sync response service fails due to provider error": {
			dto: app.RequestSyncResponseDTO{
				Version: testabilities.DefaultVersion,
				Since:   testabilities.DefaultSince,
				Topic:   testabilities.DefaultTopic,
			},
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				InitialRequest: &core.GASPInitialRequest{
					Version: testabilities.DefaultVersion,
					Since:   uint32(testabilities.DefaultSince),
				},
				Topic:                          testabilities.DefaultTopic,
				ProvideForeignSyncResponseCall: true,
				Error:                          errors.New("provider error"),
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
			response, err := service.RequestSyncResponse(context.Background(), &tc.dto)

			// then:
			var actualErr app.Error
			require.ErrorAs(t, err, &actualErr)
			require.Equal(t, tc.expectedError, actualErr)
			require.Nil(t, response)
			mock.AssertCalled()
		})
	}
}
