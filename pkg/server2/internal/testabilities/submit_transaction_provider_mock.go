package testabilities

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

// SubmitTransactionProviderMock is a mock implementation of the SubmitTransactionProvider interface
// for unit testing purposes. It allows you to control the behavior of Submit method calls and verify
// if the method was called with the expected parameters.
type SubmitTransactionProviderMock struct {
	t               *testing.T
	triggerCallback bool
	callSubmit      bool

	// expected state:
	expectedSteak *overlay.Steak
	expectedError error

	// actual state:
	called            bool
	callbackTriggered bool
	calledTaggedBEEF  overlay.TaggedBEEF
	calledSubmitMode  engine.SumbitMode
}

// Submit simulates the submission of a transaction. It will trigger the callback if the triggerCallback flag is set.
// It returns an error if the expectedError is set, otherwise returns nil.
func (s *SubmitTransactionProviderMock) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error) {
	s.t.Helper()

	if s.triggerCallback {
		s.callbackTriggered = true
		onSteakReady(s.expectedSteak)
	}

	if s.callSubmit {
		s.called = true
	}

	s.calledTaggedBEEF = taggedBEEF
	s.calledSubmitMode = mode
	return nil, s.expectedError
}

// AssertCalled is used to verify that Submit was called and that the callback was triggered correctly.
func (s *SubmitTransactionProviderMock) AssertCalled() {
	s.t.Helper()

	require.Equal(s.t, s.triggerCallback, s.callbackTriggered, "Discrepancy between expected and actual callback triggering")
	require.Equal(s.t, s.callSubmit, s.called, "Discrepancy between expected and actual Submit call")
}

// SubmitTransactionProviderMockOption is a functional option type for configuring a SubmitTransactionProviderMock.
type SubmitTransactionProviderMockOption func(*SubmitTransactionProviderMock)

// SubmitTransactionProviderMockWithSTEAK allows setting a custom steak for the mock.
func SubmitTransactionProviderMockWithSTEAK(steak *overlay.Steak) SubmitTransactionProviderMockOption {
	return func(mock *SubmitTransactionProviderMock) {
		mock.expectedSteak = steak
		mock.callSubmit = true
	}
}

// SubmitTransactionProviderMockWithError allows setting a custom error for the mock to return when Submit is called.
func SubmitTransactionProviderMockWithError(err error) SubmitTransactionProviderMockOption {
	return func(mock *SubmitTransactionProviderMock) {
		mock.expectedError = err
	}
}

// SubmitTransactionProviderMockWithTriggeredCallback configures the mock to trigger the callback when Submit is called.
func SubmitTransactionProviderMockWithTriggeredCallback() SubmitTransactionProviderMockOption {
	return func(mock *SubmitTransactionProviderMock) {
		mock.triggerCallback = true
		mock.callSubmit = true
	}
}

// NewSubmitTransactionProviderMock creates a new instance of SubmitTransactionProviderMock with the provided options.
func NewSubmitTransactionProviderMock(t *testing.T, opts ...SubmitTransactionProviderMockOption) *SubmitTransactionProviderMock {
	mock := SubmitTransactionProviderMock{t: t}
	for _, opt := range opts {
		opt(&mock)
	}
	return &mock
}
