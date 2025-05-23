package ports_test

import (
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestStartGASPSyncHandler_InvalidCase(t *testing.T) {
	// given:
	providerError := errors.New("internal start GASP sync provider error during start GASP sync handler unit test")
	expectations := testabilities.StartGASPSyncProviderMockExpectations{
		StartGASPSyncCall: true,
		Error:             providerError,
	}

	const token = "428e1f07-79b6-4901-b0a0-ec1fe815331b"
	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithStartGASPSyncProvider(testabilities.NewStartGASPSyncProviderMock(t, expectations)))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub), server2.WithAdminBearerToken(token))
	expectedResponse := testabilities.NewTestOpenapiErrorResponse(t, app.NewStartGASPSyncProviderError(providerError))

	// when:
	var actualResponse openapi.Error
	res, _ := fixture.Client().
		R().
		SetHeader(fiber.HeaderAuthorization, "Bearer "+token).
		SetError(&actualResponse).
		Post("/api/v1/admin/startGASPSync")

	// then:
	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode())
	require.Equal(t, expectedResponse, actualResponse)
	stub.AssertProvidersState()
}

func TestStartGASPSyncHandler_ValidCase(t *testing.T) {
	// given:
	const token = "428e1f07-79b6-4901-b0a0-ec1fe815331b"
	expectations := testabilities.StartGASPSyncProviderMockExpectations{
		StartGASPSyncCall: true,
		Error:             nil,
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
	require.Equal(t, fiber.StatusOK, res.StatusCode())
	require.Equal(t, ports.NewStartGASPSyncResponse(), actualResponse)
	stub.AssertProvidersState()
}
