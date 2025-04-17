package commands_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/commands"
	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/commands/testutil"
	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/stretchr/testify/require"
)

// TODO: Add missing unit tests
func Test_ArcIngestHandler_ShouldRespondsWith200AndCallsProvider(t *testing.T) {
	// given:
	payload := commands.ArcIngestRequest{
		TxID:        testutil.ValidTxId,
		MerklePath:  testutil.NewValidTestMerklePath(t),
		BlockHeight: 848372,
	}

	mock := testutil.NewMerkleProofProviderMock(nil, payload.BlockHeight)
	handler, err := commands.NewArcIngestHandler(mock)

	require.NoError(t, err)
	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL, testutil.RequestBody(t, payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// when:
	res, err := ts.Client().Do(req)

	// then:
	require.NoError(t, err)
	defer res.Body.Close()

	require.NotNil(t, res)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var actualResponse commands.ArcIngestHandlerResponse
	require.NoError(t, jsonutil.DecodeResponseBody(res, &actualResponse))

	expectedResponse := commands.NewSuccessArcIngestHandlerResponse()
	require.Equal(t, expectedResponse, actualResponse)

	mock.AssertCalled(t)
}
