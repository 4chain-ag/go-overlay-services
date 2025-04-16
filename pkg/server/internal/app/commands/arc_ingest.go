package commands

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/transaction"
)

// TODO: 1. Rewrite unit tests
// TODO: 2. Rewrite error handling impl

var (

	// ErrMissingRequiredFields is returned when required fields are missing.
	ErrMissingRequiredFields = errors.New("missing required fields: txid and merklePath are required")

	// ErrInvalidTxIDFormat is returned when the txid has an invalid format.
	ErrInvalidTxIDFormat = errors.New("invalid TxID format")

	// ErrInvalidMerklePathFormat is returned when the merkle path has an invalid format.
	ErrInvalidMerklePathFormat = errors.New("invalid Merkle path format")

	// ErrRequestTimeout is returned when the request processing exceeds the timeout.
	ErrRequestTimeout = errors.New("request processing timed out")
)

// ArcIngestHandlerResponse defines the response for the ARC ingest endpoint.
type ArcIngestHandlerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ArcIngestRequest defines the expected request structure for the ARC ingest endpoint.
type ArcIngestRequest struct {
	Txid        string `json:"txid"`
	MerklePath  string `json:"merklePath"`
	BlockHeight uint32 `json:"blockHeight"`
}

func (a *ArcIngestRequest) IsEmpty() bool {
	return a == &ArcIngestRequest{}
}

func (a *ArcIngestRequest) IsValid() error {
	switch {
	case a.IsEmpty():
		return errors.New("empty")
	case len(a.Txid) == 0:
		return ErrMissingRequiredFields
	case len(a.MerklePath) == 0:
		return ErrMissingRequiredFields
	default:
		return nil
	}
}

// NewMerkleProofProvider defines the contract for processing new merkle proofs.
// This interface allows the overlay engine to verify mined transactions and maintain
// a chain-of-custody for outputs.
type NewMerkleProofProvider interface {
	HandleNewMerkleProof(ctx context.Context, txid *chainhash.Hash, proof *transaction.MerklePath) error
}

// ArcIngestHandler processes new merkle proofs, validating requests and
// forwarding them to the overlay engine.
type ArcIngestHandler struct {
	provider         NewMerkleProofProvider
	requestBodyLimit int64
	responseTimeout  time.Duration
}

func (h *ArcIngestHandler) decodeArcIngestRequest(body io.Reader) (*ArcIngestRequest, error) {
	reader := jsonutil.BodyReader{
		Body:      body,
		ReadLimit: jsonutil.RequestBodyLimit1GB,
	}

	bb, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("body reader failure: %w", err)
	}

	var dst ArcIngestRequest
	dec := json.NewDecoder(bytes.NewBuffer(bb))
	if err := dec.Decode(&dst); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	err = dst.IsValid()
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

// ProcessMerkleProof parses and validates the merkle proof data before forwarding to the provider.
func (h *ArcIngestHandler) ProcessMerkleProof(ctx context.Context, req *ArcIngestRequest) error {
	txidBytes, err := hex.DecodeString(req.Txid)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidTxIDFormat, err.Error())
	}

	if len(txidBytes) != chainhash.HashSize {
		return fmt.Errorf("%w: invalid txid length, got %d bytes, expected %d",
			ErrInvalidTxIDFormat, len(txidBytes), chainhash.HashSize)
	}

	var txidHash chainhash.Hash
	copy(txidHash[:], txidBytes)

	if req.MerklePath == "" {
		return fmt.Errorf("%w: merkle path cannot be empty", ErrInvalidMerklePathFormat)
	}

	merklePath, err := transaction.NewMerklePathFromHex(req.MerklePath)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidMerklePathFormat, err.Error())
	}

	merklePath.BlockHeight = req.BlockHeight

	err = h.provider.HandleNewMerkleProof(ctx, &txidHash, merklePath)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return ErrRequestTimeout
		}
		return fmt.Errorf("failed to process merkle proof: %w", err)
	}

	return nil
}

// Handle processes an ARC ingest request with timeout handling and appropriate error responses.
func (h *ArcIngestHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonutil.SendHTTPResponse(w, http.StatusMethodNotAllowed, ArcIngestHandlerResponse{
			Status:  "error",
			Message: ErrInvalidHTTPMethod.Error(),
		})
		return
	}

	req, err := h.decodeArcIngestRequest(r.Body)
	if errors.Is(err, jsonutil.ErrRequestBodyTooLarge) {
		jsonutil.SendHTTPResponse(w, http.StatusRequestEntityTooLarge, ArcIngestHandlerResponse{
			Status:  "error",
			Message: jsonutil.ErrRequestBodyTooLarge.Error(),
		})
		return
	}

	if err != nil {
		jsonutil.SendHTTPResponse(w, http.StatusBadRequest, ArcIngestHandlerResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	resultChan := make(chan error, 1)
	defer close(resultChan)
	go func() { resultChan <- h.ProcessMerkleProof(r.Context(), req) }()

	select {
	case err := <-resultChan:
		if err != nil {
			statusCode := http.StatusInternalServerError
			errorMessage := "Failed to process merkle proof: " + err.Error()

			if errors.Is(err, ErrRequestTimeout) {
				statusCode = http.StatusGatewayTimeout
				errorMessage = "Request processing timed out after " + h.responseTimeout.String()
			} else if errors.Is(err, ErrInvalidTxIDFormat) || errors.Is(err, ErrInvalidMerklePathFormat) {
				statusCode = http.StatusBadRequest
			}

			jsonutil.SendHTTPResponse(w, statusCode, ArcIngestHandlerResponse{
				Status:  "error",
				Message: errorMessage,
			})
			return
		}

		jsonutil.SendHTTPResponse(w, http.StatusOK, ArcIngestHandlerResponse{
			Status:  "success",
			Message: "Transaction status updated",
		})

	case <-time.After(h.responseTimeout):
		jsonutil.SendHTTPResponse(w, http.StatusGatewayTimeout, ArcIngestHandlerResponse{
			Status:  "error",
			Message: "Request processing timed out after " + h.responseTimeout.String(),
		})
	}
}

// ArcIngestHandlerOption defines a function that can configure an ArcIngestHandler.
type ArcIngestHandlerOption func(h *ArcIngestHandler)

// WithArcResponseTimeout configures the timeout duration for merkle proof processing.
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

// NewArcIngestHandler returns an instance of an ArcIngestHandler, utilizing
// either a NewMerkleProofProvider or an engine.OverlayEngineProvider.
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
