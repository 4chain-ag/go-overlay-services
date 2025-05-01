package app

import (
	"context"
	"errors"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
)

// RequestSyncResponseProvider defines the contract that must be fulfilled

// to request a sync response from the overlay engine.

type RequestSyncResponseProvider interface {

	// ProvideForeignSyncResponse retrieves a foreign sync response based on the provided parameters.

	// It returns the response or an error if the request fails.

	ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error)
}

// RequestSyncResponseService is responsible for handling sync response requests

// using the configured RequestSyncResponseProvider.

type RequestSyncResponseService struct {
	provider RequestSyncResponseProvider
}

// RequestSyncResponse calls the configured provider's ProvideForeignSyncResponse method.

// If the provider fails, it wraps the error with ErrRequestSyncResponseProvider.

func (s *RequestSyncResponseService) RequestSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {

	response, err := s.provider.ProvideForeignSyncResponse(ctx, initialRequest, topic)

	if err != nil {

		return nil, errors.Join(err, ErrRequestSyncResponseProvider)

	}

	return response, nil

}

// NewRequestSyncResponseService creates a new instance of RequestSyncResponseService

// using the given RequestSyncResponseProvider. It panics if the provider is nil.

func NewRequestSyncResponseService(provider RequestSyncResponseProvider) *RequestSyncResponseService {

	if provider == nil {

		panic("request sync response provider is nil")

	}

	return &RequestSyncResponseService{provider: provider}

}

// ErrRequestSyncResponseProvider is returned when the RequestSyncResponseProvider

// fails to handle the sync response request.

var ErrRequestSyncResponseProvider = errors.New("failed to request sync response using provider")
