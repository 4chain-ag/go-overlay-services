package engine

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/core/gasp"
	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/lookup"
)

// OverlayEngineProvider defines the contract for the overlay engine.
// Note: The contract definition is still in development and will be updated after
// migrating the engine code.
type OverlayEngineProvider interface {
	Submit(ctx context.Context, taggedBEEF overlay.TaggedBEEF, mode SumbitMode, onSteakReady OnSteakReady) (overlay.Steak, error)
	Lookup(ctx context.Context, question *lookup.LookupQuestion) (*lookup.LookupAnswer, error)
	GetUTXOHistory(ctx context.Context, output *Output, historySelector func(beef []byte, outputIndex uint32, currentDepth uint32) bool, currentDepth uint32) (*Output, error)
	SyncAdvertisements() error
	StartGASPSync() error
	ProvideForeignSyncResponse(initialRequest *gasp.InitialRequest, topic string) (*gasp.InitialResponse, error)
	ProvideForeignGASPNode(graphId string, txid string, outputIndex uint32) (*gasp.GASPNode, error)
	ListTopicManagers() map[string]*overlay.MetaData
	ListLookupServiceProviders() map[string]*overlay.MetaData
	GetDocumentationForLookupServiceProvider(provider string) (string, error)
	GetDocumentationForTopicManager(provider string) (string, error)
}
