package app

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

const requestForeignGASPNodeServiceDescriptor = "request-foreign-gasp-node-service"

// RequestForeignGASPNodeProvider defines the interface for requesting a foreign GASP node.
type RequestForeignGASPNodeProvider interface {
	ProvideForeignGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error)
}

// RequestForeignGASPNodeService coordinates the request for a foreign GASP node.
type RequestForeignGASPNodeService struct {
	provider RequestForeignGASPNodeProvider
}

// RequestForeignGASPNode requests a foreign GASP node using the configured provider.
// Returns the GASP node on success, an error if the provider fails.
func (s *RequestForeignGASPNodeService) RequestForeignGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error) {
	node, err := s.provider.ProvideForeignGASPNode(ctx, graphID, outpoint, topic)
	if err != nil {
		return nil, NewProviderFailureError(requestForeignGASPNodeServiceDescriptor, err.Error())
	}
	return node, nil
}

// NewRequestForeignGASPNodeService creates a new RequestForeignGASPNodeService with the given provider.
// Panics if the provider is nil.
func NewRequestForeignGASPNodeService(provider RequestForeignGASPNodeProvider) *RequestForeignGASPNodeService {
	if provider == nil {
		panic("request foreign GASP node service provider is nil")
	}

	return &RequestForeignGASPNodeService{
		provider: provider,
	}
}
