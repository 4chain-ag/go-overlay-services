package ports_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestStartGASPSyncHandler_InvalidCases(t *testing.T) {

	tests := map[string]struct {
		expectedStatusCode int

		expectedResponse openapi.Error

		expectations testabilities.StartGASPSyncProviderMockExpectations
	}{

		"Start GASP sync service fails - internal error": {

			expectedStatusCode: fiber.StatusInternalServerError,

			expectedResponse: ports.StartGASPSyncInternalErrorResponse,

			expectations: testabilities.StartGASPSyncProviderMockExpectations{

				Error: errors.New("internal start GASP sync provider error during start GASP sync handler unit test"),

				StartGASPSyncCall: true,
			},
		},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {

			// given:

			const token = "428e1f07-79b6-4901-b0a0-ec1fe815331b"

			stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithStartGASPSyncProvider(testabilities.NewStartGASPSyncProviderMock(t, tc.expectations)))

			fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub), server2.WithAdminBearerToken(token))

			// when:

			var actualResponse openapi.Error

			res, _ := fixture.Client().
				R().
				SetHeader(fiber.HeaderAuthorization, "Bearer "+token).
				SetError(&actualResponse).
				Post("/api/v1/admin/startGASPSync")

			// then:

			require.Equal(t, tc.expectedStatusCode, res.StatusCode())

			require.Equal(t, &tc.expectedResponse, &actualResponse)

			stub.AssertProvidersState()

		})

	}

}

func TestStartGASPSyncHandler_ValidCase(t *testing.T) {

	// given:

	const token = "428e1f07-79b6-4901-b0a0-ec1fe815331b"

	expectations := testabilities.StartGASPSyncProviderMockExpectations{

		StartGASPSyncCall: true,

		Error: nil,
	}

	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithStartGASPSyncProvider(testabilities.NewStartGASPSyncProviderMock(t, expectations)))

	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub), server2.WithAdminBearerToken(token))

	// when:

	var actualResponse openapi.StartGASPSync

	res, _ := fixture.Client().
		R().
		SetHeader(fiber.HeaderAuthorization, "Bearer "+token).
		SetResult(&actualResponse).
		Post("/api/v1/admin/startGASPSync")

	// then:

	require.Equal(t, http.StatusOK, res.StatusCode())

	require.Equal(t, ports.StartGASPSyncSuccessResponse, actualResponse)

	stub.AssertProvidersState()

}
