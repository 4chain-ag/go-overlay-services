package app_test

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestArcIngestService_ValidCase(t *testing.T) {
	expectations := testabilities.ServiceTestMerkleProofProviderExpectations{
		ArcIngestCall:      true,
		ExpectedTxID:       testabilities.NewValidTestTxID(t).String(),
		ExpectedMerklePath: testabilities.NewValidTestMerklePath(t),
		Error:              nil,
	}

	mock := testabilities.NewServiceTestMerkleProofProviderMock(t, expectations)
	service := app.NewArcIngestService(mock)

	// when:
	err := service.HandleArcIngest(context.Background(), &app.ArcIngestDTO{
		TxID:        expectations.ExpectedTxID,
		MerklePath:  expectations.ExpectedMerklePath,
		BlockHeight: 0,
	})

	// then:
	require.NoError(t, err)
	mock.AssertCalled(t)
}

func TestArcIngestService_InvalidCase(t *testing.T) {
	tests := map[string]struct {
		dto          *app.ArcIngestDTO
		expectedErr  app.Error
		expectations testabilities.ServiceTestMerkleProofProviderExpectations
	}{
		"should fail with invalid txID format when txID is not hex": {
			dto: &app.ArcIngestDTO{
				TxID:        "not-a-hex-string",
				MerklePath:  testabilities.NewValidTestMerklePath(t),
				BlockHeight: 0,
			},
			expectedErr: app.NewInvalidTxIDFormatError(),
			expectations: testabilities.ServiceTestMerkleProofProviderExpectations{
				ArcIngestCall: false,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewServiceTestMerkleProofProviderMock(t, tc.expectations)
			service := app.NewArcIngestService(mock)

			// when:
			err := service.HandleArcIngest(context.Background(), tc.dto)

			// then:
			var actualErr app.Error
			require.ErrorAs(t, err, &actualErr)
			require.Equal(t, tc.expectedErr, actualErr)
			mock.AssertCalled(t)
		})
	}
}
