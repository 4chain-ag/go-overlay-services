package ports_test

import (
	"net/http"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestSubmitTransactionHandler_Handle_ShouldReturnBadRequestResponse(t *testing.T) {
	tests := map[string]struct {
		expectedStatusCode int
		expectedResponse   openapi.Error
		request            *http.Request
		opts               []testabilities.SubmitTransactionProviderMockOption
	}{
		"Missing x-topics header in the HTTP request": {
			expectedStatusCode: fiber.StatusBadRequest,
			request:            newSubmitTransactionRequestWithoutHeader(t),
			expectedResponse:   ports.NewRequestMissingHeaderResponse(ports.XTopicsHeader),
		},
		"Empty topics in the x-topics header in the HTTP request": {
			expectedStatusCode: fiber.StatusBadRequest,
			request:            newSubmitTransactionRequest(t, ""),
			expectedResponse:   ports.NewInvalidRequestTopicsFormatResponse(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewSubmitTransactionProviderMock(t, tc.opts...)
			engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithSubmitTransactionProvider(mock))
			serverAPI := server2.NewServerTestAdapter(server2.WithEngine(engine))

			// when:
			res, err := serverAPI.TestRequest(tc.request, -1)
			require.NoError(t, err)
			defer res.Body.Close()

			// then:
			require.Equal(t, tc.expectedStatusCode, res.StatusCode)

			var actualResponse openapi.BadRequestResponse
			testabilities.DecodeResponseBody(t, res, &actualResponse)

			require.Equal(t, &tc.expectedResponse, &actualResponse)
			mock.AssertCalled()
		})
	}
}

func TestSubmitTransactionHandler_Handle_ShouldReturnSubmitTransactionSuccessResponse(t *testing.T) {
	// given:
	steak := overlay.Steak{
		"test": &overlay.AdmittanceInstructions{
			OutputsToAdmit: []uint32{1},
		},
	}

	mock := testabilities.NewSubmitTransactionProviderMock(t,
		testabilities.SubmitTransactionProviderMockWithSTEAK(&steak),
		testabilities.SubmitTransactionProviderMockWithTriggeredCallback(),
	)
	overlayEngineStub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithSubmitTransactionProvider(mock))
	serverAPI := server2.NewServerTestAdapter(server2.WithEngine(overlayEngineStub))

	// when:
	res, err := serverAPI.TestRequest(newSubmitTransactionRequest(t, "topic1,topic2"), -1)

	// then:
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	defer res.Body.Close()

	expectedResponse := ports.NewSubmitTransactionSuccessResponse(&steak)

	var actualResponse openapi.SubmitTransactionResponse
	testabilities.DecodeResponseBody(t, res, &actualResponse)

	require.Equal(t, expectedResponse, &actualResponse)
	mock.AssertCalled()
}

// newSubmitTransactionRequest creates a new HTTP POST request to the /api/v1/submit endpoint
// with a test transaction body and sets the required headers.
// It sets the Content-Type to application/json and includes the provided topics in the X-Topics header.
// This helper is used in tests to simulate a valid transaction submission request.
func newSubmitTransactionRequest(t *testing.T, topics string) *http.Request {
	t.Helper()

	req, err := http.NewRequest(fiber.MethodPost, "/api/v1/submit", testabilities.RequestBody(t, "test transaction body"))
	require.NoError(t, err, "failed to create new HTTP request")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(ports.XTopicsHeader, topics)
	return req
}

// newSubmitTransactionRequestWithoutHeader creates a new HTTP POST request for the transaction submission endpoint,
// but without any headers. It is primarily used for testing purposes.
// The request body contains a test transaction, and any error in the request creation
// is handled with a failure in the test using require.NoError.
func newSubmitTransactionRequestWithoutHeader(t *testing.T) *http.Request {
	t.Helper()

	req, err := http.NewRequest(fiber.MethodPost, "/api/v1/submit", testabilities.RequestBody(t, "test transaction body"))
	require.NoError(t, err, "failed to create new HTTP request")
	return req
}
