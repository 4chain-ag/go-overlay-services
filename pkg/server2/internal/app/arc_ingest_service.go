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
	provider NewMerkleProofProvider

	responseTimeout time.Duration
}

// HandleArcIngest processes the ARC ingest request by passing the transaction ID and Merkle proof

// to the NewMerkleProofProvider for verification and processing.

// It returns an error if the processing fails or times out.

func (s *ArcIngestService) HandleArcIngest(ctx context.Context, txID string, merklePath string, blockHeight uint32) error {

	//TODO: break down into smaller functions

	hash, err := chainhash.NewHashFromHex(txID)

	if err != nil {

		return errors.Join(err, ErrInvalidTxIDFormat)

	}

	if len(txID) != chainhash.MaxHashStringSize {

		return ErrInvalidTxIDLength

	}

	path, err := transaction.NewMerklePathFromHex(merklePath)

	if err != nil {

		return errors.Join(err, ErrInvalidMerklePathFormat)

	}

	path.BlockHeight = blockHeight

	ctxWithTimeout, cancel := context.WithTimeout(ctx, s.responseTimeout)

	defer cancel()

	err = s.provider.HandleNewMerkleProof(ctxWithTimeout, hash, path)

	if err != nil {

		if errors.Is(err, context.DeadlineExceeded) {

			return ErrMerkleProofProcessingTimeout

		}

		if errors.Is(err, context.Canceled) {

			return ErrMerkleProofProcessingCanceled

		}

		return errors.Join(err, ErrMerkleProofProcessingFailed)

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

		provider: provider,

		responseTimeout: timeout,
	}

}

var (

	// ErrInvalidTxIDFormat is returned when the transaction ID is not in a valid format.

	ErrInvalidTxIDFormat = errors.New("invalid transaction ID format")

	// ErrInvalidTxIDLength is returned when the transaction ID does not match the expected length.

	ErrInvalidTxIDLength = errors.New("invalid transaction ID length")

	// ErrInvalidMerklePathFormat is returned when the Merkle path is malformed.

	ErrInvalidMerklePathFormat = errors.New("invalid Merkle path format")

	// ErrMerkleProofProcessingTimeout is returned when Merkle proof processing times out.

	ErrMerkleProofProcessingTimeout = errors.New("Merkle proof processing timed out")

	// ErrMerkleProofProcessingCanceled is returned when Merkle proof processing is canceled.

	ErrMerkleProofProcessingCanceled = errors.New("Merkle proof processing canceled")

	// ErrMerkleProofProcessingFailed is returned when Merkle proof processing fails.

	ErrMerkleProofProcessingFailed = errors.New("Internal server error occurred during processing")
)
