package app

import (
	"context"
	"errors"
	"time"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/transaction"
)

const DefaultArcIngestTimeout = 10 * time.Second

// NewMerkleProofProvider defines the contract for handling new Merkle proofs.
// It allows the overlay engine to verify mined transactions and maintain
// a chain-of-custody for transaction outputs.
type NewMerkleProofProvider interface {
	HandleNewMerkleProof(ctx context.Context, txid *chainhash.Hash, proof *transaction.MerklePath) error
}

// ArcIngestService coordinates the processing of ARC ingest requests, including
// validation of incoming request data, conversion into the appropriate format,
// and forwarding to the provider for processing.
type ArcIngestService struct {
	provider        NewMerkleProofProvider
	responseTimeout time.Duration
}

// HandleArcIngest processes the ARC ingest request by passing the transaction ID and Merkle proof
// to the NewMerkleProofProvider for verification and processing.
// It returns an error if the processing fails or times out.
func (s *ArcIngestService) HandleArcIngest(ctx context.Context, txID string, merklePath string, blockHeight uint32) error {
	hash, err := chainhash.NewHashFromHex(txID)
	if err != nil {
		return NewInvalidTxIDFormatError(err.Error())
	}

	if len(txID) != chainhash.MaxHashStringSize {
		return NewInvalidTxIDLengthError()
	}

	path, err := transaction.NewMerklePathFromHex(merklePath)
	if err != nil {
		return NewInvalidMerklePathFormatError(err.Error())
	}
	path.BlockHeight = blockHeight

	ctxWithTimeout, cancel := context.WithTimeout(ctx, s.responseTimeout)
	defer cancel()

	err = s.provider.HandleNewMerkleProof(ctxWithTimeout, hash, path)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return NewMerkleProofProcessingTimeoutError(s.responseTimeout)
		}
		if errors.Is(err, context.Canceled) {
			return NewMerkleProofProcessingCanceledError()
		}
		return NewMerkleProofProcessingFailedError(err.Error())
	}

	return nil
}

// NewArcIngestService creates a new ArcIngestService with the given provider and timeout.
// Panics if the provider is nil.
func NewArcIngestService(provider NewMerkleProofProvider, timeout time.Duration) *ArcIngestService {
	if provider == nil {
		panic("arc ingest service provider is nil")
	}

	return &ArcIngestService{
		provider:        provider,
		responseTimeout: timeout,
	}
}

// NewInvalidTxIDFormatError returns an Error indicating that the transaction ID is not in a valid format.
func NewInvalidTxIDFormatError(err string) Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       err,
		slug:      "Invalid transaction ID format. Please provide a valid transaction ID.",
	}
}

// NewInvalidTxIDLengthError returns an Error indicating that the transaction ID does not match the expected length.
func NewInvalidTxIDLengthError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "Invalid transaction ID length",
		slug:      "The transaction ID does not match the expected length. Please check and try again.",
	}
}

// NewInvalidMerklePathFormatError returns an Error indicating that the Merkle path is malformed.
func NewInvalidMerklePathFormatError(err string) Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       err,
		slug:      "The Merkle path format is invalid. Please provide a valid Merkle path.",
	}
}

// NewMerkleProofProcessingTimeoutError returns an Error indicating that Merkle proof processing timed out.
func NewMerkleProofProcessingTimeoutError(timeout time.Duration) Error {
	return Error{
		errorType: ErrorTypeOperationTimeout,
		err:       "Merkle proof processing timed out after " + timeout.String(),
		slug:      "The Merkle proof processing request exceeded the timeout limit.",
	}
}

// NewMerkleProofProcessingCanceledError returns an Error indicating that Merkle proof processing was canceled.
func NewMerkleProofProcessingCanceledError() Error {
	return Error{
		errorType: ErrorTypeUnknown,
		err:       "Merkle proof processing was canceled",
		slug:      "The Merkle proof processing was canceled. Please try again later.",
	}
}

// NewMerkleProofProcessingFailedError returns an Error indicating that Merkle proof processing failed.
func NewMerkleProofProcessingFailedError(err string) Error {
	return Error{
		errorType: ErrorTypeProviderFailure,
		err:       err,
		slug:      "Internal server error occurred during Merkle proof processing. Please try again later or contact support.",
	}
}
