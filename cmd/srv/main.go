package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/4chain-ag/go-overlay-services/pkg/server2"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/config"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/config/loaders"
)

func main() {
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	configPath := flag.String("config", loaders.DefaultConfigFilePath, "Path to the configuration file")
	flag.Parse()

	cfg, err := config.LoadFromPath(*configPath, "OVERLAY")
	if err != nil {
		return fmt.Errorf("load config op failed: %w", err)
	}

	ctx := context.Background()
	srv := server2.New(server2.WithConfig(cfg))
	done := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown err: %v", err)
		}
		close(done)
	}()

	err = srv.ListenAndServe(ctx)
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server listen and serve op failure: %w", err)
	}
	<-done

	return nil
}
