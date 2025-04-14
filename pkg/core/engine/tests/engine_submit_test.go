package engine_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/script"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/stretchr/testify/require"
)

func TestEngine_Submit_Success(t *testing.T) {
	// given:
	ctx := context.Background()

	sut := &engine.Engine{
		Managers: map[string]engine.TopicManager{
			"test-topic": fakeManager{
				identifyAdmissableOutputsFunc: func(ctx context.Context, beef []byte, previousCoins []uint32) (overlay.AdmittanceInstructions, error) {
					return overlay.AdmittanceInstructions{
						OutputsToAdmit: []uint32{0},
					}, nil
				},
			},
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
		Beef:   createDummyBEEF(t),
	}

	expectedSteak := overlay.Steak{
		"test-topic": &overlay.AdmittanceInstructions{
			OutputsToAdmit: []uint32{0},
		},
	}

	// when:
	steak, err := sut.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.NoError(t, err)
	require.Equal(t, expectedSteak, steak)
}

func TestEngine_Submit_InvalidBeef_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	sut := &engine.Engine{
		Managers: map[string]engine.TopicManager{
			"test-topic": fakeManager{
				identifyAdmissableOutputsFunc: func(ctx context.Context, beef []byte, previousCoins []uint32) (overlay.AdmittanceInstructions, error) {
					return overlay.AdmittanceInstructions{
						OutputsToAdmit: []uint32{0},
					}, nil
				},
			},
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
	steak, err := sut.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid-atomic-beef") // temp fix for SPV failure Submit need to be fixed by wrapping the error to use ErrorIs
	require.Nil(t, steak)
}

func TestEngine_Submit_SPVFail_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	sut := &engine.Engine{
		Managers: map[string]engine.TopicManager{
			"test-topic": fakeManager{
				identifyAdmissableOutputsFunc: func(ctx context.Context, beef []byte, previousCoins []uint32) (overlay.AdmittanceInstructions, error) {
					return overlay.AdmittanceInstructions{
						OutputsToAdmit: []uint32{0},
					}, nil
				},
			},
		},
		Storage: fakeStorage{
			findOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) (*engine.Output, error) {
				return &engine.Output{
					Outpoint: *outpoint,
					Satoshis: 1000,
					Script:   &script.Script{script.OpTRUE},
				}, nil
			},
			findOutputsFunc: func(ctx context.Context, outpoints []*overlay.Outpoint, topic *string, spent *bool, includeBEEF bool) ([]*engine.Output, error) {
				return []*engine.Output{
					{
						Outpoint: *outpoints[0],
						Satoshis: 1000,
						Script:   &script.Script{script.OpTRUE},
					},
				}, nil
			},
			doesAppliedTransactionExistFunc: func(ctx context.Context, tx *overlay.AppliedTransaction) (bool, error) {
				return false, nil
			},
			insertOutputFunc: func(ctx context.Context, output *engine.Output) error {
				return nil
			},
			markUTXOAsSpentFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
				return nil
			},
			insertAppliedTransactionFunc: func(ctx context.Context, tx *overlay.AppliedTransaction) error {
				return nil
			},
		},
		ChainTracker: fakeChainTrackerSPVFail{},
	}

	taggedBEEF := overlay.TaggedBEEF{
		Topics: []string{"test-topic"},
		Beef:   createDummyBeefWithInputs(t),
	}

	// when:
	steak, err := sut.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown txid") // temp fix for SPV failure Submit need to be fixed by wrapping the error to use ErrorIs
	require.Nil(t, steak)
}

func TestEngine_Submit_DuplicateTransaction_ShouldReturnEmptySteak(t *testing.T) {
	// given:
	ctx := context.Background()
	sut := &engine.Engine{
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
		Beef:   createDummyBEEF(t),
	}

	expectedSteak := overlay.Steak{
		"test-topic": &overlay.AdmittanceInstructions{
			OutputsToAdmit: nil,
		},
	}

	// when:
	steak, err := sut.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.NoError(t, err)
	require.Equal(t, expectedSteak, steak)
}

func TestEngine_Submit_MissingTopic_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	sut := &engine.Engine{
		Managers:     map[string]engine.TopicManager{},
		Storage:      fakeStorage{},
		ChainTracker: fakeChainTracker{},
	}
	taggedBEEF := overlay.TaggedBEEF{
		Topics: []string{"unknown-topic"},
		Beef:   createDummyBEEF(t),
	}

	// when:
	steak, err := sut.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.ErrorIs(t, err, engine.ErrUnknownTopic)
	require.Nil(t, steak)
}

func TestEngine_Submit_BroadcastFails_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	sut := &engine.Engine{
		Managers: map[string]engine.TopicManager{
			"test-topic": fakeManager{
				identifyAdmissableOutputsFunc: func(ctx context.Context, beef []byte, previousCoins []uint32) (overlay.AdmittanceInstructions, error) {
					return overlay.AdmittanceInstructions{
						OutputsToAdmit: []uint32{0},
					}, nil
				},
			},
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
		Broadcaster: fakeBroadcasterFail{
			broadcastFunc: func(tx *transaction.Transaction) (*transaction.BroadcastSuccess, *transaction.BroadcastFailure) {
				return nil, &transaction.BroadcastFailure{Description: "forced failure for testing"}
			},
			broadcastCtxFunc: func(ctx context.Context, tx *transaction.Transaction) (*transaction.BroadcastSuccess, *transaction.BroadcastFailure) {
				return nil, &transaction.BroadcastFailure{Description: "forced failure for testing"}
			},
		},
	}

	taggedBEEF := overlay.TaggedBEEF{
		Topics: []string{"test-topic"},
		Beef:   createDummyBEEF(t),
	}

	// when:
	steak, err := sut.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.Error(t, err)
	require.Nil(t, steak)
	require.EqualError(t, err, "forced failure for testing")
}

func TestEngine_Submit_OutputInsertFails_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	taggedBEEF, prevTxID := createDummyValidTaggedBEEF(t)

	sut := &engine.Engine{
		Managers: map[string]engine.TopicManager{
			"test-topic": fakeManager{
				identifyAdmissableOutputsFunc: func(ctx context.Context, beef []byte, previousCoins []uint32) (overlay.AdmittanceInstructions, error) {
					return overlay.AdmittanceInstructions{
						OutputsToAdmit: []uint32{0},
					}, nil
				},
			},
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
			deleteOutputFunc: func(ctx context.Context, outpoint *overlay.Outpoint, topic string) error {
				return nil
			},
		},
		ChainTracker: fakeChainTracker{},
	}

	// when:
	steak, err := sut.Submit(ctx, taggedBEEF, engine.SubmitModeCurrent, nil)

	// then:
	require.Error(t, err)
	require.Nil(t, steak)
	require.EqualError(t, err, "insert-failed")
}
