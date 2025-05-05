package testabilities

import (
	"context"
	"testing"

	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/stretchr/testify/require"
)

// LookupQuestionProviderMockOption defines a functional option for configuring a LookupQuestionProviderMock.
type LookupQuestionProviderMockOption func(*LookupQuestionProviderMock)

// LookupQuestionProviderMock is a mock implementation of the LookupQuestionProvider interface.
type LookupQuestionProviderMock struct {
	t *testing.T

	// expected behavior state:
	expectedAnswer *lookup.LookupAnswer
	expectedError  error
	expectedCall   bool

	// actual state:
	called         bool
	calledQuestion *lookup.LookupQuestion
}

// Lookup simulates the lookup process and returns the expected answer or error.
func (m *LookupQuestionProviderMock) Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
	m.called = true
	m.calledQuestion = question

	return m.expectedAnswer, m.expectedError
}

// AssertCalled verifies that the Lookup method was called as expected.
func (m *LookupQuestionProviderMock) AssertCalled() {
	m.t.Helper()
	require.Equal(m.t, m.expectedCall, m.called, "Discrepancy between expected and actual Lookup call")
}

// LookupQuestionProviderMockWithAnswer configures the mock to return the provided answer.
func LookupQuestionProviderMockWithAnswer(answer *lookup.LookupAnswer) LookupQuestionProviderMockOption {
	return func(m *LookupQuestionProviderMock) {
		m.expectedAnswer = answer
	}
}

// LookupQuestionProviderMockWithError configures the mock to return the provided error.
func LookupQuestionProviderMockWithError(err error) LookupQuestionProviderMockOption {
	return func(m *LookupQuestionProviderMock) {
		m.expectedError = err
	}
}

// LookupQuestionProviderMockNotCalled configures the mock to expect that Lookup should not be called.
func LookupQuestionProviderMockNotCalled() LookupQuestionProviderMockOption {
	return func(m *LookupQuestionProviderMock) {
		m.expectedCall = false
	}
}

// NewLookupQuestionProviderMock creates a new LookupQuestionProviderMock with the provided options.
func NewLookupQuestionProviderMock(t *testing.T, opts ...LookupQuestionProviderMockOption) *LookupQuestionProviderMock {
	mock := &LookupQuestionProviderMock{
		t:            t,
		expectedCall: true,
	}

	for _, opt := range opts {
		opt(mock)
	}

	return mock
}
