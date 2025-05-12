package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/stretchr/testify/require"
)
//TODO: move all mock providers to testabilities
type ServiceTestMerkleProofProvider struct {
	shouldBeCalled bool
	error          error
	calledWithTxID *chainhash.Hash
	calledWithPath *transaction.MerklePath
	called         bool
}

func (m *ServiceTestMerkleProofProvider) HandleNewMerkleProof(ctx context.Context, txid *chainhash.Hash, proof *transaction.MerklePath) error {
	m.called = true
	m.calledWithTxID = txid
	m.calledWithPath = proof
	return m.error
}

func (m *ServiceTestMerkleProofProvider) AssertCalled(t *testing.T) {
	if m.shouldBeCalled && !m.called {
		t.Error("Expected HandleNewMerkleProof to be called, but it wasn't")
	}
	if !m.shouldBeCalled && m.called {
		t.Error("Expected HandleNewMerkleProof not to be called, but it was")
	}
}

func NewServiceTestMerkleProofProvider(err error) *ServiceTestMerkleProofProvider {
	return &ServiceTestMerkleProofProvider{
		shouldBeCalled: err == nil || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) || (err != nil && !isValidationError(err)),
		error:          err,
	}
}

func isValidationError(err error) bool {
	if appErr, ok := err.(app.Error); ok {
		errorType := appErr.ErrorType()
		return errorType == app.ErrorTypeIncorrectInput
	}
	return false
}

func Test_ArcIngestService_ShouldSuccessfullyProcessValidRequest(t *testing.T) {
	// given:
	provider := NewServiceTestMerkleProofProvider(nil)
	service := app.NewArcIngestService(provider, 5*time.Second)

	// when:
	err := service.HandleArcIngest(context.Background(), testabilities.ValidTxId, testabilities.NewValidTestMerklePath(t), testabilities.DefaultBlockHeight)

	// then:
	require.NoError(t, err)
	provider.AssertCalled(t)
}

func Test_ArcIngestService_ValidationFailures(t *testing.T) {
	tests := map[string]struct {
		txID        string
		merklePath  string
		blockHeight uint32
		errorCheck  func(t *testing.T, err error)
	}{
		"should fail with invalid txID format when txID is not hex": {
			txID:        "not-a-hex-string",
			merklePath:  testabilities.NewValidTestMerklePath(t),
			blockHeight: testabilities.DefaultBlockHeight,
			errorCheck: func(t *testing.T, err error) {
				require.Error(t, err)
				appErr, ok := err.(app.Error)
				require.True(t, ok, "Expected app.Error type")
				require.Equal(t, app.ErrorTypeIncorrectInput, appErr.ErrorType())
				require.Contains(t, appErr.Slug(), "Invalid transaction ID format")
			},
		},
		"should fail with invalid txID length when txID is too short": {
			txID:        "1234abcd",
			merklePath:  testabilities.NewValidTestMerklePath(t),
			blockHeight: testabilities.DefaultBlockHeight,
			errorCheck: func(t *testing.T, err error) {
				require.Error(t, err)
				appErr, ok := err.(app.Error)
				require.True(t, ok, "Expected app.Error type")
				require.Equal(t, app.ErrorTypeIncorrectInput, appErr.ErrorType())
				require.Contains(t, appErr.Slug(), "transaction ID does not match the expected length")
			},
		},
		"should fail with invalid merkle path format when merklePath is not hex": {
			txID:        testabilities.ValidTxId,
			merklePath:  "not-a-hex-merkle-path",
			blockHeight: testabilities.DefaultBlockHeight,
			errorCheck: func(t *testing.T, err error) {
				require.Error(t, err)
				appErr, ok := err.(app.Error)
				require.True(t, ok, "Expected app.Error type")
				require.Equal(t, app.ErrorTypeIncorrectInput, appErr.ErrorType())
				require.Contains(t, appErr.Slug(), "Merkle path format is invalid")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// For validation tests, the provider should never be called
			provider := &ServiceTestMerkleProofProvider{shouldBeCalled: false}
			service := app.NewArcIngestService(provider, 5*time.Second)

			// when:
			err := service.HandleArcIngest(context.Background(), tc.txID, tc.merklePath, tc.blockHeight)

			// then:
			tc.errorCheck(t, err)
			provider.AssertCalled(t)
		})
	}
}

func Test_ArcIngestService_ProviderErrors(t *testing.T) {
	tests := map[string]struct {
		mockError  error
		errorCheck func(t *testing.T, err error)
	}{
		"should fail with timeout error when context times out": {
			mockError: context.DeadlineExceeded,
			errorCheck: func(t *testing.T, err error) {
				require.Error(t, err)
				appErr, ok := err.(app.Error)
				require.True(t, ok, "Expected app.Error type")
				require.Equal(t, app.ErrorTypeOperationTimeout, appErr.ErrorType())
				require.Contains(t, appErr.Slug(), "timeout limit")
			},
		},
		"should fail with canceled error when context is canceled": {
			mockError: context.Canceled,
			errorCheck: func(t *testing.T, err error) {
				require.Error(t, err)
				appErr, ok := err.(app.Error)
				require.True(t, ok, "Expected app.Error type")
				require.Equal(t, app.ErrorTypeUnknown, appErr.ErrorType())
				require.Contains(t, appErr.Slug(), "canceled")
			},
		},
		"should fail with processing failed error for other errors": {
			mockError: errors.New("some internal error"),
			errorCheck: func(t *testing.T, err error) {
				require.Error(t, err)
				appErr, ok := err.(app.Error)
				require.True(t, ok, "Expected app.Error type")
				require.Equal(t, app.ErrorTypeProviderFailure, appErr.ErrorType())
				require.Contains(t, appErr.Slug(), "Internal server error")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			provider := NewServiceTestMerkleProofProvider(tc.mockError)
			service := app.NewArcIngestService(provider, 5*time.Second)

			// when:
			err := service.HandleArcIngest(context.Background(), testabilities.ValidTxId, testabilities.NewValidTestMerklePath(t), testabilities.DefaultBlockHeight)

			// then:
			tc.errorCheck(t, err)
			provider.AssertCalled(t)
		})
	}
}
