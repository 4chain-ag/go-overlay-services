package app

import (
	"context"
	"errors"
)

// StartGASPSyncProvider defines the contract that must be fulfilled

// to send a GASP sync request to the overlay engine for further processing.

type StartGASPSyncProvider interface {

	// StartGASPSync triggers the GASP synchronization process.

	// It returns an error if the synchronization fails.

	StartGASPSync(ctx context.Context) error
}

// StartGASPSyncService is responsible for initiating GASP synchronization

// using the configured StartGASPSyncProvider.

type StartGASPSyncService struct {
	provider StartGASPSyncProvider
}

// StartGASPSync calls the configured provider's StartGASPSync method.

// If the provider fails, it wraps the error with ErrStartGASPSyncProvider.

func (s *StartGASPSyncService) StartGASPSync(ctx context.Context) error {

	err := s.provider.StartGASPSync(ctx)

	if err != nil {

		return errors.Join(err, ErrStartGASPSyncProvider)

	}

	return nil

}

// NewStartGASPSyncService creates a new instance of StartGASPSyncService

// using the given StartGASPSyncProvider. It panics if the provider is nil.

func NewStartGASPSyncService(provider StartGASPSyncProvider) *StartGASPSyncService {

	if provider == nil {

		panic("start GASP sync provider is nil")

	}

	return &StartGASPSyncService{provider: provider}

}

// ErrStartGASPSyncProvider is returned when the StartGASPSyncProvider fails

// to handle the GASP sync request.

var ErrStartGASPSyncProvider = errors.New("failed to start GASP sync using provider")
