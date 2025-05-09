package ports_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestSubmitTransactionHandler_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		expectedStatusCode        int
		expectedResponse          openapi.Error
		headers                   map[string]string
		body                      string
		submitTransactionMockOpts []testabilities.SubmitTransactionProviderMockOption
	}{
		"Submit transaction service fails to handle the transaction submission request - internal error": {
			expectedStatusCode: fiber.StatusInternalServerError,
			body:               "test transaction body",
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEOctetStream,
				ports.XTopicsHeader:     "topics1,topics2",
			},
			expectedResponse: ports.SubmitTransactionServiceInternalError,
			submitTransactionMockOpts: []testabilities.SubmitTransactionProviderMockOption{
				testabilities.SubmitTransactionProviderMockWithError(app.ErrSubmitTransactionProvider),
			},
		},
		"Submit transaction service fails to handle the transaction submission request - timeout error": {
			expectedStatusCode: fiber.StatusRequestTimeout,
			body:               "test transaction body",
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEOctetStream,
				ports.XTopicsHeader:     "topics1,topics2",
			},
			expectedResponse: ports.NewRequestTimeoutResponse(ports.RequestTimeout),
			submitTransactionMockOpts: []testabilities.SubmitTransactionProviderMockOption{
				testabilities.SubmitTransactionProviderMockWithError(app.ErrSubmitTransactionProviderTimeout),
				testabilities.SubmitTransactionProviderMockWithSTEAK(&overlay.Steak{}, 2*time.Second),
			},
		},
		"Missing x-topics header in the HTTP request": {
			expectedStatusCode: fiber.StatusBadRequest,
			body:               "test transaction body",
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEOctetStream,
			},
			expectedResponse: ports.NewRequestMissingHeaderResponse(ports.XTopicsHeader),
			submitTransactionMockOpts: []testabilities.SubmitTransactionProviderMockOption{
				testabilities.SubmitTransactionProviderMockNotCalled(),
			},
		},
		"Empty topics in the x-topics header in the HTTP request": {
			expectedStatusCode: fiber.StatusBadRequest,
			body:               "test transaction body",
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEOctetStream,
				ports.XTopicsHeader:     "",
			},
			expectedResponse: ports.SubmitTransactionRequestInvalidTopicsHeaderFormat,
			submitTransactionMockOpts: []testabilities.SubmitTransactionProviderMockOption{
				testabilities.SubmitTransactionProviderMockNotCalled(),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewSubmitTransactionProviderMock(t, tc.submitTransactionMockOpts...)
			engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithSubmitTransactionProvider(mock))
			fixture := server2.NewServerTestFixture(t, server2.WithEngine(engine))

			// when:
			var actualResponse openapi.BadRequestResponse

			res, _ := fixture.Client().
				R().
				SetHeaders(tc.headers).
				SetBody(tc.body).
				SetError(&actualResponse).
				Post("/api/v1/submit")

			// then:
			require.Equal(t, tc.expectedStatusCode, res.StatusCode())
			require.Equal(t, &tc.expectedResponse, &actualResponse)
			mock.AssertCalled()
		})
	}
}

func TestSubmitTransactionHandler_ValidCase(t *testing.T) {
	// given:
	steak := overlay.Steak{
		"test": &overlay.AdmittanceInstructions{
			OutputsToAdmit: []uint32{1},
		},
	}

	mock := testabilities.NewSubmitTransactionProviderMock(t, testabilities.SubmitTransactionProviderMockWithSTEAK(&steak, time.Microsecond))
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithSubmitTransactionProvider(mock))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(engine))

	headers := map[string]string{
		fiber.HeaderContentType: fiber.MIMEOctetStream,
		ports.XTopicsHeader:     "topic1,topic2",
	}

	// when:
	var actualResponse openapi.SubmitTransactionResponse

	res, _ := fixture.Client().
		R().
		SetHeaders(headers).
		SetBody("test transaction body").
		SetResult(&actualResponse).
		Post("/api/v1/submit")

	// then:
	expectedResponse := ports.NewSubmitTransactionSuccessResponse(&steak)

	require.Equal(t, http.StatusOK, res.StatusCode())
	require.Equal(t, expectedResponse, &actualResponse)
	mock.AssertCalled()
}
