package app

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
)

const requestSyncResponseServiceDescriptor = "request-sync-response-service"

// RequestSyncResponseProvider defines the interface for requesting sync responses.
type RequestSyncResponseProvider interface {
	ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error)
}

// RequestSyncResponseService coordinates foreign sync response requests.
type RequestSyncResponseService struct {
	provider RequestSyncResponseProvider
}

// RequestSyncResponse requests a foreign sync response.
func (s *RequestSyncResponseService) RequestSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {
	response, err := s.provider.ProvideForeignSyncResponse(ctx, initialRequest, topic)
	if err != nil {
		return nil, NewRequestSyncResponseProviderError(err)
	}
	return response, nil
}

// NewRequestSyncResponseService creates a new RequestSyncResponseService.
func NewRequestSyncResponseService(provider RequestSyncResponseProvider) *RequestSyncResponseService {
	if provider == nil {
		panic("request sync response provider is nil")
	}
	return &RequestSyncResponseService{provider: provider}
}

// NewRequestSyncResponseProviderError returns an Error indicating that the foreign sync
// response provider failed to process the request.
func NewRequestSyncResponseProviderError(err error) Error {
	return Error{
		errorType: ErrorTypeProviderFailure,
		err:       err.Error(),
		slug:      "Unable to process sync response request due to an error in the overlay engine.",
	}
}
