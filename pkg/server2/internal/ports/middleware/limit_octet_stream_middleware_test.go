package middleware_test

import (
	"strings"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/middleware"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestLimitOctetStreamMiddleware_ValidCases(t *testing.T) {
	const octetStreamLimit = 10

	mock := testabilities.NewSubmitTransactionProviderMock(t)
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithSubmitTransactionProvider(mock))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(engine), server2.WithOctetStreamLimit(octetStreamLimit))

	testPaths := []struct {
		endpoint string
		method   string
	}{
		{
			endpoint: "/api/v1/submit",
			method:   fiber.MethodPost,
		},
	}

	tests := map[string]struct {
		name           string
		body           string
		headers        map[string]string
		expectedStatus int
	}{
		"Request size matches octet-stream limit": {
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEOctetStream,
				ports.XTopicsHeader:     "topics1,topics2",
			},
			body:           strings.Repeat("A", octetStreamLimit),
			expectedStatus: fiber.StatusOK,
		},
		"Request size below octet-stream limit": {
			headers: map[string]string{
				fiber.HeaderContentType: fiber.MIMEOctetStream,
				ports.XTopicsHeader:     "topics1,topics2",
			},
			body:           strings.Repeat("A", 5),
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, path := range testPaths {
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				// when:
				res, _ := fixture.Client().
					R().
					SetHeaders(tc.headers).
					SetBody(tc.body).
					Execute(path.method, path.endpoint)

				// then:
				require.Equal(t, tc.expectedStatus, res.StatusCode())
			})
		}
	}
}

func TestLimitOctetStreamMiddleware_InvalidCases(t *testing.T) {
	mock := testabilities.NewSubmitTransactionProviderMock(t)
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithSubmitTransactionProvider(mock))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(engine), server2.WithOctetStreamLimit(10))

	testPaths := []struct {
		endpoint string
		method   string
	}{
		{
			endpoint: "/api/v1/submit",
			method:   fiber.MethodPost,
		},
	}

	tests := map[string]struct {
		name             string
		body             string
		headers          map[string]string
		expectedResponse openapi.Error
		expectedStatus   int
	}{
		"Request size exceeds octet-stream limit": {
			headers:          map[string]string{fiber.HeaderContentType: fiber.MIMEOctetStream},
			body:             strings.Repeat("A", 1025),
			expectedResponse: middleware.NewRequestBodyTooLargeResponse(10),
			expectedStatus:   fiber.StatusBadRequest,
		},
		"Unsupported Content-Type is rejected": {
			headers:          map[string]string{fiber.HeaderContentType: fiber.MIMEApplicationJSON},
			expectedResponse: middleware.NewUnsupportedContentTypeResponse(fiber.MIMEOctetStream),
			expectedStatus:   fiber.StatusBadRequest,
			body:             strings.Repeat("A", 10),
		},
		"Request octet-stream is empty": {
			headers:          map[string]string{fiber.HeaderContentType: fiber.MIMEOctetStream},
			expectedResponse: middleware.ResponseEmptyOctetStream,
			expectedStatus:   fiber.StatusBadRequest,
			body:             "",
		},
	}

	for _, path := range testPaths {
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				// when:
				var actual openapi.BadRequestResponse

				res, _ := fixture.Client().
					R().
					SetHeaders(tc.headers).
					SetBody(tc.body).
					SetError(&actual).
					Execute(path.method, path.endpoint)

				// then:
				require.Equal(t, tc.expectedStatus, res.StatusCode())
				require.Equal(t, tc.expectedResponse, actual, "unexpected error response from Bearer token authorization middleware")
			})
		}
	}
}
