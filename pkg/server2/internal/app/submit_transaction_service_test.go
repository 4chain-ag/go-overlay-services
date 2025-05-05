package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	testvectors "github.com/bsv-blockchain/universal-test-vectors/pkg/testabilities"
	"github.com/stretchr/testify/require"
)

func TestSubmitTransactionService_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		expectations  testabilities.SubmitTransactionProviderMockExpectations
		topics        app.TransactionTopics
		timeout       time.Duration
		txBytes       []byte
		expectedError error
	}{
		"Submit transaction service fails to handle the transaction submission - internal error": {
			topics:  app.TransactionTopics{"topic1", "topic2"},
			txBytes: DummyTxBEEF(t),
			expectations: testabilities.SubmitTransactionProviderMockExpectations{
				SubmitCall: true,
				STEAK:      nil,
				Error:      errors.New("internal submit transaction service test error"),
			},
			expectedError: app.ErrSubmitTransactionProvider,
		},
		"Submit transaction service fails to handle the transaction submission - timeout error": {
			topics:  app.TransactionTopics{"topic1", "topic2"},
			txBytes: DummyTxBEEF(t),
			timeout: time.Second,
			expectations: testabilities.SubmitTransactionProviderMockExpectations{
				SubmitCall:           true,
				TriggerCallbackAfter: 2 * time.Second,
				Error:                nil,
				STEAK:                nil,
			},
			expectedError: app.ErrSubmitTransactionProviderTimeout,
		},
		"Submit transaction service fails to handle the transaction submission - empty topics": {
			txBytes:       DummyTxBEEF(t),
			expectedError: app.ErrMissingTransactionTopics,
			expectations: testabilities.SubmitTransactionProviderMockExpectations{
				SubmitCall: false,
				Error:      nil,
				STEAK:      nil,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewSubmitTransactionProviderMock(t, tc.expectations)
			service := app.NewSubmitTransactionService(mock, tc.timeout)

			// when:
			STEAK, err := service.SubmitTransaction(context.Background(), tc.topics, tc.txBytes...)

			// then:
			require.ErrorIs(t, err, tc.expectedError)
			require.Nil(t, STEAK)
			mock.AssertCalled()
		})
	}
}

func TestSubmitTransactionService_ValidCase(t *testing.T) {
	// given:
	expectations := testabilities.SubmitTransactionProviderMockExpectations{
		STEAK: &overlay.Steak{
			"test_response": &overlay.AdmittanceInstructions{
				OutputsToAdmit: []uint32{1},
				CoinsToRetain:  []uint32{1},
				CoinsRemoved:   []uint32{1},
			},
		},
		Error:      nil,
		SubmitCall: true,
	}

	timeout := time.Second
	topics := app.TransactionTopics{"topic1", "topic2"}
	mock := testabilities.NewSubmitTransactionProviderMock(t, expectations)
	service := app.NewSubmitTransactionService(mock, timeout)

	// when:
	actualSTEAK, err := service.SubmitTransaction(context.Background(), topics)

	// then:
	require.NoError(t, err)
	require.Equal(t, expectations.STEAK, actualSTEAK)
	mock.AssertCalled()
}

// DummyTxBEEF returns a valid transaction serialized in BEEF format for use in tests.
// It creates a dummy transaction with predefined input and output values.
// The test fails immediately if the transaction cannot be serialized or results in an empty byte slice.
func DummyTxBEEF(t *testing.T) []byte {
	t.Helper()

	dummyTx := testvectors.GivenTX().
		WithInput(1000).
		WithP2PKHOutput(999).
		TX()

	bb, err := dummyTx.BEEF()
	require.NoError(t, err)
	require.NotEmpty(t, bb)
	return bb
}
