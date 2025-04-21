package engine_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/advertiser"
	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

func TestEngine_SyncAdvertisements_ShouldReturnNil_WhenAdvertiserIsNil(t *testing.T) {
	// given
	sut := &engine.Engine{
		Advertiser: nil,
	}

	// when
	err := sut.SyncAdvertisements(context.Background())

	// then
	require.NoError(t, err)
}

func TestEngine_SyncAdvertisements_ShouldNotFail_WhenCreateAdvertisementsFails(t *testing.T) {
	// given
	sut := &engine.Engine{
		Advertiser: fakeAdvertiser{
			findAllAdvertisementsFunc: func(protocol overlay.Protocol) ([]*advertiser.Advertisement, error) {
				return []*advertiser.Advertisement{}, nil
			},
			createAdvertisementsFunc: func(data []*advertiser.AdvertisementData) (overlay.TaggedBEEF, error) {
				return overlay.TaggedBEEF{}, errors.New("invalid-atomic-beef")
			},
		},
		Managers:   map[string]engine.TopicManager{"test-topic": fakeTopicManager{}},
		HostingURL: "http://localhost",
	}

	// when
	err := sut.SyncAdvertisements(context.Background())

	// then
	require.NoError(t, err)
}

func TestEngine_SyncAdvertisements_ShouldCompleteSuccessfully(t *testing.T) {
	// given
	sut := &engine.Engine{
		Advertiser: fakeAdvertiser{
			findAllAdvertisementsFunc: func(protocol overlay.Protocol) ([]*advertiser.Advertisement, error) {
				return []*advertiser.Advertisement{}, nil
			},
			createAdvertisementsFunc: func(data []*advertiser.AdvertisementData) (overlay.TaggedBEEF, error) {
				return overlay.TaggedBEEF{}, nil
			},
			revokeAdvertisementsFunc: func(data []*advertiser.Advertisement) (overlay.TaggedBEEF, error) {
				return overlay.TaggedBEEF{}, nil
			},
		},
		Managers:       map[string]engine.TopicManager{"test-topic": fakeTopicManager{}},
		LookupServices: map[string]engine.LookupService{"test-service": fakeLookupService{}},
		HostingURL:     "http://localhost",
	}

	// when
	err := sut.SyncAdvertisements(context.Background())

	// then
	require.NoError(t, err)
}

func TestEngine_SyncAdvertisements_ShouldLogAndContinue_WhenCreateOrRevokeFails(t *testing.T) {
	// given
	sut := &engine.Engine{
		Advertiser: fakeAdvertiser{
			findAllAdvertisementsFunc: func(protocol overlay.Protocol) ([]*advertiser.Advertisement, error) {
				return []*advertiser.Advertisement{}, nil
			},
			createAdvertisementsFunc: func(data []*advertiser.AdvertisementData) (overlay.TaggedBEEF, error) {
				return overlay.TaggedBEEF{}, errors.New("create failed")
			},
			revokeAdvertisementsFunc: func(data []*advertiser.Advertisement) (overlay.TaggedBEEF, error) {
				return overlay.TaggedBEEF{}, errors.New("revoke failed")
			},
		},
		Managers:       map[string]engine.TopicManager{"test-topic": fakeTopicManager{}},
		LookupServices: map[string]engine.LookupService{"test-service": fakeLookupService{}},
		HostingURL:     "http://localhost",
	}

	// when
	err := sut.SyncAdvertisements(context.Background())

	// then
	require.NoError(t, err)
}
