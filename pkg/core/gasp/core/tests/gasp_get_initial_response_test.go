package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

func TestGASP_GetInitialResponse_Success(t *testing.T) {
	// given:
	ctx := context.Background()
	request := &core.GASPInitialRequest{
		Version: 1,
		Since:   10,
	}

	expectedResponse := &core.GASPInitialResponse{
		UTXOList: []*overlay.Outpoint{
			{OutputIndex: 1},
			{OutputIndex: 2},
		},
		Since: 0,
	}

	sut := core.NewGASP(core.GASPParams{
		Version: ptr(1),
		Storage: fakeStorage{
			findKnownUTXOsFunc: func(ctx context.Context, since uint32) ([]*overlay.Outpoint, error) {
				return expectedResponse.UTXOList, nil
			},
		},
	})

	// when:
	actualResp, err := sut.GetInitialResponse(ctx, request)

	// then:
	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResp)
}

func TestGASP_GetInitialResponse_VersionMismatch_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	request := &core.GASPInitialRequest{
		Version: 99, // wrong version
		Since:   0,
	}
	sut := core.NewGASP(core.GASPParams{
		Version: ptr(1),
		Storage: fakeStorage{},
	})

	// when:
	actualResp, err := sut.GetInitialResponse(ctx, request)

	// then:
	require.IsType(t, &core.GASPVersionMismatchError{}, err)
	require.Nil(t, actualResp)
}

func TestGASP_GetInitialResponse_StorageFailure_ShouldReturnError(t *testing.T) {
	// given:
	ctx := context.Background()
	request := &core.GASPInitialRequest{
		Version: 1,
		Since:   0,
	}

	expectedErr := errors.New("forced storage error")
	sut := core.NewGASP(core.GASPParams{
		Version: ptr(1),
		Storage: fakeStorage{
			findKnownUTXOsFunc: func(ctx context.Context, since uint32) ([]*overlay.Outpoint, error) {
				return nil, expectedErr
			},
		},
	})

	// when:
	actualResp, err := sut.GetInitialResponse(ctx, request)

	// then:
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, actualResp)
}

func ptr(i int) *int {
	return &i
}
