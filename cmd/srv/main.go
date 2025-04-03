package main

import (
	"context"
	"flag"
	"net/http"

	config "github.com/4chain-ag/go-overlay-services/pkg/appconfig"
	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/gookit/slog"
)

func main() {
	configPath := flag.String("C", config.DefaultConfigFilePath, "Path to the configuration file")
	flag.Parse()

	loader := config.NewLoader("OVERLAY")
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

	if err := cfg.Validate(); err != nil {
		slog.Fatalf("Invalid configuration: %v", err)
	}

	opts := []server.HTTPOption{
		server.WithConfig(&cfg),
		server.WithMiddleware(loggingMiddleware),
		server.WithMongo(),
	}

	httpAPI, err := server.New(opts...)
	if err != nil {
		slog.Fatalf("Failed to create HTTP server: %v", err)
	}

	// Graceful shutdown handling
	ctx := context.Background()
	idleConnsClosed := httpAPI.StartWithGracefulShutdown(ctx)
	<-idleConnsClosed
	slog.Info("Server shutdown completed.")
}

// loggingMiddleware is a custom definition of the logging middleware format accepted by the HTTP API.
func loggingMiddleware(next http.Handler) http.Handler {
	slog.SetLogLevel(slog.DebugLevel)
	slog.SetFormatter(slog.NewJSONFormatter(func(f *slog.JSONFormatter) {
		f.PrettyPrint = true
	}))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.WithFields(slog.M{
			"category":    "service",
			"method":      r.Method,
			"remote-addr": r.RemoteAddr,
			"request-uri": r.RequestURI,
		})
		logger.Info("log-line")
		next.ServeHTTP(w, r)
	})
}
