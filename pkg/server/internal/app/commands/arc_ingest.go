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

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/transaction"
)

// RequestBodyLimitDefault defines the maximum allowed size for request bodies (1GB).
const requestBodyLimitDefault = 1000 * 1024 * 1024

var (
	// ErrArcIngestInvalidHTTPMethod is returned when an unsupported HTTP method is used.
	ErrArcIngestInvalidHTTPMethod = errors.New("invalid HTTP method")

	// ErrArcIngestRequestBodyRead is returned when there's an error reading the request body.
	ErrArcIngestRequestBodyRead = errors.New("failed to read request body")

	// ErrArcIngestRequestBodyTooLarge is returned when the request body exceeds the size limit.
	ErrArcIngestRequestBodyTooLarge = errors.New("request body too large")

	// ErrMissingRequiredFields is returned when required fields are missing.
	ErrMissingRequiredFields = errors.New("missing required fields: txid and merklePath are required")

	// ErrInvalidTxidFormat is returned when the txid has an invalid format.
	ErrInvalidTxidFormat = errors.New("invalid txid format")

	// ErrInvalidMerklePathFormat is returned when the merkle path has an invalid format.
	ErrInvalidMerklePathFormat = errors.New("invalid merkle path format")
	
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

// NewMerkleProofProvider defines the contract for processing new merkle proofs.
// This interface allows the overlay engine to verify mined transactions and maintain 
// a chain-of-custody for outputs.
type NewMerkleProofProvider interface {
	HandleNewMerkleProof(ctx context.Context, txid *chainhash.Hash, proof *transaction.MerklePath) error
}

// EngineProvider interface wrapper for compatibility with the overlay engine
type EngineProvider interface {
	HandleNewMerkleProof(ctx context.Context, txid *chainhash.Hash, proof *transaction.MerklePath) error
}

// OverlayEngineAdaptor adapts the OverlayEngineProvider to the NewMerkleProofProvider interface
type OverlayEngineAdaptor struct {
	Engine interface{}
}

// HandleNewMerkleProof delegates to the engine's HandleNewMerkleProof method
func (a *OverlayEngineAdaptor) HandleNewMerkleProof(ctx context.Context, txid *chainhash.Hash, proof *transaction.MerklePath) error {
	if engine, ok := a.Engine.(EngineProvider); ok {
		err := engine.HandleNewMerkleProof(ctx, txid, proof)
		if err != nil {
			return fmt.Errorf("engine HandleNewMerkleProof failed: %w", err)
		}
		return nil
	}
	return fmt.Errorf("engine does not implement HandleNewMerkleProof method")
}

// ArcIngestHandler processes new merkle proofs, validating requests and
// forwarding them to the overlay engine.
type ArcIngestHandler struct {
	provider         NewMerkleProofProvider
	requestBodyLimit int64
	responseTimeout  time.Duration
}

// DecodeAndValidateRequest reads and parses the request body, with size limit applied,
// and validates all required fields.
func (h *ArcIngestHandler) DecodeAndValidateRequest(r *http.Request) (*ArcIngestRequest, error) {
	if r.Method != http.MethodPost {
		return nil, ErrArcIngestInvalidHTTPMethod
	}

	reader := io.LimitReader(r.Body, h.requestBodyLimit+1)
	buff := make([]byte, 64*1024)
	var bodyBuffer bytes.Buffer
	var bytesRead int64

	for {
		n, err := reader.Read(buff)
		if n > 0 {
			bytesRead += int64(n)
			if bytesRead > h.requestBodyLimit {
				return nil, ErrArcIngestRequestBodyTooLarge
			}

			if _, inner := bodyBuffer.Write(buff[:n]); inner != nil {
				return nil, ErrArcIngestRequestBodyRead
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, ErrArcIngestRequestBodyRead
		}
	}

	var req ArcIngestRequest
	if err := json.Unmarshal(bodyBuffer.Bytes(), &req); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	if req.Txid == "" || req.MerklePath == "" {
		return nil, ErrMissingRequiredFields
	}

	return &req, nil
}

// ProcessMerkleProof parses and validates the merkle proof data before forwarding to the provider.
func (h *ArcIngestHandler) ProcessMerkleProof(ctx context.Context, req *ArcIngestRequest) error {
	txidBytes, err := hex.DecodeString(req.Txid)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidTxidFormat, err.Error())
	}
	
	if len(txidBytes) != chainhash.HashSize {
		return fmt.Errorf("%w: invalid txid length, got %d bytes, expected %d", 
			ErrInvalidTxidFormat, len(txidBytes), chainhash.HashSize)
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
	req, err := h.DecodeAndValidateRequest(r)
	if errors.Is(err, ErrArcIngestInvalidHTTPMethod) {
		jsonutil.SendHTTPResponse(w, http.StatusMethodNotAllowed, ArcIngestHandlerResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if errors.Is(err, ErrArcIngestRequestBodyTooLarge) {
		jsonutil.SendHTTPResponse(w, http.StatusRequestEntityTooLarge, ArcIngestHandlerResponse{
			Status:  "error",
			Message: err.Error(),
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

	ctx, cancel := context.WithTimeout(r.Context(), h.responseTimeout)
	defer cancel()

	resultChan := make(chan error, 1)
	
	go func() {
		resultChan <- h.ProcessMerkleProof(ctx, req)
	}()
	
	select {
	case err := <-resultChan:
		if err != nil {
			statusCode := http.StatusInternalServerError
			errorMessage := "Failed to process merkle proof: " + err.Error()
			
			if errors.Is(err, ErrRequestTimeout) {
				statusCode = http.StatusGatewayTimeout
				errorMessage = "Request processing timed out after " + h.responseTimeout.String()
			} else if errors.Is(err, ErrInvalidTxidFormat) || errors.Is(err, ErrInvalidMerklePathFormat) {
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
func NewArcIngestHandler(providerOrAdaptor interface{}, opts ...ArcIngestHandlerOption) (*ArcIngestHandler, error) {
    if providerOrAdaptor == nil {
        return nil, fmt.Errorf("provider is nil")
    }

    var provider NewMerkleProofProvider

    if p, ok := providerOrAdaptor.(NewMerkleProofProvider); ok {
        provider = p
    } else if engineProvider, ok := providerOrAdaptor.(engine.OverlayEngineProvider); ok {
        provider = &OverlayEngineAdaptor{
            Engine: engineProvider,
        }
    } else {
        return nil, fmt.Errorf("provider must implement either NewMerkleProofProvider or engine.OverlayEngineProvider")
    }

    h := &ArcIngestHandler{
        provider:         provider,
        requestBodyLimit: requestBodyLimitDefault,
        responseTimeout:  10 * time.Second,
    }
    for _, opt := range opts {
        opt(h)
    }
    return h, nil
}

// NewOverlayEngineAdaptor creates a new adaptor for an overlay engine
func NewOverlayEngineAdaptor(engine interface{}) *OverlayEngineAdaptor {
    return &OverlayEngineAdaptor{
        Engine: engine,
    }
}
