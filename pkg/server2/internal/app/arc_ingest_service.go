package app

import (
	"context"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/transaction"
)

// ArcIngestDTO defines the expected structure for the ARC ingest request body,
// containing the transaction ID, Merkle path, and block height. This structure
// is used to validate and process incoming ARC ingest requests.
type ArcIngestDTO struct {
	TxID        string
	MerklePath  string
	BlockHeight uint32
}

// ArcIngestProvider defines the contract for handling new Merkle proofs.
// It allows the overlay engine to verify mined transactions and maintain
// a chain-of-custody for transaction outputs.
type ArcIngestProvider interface {
	HandleNewMerkleProof(ctx context.Context, txid *chainhash.Hash, proof *transaction.MerklePath) error
}

// ArcIngestService coordinates the processing of ARC ingest requests, including
// validation of incoming request data, conversion into the appropriate format,
// and forwarding to the provider for processing.
type ArcIngestService struct {
	provider ArcIngestProvider
}

// HandleArcIngest processes the ARC ingest request by passing the transaction ID and Merkle proof
// to the ArcIngestProvider for verification and processing.
// It returns an error if the processing fails or times out.
func (s *ArcIngestService) HandleArcIngest(ctx context.Context, dto *ArcIngestDTO) error {
	hash, err := chainhash.NewHashFromHex(dto.TxID)
	if err != nil {
		return NewInvalidTxIDFormatError()
	}

	path, err := transaction.NewMerklePathFromHex(dto.MerklePath)
	if err != nil {
		return NewInvalidMerklePathFormatError()
	}

	path.BlockHeight = dto.BlockHeight

	err = s.provider.HandleNewMerkleProof(ctx, hash, path)
	if err != nil {
		return NewMerkleProofProcessingFailedError(err.Error())
	}

	return nil
}

// NewArcIngestService creates a new ArcIngestService with the given provider.
// Panics if the provider is nil.
func NewArcIngestService(provider ArcIngestProvider) *ArcIngestService {
	if provider == nil {
		panic("arc ingest service provider is nil")
	}

	return &ArcIngestService{
		provider: provider,
	}
}

// NewInvalidMerklePathFormatError returns an Error indicating that the Merkle path is malformed.
func NewInvalidMerklePathFormatError() Error {
	const str = "Invalid Merkle path format. Please provide a valid Merkle path."
	return NewIncorrectInputError(str, str)
}

// NewMerkleProofProcessingFailedError returns an Error indicating that Merkle proof processing failed.
func NewMerkleProofProcessingFailedError(err string) Error {
	return Error{
		errorType: ErrorTypeProviderFailure,
		err:       err,
		slug:      "Internal server error occurred during Merkle proof processing. Please try again later or contact support.",
	}
}

// NewInvalidTxIDFormatError returns an Error indicating that the transaction ID is not in a valid format.
func NewInvalidTxIDFormatError() Error {
	const str = "Invalid transaction ID format. Please provide a valid transaction ID."
	return NewIncorrectInputError(str, str)
}
