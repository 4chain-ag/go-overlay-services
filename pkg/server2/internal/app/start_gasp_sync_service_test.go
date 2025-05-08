package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/stretchr/testify/require"
)

func TestStartGASPSyncService_InvalidCases(t *testing.T) {

	tests := map[string]struct {
		expectations testabilities.StartGASPSyncProviderMockExpectations

		expectedErrorType app.ErrorType
	}{

		"Start GASP sync service fails - internal error": {

			expectations: testabilities.StartGASPSyncProviderMockExpectations{

				StartGASPSyncCall: true,

				Error: errors.New("internal start GASP sync service test error"),
			},

			expectedErrorType: app.ErrorTypeProviderFailure,
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			// given:

			mock := testabilities.NewStartGASPSyncProviderMock(t, tc.expectations)

			service := app.NewStartGASPSyncService(mock)

			// when:

			err := service.StartGASPSync(context.Background())

			// then:

			var appErr app.Error

			require.ErrorAs(t, err, &appErr)

			require.Equal(t, tc.expectedErrorType, appErr.ErrorType())

			mock.AssertCalled()

		})

	}

}

func TestStartGASPSyncService_ValidCase(t *testing.T) {

	// given:

	expectations := testabilities.StartGASPSyncProviderMockExpectations{

		StartGASPSyncCall: true,

		Error: nil,
	}

	mock := testabilities.NewStartGASPSyncProviderMock(t, expectations)

	service := app.NewStartGASPSyncService(mock)

	// when:

	err := service.StartGASPSync(context.Background())

	// then:

	require.NoError(t, err)

	mock.AssertCalled()

}
