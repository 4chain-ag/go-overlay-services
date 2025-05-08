package main

import (
	"context"
	"flag"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/config"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/config/loaders"
	"github.com/gookit/slog"
)

func main() {
	configPath := flag.String("config", loaders.DefaultConfigFilePath, "Path to the configuration file")
	flag.Parse()

	loader := loaders.NewLoader(config.NewDefault, "OVERLAY")
	if err := loader.SetConfigFilePath(*configPath); err != nil {
		slog.Fatalf("Invalid config file path: %v", err)
	}

	cfg, err := loader.Load()
	if err != nil {
		slog.Fatalf("failed to load config: %v", err)
	}

	if err := config.PrettyPrintAs(cfg, "json"); err != nil {
		slog.Fatalf("failed to pretty print config: %v", err)
	}

	srv := server2.New(server2.WithConfig(&cfg.Server))
	<-srv.ListenAndServe(context.Background())
}
