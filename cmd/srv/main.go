package main

import (
	"flag"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/config"
	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/4chain-ag/go-overlay-services/pkg/server/config"
	"github.com/4chain-ag/go-overlay-services/pkg/server/mongo"
	"github.com/gookit/slog"
)

func main() {
	configPath := flag.String("C", config.DefaultConfigFilePath, "Path to the configuration file")
	flag.Parse()

	loader := config.NewLoader("OVERLAY")
	if err := loader.SetConfigFilePath(*configPath); err != nil {
		slog.Fatalf("Invalid config file path: %v", err)
	}

	if err := loader.PrettyPrintAs("json"); err != nil {
		slog.Fatalf("failed to pretty print config: %v", err)
	}

	cfg, err := loader.Load()
	if err != nil {
		slog.Fatalf("failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		slog.Fatalf("Invalid configuration: %v", err)
	}

	var mongoClient *mongo.Client
	if cfg.Mongo.URI != "" {
		mongoConn, err := mongo.New(cfg)
		if err != nil {
			slog.Fatalf("failed to connect to MongoDB: %v", err)
		}
		mongoClient = mongoConn
	}

	opts := []server.HTTPOption{
		server.WithConfig(&cfg),
		server.WithMiddleware(loggingMiddleware),
	}

	if mongoClient != nil {
		opts = append(opts, server.WithMongo(mongoClient.Client()))
	}

	httpAPI := server.New(opts...)

	if err := httpAPI.ListenAndServe(); err != nil {
		slog.Fatalf("HTTP server failed: %v", err)
	}
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
