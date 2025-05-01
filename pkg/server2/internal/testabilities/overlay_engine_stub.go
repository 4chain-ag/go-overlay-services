package testabilities

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp/core"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
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
func WithSyncAdvertisementsProvider(provider app.SyncAdvertisementsProvider) TestOverlayEngineStubOption {
	return func(stub *TestOverlayEngineStub) {
		stub.syncAdvertisementsProvider = provider
	}
}

// WithSubmitTransactionProvider allows setting a custom SubmitTransactionProvider in a TestOverlayEngineStub.
// This can be used to mock transaction submission behavior during tests.
func WithSubmitTransactionProvider(provider app.SubmitTransactionProvider) TestOverlayEngineStubOption {
	return func(stub *TestOverlayEngineStub) {
		stub.submitTransactionProvider = provider
	}
}

// WithStartGASPSyncProvider allows setting a custom StartGASPSyncProvider in a TestOverlayEngineStub.
// This can be used to mock GASP sync behavior during tests.
func WithStartGASPSyncProvider(provider app.StartGASPSyncProvider) TestOverlayEngineStubOption {
	return func(stub *TestOverlayEngineStub) {
		stub.startGASPSyncProvider = provider
	}
}

// WithRequestSyncResponseProvider allows setting a custom RequestSyncResponseProvider in a TestOverlayEngineStub.
// This can be used to mock request sync response behavior during tests.
func WithRequestSyncResponseProvider(provider app.RequestSyncResponseProvider) TestOverlayEngineStubOption {
	return func(stub *TestOverlayEngineStub) {
		stub.requestSyncResponseProvider = provider
	}
}

// WithRequestForeignGASPNodeProvider allows setting a custom RequestForeignGASPNodeProvider in a TestOverlayEngineStub.
// This can be used to mock foreign GASP node request behavior during tests.
func WithRequestForeignGASPNodeProvider(provider app.RequestForeignGASPNodeProvider) TestOverlayEngineStubOption {
	return func(stub *TestOverlayEngineStub) {
		stub.requestForeignGASPNodeProvider = provider
	}
}

// TestOverlayEngineStub is a test implementation of the engine.OverlayEngineProvider interface.
// It is used to mock engine behavior in unit tests, allowing the simulation of various engine actions
// like submitting transactions and synchronizing advertisements.
type TestOverlayEngineStub struct {
	t                              *testing.T
	syncAdvertisementsProvider     app.SyncAdvertisementsProvider
	submitTransactionProvider      app.SubmitTransactionProvider
	startGASPSyncProvider          app.StartGASPSyncProvider
	requestSyncResponseProvider    app.RequestSyncResponseProvider
	requestForeignGASPNodeProvider app.RequestForeignGASPNodeProvider
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

// ProvideForeignGASPNode returns a foreign GASP node.
// It delegates to the requestForeignGASPNodeProvider if one is configured.
func (t TestOverlayEngineStub) ProvideForeignGASPNode(ctx context.Context, graphId *overlay.Outpoint, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error) {
	t.t.Helper()

	return t.requestForeignGASPNodeProvider.ProvideForeignGASPNode(ctx, graphId, outpoint, topic)
}

// ProvideForeignSyncResponse returns a foreign sync response.
// It delegates to the requestSyncResponseProvider if one is configured.
func (t TestOverlayEngineStub) ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {
	t.t.Helper()

	return t.requestSyncResponseProvider.ProvideForeignSyncResponse(ctx, initialRequest, topic)
}

// StartGASPSync starts the GASP synchronization process.
// It delegates to the startGASPSyncProvider if one is configured.
func (t TestOverlayEngineStub) StartGASPSync(ctx context.Context) error {
	t.t.Helper()

	return t.startGASPSyncProvider.StartGASPSync(ctx)
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
		t:                              t,
		submitTransactionProvider:      submitTransactionProviderAlwaysSuccessStub{},
		syncAdvertisementsProvider:     syncAdvertisementsProviderAlwaysSuccessStub{},
		startGASPSyncProvider:          startGASPSyncProviderAlwaysSuccessStub{},
		requestSyncResponseProvider:    requestSyncResponseProviderAlwaysSuccessStub{},
		requestForeignGASPNodeProvider: requestForeignGASPNodeProviderAlwaysSuccessStub{},
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

// startGASPSyncProviderAlwaysSuccessStub is a mock implementation of StartGASPSyncProvider that always succeeds.
// It is used as the default StartGASPSyncProvider in the TestOverlayEngineStub.
type startGASPSyncProviderAlwaysSuccessStub struct{}

// StartGASPSync simulates a successful GASP sync request call.
// It always returns nil, indicating that the synchronization was successful.
func (startGASPSyncProviderAlwaysSuccessStub) StartGASPSync(ctx context.Context) error {
	return nil
}

// requestSyncResponseProviderAlwaysSuccessStub is a mock implementation of RequestSyncResponseProvider that always succeeds.
// It is used as the default RequestSyncResponseProvider in the TestOverlayEngineStub.
type requestSyncResponseProviderAlwaysSuccessStub struct{}

// ProvideForeignSyncResponse simulates a successful foreign sync response request call.
// It always returns an empty GASPInitialResponse, indicating that the request was successful.
func (requestSyncResponseProviderAlwaysSuccessStub) ProvideForeignSyncResponse(ctx context.Context, initialRequest *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {
	return &core.GASPInitialResponse{
		UTXOList: []*overlay.Outpoint{},
		Since:    initialRequest.Since,
	}, nil
}

// requestForeignGASPNodeProviderAlwaysSuccessStub is a mock implementation of RequestForeignGASPNodeProvider that always succeeds.
// It is used as the default RequestForeignGASPNodeProvider in the TestOverlayEngineStub.
type requestForeignGASPNodeProviderAlwaysSuccessStub struct{}

// ProvideForeignGASPNode simulates a successful foreign GASP node request call.
// It always returns an empty GASPNode, indicating that the request was successful.
func (requestForeignGASPNodeProviderAlwaysSuccessStub) ProvideForeignGASPNode(ctx context.Context, graphId *overlay.Outpoint, outpoint *overlay.Outpoint, topic string) (*core.GASPNode, error) {
	return &core.GASPNode{}, nil
}
