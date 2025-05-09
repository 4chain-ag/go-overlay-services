package testabilities

import (
	"context"
	"testing"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

// SubmitTransactionProviderMock is a mock implementation of the SubmitTransactionProvider interface,
// used for unit testing purposes. It simulates the behavior of the Submit method, allowing you to
// configure and verify expected behavior such as whether Submit was called, whether a callback was triggered,
// and what error (if any) should be returned.
//
// Default Behavior:
//   - expectedSubmitCall: true         // The mock expects Submit to be called.
//   - expectedCallbackTriggering: true // The mock expects the callback to be invoked.
//   - expectedError: nil               // No error will be returned by default.
//   - expectedSteak: empty Steak       // A zero-value Steak will be passed to the callback.
//   - callbackSleep: 1µs               // The callback will be triggered after 1 microsecond.
//
// These defaults can be overridden using functional options when creating the mock.
type SubmitTransactionProviderMock struct {
	t *testing.T

	// expected behavior state:
	expectedSteak              *overlay.Steak
	expectedError              error
	expectedSubmitCall         bool
	expectedCallbackTriggering bool

	// actual state:
	called           bool
	callbackInvoked  bool
	callbackSleep    time.Duration
	calledTaggedBEEF overlay.TaggedBEEF
	calledSubmitMode engine.SumbitMode
}

// Submit simulates the submission of a transaction.
// It records the call parameters and—if expectedCallbackTriggering is true—invokes the callback
// after a delay defined by callbackSleep. Returns expectedError if one is configured,
// otherwise returns a zero-value overlay.Steak.
func (s *SubmitTransactionProviderMock) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, callback engine.OnSteakReady) (overlay.Steak, error) {
	s.t.Helper()

	s.called = true
	s.calledTaggedBEEF = taggedBEEF
	s.calledSubmitMode = mode

	if s.expectedError != nil {
		return nil, s.expectedError
	}

	if s.expectedCallbackTriggering {
		time.AfterFunc(s.callbackSleep, func() {
			callback(s.expectedSteak)
			s.callbackInvoked = true
		})
	}

	return overlay.Steak{}, nil
}

// AssertCalled verifies that the Submit method was called as expected and that the callback
// was triggered if configured. It fails the test if the actual behavior deviates from expectations.
func (s *SubmitTransactionProviderMock) AssertCalled() {
	s.t.Helper()

	require.Equal(s.t, s.expectedCallbackTriggering, s.callbackInvoked, "Discrepancy between expected and actual callback triggering")
	require.Equal(s.t, s.expectedSubmitCall, s.called, "Discrepancy between expected and actual Submit call")
}

// SubmitTransactionProviderMockOption defines a functional option for configuring a SubmitTransactionProviderMock.
type SubmitTransactionProviderMockOption func(*SubmitTransactionProviderMock)

// SubmitTransactionProviderMockWithSTEAK configures the mock to invoke the callback
// with the specified Steak after the given timeout duration.
func SubmitTransactionProviderMockWithSTEAK(steak *overlay.Steak, timeout time.Duration) SubmitTransactionProviderMockOption {
	return func(mock *SubmitTransactionProviderMock) {
		mock.expectedSteak = steak
		mock.callbackSleep = timeout
	}
}

// SubmitTransactionProviderMockNotCalled configures the mock to expect that Submit
// should not be called and the callback should not be triggered.
func SubmitTransactionProviderMockNotCalled() SubmitTransactionProviderMockOption {
	return func(mock *SubmitTransactionProviderMock) {
		mock.expectedSubmitCall = false
		mock.expectedCallbackTriggering = false
		mock.expectedError = nil
	}
}

// SubmitTransactionProviderMockWithError configures the mock to return the specified error
// when Submit is called. It also disables callback triggering by default.
func SubmitTransactionProviderMockWithError(err error) SubmitTransactionProviderMockOption {
	return func(mock *SubmitTransactionProviderMock) {
		mock.expectedError = err
		mock.expectedCallbackTriggering = false
	}
}

// NewSubmitTransactionProviderMock creates a new instance of SubmitTransactionProviderMock
// with the provided testing object and functional options to override default behavior.
func NewSubmitTransactionProviderMock(t *testing.T, opts ...SubmitTransactionProviderMockOption) *SubmitTransactionProviderMock {
	mock := &SubmitTransactionProviderMock{
		t:                          t,
		callbackSleep:              time.Microsecond,
		expectedSubmitCall:         true,
		expectedCallbackTriggering: true,
		expectedError:              nil,
		expectedSteak:              &overlay.Steak{},
	}

	for _, opt := range opts {
		opt(mock)
	}
	return mock
}
