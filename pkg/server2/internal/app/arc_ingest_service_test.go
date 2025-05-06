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
	return errors.Is(err, app.ErrInvalidTxIDFormat) ||
		errors.Is(err, app.ErrInvalidTxIDLength) ||
		errors.Is(err, app.ErrInvalidMerklePathFormat)
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
		expectedErr error
	}{
		"should fail with ErrInvalidTxIDFormat when txID is not hex": {
			txID:        "not-a-hex-string",
			merklePath:  testabilities.NewValidTestMerklePath(t),
			blockHeight: testabilities.DefaultBlockHeight,
			expectedErr: app.ErrInvalidTxIDFormat,
		},
		"should fail with ErrInvalidTxIDLength when txID is too short": {
			txID:        "1234abcd",
			merklePath:  testabilities.NewValidTestMerklePath(t),
			blockHeight: testabilities.DefaultBlockHeight,
			expectedErr: app.ErrInvalidTxIDLength,
		},
		"should fail with ErrInvalidMerklePathFormat when merklePath is not hex": {
			txID:        testabilities.ValidTxId,
			merklePath:  "not-a-hex-merkle-path",
			blockHeight: testabilities.DefaultBlockHeight,
			expectedErr: app.ErrInvalidMerklePathFormat,
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
			require.Error(t, err)
			require.ErrorIs(t, err, tc.expectedErr)
			provider.AssertCalled(t)
		})
	}
}

func Test_ArcIngestService_ProviderErrors(t *testing.T) {
	tests := map[string]struct {
		mockError   error
		expectedErr error
	}{
		"should fail with ErrMerkleProofProcessingTimeout when context times out": {
			mockError:   context.DeadlineExceeded,
			expectedErr: app.ErrMerkleProofProcessingTimeout,
		},
		"should fail with ErrMerkleProofProcessingCanceled when context is canceled": {
			mockError:   context.Canceled,
			expectedErr: app.ErrMerkleProofProcessingCanceled,
		},
		"should fail with ErrMerkleProofProcessingFailed for other errors": {
			mockError:   errors.New("some internal error"),
			expectedErr: app.ErrMerkleProofProcessingFailed,
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
			require.Error(t, err)
			require.ErrorIs(t, err, tc.expectedErr)
			provider.AssertCalled(t)
		})
	}
}
