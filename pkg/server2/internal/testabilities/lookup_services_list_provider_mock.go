package testabilities

import (
	"testing"

	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

// LookupServicesListProviderMockExpectations defines the expected behavior of the LookupServicesListProviderMock during a test.
type LookupServicesListProviderMockExpectations struct {
	// MetadataList is the mock lookup service providers that will be returned.
	MetadataList map[string]*overlay.MetaData

	// ListLookupServiceProvidersCall indicates whether the ListLookupServiceProviders method is expected to be called during the test.
	ListLookupServiceProvidersCall bool
}

// DefaultLookupServicesListProviderMockExpectations provides default expectations for LookupServicesListProviderMock.
var DefaultLookupServicesListProviderMockExpectations = LookupServicesListProviderMockExpectations{
	MetadataList:                   map[string]*overlay.MetaData{},
	ListLookupServiceProvidersCall: true,
}

// LookupServicesListProviderMock is a mock implementation of a lookup services list provider,
// used for testing the behavior of components that depend on lookup service provider listing.
type LookupServicesListProviderMock struct {
	t *testing.T

	// expectations defines the expected behavior and outcomes for this mock.
	expectations LookupServicesListProviderMockExpectations

	// called is true if the ListLookupServiceProviders method was called.
	called bool
}

// ListLookupServiceProviders returns the predefined list of lookup service providers.
func (m *LookupServicesListProviderMock) ListLookupServiceProviders() map[string]*overlay.MetaData {
	m.t.Helper()
	m.called = true
	return m.expectations.MetadataList
}

// AssertCalled verifies that the ListLookupServiceProviders method was called if it was expected to be.
func (m *LookupServicesListProviderMock) AssertCalled() {
	m.t.Helper()
	require.Equal(m.t, m.expectations.ListLookupServiceProvidersCall, m.called, "Discrepancy between expected and actual ListLookupServiceProviders call")
}

// NewLookupServicesListProviderMock creates a new instance of LookupServicesListProviderMock with the given expectations.
func NewLookupServicesListProviderMock(t *testing.T, expectations LookupServicesListProviderMockExpectations) *LookupServicesListProviderMock {
	mock := &LookupServicesListProviderMock{
		t:            t,
		expectations: expectations,
	}
	return mock
}

// LookupListProviderAlwaysEmpty is a mock that always returns an empty lookup services list.
type LookupListProviderAlwaysEmpty struct{}

// ListLookupServiceProviders returns an empty map of lookup service providers.
func (*LookupListProviderAlwaysEmpty) ListLookupServiceProviders() map[string]*overlay.MetaData {
	return map[string]*overlay.MetaData{}
}

// LookupListProviderAlwaysSuccess is a mock that always returns a predefined list of lookup service providers.
type LookupListProviderAlwaysSuccess struct{}

// ListLookupServiceProviders returns a predefined map of lookup service providers.
func (*LookupListProviderAlwaysSuccess) ListLookupServiceProviders() map[string]*overlay.MetaData {
	return map[string]*overlay.MetaData{
		"provider1": {
			Description: "Description 1",
			Icon:        "https://example.com/icon.png",
			Version:     "1.0.0",
			InfoUrl:     "https://example.com/info",
		},
		"provider2": {
			Description: "Description 2",
			Icon:        "https://example.com/icon2.png",
			Version:     "2.0.0",
			InfoUrl:     "https://example.com/info2",
		},
	}
}
