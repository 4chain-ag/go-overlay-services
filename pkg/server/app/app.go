package app

import (
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/queries"
)

// Commands aggregate all the supported commands by the overlay API.
type Commands struct {
	SubmitTransactionHandler      *commands.SubmitTransactionHandler
	SyncAdvertismentsHandler      *commands.SyncAdvertisementsHandler
	StartGASPSyncHandler          *commands.StartGASPSyncHandler
	RequestForeignGASPNodeHandler *commands.RequestForeignGASPNodeHandler
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
func New(provider engine.OverlayEngineProvider) *Application {
	if provider == nil {
		panic("overlay engine provider is nil")
	}

	submitHandler, err := commands.NewSubmitTransactionCommandHandler(provider)
	if err != nil {
		panic(fmt.Sprintf("failed to create submit transaction handler: %v", err))
	}

	return &Application{
		Commands: &Commands{
			SubmitTransactionHandler:      submitHandler,
			SyncAdvertismentsHandler:      commands.NewSyncAdvertisementsCommandHandler(provider),
			StartGASPSyncHandler:          commands.NewStartGASPSyncHandler(provider),
			RequestForeignGASPNodeHandler: commands.NewRequestForeignGASPNodeHandler(provider),
		},
		Queries: &Queries{
			TopicManagerDocumentationHandler: queries.NewTopicManagerDocumentationHandler(provider),
		},
	}
}
