package main

import (
	"context"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/gookit/slog"
)

func main() {
	ctx := context.Background()
	srv := server2.New()

	<-srv.ListenAndServe(ctx)
	slog.Info("Server shutdown completed.")
}
