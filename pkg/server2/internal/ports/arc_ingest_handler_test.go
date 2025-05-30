package ports_test

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestArcIngestHandler_InvalidJSONCase(t *testing.T) {
	// given:
	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithArcIngestProvider(testabilities.NewServiceTestMerkleProofProviderMock(t, testabilities.ServiceTestMerkleProofProviderExpectations{ArcIngestCall: false})))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	// when:
	var actualResponse openapi.Error

	res, _ := fixture.Client().
		R().
		SetHeader("Content-Type", "application/json").
		SetBody("INVALID_JSON").
		SetError(&actualResponse).
		Post("/api/v1/arc-ingest")

	// then:
	require.Equal(t, fiber.StatusBadRequest, res.StatusCode())
	require.Equal(t, "Invalid request format", actualResponse.Message)
	stub.AssertProvidersState()
}

func TestArcIngestHandler_ValidCase(t *testing.T) {
	// given:
	expectations := testabilities.ServiceTestMerkleProofProviderExpectations{
		ArcIngestCall:      true,
		ExpectedTxID:       testabilities.NewValidTestTxID(t).String(),
		ExpectedMerklePath: testabilities.NewValidTestMerklePath(t),
		Error:              nil,
	}

	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithArcIngestProvider(testabilities.NewServiceTestMerkleProofProviderMock(t, expectations)))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

	payload := openapi.ArcIngestBody{
		Txid:        expectations.ExpectedTxID,
		MerklePath:  expectations.ExpectedMerklePath,
		BlockHeight: 0,
	}

	expectedResponse := ports.NewArcIngestSuccessResponse()

	// when:
	var actualResponse openapi.ArcIngest

	res, _ := fixture.Client().
		R().
		SetHeader("Content-Type", fiber.MIMEApplicationJSON).
		SetBody(payload).
		SetResult(&actualResponse).
		Post("/api/v1/arc-ingest")

	// then:
	require.Equal(t, fiber.StatusOK, res.StatusCode())
	require.Equal(t, expectedResponse, &actualResponse)
	stub.AssertProvidersState()
}
