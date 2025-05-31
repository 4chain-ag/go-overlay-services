package app_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/testabilities"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/stretchr/testify/require"
)

func TestLookupQuestionService_ValidCase(t *testing.T) {
	// given:
	expectations := testabilities.LookupQuestionProviderMockExpectations{
		Answer: &lookup.LookupAnswer{
			Type:   lookup.AnswerTypeFreeform,
			Result: map[string]interface{}{"test": "value"},
		},
		LookupQuestionCall: true,
	}

	question := &lookup.LookupQuestion{
		Service: "test-service",
		Query:   json.RawMessage(`{}`),
	}

	mock := testabilities.NewLookupQuestionProviderMock(t, expectations)
	service := app.NewLookupQuestionService(mock)

	// when:
	answer, err := service.LookupQuestion(context.Background(), question)

	// then:
	require.NoError(t, err)
	require.Equal(t, expectations.Answer, answer)
	mock.AssertCalled()
}

func TestLookupQuestionService_InvalidCases(t *testing.T) {
	tests := map[string]struct {
		expectations  testabilities.LookupQuestionProviderMockExpectations
		question      *lookup.LookupQuestion
		expectedError app.Error
	}{
		"LookupQuestion should return error when question is nil": {
			expectations: testabilities.LookupQuestionProviderMockExpectations{
				LookupQuestionCall: false,
			},
			question:      nil,
			expectedError: app.NewInvalidLookupQuestionError(),
		},
		"LookupQuestion should return error when service is empty": {
			expectations: testabilities.LookupQuestionProviderMockExpectations{
				LookupQuestionCall: false,
			},
			question: &lookup.LookupQuestion{
				Service: "",
				Query:   json.RawMessage(`{}`),
			},
			expectedError: app.NewLookupQuestionMissingServiceFieldError(),
		},
		"LookupQuestion should return error when query is nil": {
			expectations: testabilities.LookupQuestionProviderMockExpectations{
				LookupQuestionCall: false,
			},
			question: &lookup.LookupQuestion{
				Service: "test-service",
				Query:   nil,
			},
			expectedError: app.NewLookupQuestionMissingQueryFieldError(),
		},
		"LookupQuestion should return error from provider": {
			expectations: testabilities.LookupQuestionProviderMockExpectations{
				LookupQuestionCall: true,
				Error:              errors.New("provider error"),
			},
			question: &lookup.LookupQuestion{
				Service: "test-service",
				Query:   json.RawMessage(`{}`),
			},
			expectedError: app.NewLookupQuestionProviderError(errors.New("provider error")),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			mock := testabilities.NewLookupQuestionProviderMock(t, tc.expectations)
			service := app.NewLookupQuestionService(mock)

			// when:
			answer, err := service.LookupQuestion(context.Background(), tc.question)

			// then:
			var actualErr app.Error
			require.ErrorAs(t, err, &actualErr)
			require.Equal(t, tc.expectedError, actualErr)
			require.Nil(t, answer)
			mock.AssertCalled()
		})
	}
}
