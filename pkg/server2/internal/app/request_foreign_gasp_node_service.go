package app

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

type RequestForeignGASPNodeDTO struct {
	GraphID     string
	TxID        string
	OutputIndex uint32
	Topic       string
}

// RequestForeignGASPNodeProvider defines the interface for requesting a foreign GASP node.
type RequestForeignGASPNodeProvider interface {
	ProvideForeignGASPNode(ctx context.Context, graphID, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error)
}

// RequestForeignGASPNodeService coordinates the request for a foreign GASP node.
type RequestForeignGASPNodeService struct {
	provider RequestForeignGASPNodeProvider
}

// RequestForeignGASPNode takes string representations of graphID and txID,
// validates and converts them to appropriate types, and requests a foreign GASP node.
// Returns the GASP node on success, an error if the provider fails.
func (s *RequestForeignGASPNodeService) RequestForeignGASPNode(ctx context.Context, dto RequestForeignGASPNodeDTO) (*core.GASPNode, error) {
	if dto.Topic == "" {
		return nil, NewRequestForeignGASPNodeMissingTopicError()
	}

	outpoint := &overlay.Outpoint{
		OutputIndex: dto.OutputIndex,
	}
	txid, err := chainhash.NewHashFromHex(dto.TxID)
	if err != nil {
		return nil, NewRequestForeignGASPNodeInvalidTxIDError()
	}
	outpoint.Txid = *txid

	graphID, err := overlay.NewOutpointFromString(dto.GraphID)
	if err != nil {
		return nil, NewRequestForeignGASPNodeInvalidGraphIDError()
	}

	node, err := s.provider.ProvideForeignGASPNode(ctx, graphID, outpoint, dto.Topic)
	if err != nil {
		return nil, NewRequestForeignGASPNodeProviderError(err)
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

// NewRequestForeignGASPNodeProviderError returns an Error indicating that the provider
// failed to retrieve a foreign GASP node.
func NewRequestForeignGASPNodeProviderError(err error) Error {
	return Error{
		errorType: ErrorTypeProviderFailure,
		err:       err.Error(),
		slug:      "Unable to retrieve foreign GASP node due to an error in the provider.",
	}
}

// NewRequestForeignGASPNodeMissingTopicError returns an Error indicating that the topic is missing.
func NewRequestForeignGASPNodeMissingTopicError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "topic is required",
		slug:      "The topic parameter is required for requesting a foreign GASP node.",
	}
}

// NewRequestForeignGASPNodeInvalidTxIDError returns an Error indicating that the txID is invalid.
func NewRequestForeignGASPNodeInvalidTxIDError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "invalid txID format",
		slug:      "The submitted txID is not a valid transaction hash.",
	}
}

// NewRequestForeignGASPNodeInvalidGraphIDError returns an Error indicating that the graphID is invalid.
func NewRequestForeignGASPNodeInvalidGraphIDError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "invalid graphID format",
		slug:      "The submitted graphID is not in a valid format (expected: txID.outputIndex).",
	}
}
