package app

import (
	"context"
	"errors"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

const DefaultSubmitTransactionTimeout = 5 * time.Second

// SubmitTransactionProvider defines the interface for sending a tagged transaction
// to the overlay engine for processing.
type SubmitTransactionProvider interface {
	Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error)
}

// SubmitTransactionService coordinates the transaction submission process
// using a SubmitTransactionProvider with a configurable timeout for awaiting a response.
type SubmitTransactionService struct {
	provider          SubmitTransactionProvider
	submitCallTimeout time.Duration
}

// SubmitTransaction submits a transaction to the overlay engine using the configured provider.
// It validates the provided topics, sends the transaction, and waits for a response (STEAK).
// Returns a non-nil *overlay.Steak on success.An error if topics are missing, invalid, the provider fails, or a timeout occurs.
func (s *SubmitTransactionService) SubmitTransaction(ctx context.Context, topics Topics, bytes ...byte) (*overlay.Steak, error) {
	err := topics.Verify()
	if err != nil {
		return nil, err
	}

	ch := make(chan *overlay.Steak, 1)
	_, err = s.provider.Submit(ctx, overlay.TaggedBEEF{Beef: bytes, Topics: topics}, engine.SubmitModeCurrent, func(steak *overlay.Steak) {
		ch <- steak
	})
	if err != nil {
		return nil, errors.Join(err, ErrSubmitTransactionProvider)
	}

	select {
	case steak := <-ch:
		return steak, nil
	case <-time.After(s.submitCallTimeout):
		return nil, ErrSubmitTransactionProviderTimeout
	}
}

// NewSubmitTransactionService creates a new SubmitTransactionService with the given provider and timeout.
// Panics if the provider is nil.
func NewSubmitTransactionService(provider SubmitTransactionProvider, timeout time.Duration) *SubmitTransactionService {
	if provider == nil {
		panic("submit transaction service provider is nil")
	}

	return &SubmitTransactionService{
		provider:          provider,
		submitCallTimeout: timeout,
	}
}

// Topics represents a list of topics that must be provided when submitting a transaction.
type Topics []string

// Verify ensures the topic list is non-empty and that each topic is non-blank.
// Returns ErrMissingTopics or ErrInvalidTopicFormat on failure.
func (tt Topics) Verify() error {
	if len(tt) == 0 {
		return ErrMissingTopics
	}

	for _, t := range tt {
		if len(t) == 0 { // TODO: Add more robust topic format check.
			return ErrInvalidTopicFormat
		}
	}
	return nil
}

var (
	// ErrSubmitTransactionProvider indicates a failure when submitting the transaction using the provider.
	ErrSubmitTransactionProvider = errors.New("failed to submit transaction using provider")

	// ErrSubmitTransactionProviderTimeout is returned if the provider does not respond within the configured timeout.
	ErrSubmitTransactionProviderTimeout = errors.New("submit transaction timeout occurred")

	// ErrMissingTopics is returned when no topics are provided.
	ErrMissingTopics = errors.New("provided topics cannot be an empty slice")

	// ErrMissingTransactionBytes is returned when the transaction data is empty.
	ErrMissingTransactionBytes = errors.New("provided tx bytes data cannot be an empty slice")

	// ErrInvalidTopicFormat is returned when a topic is empty or malformed.
	ErrInvalidTopicFormat = errors.New("invalid topic header format")
)
