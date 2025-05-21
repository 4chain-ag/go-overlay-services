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

// ProviderStateAsserter is an interface for asserting internal state after a test run.
type ProviderStateAsserter interface {
	AssertCalled()
}

// SyncAdvertisementsProvider extends app.SyncAdvertisementsProvider with the ability
// to assert whether it was called during a test.
type SyncAdvertisementsProvider interface {
	app.SyncAdvertisementsProvider
	ProviderStateAsserter
}

// SubmitTransactionProvider extends app.SubmitTransactionProvider with the ability
// to assert whether it was called during a test.
type SubmitTransactionProvider interface {
	app.SubmitTransactionProvider
	ProviderStateAsserter
}

type LookupServiceDocumentationProvider interface {
	app.LookupServiceDocumentationProvider
	ProviderStateAsserter
}

// TestOverlayEngineStubOption is a functional option type used to configure a TestOverlayEngineStub.
// It allows setting custom behaviors for different parts of the TestOverlayEngineStub.
type TestOverlayEngineStubOption func(*TestOverlayEngineStub)

// WithSubmitTransactionProvider allows setting a custom SubmitTransactionProvider in a TestOverlayEngineStub.
// This can be used to mock transaction submission behavior during tests.
func WithSubmitTransactionProvider(provider SubmitTransactionProvider) TestOverlayEngineStubOption {
	return func(stub *TestOverlayEngineStub) {
		stub.submitTransactionProvider = provider
	}
}

func WithLookupDocumentationProvider(provider LookupServiceDocumentationProvider) TestOverlayEngineStubOption {
	return func(stub *TestOverlayEngineStub) {
		stub.lookupDocumentationProvider = provider
	}
}

// WithSyncAdvertisementsProvider allows setting a custom SyncAdvertisementsProvider in a TestOverlayEngineStub.
// This can be used to mock advertisement synchronization behavior during tests.
func WithSyncAdvertisementsProvider(provider SyncAdvertisementsProvider) TestOverlayEngineStubOption {
	return func(stub *TestOverlayEngineStub) {
		stub.syncAdvertisementsProvider = provider
	}
}

// TestOverlayEngineStub is a test implementation of the engine.OverlayEngineProvider interface.
// It is used to mock engine behavior in unit tests, allowing the simulation of various engine actions
// like submitting transactions and synchronizing advertisements.
type TestOverlayEngineStub struct {
	t                           *testing.T
	submitTransactionProvider   SubmitTransactionProvider
	lookupDocumentationProvider LookupServiceDocumentationProvider
	syncAdvertisementsProvider SyncAdvertisementsProvider
}

// GetDocumentationForLookupServiceProvider returns documentation for a lookup service provider (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (s *TestOverlayEngineStub) GetDocumentationForLookupServiceProvider(provider string) (string, error) {
	s.t.Helper()

	return s.lookupDocumentationProvider.GetDocumentationForLookupServiceProvider(provider)
}

// GetDocumentationForTopicManager returns documentation for a topic manager (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (s *TestOverlayEngineStub) GetDocumentationForTopicManager(provider string) (string, error) {
	panic("unimplemented")
}

// GetUTXOHistory retrieves UTXO history for the given output (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (s *TestOverlayEngineStub) GetUTXOHistory(ctx context.Context, outpus *engine.Output, historySelector func(beef []byte, outputIndex uint32, currentDepth uint32) bool, currentDepth uint32) (*engine.Output, error) {
	panic("unimplemented")
}

// HandleNewMerkleProof processes a new Merkle proof for a transaction (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (s *TestOverlayEngineStub) HandleNewMerkleProof(ctx context.Context, txid *chainhash.Hash, proof *transaction.MerklePath) error {
	panic("unimplemented")
}

// ListLookupServiceProviders lists the available lookup service providers (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (s *TestOverlayEngineStub) ListLookupServiceProviders() map[string]*overlay.MetaData {
	panic("unimplemented")
}

// ListTopicManagers lists the available topic managers (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (s *TestOverlayEngineStub) ListTopicManagers() map[string]*overlay.MetaData {
	panic("unimplemented")
}

// Lookup performs a lookup query based on the provided LookupQuestion (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (s *TestOverlayEngineStub) Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error) {
	panic("unimplemented")
}

// ProvideForeignGASPNode returns a foreign GASP node (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (s *TestOverlayEngineStub) ProvideForeignGASPNode(ctx context.Context, graphId *overlay.Outpoint, outpoins *overlay.Outpoint, topic string) (*core.GASPNode, error) {
	panic("unimplemented")
}

// ProvideForeignSyncResponse returns a foreign sync response (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (s *TestOverlayEngineStub) ProvideForeignSyncResponse(ctx context.Context, initialRequess *core.GASPInitialRequest, topic string) (*core.GASPInitialResponse, error) {
	panic("unimplemented")
}

// StartGASPSync starts the GASP synchronization process (unimplemented).
// This is a placeholder function meant to be overridden in actual implementations.
func (s *TestOverlayEngineStub) StartGASPSync(ctx context.Context) error {
	panic("unimplemented")
}

// Submit processes a transaction submission and returns a steak or error based on the provided inputs.
// It calls the Submit method of the configured SubmitTransactionProvider and handles the steak callback.
func (s *TestOverlayEngineStub) Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode engine.SumbitMode, onSteakReady engine.OnSteakReady) (overlay.Steak, error) {
	s.t.Helper()

	return s.submitTransactionProvider.Submit(ctx, taggedBEEF, mode, onSteakReady)
}

// SyncAdvertisements synchronizes advertisements using the configured SyncAdvertisementsProvider.
// It calls the SyncAdvertisements method of the provider and handles the result.
func (s *TestOverlayEngineStub) SyncAdvertisements(ctx context.Context) error {
	s.t.Helper()

	return s.syncAdvertisementsProvider.SyncAdvertisements(ctx)
}

// AssertProvidersState asserts that all configured providers were used as expected.
func (s *TestOverlayEngineStub) AssertProvidersState() {
	s.t.Helper()

	providers := []ProviderStateAsserter{
		s.submitTransactionProvider,
		s.lookupDocumentationProvider,
		s.syncAdvertisementsProvider,
	}
	for _, p := range providers {
		p.AssertCalled()
	}
}

// NewTestOverlayEngineStub creates and returns a new instance of TestOverlayEngineStub with the provided options.
// The options allow for configuring custom providers for transaction submission and advertisement synchronization.
func NewTestOverlayEngineStub(t *testing.T, opts ...TestOverlayEngineStubOption) *TestOverlayEngineStub {
	stub := TestOverlayEngineStub{
		t:                           t,
		submitTransactionProvider:   NewSubmitTransactionProviderMock(t, SubmitTransactionProviderMockExpectations{SubmitCall: false}),
		lookupDocumentationProvider: NewLookupServiceDocumentationProviderMock(t, LookupServiceDocumentationProviderMockExpectations{DocumentationCall: false}),
		syncAdvertisementsProvider: NewSyncAdvertisementsProviderMock(t, SyncAdvertisementsProviderMockExpectations{SyncAdvertisementsCall: false}),
	}

	for _, opt := range opts {
		opt(&stub)
	}
	return &stub
}
