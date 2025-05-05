package app_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/stretchr/testify/require"
)

type mockLookupProvider struct {
	answer *lookup.LookupAnswer

	err error
}

func (m *mockLookupProvider) Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {

	return m.answer, m.err

}

func TestLookupQuestionService(t *testing.T) {

	t.Run("should return error when question is nil", func(t *testing.T) {

		// given

		provider := &mockLookupProvider{}

		service := app.NewLookupQuestionService(provider)

		// when

		answer, err := service.LookupQuestion(context.Background(), nil)

		// then

		require.Error(t, err)

		require.Nil(t, answer)

		require.Equal(t, app.ErrInvalidLookupQuestion, err)

	})

	t.Run("should return error when service is empty", func(t *testing.T) {

		// given

		provider := &mockLookupProvider{}

		service := app.NewLookupQuestionService(provider)

		question := &lookup.LookupQuestion{

			Service: "",

			Query: json.RawMessage(`{}`),
		}

		// when

		answer, err := service.LookupQuestion(context.Background(), question)

		// then

		require.Error(t, err)

		require.Nil(t, answer)

		require.Equal(t, app.ErrMissingServiceField, err)

	})

	t.Run("should return error from provider", func(t *testing.T) {

		// given

		expectedErr := errors.New("provider error")

		provider := &mockLookupProvider{err: expectedErr}

		service := app.NewLookupQuestionService(provider)

		question := &lookup.LookupQuestion{

			Service: "test-service",

			Query: json.RawMessage(`{}`),
		}

		// when

		answer, err := service.LookupQuestion(context.Background(), question)

		// then

		require.Error(t, err)

		require.Nil(t, answer)

		require.Equal(t, expectedErr, err)

	})

	t.Run("should return answer from provider", func(t *testing.T) {

		// given

		expectedAnswer := &lookup.LookupAnswer{

			Type: lookup.AnswerTypeFreeform,

			Result: map[string]interface{}{"test": "value"},
		}

		provider := &mockLookupProvider{answer: expectedAnswer}

		service := app.NewLookupQuestionService(provider)

		question := &lookup.LookupQuestion{

			Service: "test-service",

			Query: json.RawMessage(`{}`),
		}

		// when

		answer, err := service.LookupQuestion(context.Background(), question)

		// then

		require.NoError(t, err)

		require.Equal(t, expectedAnswer, answer)

	})

}
