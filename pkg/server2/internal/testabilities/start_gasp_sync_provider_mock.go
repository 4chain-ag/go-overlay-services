package testabilities

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ErrMockStartGASPSync is the error returned by StartGASPSyncProviderMock when shouldFail is true.
var ErrMockStartGASPSync = assert.AnError

// StartGASPSyncProviderMock is a mock implementation of the StartGASPSyncProvider interface.
// It's used for testing the StartGASPSyncHandler.
type StartGASPSyncProviderMock struct {
	t          *testing.T
	shouldFail bool
	called     bool
}

// StartGASPSync implements the StartGASPSyncProvider interface.
// It records that the method was called and returns an error if shouldFail is true.
func (m *StartGASPSyncProviderMock) StartGASPSync(ctx context.Context) error {
	m.called = true
	if m.shouldFail {
		return ErrMockStartGASPSync
	}
	return nil
}

// AssertCalled verifies that StartGASPSync was called.
func (m *StartGASPSyncProviderMock) AssertCalled(t *testing.T) {
	assert.True(t, m.called, "Expected StartGASPSync to be called")
}

// NewStartGASPSyncProviderMock creates a new StartGASPSyncProviderMock with the specified behavior.
func NewStartGASPSyncProviderMock(t *testing.T, shouldFail bool) *StartGASPSyncProviderMock {
	return &StartGASPSyncProviderMock{
		t:          t,
		shouldFail: shouldFail,
		called:     false,
	}
}
