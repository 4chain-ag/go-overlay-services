package app

import (
	"bytes"
	"context"
	"errors"
	"io"
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
	provider                         SubmitTransactionProvider
	submitTransactionProviderTimeout time.Duration
	readerLimit                      int64
}

// SubmitTransaction validates the topics, creates a TaggedBEEF from the request body,
// and submits it using the configured provider.
// Returns the resulting Steak on success.
// Possible errors include:
//   - ErrMissingTopics or ErrInvalidTopicFormat (invalid topics)
//   - ErrReaderLimitExceeded (body exceeds size limit)
//   - ErrReaderBytesRead (read failure)
//   - ErrSubmitTransactionProvider (submission failure)
//   - ErrSubmitTransactionProviderTimeout (timeout waiting for provider callback)
func (s *SubmitTransactionService) SubmitTransaction(ctx context.Context, topics Topics, body ...byte) (*overlay.Steak, error) {
	if len(body) == 0 {
		return nil, ErrMissingTransactionBytes
	}

	err := topics.Verify()
	if err != nil {
		return nil, err
	}

	taggedBEEF, err := s.CreateTaggedBEEF(topics, body...)
	if errors.Is(err, ErrReaderLimitExceeded) {
		return nil, ErrReaderLimitExceeded
	}
	if err != nil {
		return nil, ErrReaderBytesRead
	}

	ch := make(chan *overlay.Steak, 1)
	_, err = s.provider.Submit(ctx, *taggedBEEF, engine.SubmitModeCurrent, func(steak *overlay.Steak) { ch <- steak })
	if err != nil {
		return nil, errors.Join(err, ErrSubmitTransactionProvider)
	}

	select {
	case steak := <-ch:
		return steak, nil
	case <-time.After(s.submitTransactionProviderTimeout):
		return nil, ErrSubmitTransactionProviderTimeout
	}
}

// CreateTaggedBEEF reads the given body with a size limit and wraps it with the provided topics
// into a TaggedBEEF. Returns ErrReaderLimitExceeded if the body exceeds the limit,
// or ErrReaderBytesRead on read/write failures.
func (s *SubmitTransactionService) CreateTaggedBEEF(topics []string, body ...byte) (*overlay.TaggedBEEF, error) {
	reader := io.LimitReader(bytes.NewBuffer(body), s.readerLimit+1)
	buff := make([]byte, 64*1024)
	var dst bytes.Buffer
	var bytesRead int64

	for {
		n, err := reader.Read(buff)
		if n > 0 {
			bytesRead += int64(n)
			if bytesRead > s.readerLimit {
				return nil, ErrReaderLimitExceeded
			}

			if _, inner := dst.Write(buff[:n]); inner != nil {
				return nil, errors.Join(inner, ErrReaderBytesRead)
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, errors.Join(err, ErrReaderBytesRead)
		}
	}

	return &overlay.TaggedBEEF{Beef: dst.Bytes(), Topics: topics}, nil
}

// NewSubmitTransactionService creates a new SubmitTransactionService instance using the given provider.
// Panics if the provider is nil.
func NewSubmitTransactionService(provider SubmitTransactionProvider) *SubmitTransactionService {
	if provider == nil {
		panic("submit transaction service provider is nil")
	}

	return &SubmitTransactionService{
		provider:                         provider,
		submitTransactionProviderTimeout: 5 * time.Second,
		readerLimit:                      1000 * 1024 * 1024, // 1GB
	}
}

var (
	// ErrReaderBytesRead indicates a failure while reading input data.
	ErrReaderBytesRead = errors.New("failed to read input data")

	// ErrReaderLimitExceeded is returned when the read exceeds the allowed byte limit.
	ErrReaderLimitExceeded = errors.New("input data too large")

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
