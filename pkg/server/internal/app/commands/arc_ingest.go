package commands

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/transaction"
)

var (
	// ErrInvalidTxIDFormat is returned when the transaction ID is not in a valid format (e.g., not hexadecimal).
	ErrInvalidTxIDFormat = errors.New("invalid transaction ID format")
	// ErrInvalidTxIDLength is returned when the transaction ID does not match the expected length (typically 64 characters for a SHA-256 hash).
	ErrInvalidTxIDLength = errors.New("invalid transaction ID length")
	// ErrInvalidMerklePathFormat is returned when the Merkle path is malformed or does not conform to the expected structure.
	ErrInvalidMerklePathFormat = errors.New("invalid Merkle path format")
	// ErrMissingRequiredRequestFieldsDefinition is returned when the request body is missing
	// required fields, such as the transaction ID and Merkle path.
	ErrMissingRequiredRequestFieldsDefinition = errors.New("missing required fields: txid, merkle path")
	// ErrMissingRequiredTxIDFieldDefinition is returned when the request body is missing
	// the required transaction ID field.
	ErrMissingRequiredTxIDFieldDefinition = errors.New("missing required field: txid")
	// ErrMissingRequiredMerklePathFieldDefinition is returned when the request body is missing
	// the required Merkle path field.
	ErrMissingRequiredMerklePathFieldDefinition = errors.New("missing required field: merkle path")
)

// ArcIngestHandlerResponse represents the response format for the ArcIngestHandler,
// containing the status of the operation and a message providing additional context.
type ArcIngestHandlerResponse struct {
	Status  string `json:"status"`  // The status of the request (e.g., "success", "error")
	Message string `json:"message"` // A message providing additional information about the result
}

// ArcIngestRequest defines the expected structure for the ARC ingest request body,
// containing the transaction ID, Merkle path, and block height. This structure
// is used to validate and process incoming ARC ingest requests.
type ArcIngestRequest struct {
	TxID        string `json:"txid"`        // Transaction ID in hexadecimal format
	MerklePath  string `json:"merklePath"`  // Merkle path as a hex string
	BlockHeight uint32 `json:"blockHeight"` // Block height associated with the Merkle path
}

// IsMerklePathEmpty checks if the MerklePath field is empty, indicating missing data.
func (r *ArcIngestRequest) IsMerklePathEmpty() bool {
	return r.MerklePath == ""
}

// IsTxIDEmpty checks if the TxID field is empty, indicating missing data.
func (r *ArcIngestRequest) IsTxIDEmpty() bool {
	return r.TxID == ""
}

// IsEmpty checks if all fields in the ArcIngestRequest are zero or empty values,
// signifying an invalid or incomplete request.
func (r *ArcIngestRequest) IsEmpty() bool {
	return r.TxID == "" && r.MerklePath == "" && r.BlockHeight == 0
}

// Validate checks if all required fields are present and valid. It returns an
// error if any of the required fields are missing or improperly formatted.
func (r *ArcIngestRequest) Validate() error {
	if r.IsEmpty() {
		return ErrMissingRequiredRequestFieldsDefinition
	}
	if r.IsTxIDEmpty() {
		return ErrMissingRequiredTxIDFieldDefinition
	}
	if r.IsMerklePathEmpty() {
		return ErrMissingRequiredMerklePathFieldDefinition
	}
	return nil
}

// MerklePathStruct parses the MerklePath hex string and returns the Merkle path
// structure. It also sets the associated block height on the resulting structure.
func (r *ArcIngestRequest) MerklePathStruct() (*transaction.MerklePath, error) {
	path, err := transaction.NewMerklePathFromHex(r.MerklePath)
	if err != nil {
		return nil, errors.Join(err, ErrInvalidMerklePathFormat)
	}
	path.BlockHeight = r.BlockHeight
	return path, nil
}

// TxIDHash converts the transaction ID (TxID) from a hex string to a chainhash.Hash
// representation. It returns an error if the string is improperly formatted or has
// an invalid length.
func (r *ArcIngestRequest) TxIDHash() (chainhash.Hash, error) {
	bb, err := hex.DecodeString(r.TxID)
	if err != nil {
		return chainhash.Hash{}, errors.Join(err, ErrInvalidTxIDFormat)
	}
	if len(bb) != chainhash.HashSize {
		return chainhash.Hash{}, ErrInvalidTxIDLength
	}
	var hash chainhash.Hash
	copy(hash[:], bb)
	return hash, nil
}

// NewMerkleProofProvider defines the contract for handling new Merkle proofs.
// It allows the overlay engine to verify mined transactions and maintain
// a chain-of-custody for transaction outputs.
type NewMerkleProofProvider interface {
	HandleNewMerkleProof(ctx context.Context, txid *chainhash.Hash, proof *transaction.MerklePath) error
}

// ArcIngestHandler orchestrates the processing of ARC ingest requests, including
// validation of incoming request bodies, conversion of the data into the
// appropriate format, and forwarding the data to the overlay engine for processing.
type ArcIngestHandler struct {
	provider         NewMerkleProofProvider
	requestBodyLimit int64
	responseTimeout  time.Duration
}

// decode reads and decodes the request body into the provided destination
// struct. It ensures the request body size is within the allowed limit.
func (h *ArcIngestHandler) decode(body io.Reader, dst *ArcIngestRequest) error {
	reader := jsonutil.LimitedBodyReader{
		Body:      body,
		ReadLimit: jsonutil.RequestBodyLimit1GB,
	}
	bb, err := reader.Read()
	if err != nil {
		return fmt.Errorf("body reader failure: %w", err)
	}
	err = jsonutil.DecodeBytes(bb, dst)
	if err != nil {
		return fmt.Errorf("decode bytes failure: %w", err)
	}
	return nil
}

// buildHandleNewMerkleProofArguments decodes the request body, validates the
// fields, and returns the Merkle proof data (TxID and MerklePath) for further
// processing by the provider.
func (h *ArcIngestHandler) buildHandleNewMerkleProofArguments(body io.Reader) (chainhash.Hash, *transaction.MerklePath, error) {
	var dst ArcIngestRequest
	err := h.decode(body, &dst)
	if err != nil {
		return chainhash.Hash{}, nil, err
	}
	err = dst.Validate()
	if err != nil {
		return chainhash.Hash{}, nil, err
	}
	txIDHash, err := dst.TxIDHash()
	if err != nil {
		return chainhash.Hash{}, nil, err
	}
	merklePath, err := dst.MerklePathStruct()
	if err != nil {
		return chainhash.Hash{}, nil, err
	}
	return txIDHash, merklePath, nil
}

// Handle processes an ARC ingest request, handling the validation, converting
// the request data into the correct format, and invoking the provider to handle
// the new Merkle proof. It also manages timeout and error responses.
func (h *ArcIngestHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonutil.SendHTTPResponse(w, http.StatusMethodNotAllowed, NewFailureArcIngestHandlerResponse(ErrInvalidHTTPMethod.Error()))
		return
	}
	txIDHash, merklePath, err := h.buildHandleNewMerkleProofArguments(r.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("[ArcIngest] Failed to build Merkle proof arguments: %v", err))

		switch {
		case errors.Is(err, jsonutil.ErrBodyReaderFailure),
			errors.Is(err, jsonutil.JSONDecoderFailure):
			jsonutil.SendHTTPResponse(w, http.StatusInternalServerError, NewInternalFailureArcIngestHandlerResponse())
			return
		case errors.Is(err, ErrInvalidMerklePathFormat):
			jsonutil.SendHTTPResponse(w, http.StatusBadRequest, NewFailureArcIngestHandlerResponse(ErrInvalidMerklePathFormat.Error()))
			return
		case errors.Is(err, ErrInvalidTxIDFormat):
			jsonutil.SendHTTPResponse(w, http.StatusBadRequest, NewFailureArcIngestHandlerResponse(ErrInvalidTxIDFormat.Error()))
			return
		default:
			jsonutil.SendHTTPResponse(w, http.StatusBadRequest, NewFailureArcIngestHandlerResponse(err.Error()))
			return
		}
	}
	errCh := make(chan error, 1)
	defer close(errCh)
	go func() { errCh <- h.provider.HandleNewMerkleProof(r.Context(), &txIDHash, merklePath) }()
	select {
	case err := <-errCh:
		if err == nil {
			jsonutil.SendHTTPResponse(w, http.StatusOK, NewSuccessArcIngestHandlerResponse())
			return
		}
		slog.Error(fmt.Sprintf("[ArcIngest] Merkle proof processing failed: %v", err))

		jsonutil.SendHTTPResponse(w, http.StatusInternalServerError, NewFailureArcIngestHandlerResponse(http.StatusText(http.StatusInternalServerError)))
	case <-time.After(h.responseTimeout):
		jsonutil.SendHTTPResponse(w, http.StatusGatewayTimeout, NewFailureArcIngestHandlerResponse(http.StatusText(http.StatusGatewayTimeout)))
	}
}

// ArcIngestHandlerOption defines a function that configures an ArcIngestHandler.
type ArcIngestHandlerOption func(h *ArcIngestHandler)

// WithArcResponseTimeout configures the timeout duration for Merkle proof processing.
func WithArcResponseTimeout(d time.Duration) ArcIngestHandlerOption {
	return func(h *ArcIngestHandler) {
		h.responseTimeout = d
	}
}

// WithArcRequestBodyLimit configures the maximum allowed size for request bodies.
func WithArcRequestBodyLimit(limit int64) ArcIngestHandlerOption {
	return func(h *ArcIngestHandler) {
		h.requestBodyLimit = limit
	}
}

// NewArcIngestHandler creates a new instance of an ArcIngestHandler, utilizing
// the provided Merkle proof provider and any optional configurations.
func NewArcIngestHandler(provider NewMerkleProofProvider, opts ...ArcIngestHandlerOption) (*ArcIngestHandler, error) {
	if provider == nil {
		return nil, fmt.Errorf("provider is nil")
	}
	h := ArcIngestHandler{
		provider:         provider,
		requestBodyLimit: jsonutil.RequestBodyLimit1GB,
		responseTimeout:  10 * time.Second,
	}
	for _, opt := range opts {
		opt(&h)
	}
	return &h, nil
}

// NewSuccessArcIngestHandlerResponse creates a success response for the ArcIngestHandler,
// indicating that the transaction status has been successfully updated.
func NewSuccessArcIngestHandlerResponse() ArcIngestHandlerResponse {
	return ArcIngestHandlerResponse{Status: "success", Message: "Transaction status updated"}
}

// NewFailureArcIngestHandlerResponse creates a failure response for the ArcIngestHandler,
// with the provided message describing the error that occurred during the process.
func NewFailureArcIngestHandlerResponse(message string) ArcIngestHandlerResponse {
	return ArcIngestHandlerResponse{
		Status:  "error",
		Message: message,
	}
}

// NewInternalFailureArcIngestHandlerResponse creates a failure response for the ArcIngestHandler,
// indicating an internal server error, with the default error message for HTTP 500.
func NewInternalFailureArcIngestHandlerResponse() ArcIngestHandlerResponse {
	return ArcIngestHandlerResponse{
		Status:  "error",
		Message: http.StatusText(http.StatusInternalServerError),
	}
}
