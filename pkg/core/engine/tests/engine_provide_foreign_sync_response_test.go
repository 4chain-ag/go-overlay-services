package engine_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

func TestEngine_ProvideForeignSyncResponse_ShouldReturnUTXOList(t *testing.T) {
	t.Parallel()

	// given
	expectedOutpoint := &overlay.Outpoint{
		Txid:        fakeTxID(),
		OutputIndex: 1,
	}
	e := &engine.Engine{
		Storage: fakeStorage{
			findUTXOsForTopicFunc: func(ctx context.Context, topic string, since uint32, includeBEEF bool) ([]*engine.Output, error) {
				return []*engine.Output{
					{Outpoint: *expectedOutpoint},
				}, nil
			},
		},
	}

	// when
	resp, err := e.ProvideForeignSyncResponse(context.Background(), &core.GASPInitialRequest{Since: 0}, "test-topic")

	// then
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.UTXOList, 1)
	require.Equal(t, expectedOutpoint, resp.UTXOList[0])
}

func TestEngine_ProvideForeignSyncResponse_ShouldReturnError_WhenStorageFails(t *testing.T) {
	t.Parallel()

	// given
	expectedError := errors.New("storage failed")
	e := &engine.Engine{
		Storage: fakeStorage{
			findUTXOsForTopicFunc: func(ctx context.Context, topic string, since uint32, includeBEEF bool) ([]*engine.Output, error) {
				return nil, expectedError
			},
		},
	}

	// when
	resp, err := e.ProvideForeignSyncResponse(context.Background(), &core.GASPInitialRequest{Since: 0}, "test-topic")

	// then
	require.Error(t, err)
	require.Nil(t, resp)
	require.Equal(t, expectedError, err)
}
