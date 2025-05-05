package ports_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestLookupQuestionHandler_Handle_ShouldReturnBadRequestResponse(t *testing.T) {
	tests := map[string]struct {
		expectedStatusCode     int
		expectedResponse       openapi.Error
		body                   interface{}
		lookupProviderMockOpts []testabilities.LookupQuestionProviderMockOption
	}{
		"Invalid request body - malformed JSON": {
			expectedStatusCode: fiber.StatusBadRequest,
			body:               `{invalid json`,
			expectedResponse:   ports.NewInvalidRequestBodyResponse(),
			lookupProviderMockOpts: []testabilities.LookupQuestionProviderMockOption{
				testabilities.LookupQuestionProviderMockNotCalled(),
			},
		},
		"Missing service field in request body": {
			expectedStatusCode: fiber.StatusBadRequest,
			body:               map[string]interface{}{"query": map[string]string{"test": "value"}},
			expectedResponse:   ports.NewMissingServiceFieldResponse(),
			lookupProviderMockOpts: []testabilities.LookupQuestionProviderMockOption{
				testabilities.LookupQuestionProviderMockWithError(app.ErrMissingServiceField),
			},
		},
		"Empty service field in request body": {
			expectedStatusCode: fiber.StatusBadRequest,
			body:               map[string]interface{}{"service": "", "query": map[string]string{"test": "value"}},
			expectedResponse:   ports.NewMissingServiceFieldResponse(),
			lookupProviderMockOpts: []testabilities.LookupQuestionProviderMockOption{
				testabilities.LookupQuestionProviderMockWithError(app.ErrMissingServiceField),
			},
		},
		"Provider returns error": {
			expectedStatusCode: fiber.StatusInternalServerError,
			body:               map[string]interface{}{"service": "test-service", "query": map[string]string{"test": "value"}},
			expectedResponse:   ports.NewLookupQuestionProviderErrorResponse(),
			lookupProviderMockOpts: []testabilities.LookupQuestionProviderMockOption{
				testabilities.LookupQuestionProviderMockWithError(errors.New("provider error")),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewLookupQuestionProviderMock(t, tc.lookupProviderMockOpts...)
			engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupQuestionProvider(mock))
			fixture := server2.NewServerTestFixture(t, server2.WithEngine(engine))

			// when:
			var actualResponse openapi.BadRequestResponse
			var requestBody []byte

			if jsonBody, ok := tc.body.(string); ok {
				requestBody = []byte(jsonBody)
			} else {
				requestBody, _ = json.Marshal(tc.body)
			}

			res, _ := fixture.Client().
				R().
				SetHeader("Content-Type", "application/json").
				SetBody(requestBody).
				SetError(&actualResponse).
				Post("/api/v1/lookup")

			// then:
			require.Equal(t, tc.expectedStatusCode, res.StatusCode())

			if tc.expectedStatusCode >= 400 {
				if strings.Contains(string(res.Body()), "message") {
					require.Equal(t, &tc.expectedResponse, &actualResponse)
				}
			}

			mock.AssertCalled()
		})
	}
}

func TestLookupQuestionHandler_Handle_ShouldReturnLookupQuestionSuccessResponse(t *testing.T) {
	// given:
	expectedAnswer := &lookup.LookupAnswer{
		Type:   lookup.AnswerTypeFreeform,
		Result: map[string]interface{}{"test": "value"},
	}

	mock := testabilities.NewLookupQuestionProviderMock(t, testabilities.LookupQuestionProviderMockWithAnswer(expectedAnswer))
	engine := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupQuestionProvider(mock))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(engine))

	requestBody := map[string]interface{}{
		"service": "test-service",
		"query":   map[string]string{"test": "query"},
	}

	// when:
	var actualResponse openapi.LookupAnswer

	res, _ := fixture.Client().
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		SetResult(&actualResponse).
		Post("/api/v1/lookup")

	// then:
	expectedResponse := ports.NewLookupQuestionSuccessResponse(expectedAnswer)

	require.Equal(t, http.StatusOK, res.StatusCode())
	require.Equal(t, expectedResponse.Answer, actualResponse.Answer)
	mock.AssertCalled()
}
