package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// DefaultSubmitTransactionTimeout defines the default maximum duration allowed
// for submitting a transaction before timing out.
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
func (s *SubmitTransactionService) SubmitTransaction(ctx context.Context, topics TransactionTopics, txBytes ...byte) (*overlay.Steak, error) {
	err := topics.Verify()
	if err != nil {
		return nil, err
	}

	ch := make(chan *overlay.Steak, 1)
	_, err = s.provider.Submit(ctx, overlay.TaggedBEEF{Beef: txBytes, Topics: topics}, engine.SubmitModeCurrent, func(steak *overlay.Steak) {
		ch <- steak
	})
	if err != nil {
		return nil, NewSubmitTransactionProviderError(err)
	}

	select {
	case steak := <-ch:
		return steak, nil
	case <-time.After(s.submitCallTimeout):
		return nil, NewSubmitTransactionServiceTimeoutError(s.submitCallTimeout)
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

// TransactionTopics represents a list of topics that must be provided when submitting a transaction.
type TransactionTopics []string

// Verify ensures the topic list is non-empty and that each topic is non-blank.
// Returns ErrMissingTransactionTopics or ErrInvalidTransactionTopicFormat on failure.
func (tt TransactionTopics) Verify() error {
	if len(tt) == 0 {
		return NewEmptyTransactionTopicsError()
	}

	for i, t := range tt {
		t = strings.TrimSpace(t)
		if len(t) == 0 || len(t) == 1 { // TODO: Add more robust topic format check.
			return NewErrInvalidTopicFormatError(i)
		}
	}

	return nil
}

// NewEmptyTransactionTopicsError returns an Error indicating that the topics slice is empty,
// which is invalid input when submitting a transaction.
func NewEmptyTransactionTopicsError() Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       "provided topics cannot be an empty slice",
		slug:      "At least one topic must be provided in the correct string format. Empty topic values are not allowed.",
	}
}

// NewErrInvalidTopicFormatError returns an Error indicating that a specific topic,
// identified by its index, is in an invalid format.
func NewErrInvalidTopicFormatError(i int) Error {
	return Error{
		errorType: ErrorTypeIncorrectInput,
		err:       fmt.Sprintf("invalid topic header format for topic no. %d", i+1),
		slug:      "One or more topics are in an invalid format. Empty string values are not allowed.",
	}
}

// NewSubmitTransactionProviderError returns an Error indicating that the overlay engine
// failed to process a submitted transaction octet-stream.
func NewSubmitTransactionProviderError(err error) Error {
	return Error{
		errorType: ErrorTypeProviderFailure,
		err:       err.Error(),
		slug:      "Unable to process submitted transaction octet-stream due to an error in the overlay engine.",
	}
}

// NewSubmitTransactionServiceTimeoutError returns an Error indicating that the submit transaction
// request exceeded the configured timeout limit.
func NewSubmitTransactionServiceTimeoutError(timeout time.Duration) Error {
	return Error{
		errorType: ErrorTypeOperationTimeout,
		err:       fmt.Sprintf("submit transaction timeout occurred - limit set to %f seconds", timeout.Seconds()),
		slug:      "The submitted request exceeded the timeout limit",
	}
}
