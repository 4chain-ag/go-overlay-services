package overlayhttp_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/overlayhttp"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/overlayhttp/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/stretchr/testify/require"
)

func TestSubmitTransactionHandler_Handle(t *testing.T) {
	// given:
	steak := overlay.Steak{
		"test": &overlay.AdmittanceInstructions{
			OutputsToAdmit: []uint32{1},
		},
	}

	opts := []testabilities.SubmitTransactionProviderMockOption{
		testabilities.SubmitTransactionProviderMockWithSTEAK(&steak),
		testabilities.SubmitTransactionProviderMockWithTriggeredCallback(),
	}

	mock := testabilities.NewSubmitTransactionProviderMock(t, opts...)
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithSubmitTransactionProvider(mock))
	serverAPI := &openapi.ServerInterfaceWrapper{Handler: overlayhttp.NewServerHandlers("", engine)}

	httpHandler := adaptor.FiberHandler(serverAPI.SubmitTransaction)
	ts := httptest.NewServer(httpHandler)
	defer ts.Close()

	requestBody := []byte("test transaction body")

	// Using comma-separated topics
	topics := "topic1,topic2"

	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewBuffer(requestBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(overlayhttp.XTopicsHeader, topics)

	// when:
	res, err := ts.Client().Do(req)

	// then:
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	var actual openapi.SubmitTransactionResponse
	testabilities.DecodeResponseBody(t, res, &actual)

	expected := overlayhttp.NewSubmitTransactionSuccessResponse(&steak)
	require.Equal(t, expected, &actual)

	mock.AssertCalled()
}
