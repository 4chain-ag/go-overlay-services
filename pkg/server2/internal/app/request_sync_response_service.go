package app

import (
	"context"
	"math"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
)

type RequestSyncResponseDTO struct {
	Version int
	Since   int
	Topic   string
}

// RequestSyncResponseProvider defines the interface for requesting sync responses.
type RequestSyncResponseProvider interface {
	ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error)
}

// RequestSyncResponseService coordinates foreign sync response requests.
type RequestSyncResponseService struct {
	provider RequestSyncResponseProvider
}

// RequestSyncResponse requests a foreign sync response.
func (s *RequestSyncResponseService) RequestSyncResponse(ctx context.Context, dto *RequestSyncResponseDTO) (*core.GASPInitialResponse, error) {
	if dto.Topic == "" {
		return nil, NewRequestSyncResponseInvalidInputError()
	}

	version := dto.Version
	if version <= 0 {
		return nil, NewRequestSyncResponseInvalidVersionError()
	}

	since := dto.Since
	if since < 0 || since > math.MaxUint32 {
		return nil, NewRequestSyncResponseInvalidSinceError()
	}

	response, err := s.provider.ProvideForeignSyncResponse(ctx, &core.GASPInitialRequest{Version: version, Since: uint32(since)}, dto.Topic)
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

// NewRequestSyncResponseInvalidInputError returns an Error indicating that the topic is empty.
func NewRequestSyncResponseInvalidInputError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "topic cannot be empty",
		slug:      "A valid topic must be provided to request a sync response.",
	}
}

// NewRequestSyncResponseInvalidRequestError returns an Error indicating that the initial request is nil.
func NewRequestSyncResponseInvalidRequestError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "initial request cannot be nil",
		slug:      "A valid initial request must be provided to request a sync response.",
	}
}

// NewRequestSyncResponseInvalidVersionError returns an Error indicating that the initial request version is invalid.
func NewRequestSyncResponseInvalidVersionError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "initial request version must be greater than 0",
		slug:      "A valid version must be provided for the initial request to request a sync response.",
	}
}

// NewRequestSyncResponseInvalidSinceError returns an Error indicating that the initial request since is invalid.
func NewRequestSyncResponseInvalidSinceError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "initial request since must be greater than 0 and less than or equal to 4294967295",
		slug:      "A valid since value between 1 and 4294967295 must be provided for the initial request to request a sync response.",
	}
}

// NewRequestSyncResponseInvalidJSONError returns an Error indicating that the JSON input is invalid.
func NewRequestSyncResponseInvalidJSONError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "invalid JSON format",
		slug:      "The request body must contain valid JSON.",
	}
}
