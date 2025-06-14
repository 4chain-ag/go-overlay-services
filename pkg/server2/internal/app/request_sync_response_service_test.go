package app_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/stretchr/testify/require"
)

func TestRequestSyncResponseService_ValidCase(t *testing.T) {
	// given:
	expectations := testabilities.RequestSyncResponseProviderMockExpectations{
		ProvideForeignSyncResponseCall: true,
		InitialRequest: &core.GASPInitialRequest{
			Version: testabilities.DefaultVersion,
			Since:   testabilities.DefaultSince,
		},
		Topic: testabilities.DefaultTopic,
		Response: &core.GASPInitialResponse{
			Since: testabilities.DefaultSince,
			UTXOList: []*transaction.Outpoint{
				{
					Txid:  *testabilities.DummyTxHash(t, "03895fb984362a4196bc9931629318fcbb2aeba7c6293638119ea653fa31d119"),
					Index: 0,
				},
				{
					Txid:  *testabilities.DummyTxHash(t, "27c8f37851aabc468d3dbb6bf0789dc398a602dcb897ca04e7815d939d621595"),
					Index: 1,
				},
				{
					Txid:  *testabilities.DummyTxHash(t, "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"),
					Index: 2,
				},
			},
		},
	}

	expectedDTO := app.NewRequestSyncResponseDTO(expectations.Response)
	provider := testabilities.NewRequestSyncResponseProviderMock(t, expectations)
	service := app.NewRequestSyncResponseService(provider)

	// when:
	actualDTO, err := service.RequestSyncResponse(
		t.Context(),
		testabilities.DefaultTopic,
		testabilities.DefaultVersion,
		testabilities.DefaultSince)

	// then:
	require.NoError(t, err)
	require.Equal(t, expectedDTO, actualDTO)
	provider.AssertCalled()
}

func TestRequestSyncResponseService_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		version       app.Version
		since         app.Since
		topic         app.Topic
		expectations  testabilities.RequestSyncResponseProviderMockExpectations
		expectedError app.Error
	}{
		"Request sync response service fails to handle the sync request - empty topic": {
			version: testabilities.DefaultVersion,
			since:   testabilities.DefaultSince,
			topic:   "",

			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				InitialRequest:                 nil,
				Topic:                          "",
				ProvideForeignSyncResponseCall: false,
			},
			expectedError: app.NewIncorrectInputWithFieldError("topic"),
		},
		"Request sync response service fails to handle the sync request - negative version": {
			version: -1,
			since:   testabilities.DefaultSince,
			topic:   testabilities.DefaultTopic,
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				InitialRequest:                 nil,
				Topic:                          "",
				ProvideForeignSyncResponseCall: false,
			},
			expectedError: app.NewIncorrectInputWithFieldError("version"),
		},
		"Request sync response service fails to handle the sync request - internal provider error": {
			version: testabilities.DefaultVersion,
			since:   testabilities.DefaultSince,
			topic:   testabilities.DefaultTopic,
			expectations: testabilities.RequestSyncResponseProviderMockExpectations{
				InitialRequest: &core.GASPInitialRequest{
					Version: testabilities.DefaultVersion,
					Since:   uint32(testabilities.DefaultSince),
				},
				Topic:                          testabilities.DefaultTopic,
				ProvideForeignSyncResponseCall: true,
				Error:                          testabilities.ErrTestNoopOpFailure,
			},
			expectedError: app.NewRequestSyncResponseProviderError(testabilities.ErrTestNoopOpFailure),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewRequestSyncResponseProviderMock(t, tc.expectations)
			service := app.NewRequestSyncResponseService(mock)

			// when:
			response, err := service.RequestSyncResponse(
				t.Context(),
				tc.topic,
				tc.version,
				tc.since,
			)

			// then:
			var actualErr app.Error
			require.ErrorAs(t, err, &actualErr)
			require.Equal(t, tc.expectedError, actualErr)
			require.Nil(t, response)
			mock.AssertCalled()
		})
	}
}
