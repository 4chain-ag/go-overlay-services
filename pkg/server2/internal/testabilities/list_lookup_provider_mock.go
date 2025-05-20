package testabilities

import (
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

// Standard metadata maps for lookup service provider tests
var (
	// EmptyMetadata is an empty map of metadata
	LookupListEmptyMetadata = map[string]*overlay.MetaData{}

	// LookupDefaultMetadata contains standard metadata for testing lookup providers
	LookupDefaultMetadata = map[string]*overlay.MetaData{
		"lookup_provider1": {
			Description: "Description 1",
			Icon:        "https://example.com/icon.png",
			Version:     "1.0.0",
			InfoUrl:     "https://example.com/info",
		},
		"lookup_provider2": {
			Description: "Description 2",
			Icon:        "https://example.com/icon2.png",
			Version:     "2.0.0",
			InfoUrl:     "https://example.com/info2",
		},
	}
)

// Standard expected responses for lookup service providers that can be used for testing
var (
	// EmptyLookupListExpectedResponse is an empty response
	EmptyLookupListExpectedResponse = app.LookupListResponse{}

	// DefaultLookupListExpectedResponse contains the standard expected response matching DefaultMetadata
	DefaultLookupListExpectedResponse = app.LookupListResponse{
		"lookup_provider1": app.LookupServiceProviderMetadata{
			Name:             "lookup_provider1",
			ShortDescription: "Description 1",
			IconURL:          ptr.To("https://example.com/icon.png"),
			Version:          ptr.To("1.0.0"),
			InformationURL:   ptr.To("https://example.com/info"),
		},
		"lookup_provider2": app.LookupServiceProviderMetadata{
			Name:             "lookup_provider2",
			ShortDescription: "Description 2",
			IconURL:          ptr.To("https://example.com/icon2.png"),
			Version:          ptr.To("2.0.0"),
			InformationURL:   ptr.To("https://example.com/info2"),
		},
	}
)

// LookupListProviderMockExpectations defines the expected behavior of the LookupListProviderMock during a test.
type LookupListProviderMockExpectations struct {
	// MetadataList is the mock lookup service providers that will be returned.
	MetadataList map[string]*overlay.MetaData

	// ListLookupServiceProvidersCall indicates whether the ListLookupServiceProviders method is expected to be called during the test.
	ListLookupServiceProvidersCall bool
}

// LookupListProviderMock is a mock implementation of a lookup service provider list provider,
// used for testing the behavior of components that depend on lookup service provider listing.
type LookupListProviderMock struct {
	t *testing.T

	// expectations defines the expected behavior and outcomes for this mock.
	expectations LookupListProviderMockExpectations

	// called is true if the ListLookupServiceProviders method was called.
	called bool
}

// ListLookupServiceProviders returns the predefined list of lookup service providers.
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

// LookupListProviderAlwaysEmpty is a mock that always returns an empty lookup service provider list.
type LookupListProviderAlwaysEmpty struct{}

// ListLookupServiceProviders returns an empty map of lookup service providers.
func (*LookupListProviderAlwaysEmpty) ListLookupServiceProviders() map[string]*overlay.MetaData {
	return LookupListEmptyMetadata
}

// LookupListProviderAlwaysSuccess is a mock that always returns a predefined list of lookup service providers.
type LookupListProviderAlwaysSuccess struct{}

// ListLookupServiceProviders returns a predefined map of lookup service providers.
func (*LookupListProviderAlwaysSuccess) ListLookupServiceProviders() map[string]*overlay.MetaData {
	return LookupDefaultMetadata
}
