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
		return NewStartGASPSyncProviderError(err)
	}
	return nil
}

// NewStartGASPSyncService creates a new StartGASPSyncService with the given provider.
// Returns an error if the provider is nil.
func NewStartGASPSyncService(provider StartGASPSyncProvider) (*StartGASPSyncService, error) {
	if provider == nil {
		return nil, NewStartGASPSyncNilProviderError()
	}

	return &StartGASPSyncService{
		provider: provider,
	}, nil
}

// NewStartGASPSyncProviderError returns an Error indicating that the configured provider
// failed to process a GASP sync request.
func NewStartGASPSyncProviderError(err error) Error {
	return Error{
		errorType: ErrorTypeProviderFailure,
		err:       err.Error(),
		slug:      "Unable to synchronize GASP due to an internal error. Please try again later or contact the support team.",
	}
}

// NewStartGASPSyncNilProviderError returns an Error indicating that the required provider was nil,
// which is invalid input when creating the service.
func NewStartGASPSyncNilProviderError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "start GASP sync service provider cannot be nil",
		slug:      "The required provider was not properly initialized",
	}
}
