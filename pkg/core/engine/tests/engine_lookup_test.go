package engine_test

import (
	"context"
	"testing"

	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/stretchr/testify/require"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
)

func TestEngine_Lookup_ShouldReturnError_WhenServiceUnknown(t *testing.T) {
	t.Parallel()

	// given
	e := &engine.Engine{
		LookupServices: make(map[string]engine.LookupService),
	}

	// when
	answer, err := e.Lookup(context.Background(), &lookup.LookupQuestion{Service: "non-existing"})

	// then
	require.Error(t, err)
	require.Nil(t, answer)
	require.Equal(t, engine.ErrUnknownTopic, err)
}

func TestEngine_Lookup_ShouldReturnError_WhenServiceLookupFails(t *testing.T) {
	t.Parallel()

	// given
	e := &engine.Engine{
		LookupServices: map[string]engine.LookupService{
			"test": fakeLookupService{
				lookupFunc: func(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
					return nil, errFakeLookup
				},
			},
		},
	}

	// when
	answer, err := e.Lookup(context.Background(), &lookup.LookupQuestion{Service: "test"})

	// then
	require.Error(t, err)
	require.Nil(t, answer)
	require.Equal(t, errFakeLookup, err)
}

func TestEngine_Lookup_ShouldReturnDirectResult_WhenAnswerTypeIsFreeform(t *testing.T) {
	t.Parallel()

	// given
	e := &engine.Engine{
		LookupServices: map[string]engine.LookupService{
			"test": fakeLookupService{
				lookupFunc: func(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
					return &lookup.LookupAnswer{Type: lookup.AnswerTypeFreeform}, nil
				},
			},
		},
	}

	// when
	answer, err := e.Lookup(context.Background(), &lookup.LookupQuestion{Service: "test"})

	// then
	require.NoError(t, err)
	require.NotNil(t, answer)
	require.Equal(t, lookup.AnswerTypeFreeform, answer.Type)
}

func TestEngine_Lookup_ShouldReturnDirectResult_WhenAnswerTypeIsOutputList(t *testing.T) {
	t.Parallel()

	// given
	e := &engine.Engine{
		LookupServices: map[string]engine.LookupService{
			"test": fakeLookupService{
				lookupFunc: func(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
					return &lookup.LookupAnswer{Type: lookup.AnswerTypeOutputList}, nil
				},
			},
		},
	}

	// when
	answer, err := e.Lookup(context.Background(), &lookup.LookupQuestion{Service: "test"})

	// then
	require.NoError(t, err)
	require.NotNil(t, answer)
	require.Equal(t, lookup.AnswerTypeOutputList, answer.Type)
}

func TestEngine_Lookup_ShouldHydrateOutputs_WhenFormulasProvided(t *testing.T) {
	t.Parallel()

	// given
	ctx := context.Background()
	expectedBeef := []byte("hydrated beef")

	e := &engine.Engine{
		LookupServices: map[string]engine.LookupService{
			"test": fakeLookupService{
				lookupFunc: func(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
					return &lookup.LookupAnswer{
						Type: lookup.AnswerTypeFormula,
						Formulas: []lookup.LookupFormula{
							{Outpoint: &overlay.Outpoint{Txid: fakeTxID(), OutputIndex: 0}},
						},
					}, nil
				},
			},
		},
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				return &engine.Output{
					Outpoint: *outpoint,
					Beef:     expectedBeef,
				}, nil
			},
		},
	}

	// when
	answer, err := e.Lookup(ctx, &lookup.LookupQuestion{Service: "test"})

	// then
	require.NoError(t, err)
	require.NotNil(t, answer)
	require.Equal(t, lookup.AnswerTypeOutputList, answer.Type)
	require.Len(t, answer.Outputs, 1)
	require.Equal(t, expectedBeef, answer.Outputs[0].Beef)
}
