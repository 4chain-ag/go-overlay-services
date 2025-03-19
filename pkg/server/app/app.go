package app

import (
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/queries"
)

type Commands struct {
	SubmitTransactionHandler *commands.SubmitTransactionHandler
	SyncAdvertismentsHandler *commands.SyncAdvertismentsHandler
}

type Queries struct {
	TopicManagerDocumentationHandler *queries.TopicManagerDocumentationHandler
}

type Application struct {
	Commands *Commands
	Queries  *Queries
}
