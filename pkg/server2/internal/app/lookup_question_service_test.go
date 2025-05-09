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

func TestLookupQuestionService(t *testing.T) {
	t.Run("should return error when question is nil", func(t *testing.T) {
		// given
		provider := &testabilities.SimpleLookupQuestionProvider{}
		service := app.NewLookupQuestionService(provider)

		// when
		answer, err := service.LookupQuestion(context.Background(), nil)

		// then
		require.Error(t, err)
		require.Nil(t, answer)

		var appErr app.Error
		require.ErrorAs(t, err, &appErr)
		require.Equal(t, app.ErrorTypeIncorrectInput, appErr.ErrorType())
		require.Contains(t, appErr.Error(), "lookup question cannot be nil")
	})

	t.Run("should return error when service is empty", func(t *testing.T) {
		// given
		provider := &testabilities.SimpleLookupQuestionProvider{}
		service := app.NewLookupQuestionService(provider)
		question := &lookup.LookupQuestion{
			Service: "",
			Query:   json.RawMessage(`{}`),
		}

		// when
		answer, err := service.LookupQuestion(context.Background(), question)

		// then
		require.Error(t, err)
		require.Nil(t, answer)

		var appErr app.Error
		require.ErrorAs(t, err, &appErr)
		require.Equal(t, app.ErrorTypeIncorrectInput, appErr.ErrorType())
		require.Contains(t, appErr.Error(), "missing required service field")
	})

	t.Run("should return error from provider", func(t *testing.T) {
		// given
		expectedErr := errors.New("provider error")
		provider := &testabilities.SimpleLookupQuestionProvider{Err: expectedErr}
		service := app.NewLookupQuestionService(provider)
		question := &lookup.LookupQuestion{
			Service: "test-service",
			Query:   json.RawMessage(`{}`),
		}

		// when
		answer, err := service.LookupQuestion(context.Background(), question)

		// then
		require.Error(t, err)
		require.Nil(t, answer)

		var appErr app.Error
		require.ErrorAs(t, err, &appErr)
		require.Equal(t, app.ErrorTypeProviderFailure, appErr.ErrorType())
		require.Contains(t, appErr.Error(), "provider error")
	})

	t.Run("should return answer from provider", func(t *testing.T) {
		// given
		expectedAnswer := &lookup.LookupAnswer{
			Type:   lookup.AnswerTypeFreeform,
			Result: map[string]interface{}{"test": "value"},
		}
		provider := &testabilities.SimpleLookupQuestionProvider{Answer: expectedAnswer}
		service := app.NewLookupQuestionService(provider)
		question := &lookup.LookupQuestion{
			Service: "test-service",
			Query:   json.RawMessage(`{}`),
		}

		// when
		answer, err := service.LookupQuestion(context.Background(), question)

		// then
		require.NoError(t, err)
		require.Equal(t, expectedAnswer, answer)
	})
}
