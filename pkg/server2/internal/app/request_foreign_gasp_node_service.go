package app

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/chainhash"
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
	// Validate input parameters
	if topic == "" {
		return nil, NewMissingTopicError()
	}

	if graphID == nil {
		return nil, NewMissingGraphIDError()
	}

	if outpoint == nil {
		return nil, NewMissingOutpointError()
	}

	node, err := s.provider.ProvideForeignGASPNode(ctx, graphID, outpoint, topic)
	if err != nil {
		return nil, NewRequestForeignGASPNodeProviderError(err)
	}
	return node, nil
}

// RequestForeignGASPNodeWithStrings takes string representations of graphID and txID,
// validates and converts them to appropriate types, and requests a foreign GASP node.
func (s *RequestForeignGASPNodeService) RequestForeignGASPNodeWithStrings(
	ctx context.Context,
	graphIDStr string,
	txIDStr string,
	outputIndex uint32,
	topic string,
) (*core.GASPNode, error) {
	if topic == "" {
		return nil, NewMissingTopicError()
	}

	outpoint := &overlay.Outpoint{
		OutputIndex: outputIndex,
	}
	txid, err := chainhash.NewHashFromHex(txIDStr)
	if err != nil {
		return nil, NewInvalidTxIDError()
	}
	outpoint.Txid = *txid

	graphID, err := overlay.NewOutpointFromString(graphIDStr)
	if err != nil {
		return nil, NewInvalidGraphIDError()
	}

	return s.RequestForeignGASPNode(ctx, graphID, outpoint, topic)
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

// NewMissingTopicError returns an Error indicating that the topic is missing.
func NewMissingTopicError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "topic is required",
		slug:      "The topic parameter is required for requesting a foreign GASP node.",
	}
}

// NewMissingGraphIDError returns an Error indicating that the graphID is missing.
func NewMissingGraphIDError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "graphID is required",
		slug:      "The graphID parameter is required for requesting a foreign GASP node.",
	}
}

// NewMissingOutpointError returns an Error indicating that the outpoint is missing.
func NewMissingOutpointError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "outpoint is required",
		slug:      "The outpoint parameter is required for requesting a foreign GASP node.",
	}
}

// NewInvalidTxIDError returns an Error indicating that the txID is invalid.
func NewInvalidTxIDError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "invalid txID format",
		slug:      "The submitted txID is not a valid transaction hash.",
	}
}

// NewInvalidGraphIDError returns an Error indicating that the graphID is invalid.
func NewInvalidGraphIDError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "invalid graphID format",
		slug:      "The submitted graphID is not in a valid format (expected: txID.outputIndex).",
	}
}
