package testabilities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type LookupServiceDocumentationProviderMockExpectations struct {
	DocumentationCall bool

	Error error

	Documentation string
}

var DefaultLookupServiceDocumentationProviderMockExpectations = LookupServiceDocumentationProviderMockExpectations{

	DocumentationCall: true,

	Error: nil,

	Documentation: "# Test Documentation\nThis is a test markdown document.",
}

// LookupServiceDocumentationProviderMock is a simple mock implementation for testing

type LookupServiceDocumentationProviderMock struct {
	t *testing.T

	expectations LookupServiceDocumentationProviderMockExpectations

	called bool
}

// GetDocumentationForLookupServiceProvider simulates a documentation retrieval operation

func (m *LookupServiceDocumentationProviderMock) GetDocumentationForLookupServiceProvider(lookupServiceName string) (string, error) {

	m.t.Helper()

	m.called = true

	if m.expectations.Error != nil {

		return "", m.expectations.Error

	}

	return m.expectations.Documentation, nil

}

func (m *LookupServiceDocumentationProviderMock) AssertCalled() {

	m.t.Helper()

	require.Equal(m.t, m.expectations.DocumentationCall, m.called, "Discrepancy between expected and actual DocumentationCall")

}

func NewLookupServiceDocumentationProviderMock(t *testing.T, expectations LookupServiceDocumentationProviderMockExpectations) *LookupServiceDocumentationProviderMock {

	return &LookupServiceDocumentationProviderMock{

		t: t,

		expectations: expectations,
	}

}
