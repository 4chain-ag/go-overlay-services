package ports_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/middleware"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestSubmitTransactionHandler_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		expectedStatusCode int
		expectedResponse   openapi.Error
		headers            map[string]string
		body               string
		expectations       testabilities.SubmitTransactionProviderMockExpectations
	}{
		"Submit transaction service fails to handle the transaction submission request - internal error": {
			expectedStatusCode: fiber.StatusInternalServerError,
			body:               "test transaction body",
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEOctetStream,
				ports.XTopicsHeader:     "topics1,topics2",
			},
			expectedResponse: ports.SubmitTransactionServiceInternalError,
			expectations: testabilities.SubmitTransactionProviderMockExpectations{
				Error:      errors.New("internal submit transaction provider error during submit transaction handler unit test"),
				SubmitCall: true,
			},
		},
		"Submit transaction service fails to handle the transaction submission request - timeout error": {
			expectedStatusCode: fiber.StatusRequestTimeout,
			body:               "test transaction body",
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEOctetStream,
				ports.XTopicsHeader:     "topics1,topics2",
			},
			expectedResponse: middleware.NewRequestTimeoutResponse(time.Second),
			expectations: testabilities.SubmitTransactionProviderMockExpectations{
				SubmitCall:           true,
				TriggerCallbackAfter: 2 * time.Second,
			},
		},
		"Missing x-topics header in the HTTP request": {
			expectedStatusCode: fiber.StatusBadRequest,
			body:               "test transaction body",
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEOctetStream,
			},
			expectedResponse: ports.NewRequestMissingHeaderResponse(ports.XTopicsHeader),
			expectations: testabilities.SubmitTransactionProviderMockExpectations{
				SubmitCall: false,
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
			expectations: testabilities.SubmitTransactionProviderMockExpectations{
				SubmitCall: false,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithSubmitTransactionProvider(testabilities.NewSubmitTransactionProviderMock(t, tc.expectations)))
			fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub), server2.WithSubmitTransactionHandlerResponseTime(stub, time.Second))

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
			stub.AssertProvidersState()
		})
	}
}

func TestSubmitTransactionHandler_ValidCase(t *testing.T) {
	// given:
	expectations := testabilities.SubmitTransactionProviderMockExpectations{
		SubmitCall: true,
		STEAK: &overlay.Steak{
			"test": &overlay.AdmittanceInstructions{
				OutputsToAdmit: []uint32{1},
			},
		},
	}

	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithSubmitTransactionProvider(testabilities.NewSubmitTransactionProviderMock(t, expectations)))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

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
	expectedResponse := ports.NewSubmitTransactionSuccessResponse(expectations.STEAK)

	require.Equal(t, http.StatusOK, res.StatusCode())
	require.Equal(t, expectedResponse, &actualResponse)
	stub.AssertProvidersState()
}
