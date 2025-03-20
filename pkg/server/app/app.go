package app

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/queries"
)

// OverlayEngineProvider defines the contract for the overlay engine.
// Note: The contract definition is still in development and will be updated after
// migrating the engine code.
type OverlayEngineProvider interface {
	SubmitTransaction(ctx context.Context) error
	SyncAdvertisments(ctx context.Context) error
	GetTopicManagerDocumentation(ctx context.Context) error
}

// Commands aggregate all the supported commands by the overlay API.
type Commands struct {
	SubmitTransactionHandler *commands.SubmitTransactionHandler
	SyncAdvertismentsHandler *commands.SyncAdvertismentsHandler
}

// Queries aggregate all the supported queries by the overlay API.
type Queries struct {
	TopicManagerDocumentationHandler *queries.TopicManagerDocumentationHandler
}

// Application aggregates queries and commands supported by the overlay API.
type Application struct {
	Commands *Commands
	Queries  *Queries
}

// New returns an instance of an Application with intialized commands and queries
// utilizing an implementation of OverlayEngineProvider. If the provided argument is nil, it triggers a panic.
func New(provider OverlayEngineProvider) *Application {
	if provider == nil {
		panic("overlay engine provider is nil")
	}
	return &Application{
		Commands: &Commands{
			SubmitTransactionHandler: commands.NewSubmitTransactionCommandHandler(provider),
			SyncAdvertismentsHandler: commands.NewSyncAdvertismentsHandler(provider),
		},
		Queries: &Queries{
			TopicManagerDocumentationHandler: queries.NewTopicManagerDocumentationHandler(provider),
		},
	}
}
