package testabilities

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
)

// Standard metadata maps that can be used for testing
var (
	// LookupListEmptyMetadata is an empty metadata map
	LookupListEmptyMetadata = map[string]*overlay.MetaData{}

	// LookupListDefaultMetadata contains standard metadata for testing
	LookupListDefaultMetadata = map[string]*overlay.MetaData{
		"lookup_service1": {
			Description: "Description 1",
			Icon:        "https://example.com/icon.png",
			Version:     "1.0.0",
			InfoUrl:     "https://example.com/info",
		},
		"lookup_service2": {
			Description: "Description 2",
			Icon:        "https://example.com/icon2.png",
			Version:     "2.0.0",
			InfoUrl:     "https://example.com/info2",
		},
	}
)

// Standard expected responses that can be used for testing
var (
	// LookupListEmptyExpectedResponse is an empty response
	LookupListEmptyExpectedResponse = app.LookupServiceProviders{}

	// LookupListDefaultExpectedResponse contains the standard expected response matching LookupListDefaultMetadata
	LookupListDefaultExpectedResponse = app.LookupServiceProviders{
		"lookup_service1": app.LookupMetadata{
			Name:             "lookup_service1",
			ShortDescription: "Description 1",
			IconURL:          "https://example.com/icon.png",
			Version:          "1.0.0",
			InformationURL:   "https://example.com/info",
		},
		"lookup_service2": app.LookupMetadata{
			Name:             "lookup_service2",
			ShortDescription: "Description 2",
			IconURL:          "https://example.com/icon2.png",
			Version:          "2.0.0",
			InformationURL:   "https://example.com/info2",
		},
	}
)

// LookupListProviderMockExpectations defines the expected behavior of the LookupListProviderMock during a test.
type LookupListProviderMockExpectations struct {
	MetadataList                   map[string]*overlay.MetaData
	ListLookupServiceProvidersCall bool
}

// LookupListProviderMock is a mock implementation of a lookup service list provider,
// used for testing the behavior of components that depend on lookup service listing.
type LookupListProviderMock struct {
	t            *testing.T
	expectations LookupListProviderMockExpectations
	called       bool
}

// ListLookupServiceProviders returns the predefined list of lookup services.
func (m *LookupListProviderMock) ListLookupServiceProviders() map[string]*overlay.MetaData {
	m.t.Helper()
	m.called = true
	return m.expectations.MetadataList
}

// AssertCalled verifies that the ListLookupServiceProviders method was called if it was expected to be.
func (m *LookupListProviderMock) AssertCalled() {
	m.t.Helper()
	require.Equal(m.t, m.expectations.ListLookupServiceProvidersCall, m.called, "Discrepancy between expected and actual ListLookupServiceProviders call")
}

// NewLookupListProviderMock creates a new instance of LookupListProviderMock with the given expectations.
func NewLookupListProviderMock(t *testing.T, expectations LookupListProviderMockExpectations) *LookupListProviderMock {
	return &LookupListProviderMock{
		t:            t,
		expectations: expectations,
	}
}
