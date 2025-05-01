package app

import (
	"context"
	"errors"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// RequestForeignGASPNodeProvider defines the contract that must be fulfilled
// to send a foreign GASP node request to the overlay engine.
type RequestForeignGASPNodeProvider interface {
	// ProvideForeignGASPNode retrieves a foreign GASP node based on the provided parameters.
	// It returns the node or an error if the request fails.
	ProvideForeignGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error)
}

// RequestForeignGASPNodeService is responsible for handling foreign GASP node requests
// using the configured RequestForeignGASPNodeProvider.
type RequestForeignGASPNodeService struct {
	provider RequestForeignGASPNodeProvider
}

// RequestForeignGASPNode calls the configured provider's ProvideForeignGASPNode method.
// If the provider fails, it wraps the error with ErrRequestForeignGASPNodeProvider.
func (s *RequestForeignGASPNodeService) RequestForeignGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error) {
	node, err := s.provider.ProvideForeignGASPNode(ctx, graphID, outpoint, topic)
	if err != nil {
		return nil, errors.Join(err, ErrRequestForeignGASPNodeProvider)
	}
	return node, nil
}

// NewRequestForeignGASPNodeService creates a new instance of RequestForeignGASPNodeService
// using the given RequestForeignGASPNodeProvider. It panics if the provider is nil.
func NewRequestForeignGASPNodeService(provider RequestForeignGASPNodeProvider) *RequestForeignGASPNodeService {
	if provider == nil {
		panic("request foreign GASP node provider is nil")
	}

	return &RequestForeignGASPNodeService{provider: provider}
}

// ErrRequestForeignGASPNodeProvider is returned when the RequestForeignGASPNodeProvider
// fails to handle the foreign GASP node request.
var ErrRequestForeignGASPNodeProvider = errors.New("failed to request foreign GASP node using provider") 
