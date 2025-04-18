package engine_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngine_GetUTXOHistory_ShouldReturnImmediateOutput_WhenSelectorIsNil(t *testing.T) {
	// given
	output := &engine.Output{Beef: []byte("beef")}
	sut := &engine.Engine{}

	// when
	result, err := sut.GetUTXOHistory(context.Background(), output, nil, 0)

	// then
	require.NoError(t, err)
	assert.Equal(t, output, result)
}

func TestEngine_GetUTXOHistory_ShouldReturnNil_WhenSelectorReturnsFalse(t *testing.T) {
	// given
	output := &engine.Output{Beef: []byte("beef")}
	sut := &engine.Engine{}

	historySelector := func(beef []byte, outputIndex uint32, currentDepth uint32) bool {
		return false
	}

	// when
	result, err := sut.GetUTXOHistory(context.Background(), output, historySelector, 0)

	// then
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestEngine_GetUTXOHistory_ShouldReturnOutput_WhenNoOutputsConsumed(t *testing.T) {
	// given
	output := &engine.Output{
		Beef:            []byte("beef"),
		OutputsConsumed: nil,
	}
	sut := &engine.Engine{}

	historySelector := func(beef []byte, outputIndex uint32, currentDepth uint32) bool {
		return true
	}

	// when
	result, err := sut.GetUTXOHistory(context.Background(), output, historySelector, 0)

	// then
	require.NoError(t, err)
	assert.Equal(t, output, result)
}

func TestEngine_GetUTXOHistory_ShouldTravelRecursively_WhenOutputsConsumedPresent(t *testing.T) {
	// given
	ctx := context.Background()

	parentOutpoint := &overlay.Outpoint{Txid: fakeTxID(t), OutputIndex: 0}
	childOutpoint := &overlay.Outpoint{Txid: fakeTxID(t), OutputIndex: 1}

	childBeef := createDummyBEEF(t)
	parentBeef := createDummyBEEF(t)

	childOutput := &engine.Output{
		Outpoint: *childOutpoint,
		Beef:     childBeef,
	}
	parentOutput := &engine.Output{
		Outpoint:        *parentOutpoint,
		Beef:            parentBeef,
		OutputsConsumed: []*overlay.Outpoint{childOutpoint},
	}

	sut := &engine.Engine{
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				if outpoint.String() == childOutpoint.String() {
					return childOutput, nil
				}
				return nil, errors.New("unexpected output")
			},
		},
	}

	historySelector := func(beef []byte, outputIndex uint32, currentDepth uint32) bool {
		return true
	}

	// when
	result, err := sut.GetUTXOHistory(ctx, parentOutput, historySelector, 0)

	// then
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Beef)
}

func TestEngine_GetUTXOHistory_ShouldReturnError_WhenStorageFails(t *testing.T) {
	// given
	ctx := context.Background()

	parentOutpoint := &overlay.Outpoint{Txid: fakeTxID(t), OutputIndex: 0}
	childOutpoint := &overlay.Outpoint{Txid: fakeTxID(t), OutputIndex: 1}

	parentOutput := &engine.Output{
		Outpoint:        *parentOutpoint,
		Beef:            []byte("parent beef"),
		OutputsConsumed: []*overlay.Outpoint{childOutpoint},
	}

	sut := &engine.Engine{
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				return nil, errors.New("storage error")
			},
		},
	}

	historySelector := func(beef []byte, outputIndex uint32, currentDepth uint32) bool {
		return true
	}

	// when
	result, err := sut.GetUTXOHistory(ctx, parentOutput, historySelector, 0)

	// then
	require.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "storage error")
}
