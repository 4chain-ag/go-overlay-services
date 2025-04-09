package server

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
)

// NoopEngineProvider is a custom test overlay engine implementation. This is only a temporary solution and will be removed
// after migrating the engine code. Currently, it functions as mock for the overlay HTTP server.
type NoopEngineProvider struct{}

// Submit is a no-op call that always returns an empty STEAK with nil error.
func (*NoopEngineProvider) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error) {
	onSteakReady(&overlay.Steak{
		"test_key": &overlay.AdmittanceInstructions{
			OutputsToAdmit: []uint32{1},
		}})
	return overlay.Steak{}, nil
}

// SyncAdvertisements is a no-op call that always returns a nil error.
func (*NoopEngineProvider) SyncAdvertisements(ctx context.Context) error { return nil }

// GetTopicManagerDocumentation is a no-op call that always returns a nil error.
func (*NoopEngineProvider) GetTopicManagerDocumentation(ctx context.Context) error { return nil }

// Lookup is a no-op call that always returns an empty lookup answer with nil error.
func (*NoopEngineProvider) Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
	return &lookup.LookupAnswer{}, nil
}

// GetUTXOHistory is a no-op call that always returns an empty engine output with nil error.
func (*NoopEngineProvider) GetUTXOHistory(ctx context.Context, output *engine.Output, historySelector func(beef []byte, outputIndex uint32, currentDepth uint32) bool, currentDepth uint32) (*engine.Output, error) {
	return &engine.Output{}, nil
}

// StartGASPSync is a no-op call that always returns a nil error.
func (*NoopEngineProvider) StartGASPSync(ctx context.Context) error { return nil }

// ProvideForeignSyncResponse is a no-op call that always returns an empty initial GASP response with nil error.
func (*NoopEngineProvider) ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {
	return &core.GASPInitialResponse{}, nil
}

// ProvideForeignGASPNode is a no-op call that always returns an empty GASP node with nil error.
func (*NoopEngineProvider) ProvideForeignGASPNode(ctx context.Context, graphId, outpoint *overlay.Outpoint) (*core.GASPNode, error) {
	return &core.GASPNode{}, nil
}

// ListTopicManagers is a no-op call that always returns an empty topic managers map with nil error.
func (*NoopEngineProvider) ListTopicManagers() map[string]*overlay.MetaData {
	return map[string]*overlay.MetaData{}
}

// ListLookupServiceProviders is a no-op call that always returns an empty lookup service providers map with nil error.
func (*NoopEngineProvider) ListLookupServiceProviders() map[string]*overlay.MetaData {
	return map[string]*overlay.MetaData{}
}

// GetDocumentationForLookupServiceProvider is a no-op call that always returns an empty string with nil error.
func (*NoopEngineProvider) GetDocumentationForLookupServiceProvider(provider string) (string, error) {
	return "", nil
}

// GetDocumentationForTopicManager is a no-op call that always returns an empty string with nil error.
func (*NoopEngineProvider) GetDocumentationForTopicManager(provider string) (string, error) {
	return "", nil
}

// NewNoopEngineProvider returns an OverlayEngineProvider implementation
// and checks whether the engine contract matches the implemented method set.
func NewNoopEngineProvider() engine.OverlayEngineProvider {
	return &NoopEngineProvider{}
}
