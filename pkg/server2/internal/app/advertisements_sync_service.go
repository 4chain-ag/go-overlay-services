package app

import (
	"context"
	"errors"
)

// SyncAdvertisementsProvider defines the contract that must be fulfilled
// to send a synchronize advertisements request to the overlay engine for further processing.
type SyncAdvertisementsProvider interface {
	// SyncAdvertisements triggers the advertisement synchronization process.
	// It returns an error if the synchronization fails.
	SyncAdvertisements(ctx context.Context) error
}

// AdvertisementsSyncServcie is responsible for synchronizing advertisements
// using the configured SyncAdvertisementsProvider.
type AdvertisementsSyncServcie struct {
	provider SyncAdvertisementsProvider
}

// SyncAdvertisements calls the configured provider's SyncAdvertisements method.
// If the provider fails, it wraps the error with ErrSyncAdvertisementsProvider.
func (a *AdvertisementsSyncServcie) SyncAdvertisements(ctx context.Context) error {
	err := a.provider.SyncAdvertisements(ctx)
	if err != nil {
		return errors.Join(err, ErrSyncAdvertisementsProvider)
	}
	return nil
}

// NewAdvertisementsSyncServcie creates a new instance of AdvertisementsSyncServcie
// using the given SyncAdvertisementsProvider. It panics if the provider is nil.
func NewAdvertisementsSyncServcie(provider SyncAdvertisementsProvider) *AdvertisementsSyncServcie {
	if provider == nil {
		panic("sync advertisements provider is nil")
	}

	return &AdvertisementsSyncServcie{provider: provider}
}

// ErrSyncAdvertisementsProvider is returned when the SyncAdvertisementsProvider fails
// to handle the synchronize advertisements request.
var ErrSyncAdvertisementsProvider = errors.New("failed to sync advertisements using provider")
