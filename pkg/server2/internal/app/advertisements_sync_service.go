package app

import (
	"context"
)

// syncAdvertisementsProviderDescriptor is the service descriptor label used for identifying
// the submit transaction service in logs, metrics, or tracing contexts.
const syncAdvertisementsProviderDescriptor = "advertisements-sync-service"

// SyncAdvertisementsProvider defines the contract that must be fulfilled
// to send a synchronize advertisements request to the overlay engine for further processing.
type SyncAdvertisementsProvider interface {
	// SyncAdvertisements triggers the advertisement synchronization process.
	// It returns an error if the synchronization fails.
	SyncAdvertisements(ctx context.Context) error
}

// AdvertisementsSyncService is responsible for synchronizing advertisements
// using the configured SyncAdvertisementsProvider.
type AdvertisementsSyncService struct {
	provider SyncAdvertisementsProvider
}

// SyncAdvertisements delegates the advertisement synchronization task to the configured provider.
// It calls the provider's SyncAdvertisements method with the given context.
// If the provider returns an error, it wraps the error as a ProviderFailureError
// using the syncAdvertisementsProviderDescriptor label for context.
func (a *AdvertisementsSyncService) SyncAdvertisements(ctx context.Context) error {
	err := a.provider.SyncAdvertisements(ctx)
	if err != nil {
		return NewProviderFailureError(syncAdvertisementsProviderDescriptor, err.Error())
	}
	return nil
}

// NewAdvertisementsSyncService creates a new instance of AdvertisementsSyncServcie
// using the given SyncAdvertisementsProvider. It panics if the provider is nil.
func NewAdvertisementsSyncService(provider SyncAdvertisementsProvider) *AdvertisementsSyncService {
	if provider == nil {
		panic("sync advertisements provider is nil")
	}

	return &AdvertisementsSyncService{provider: provider}
}
