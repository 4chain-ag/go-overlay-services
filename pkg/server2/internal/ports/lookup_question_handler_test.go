package ports_test

import (
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestLookupQuestionHandler_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		expectedStatusCode int
		body               interface{}
		expectedResponse   openapi.Error
		expectations       testabilities.LookupQuestionProviderMockExpectations
	}{
		"Lookup Question service fails with invalid request body malformed JSON": {
			expectedStatusCode: fiber.StatusBadRequest,
			body:               `{invalid json`,
			expectedResponse:   testabilities.NewLookupQuestionInvalidRequestBodyResponse(),
			expectations: testabilities.LookupQuestionProviderMockExpectations{
				LookupQuestionCall: false,
			},
		},
		"Lookup Question service fails with missing service field in request body": {
			expectedStatusCode: fiber.StatusBadRequest,
			body:               map[string]interface{}{"query": map[string]string{"test": "value"}},
			expectedResponse:   testabilities.NewLookupQuestionMissingServiceFieldResponse(),
			expectations: testabilities.LookupQuestionProviderMockExpectations{
				LookupQuestionCall: false,
			},
		},
		"Lookup Question service fails with empty service field in request body": {
			expectedStatusCode: fiber.StatusBadRequest,
			body:               map[string]interface{}{"service": "", "query": map[string]string{"test": "value"}},
			expectedResponse:   testabilities.NewLookupQuestionMissingServiceFieldResponse(),
			expectations: testabilities.LookupQuestionProviderMockExpectations{
				LookupQuestionCall: false,
			},
		},
		"Lookup Question service fails with provider error": {
			expectedStatusCode: fiber.StatusInternalServerError,
			body:               map[string]interface{}{"service": "test-service", "query": map[string]string{"test": "value"}},
			expectedResponse:   testabilities.NewLookupQuestionProviderErrorResponse(),
			expectations: testabilities.LookupQuestionProviderMockExpectations{
				LookupQuestionCall: true,
				Error:              errors.New("provider error"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupQuestionProvider(testabilities.NewLookupQuestionProviderMock(t, tc.expectations)))
			fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

			// when:
			var actualResponse openapi.Error

			res, _ := fixture.Client().
				R().
				SetHeader("Content-Type", "application/json").
				SetBody(tc.body).
				SetError(&actualResponse).
				Post("/api/v1/lookup")

			// then:
			require.Equal(t, tc.expectedStatusCode, res.StatusCode())
			require.Equal(t, tc.expectedResponse, actualResponse)
			stub.AssertProvidersState()
		})
	}
}

func TestLookupQuestionHandler_ValidCase(t *testing.T) {
	// given:
	expectations := testabilities.LookupQuestionProviderMockExpectations{
		LookupQuestionCall: true,
		Answer: &lookup.LookupAnswer{
			Type:   lookup.AnswerTypeFreeform,
			Result: map[string]interface{}{"test": "value"},
		},
	}

	stub := testabilities.NewTestOverlayEngineStub(t, testabilities.WithLookupQuestionProvider(testabilities.NewLookupQuestionProviderMock(t, expectations)))
	fixture := server2.NewServerTestFixture(t, server2.WithEngine(stub))

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
	expectedResponse := ports.NewLookupQuestionSuccessResponse(expectations.Answer)
	require.Equal(t, fiber.StatusOK, res.StatusCode())
	require.Equal(t, expectedResponse, &actualResponse)
	stub.AssertProvidersState()
}
