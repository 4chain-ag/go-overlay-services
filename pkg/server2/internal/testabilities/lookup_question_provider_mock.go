package testabilities

import (
	"context"
	"errors"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/stretchr/testify/require"
)

// Helper functions to create standardized error responses for testing
// These match the error messages from app.Error types

// NewMissingServiceFieldResponse creates a standard error response for missing service field.
func NewMissingServiceFieldResponse() openapi.Error {
	return openapi.Error{
		Message: "The service field is required in the lookup question request.",
	}
}

// NewLookupQuestionProviderErrorResponse creates a standard error response for lookup provider errors.
func NewLookupQuestionProviderErrorResponse() openapi.Error {
	return openapi.Error{
		Message: "Unable to process lookup question due to an error in the overlay engine.",
	}
}

// NewInvalidRequestBodyResponse creates a standard error response for invalid request body.
func NewInvalidRequestBodyResponse() openapi.Error {
	return openapi.Error{
		Message: "Invalid request body format or structure. Please check the API documentation for the correct format.",
	}
}

// SimpleLookupQuestionProvider is a minimal implementation of the LookupQuestionProvider interface
// for simple test cases.
type SimpleLookupQuestionProvider struct {
	Answer *lookup.LookupAnswer
	Err    error
}

// Lookup implements the LookupQuestionProvider interface by simply returning the pre-configured
// answer and error.
func (m *SimpleLookupQuestionProvider) Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
	return m.Answer, m.Err
}

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

// LookupQuestionProviderMockWithAppError configures the mock to return an app.Error.
func LookupQuestionProviderMockWithAppError(appErr app.Error) LookupQuestionProviderMockOption {
	return func(m *LookupQuestionProviderMock) {
		m.expectedError = appErr
	}
}

// LookupQuestionProviderMockWithMissingServiceField configures the mock to return a missing service field error.
func LookupQuestionProviderMockWithMissingServiceField() LookupQuestionProviderMockOption {
	return func(m *LookupQuestionProviderMock) {
		m.expectedError = app.NewMissingServiceFieldError()
	}
}

// LookupQuestionProviderMockWithInvalidLookupQuestion configures the mock to return an invalid lookup question error.
func LookupQuestionProviderMockWithInvalidLookupQuestion() LookupQuestionProviderMockOption {
	return func(m *LookupQuestionProviderMock) {
		m.expectedError = app.NewInvalidLookupQuestionError()
	}
}

// LookupQuestionProviderMockWithProviderError configures the mock to return a provider error.
func LookupQuestionProviderMockWithProviderError(message string) LookupQuestionProviderMockOption {
	return func(m *LookupQuestionProviderMock) {
		m.expectedError = app.NewLookupQuestionProviderError(errors.New(message))
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
