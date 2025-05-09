package testabilities

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// SyncAdvertisementsProviderMockOption defines a functional option for configuring a SyncAdvertisementsProviderMock.
type SyncAdvertisementsProviderMockOption func(*SyncAdvertisementsProviderMock)

// SyncAdvertisementsProviderMockNotCalled configures the mock to expect that SyncAdvertisements should not be called.
func SyncAdvertisementsProviderMockNotCalled() SyncAdvertisementsProviderMockOption {
	return func(mock *SyncAdvertisementsProviderMock) {
		mock.expectedErr = nil
		mock.expectedSyncAdvertisementsCall = false
	}
}

// SyncAdvertisementsProviderMockWithError configures the mock to return the specified error
// when SyncAdvertisements is called.
func SyncAdvertisementsProviderMockWithError(err error) SyncAdvertisementsProviderMockOption {
	return func(mock *SyncAdvertisementsProviderMock) {
		mock.expectedErr = err
	}
}

// SyncAdvertisementsProviderMock is a mock implementation of the SyncAdvertisementsProvider interface,
// used for unit testing purposes. It simulates the behavior of the SyncAdvertisements method, allowing you to
// configure and verify expected behavior such as whether SyncAdvertisements was called, and what error (if any) should be returned.
//
// Default Behavior:
//   - expectedSyncAdvertisementsCall: true         // The mock expects expectedSyncAdvertisementsCall to be called.
//   - expectedError: nil                          // No error will be returned by default.
//
// These defaults can be overridden using functional options when creating the mock.
type SyncAdvertisementsProviderMock struct {
	t *testing.T
	// expected behavior state:
	expectedErr                    error
	expectedSyncAdvertisementsCall bool

	// actual state:
	called bool
}

// SyncAdvertisements simulates the synchronize advertisements request.
// It records the call parameters. Returns expectedError if one is configured, otherwise returns a nil error.
func (s *SyncAdvertisementsProviderMock) SyncAdvertisements(ctx context.Context) error {
	s.called = true
	if s.expectedErr != nil {
		return s.expectedErr
	}
	return nil
}

// AssertCalled verifies that the Submit method was called as expected and that the callback
// was triggered if configured. It fails the test if the actual behavior deviates from expectations.
func (s *SyncAdvertisementsProviderMock) AssertCalled() {
	s.t.Helper()

	require.Equal(s.t, s.expectedSyncAdvertisementsCall, s.called, "Discrepancy between expected and actual SyncAdvertisements call")
}

// NewSubmitTransactionProviderMock creates a new instance of SubmitTransactionProviderMock
// with the provided testing object and functional options to override default behavior.
func NewSyncAdvertisementsProviderMock(t *testing.T, opts ...SyncAdvertisementsProviderMockOption) *SyncAdvertisementsProviderMock {
	mock := &SyncAdvertisementsProviderMock{
		t:                              t,
		expectedErr:                    nil,
		expectedSyncAdvertisementsCall: true,
	}

	for _, opt := range opts {
		opt(mock)
	}
	return mock
}
