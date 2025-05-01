package app

import (
	"context"
	"errors"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// SubmitTransactionProvider defines the contract that must be fulfilled
// to send a transaction request to the overlay engine for further processing.
type SubmitTransactionProvider interface {
	Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error)
}

// SubmitTransactionService provides functionality for validating and submitting transactions
// with topic-based tagging and a size-limited request body.
type SubmitTransactionService struct {
	engine            SubmitTransactionProvider
	submitCallTimeout time.Duration
}

func (s *SubmitTransactionService) SubmitTransaction(ctx context.Context, topics Topics, bytes ...byte) (*overlay.Steak, error) {
	err := topics.Verify()
	if err != nil {
		return nil, err
	}

	reader := &LimitedBytesReader{
		Bytes:     bytes,
		ReadLimit: ReadBodyLimit1GB,
	}

	readBytes, err := reader.Read()
	if err != nil {
		if errors.Is(err, ErrReaderMissingBytes) {
			return nil, ErrMissingTransactionBytes
		}

		if errors.Is(err, ErrReaderLimitExceeded) {
			return nil, ErrReaderLimitExceeded
		}
		return nil, ErrReaderBytesRead
	}

	ch := make(chan *overlay.Steak, 1)
	_, err = s.engine.Submit(ctx, overlay.TaggedBEEF{Beef: readBytes, Topics: topics}, engine.SubmitModeCurrent, func(steak *overlay.Steak) { ch <- steak })
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

// NewSubmitTransactionService creates a new SubmitTransactionService instance using the given provider.
// Panics if the provider is nil.
func NewSubmitTransactionService(provider SubmitTransactionProvider) *SubmitTransactionService {
	if provider == nil {
		panic("submit transaction service provider is nil")
	}

	return &SubmitTransactionService{
		engine:            provider,
		submitCallTimeout: 5 * time.Second,
	}
}

var (
	// ErrSubmitTransactionProvider is returned when the SubmitTransactionProvider fails to handle the transaction submission request.
	ErrSubmitTransactionProvider = errors.New("failed to submit transaction using provider")

	// ErrSubmitTransactionProviderTimeout is returned when the transaction submission request times out.
	ErrSubmitTransactionProviderTimeout = errors.New("submit transaction timeout occurred")

	// ErrMissingTopics is returned when an empty topics slice is provided as an argument.
	ErrMissingTopics = errors.New("provided topics cannot be an empty slice")

	// ErrMissingTransactionBytes is returned when an empty tx bytes slice is provided as an argument.
	ErrMissingTransactionBytes = errors.New("provided tx bytes data cannot be an empty slice")

	// ErrInvalidTopicFormat is returned when the topic header has an invalid format.
	ErrInvalidTopicFormat = errors.New("invalid topic header format")
)

// Topics represents a list of required topics for submitting a transaction.
type Topics []string

// Verify validates the Topics slice.
// It returns ErrMissingTopics if the slice is empty,
// or ErrInvalidTopicFormat if any topic in the slice is an empty string.
func (tt Topics) Verify() error {
	if len(tt) == 0 {
		return ErrMissingTopics
	}

	for _, t := range tt {
		if len(t) == 0 {
			return ErrInvalidTopicFormat
		}
	}
	return nil
}
