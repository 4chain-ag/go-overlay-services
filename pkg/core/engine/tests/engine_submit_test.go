package engine_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/script"
	"github.com/stretchr/testify/require"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
)

func TestEngine_Submit_Success(t *testing.T) {
	// given:
	ctx := context.Background()
	e := &engine.Engine{
		Managers: map[string]engine.TopicManager{
			"test-topic": fakeManager{},
		},
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				return &engine.Output{}, nil
			},
			doesAppliedTransactionExistFunc: func(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
				return false, nil
			},
			markUTXOAsSpentFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
				return nil
			},
			insertOutputFunc: func(ctx context.Context, output *engine.Output) error {
				return nil
			},
			insertAppliedTransactionFunc: func(ctx context.Context, tx *overlay.AppliedTransaction) error {
				return nil
			},
		},
		ChainTracker: fakeChainTracker{},
	}
	taggedBEEF := overlay.TaggedBEEF{
		Topics: []string{"test-topic"},
		Beef:   createDummyBeef(t),
	}

	// when:
	steak, err := e.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.NoError(t, err)
	require.NotNil(t, steak)
}

func TestEngine_Submit_InvalidBeef_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	e := &engine.Engine{
		Managers: map[string]engine.TopicManager{
			"test-topic": fakeManager{},
		},
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				return &engine.Output{}, nil
			},
			doesAppliedTransactionExistFunc: func(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
				return false, nil
			},
		},
		ChainTracker: fakeChainTracker{},
	}
	taggedBEEF := overlay.TaggedBEEF{
		Topics: []string{"test-topic"},
		Beef:   []byte{0xFF}, // invalid beef
	}

	// when:
	steak, err := e.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.Error(t, err)
	require.Nil(t, steak)
}

func TestEngine_Submit_SPVFail_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	e := &engine.Engine{
		Managers: map[string]engine.TopicManager{
			"test-topic": fakeManager{},
		},
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				return &engine.Output{}, nil
			},
			doesAppliedTransactionExistFunc: func(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
				return false, nil
			},
		},
		ChainTracker: fakeChainTrackerSPVFail{},
	}
	taggedBEEF := overlay.TaggedBEEF{
		Topics: []string{"test-topic"},
		Beef:   createDummyBeef(t),
	}

	// when:
	steak, err := e.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.Error(t, err)
	require.Nil(t, steak)
}

func TestEngine_Submit_DuplicateTransaction_ShouldReturnEmptySteak(t *testing.T) {
	// given:
	ctx := context.Background()
	e := &engine.Engine{
		Managers: map[string]engine.TopicManager{
			"test-topic": fakeManager{},
		},
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				return &engine.Output{}, nil
			},
			doesAppliedTransactionExistFunc: func(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
				return true, nil
			},
			markUTXOAsSpentFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
				return nil
			},
			insertAppliedTransactionFunc: func(ctx context.Context, tx *overlay.AppliedTransaction) error {
				return nil
			},
			insertOutputFunc: func(ctx context.Context, output *engine.Output) error {
				return nil
			},
		},
		ChainTracker: fakeChainTracker{},
	}
	taggedBEEF := overlay.TaggedBEEF{
		Topics: []string{"test-topic"},
		Beef:   createDummyBeef(t),
	}

	// when:
	steak, err := e.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.NoError(t, err)
	require.NotNil(t, steak)
	require.Empty(t, steak["test-topic"].OutputsToAdmit)
}

func TestEngine_Submit_MissingTopic_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	e := &engine.Engine{
		Managers: map[string]engine.TopicManager{
			// no managers, missing topic
		},
		Storage: fakeStorage{},
		ChainTracker: fakeChainTracker{},
	}
	taggedBEEF := overlay.TaggedBEEF{
		Topics: []string{"unknown-topic"},
		Beef:   createDummyBeef(t),
	}

	// when:
	steak, err := e.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.Error(t, err)
	require.Nil(t, steak)
	require.Equal(t, engine.ErrUnknownTopic, err)
}

func TestEngine_Submit_BroadcastFails_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	e := &engine.Engine{
		Managers: map[string]engine.TopicManager{
			"test-topic": fakeManager{},
		},
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				return &engine.Output{}, nil
			},
			doesAppliedTransactionExistFunc: func(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
				return false, nil
			},
			markUTXOAsSpentFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
				return nil
			},
			insertAppliedTransactionFunc: func(ctx context.Context, tx *overlay.AppliedTransaction) error {
				return nil
			},
			insertOutputFunc: func(ctx context.Context, output *engine.Output) error {
				return nil
			},
		},
		ChainTracker: fakeChainTracker{},
		Broadcaster:  fakeBroadcasterFail{},
	}
	taggedBEEF := overlay.TaggedBEEF{
		Topics: []string{"test-topic"},
		Beef:   createDummyBeef(t),
	}

	// when:
	steak, err := e.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.Error(t, err)
	require.Nil(t, steak)
	require.EqualError(t, err, "forced failure for testing")
}

func TestEngine_Submit_OutputInsertFails_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	taggedBEEF, prevTxID := createDummyValidTaggedBEEF(t)

	e := &engine.Engine{
		Managers: map[string]engine.TopicManager{
			"test-topic": fakeManager{},
		},
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				return &engine.Output{
					Outpoint: overlay.Outpoint{
						Txid:        *prevTxID,
						OutputIndex: 0,
					},
					Satoshis: 1000,
					Script:   &script.Script{script.OpTRUE},
					Topic:    "test-topic",
				}, nil
			},
			doesAppliedTransactionExistFunc: func(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
				return false, nil
			},
			markUTXOAsSpentFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
				return nil
			},
			insertAppliedTransactionFunc: func(ctx context.Context, tx *overlay.AppliedTransaction) error {
				return nil
			},
			insertOutputFunc: func(ctx context.Context, output *engine.Output) error {
				return errors.New("insert-failed")
			},
			updateConsumedByFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic string, consumedBy []*overlay.Outpoint) error {
				return nil
			},
		},		
		ChainTracker: fakeChainTracker{},
	}

	// when:
	steak, err := e.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.Error(t, err)
	require.Nil(t, steak)
	require.EqualError(t, err, "insert-failed")
}
