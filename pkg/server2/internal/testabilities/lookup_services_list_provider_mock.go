package testabilities

import (
	"github.com/bsv-blockchain/go-sdk/overlay"
)

// LookupMetadataMock represents mock metadata for a lookup service provider.
type LookupMetadataMock struct {
	Description string
	Icon        string
	Version     string
	InfoUrl     string
}

// WithLookupServicesList returns a TestOverlayEngineStubOption that configures the engine
// to return a specific list of lookup service providers.
func WithLookupServicesList(metadata map[string]*LookupMetadataMock) TestOverlayEngineStubOption {
	return func(stub *TestOverlayEngineStub) {
		stub.listLookupServicesProvider = &mockLookupServicesListProvider{
			metadata: metadata,
		}
	}
}

// WithEmptyLookupServicesList returns a TestOverlayEngineStubOption that configures the engine
// to return an empty list of lookup service providers.
func WithEmptyLookupServicesList() TestOverlayEngineStubOption {
	return WithLookupServicesList(map[string]*LookupMetadataMock{})
}

// mockLookupServicesListProvider is a mock implementation of the lookup services list provider.
type mockLookupServicesListProvider struct {
	metadata map[string]*LookupMetadataMock
}

// ListLookupServiceProviders returns a mock list of lookup service providers.
func (m *mockLookupServicesListProvider) ListLookupServiceProviders() map[string]*overlay.MetaData {
	result := make(map[string]*overlay.MetaData, len(m.metadata))

	for name, mock := range m.metadata {
		if mock == nil {
			result[name] = nil
			continue
		}

		result[name] = &overlay.MetaData{
			Description: mock.Description,
			Icon:        mock.Icon,
			Version:     mock.Version,
			InfoUrl:     mock.InfoUrl,
		}
	}

	return result
}
