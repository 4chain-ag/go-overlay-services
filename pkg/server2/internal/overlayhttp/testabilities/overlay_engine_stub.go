package testabilities

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/overlayhttp"
	"github.com/bsv-blockchain/go-sdk/chainhash"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
	"github.com/bsv-blockchain/go-sdk/transaction"
)

// TestOverlayEngineStubOption is a functional option type used to configure a TestOverlayEngineStub.
// It allows setting custom behaviors for different parts of the TestOverlayEngineStub.
type TestOverlayEngineStubOption func(*TestOverlayEngineStub)

// WithSyncAdvertisementsProvider allows setting a custom SyncAdvertisementsProvider in a TestOverlayEngineStub.
// This can be used to mock advertisement synchronization behavior during tests.
func WithSyncAdvertisementsProvider(provider overlayhttp.SyncAdvertisementsProvider) TestOverlayEngineStubOption {
	return func(engine *TestOverlayEngineStub) {
		engine.syncAdvertisementsProvider = provider
	}
}

// WithSubmitTransactionProvider allows setting a custom SubmitTransactionProvider in a TestOverlayEngineStub.
// This can be used to mock transaction submission behavior during tests.
func WithSubmitTransactionProvider(provider overlayhttp.SubmitTransactionProvider) TestOverlayEngineStubOption {
	return func(engine *TestOverlayEngineStub) {
		engine.submitTransactionProvider = provider
	}
}

// TestOverlayEngineStub is a test implementation of the engine.OverlayEngineProvider interface.
// It is used to mock engine behavior in unit tests, allowing the simulation of various engine actions
// like submitting transactions and synchronizing advertisements.
type TestOverlayEngineStub struct {
	t                          *testing.T
	syncAdvertisementsProvider overlayhttp.SyncAdvertisementsProvider
	submitTransactionProvider  overlayhttp.SubmitTransactionProvider
}

// GetDocumentationForLookupServiceProvider returns documentation for a lookup service provider (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (t TestOverlayEngineStub) GetDocumentationForLookupServiceProvider(provider string) (string, error) {
	panic("unimplemented")
}

// GetDocumentationForTopicManager returns documentation for a topic manager (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (t TestOverlayEngineStub) GetDocumentationForTopicManager(provider string) (string, error) {
	panic("unimplemented")
}

// GetUTXOHistory retrieves UTXO history for the given output (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (t TestOverlayEngineStub) GetUTXOHistory(ctx context.Context, output *engine.Output, historySelector func(beef []byte, outputIndex uint32, currentDepth uint32) bool, currentDepth uint32) (*engine.Output, error) {
	panic("unimplemented")
}

// HandleNewMerkleProof processes a new Merkle proof for a transaction (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (t TestOverlayEngineStub) HandleNewMerkleProof(ctx context.Context, txid *chainhash.Hash, proof *transaction.MerklePath) error {
	panic("unimplemented")
}

// ListLookupServiceProviders lists the available lookup service providers (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (t TestOverlayEngineStub) ListLookupServiceProviders() map[string]*overlay.MetaData {
	panic("unimplemented")
}

// ListTopicManagers lists the available topic managers (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (t TestOverlayEngineStub) ListTopicManagers() map[string]*overlay.MetaData {
	panic("unimplemented")
}

// Lookup performs a lookup query based on the provided LookupQuestion (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (t TestOverlayEngineStub) Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
	panic("unimplemented")
}

// ProvideForeignGASPNode returns a foreign GASP node (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (t TestOverlayEngineStub) ProvideForeignGASPNode(ctx context.Context, graphId *overlay.Outpoint, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error) {
	panic("unimplemented")
}

// ProvideForeignSyncResponse returns a foreign sync response (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (t TestOverlayEngineStub) ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {
	panic("unimplemented")
}

// StartGASPSync starts the GASP synchronization process (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (t TestOverlayEngineStub) StartGASPSync(ctx context.Context) error {
	panic("unimplemented")
}

// Submit processes a transaction submission and returns a steak or error based on the provided inputs.
// It calls the Submit method of the configured SubmitTransactionProvider and handles the steak callback.
func (t TestOverlayEngineStub) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error) {
	t.t.Helper()

	return t.submitTransactionProvider.Submit(ctx, taggedBEEF, mode, onSteakReady)
}

// SyncAdvertisements synchronizes advertisements using the configured SyncAdvertisementsProvider.
// It calls the SyncAdvertisements method of the provider and handles the result.
func (t TestOverlayEngineStub) SyncAdvertisements(ctx context.Context) error {
	t.t.Helper()

	return t.syncAdvertisementsProvider.SyncAdvertisements(ctx)
}

// NewTestOverlayEngineStub creates and returns a new instance of TestOverlayEngineStub with the provided options.
// The options allow for configuring custom providers for transaction submission and advertisement synchronization.
func NewTestOverlayEngineStub(t *testing.T, opts ...TestOverlayEngineStubOption) engine.OverlayEngineProvider {
	engine := TestOverlayEngineStub{
		t:                          t,
		submitTransactionProvider:  submitTransactionProviderAlwaysSuccessStub{},
		syncAdvertisementsProvider: syncAdvertisementsProviderAlwaysSuccessStub{},
	}

	for _, opt := range opts {
		opt(&engine)
	}
	return engine
}

// submitTransactionProviderAlwaysSuccessStub is a mock implementation of SubmitTransactionProvider that always succeeds.
// It is used as the default SubmitTransactionProvider in the TestOverlayEngineStub.
type submitTransactionProviderAlwaysSuccessStub struct{ ExpectedSteak overlay.Steak }

// Submit simulates a successful transaction submission, triggering the onSteakReady callback with the expected steak.
func (s submitTransactionProviderAlwaysSuccessStub) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error) {
	onSteakReady(&s.ExpectedSteak)
	return nil, nil
}

// syncAdvertisementsProviderAlwaysSuccessStub is a mock implementation of SyncAdvertisementsProvider that always succeeds.
// It is used as the default SyncAdvertisementsProvider in the TestOverlayEngineStub.
type syncAdvertisementsProviderAlwaysSuccessStub struct{}

// SyncAdvertisements simulates a successful advertisements synchronization request call.
// It always returns nil, indicating that the synchronization was successful.
func (syncAdvertisementsProviderAlwaysSuccessStub) SyncAdvertisements(ctx context.Context) error {
	return nil
}
