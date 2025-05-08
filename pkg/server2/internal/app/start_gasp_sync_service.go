package app

import (
	"context"
)

// startGASPSyncServiceDescriptor is the service descriptor label used for identifying
// the start GASP sync service in logs, metrics, or tracing contexts.
const startGASPSyncServiceDescriptor = "start-gasp-sync-service"

// StartGASPSyncProvider defines the interface for triggering GASP sync.
type StartGASPSyncProvider interface {
	StartGASPSync(ctx context.Context) error
}

// StartGASPSyncService coordinates the GASP synchronization process.
type StartGASPSyncService struct {
	provider StartGASPSyncProvider
}

// StartGASPSync initiates the GASP synchronization process using the configured provider.
// Returns nil on success, an error if the provider fails.
func (s *StartGASPSyncService) StartGASPSync(ctx context.Context) error {
	if err := s.provider.StartGASPSync(ctx); err != nil {
		return NewProviderFailureError(startGASPSyncServiceDescriptor, err.Error())
	}
	return nil
}

// NewStartGASPSyncService creates a new StartGASPSyncService with the given provider.
// Panics if the provider is nil.
func NewStartGASPSyncService(provider StartGASPSyncProvider) *StartGASPSyncService {
	if provider == nil {
		panic("start GASP sync service provider is nil")
	}

	return &StartGASPSyncService{
		provider: provider,
	}
}
