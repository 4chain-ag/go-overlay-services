package server2

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/adapters"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Config holds the configuration settings for the HTTP server
type Config struct {
	// AppName is the name of the application.
	AppName string `mapstructure:"app_name"`

	// Port is the TCP port on which the server will listen.
	Port int `mapstructure:"port"`

	// Addr is the address the server will bind to.
	Addr string `mapstructure:"addr"`

	// ServerHeader is the value of the Server header returned in HTTP responses.
	ServerHeader string `mapstructure:"server_header"`

	// AdminBearerToken is the token required to access admin-only endpoints.
	AdminBearerToken string `mapstructure:"admin_bearer_token"`

	// OctetStreamLimit  defines the maximum allowed bytes read size (in bytes).
	// This limit by default is set to 1GB to protect against excessively large payloads.
	OctetStreamLimit int64 `mapstructure:"octet_stream_limit"`
}

// DefaultConfig provides a default configuration with reasonable values for local development.
var DefaultConfig = Config{
	AppName:          "Overlay API v0.0.0",
	Port:             3000,
	Addr:             "localhost",
	ServerHeader:     "Overlay API",
	AdminBearerToken: uuid.NewString(),
	OctetStreamLimit: middleware.ReadBodyLimit1GB,
}

// ServerOption defines a functional option for configuring an HTTP server.
// These options allow for flexible setup of middlewares and configurations.
type ServerOption func(*ServerHTTP)

// WithMiddleware adds a Fiber middleware handler to the HTTP server configuration.
// It returns a ServerOption that appends the given middleware to the server's middleware stack.
func WithMiddleware(f fiber.Handler) ServerOption {
	return func(s *ServerHTTP) {
		s.middleware = append(s.middleware, f)
	}
}

// WithEngine sets the overlay engine provider for the HTTP server.
// It configures the ServerHTTP handlers to use the provided engine implementation.
func WithEngine(provider engine.OverlayEngineProvider) ServerOption {
	return func(s *ServerHTTP) {
		s.submitTransactionHandler = ports.NewSubmitTransactionHandler(provider, ports.RequestTimeout)
		s.syncAdvertisementsHandler = ports.NewSyncAdvertisementsHandler(provider)
		s.lookupDocumentationHandler = ports.NewLookupProviderDocumentationHandler(provider)
	}
}

// WithSubmitTransactionHandlerResponseTime sets the submit transaction handler's response time threshold
// for the HTTP server. This timeout defines how long the handler waits before returning a request timeout response.
func WithSubmitTransactionHandlerResponseTime(provider app.SubmitTransactionProvider, timeout time.Duration) ServerOption {
	return func(s *ServerHTTP) {
		s.submitTransactionHandler = ports.NewSubmitTransactionHandler(provider, timeout)
	}
}

// WithOctetStreamLimit returns a ServerOption that sets the maximum allowed size (in bytes)
// for incoming requests with Content-Type: application/octet-stream.
// This is useful for controlling memory usage when clients upload large binary payloads.
//
// Example: To limit uploads to 512MB:
//
//	WithOctetStreamLimit(512 * 1024 * 1024)
func WithOctetStreamLimit(limit int64) ServerOption {
	return func(s *ServerHTTP) {
		s.cfg.OctetStreamLimit = limit
	}
}

// WithAdminBearerToken sets the admin bearer token used for authenticating
// admin routes on the HTTP server.
// It returns a ServerOption that applies this configuration to ServerHTTP.
func WithAdminBearerToken(token string) ServerOption {
	return func(s *ServerHTTP) {
		s.cfg.AdminBearerToken = token
	}
}

// WithConfig sets the configuration for the HTTP server using the provided Config.
// It initializes a new Fiber application with the specified server settings.
// Returns a ServerOption to apply during server setup.
func WithConfig(cfg *Config) ServerOption {
	return func(s *ServerHTTP) {
		s.cfg = cfg
		s.app = fiber.New(fiber.Config{
			CaseSensitive: true,
			StrictRouting: true,
			ServerHeader:  cfg.ServerHeader,
			AppName:       cfg.AppName,
		})
	}
}

// ServerHTTP represents the HTTP server instance, including configuration,
// Fiber app instance, middleware stack, and registered request handlers.
type ServerHTTP struct {
	cfg        *Config         // cfg holds the server configuration settings.
	app        *fiber.App      // app is the Fiber application instance serving HTTP requests.
	middleware []fiber.Handler // middleware is a list of Fiber middleware functions to be applied globally.

	// Handlers for processing incoming HTTP requests:
	submitTransactionHandler   *ports.SubmitTransactionHandler           // submitTransactionHandler handles transaction submission requests.
	syncAdvertisementsHandler  *ports.SyncAdvertisementsHandler          // syncAdvertisementsHandler handles advertisement sync requests.
	lookupDocumentationHandler *ports.LookupProviderDocumentationHandler // lookupDocumentationHandler handles lookup service documentation requests.
}

// SocketAddr builds the address string for binding.
func (s *ServerHTTP) SocketAddr() string {
	return fmt.Sprintf("%s:%d", s.cfg.Addr, s.cfg.Port)
}

// ListenAndServe starts the HTTP server and listens for termination signals.
// It returns a channel that will be closed once the shutdown is complete.
func (s *ServerHTTP) ListenAndServe(ctx context.Context) <-chan struct{} {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

		select {
		case <-sigint:
			slog.Info("Shutdown signal received. Cleaning up...")
		case <-ctx.Done():
			slog.Info("Shutdown context canceled. Cleaning up...")
		}

		if err := s.app.ShutdownWithContext(ctx); err != nil {
			slog.Error("HTTP shutdown error", slog.Attr{Key: "server_shutdown_err", Value: slog.StringValue(err.Error())})
		}

		close(idleConnsClosed)
	}()

	go func() {
		err := s.app.Listen(s.SocketAddr())
		isNotErrServerClosed := err != nil && !errors.Is(err, http.ErrServerClosed)
		if isNotErrServerClosed {
			slog.Error("HTTP shutdown error", slog.Attr{Key: "server_listen_err", Value: slog.StringValue(err.Error())})
		}
	}()

	return idleConnsClosed
}

// New creates and configures a new instance of ServerHTTP.
// It initializes the application with default settings and middleware, registers OpenAPI handlers,
// sets up transaction submission and advertisement synchronization handlers using the provided OverlayEngineProvider,
// and applies any optional functional configuration options passed via opts.
//
// The returned ServerHTTP instance is ready to be started by calling .Listen(...) or integrated into tests.
func New(opts ...ServerOption) *ServerHTTP {
	noop := adapters.NewNoopEngineProvider()
	srv := &ServerHTTP{
		submitTransactionHandler:   ports.NewSubmitTransactionHandler(noop, app.DefaultSubmitTransactionTimeout),
		syncAdvertisementsHandler:  ports.NewSyncAdvertisementsHandler(noop),
		lookupDocumentationHandler: ports.NewLookupProviderDocumentationHandler(noop),
		cfg:                        &DefaultConfig,
		app: fiber.New(fiber.Config{
			CaseSensitive: true,
			StrictRouting: true,
			ServerHeader:  "Overlay API",
			AppName:       "Overlay API v0.0.0",
		}),
		middleware: middleware.BasicMiddlewareGroup(),
	}

	for _, m := range srv.middleware {
		srv.app.Use(m)
	}

	for _, opt := range opts {
		opt(srv)
	}

	srv.registerRoutes()
	return srv
}

func (s *ServerHTTP) registerRoutes() {
	api := s.app.Group("/api")

	v1 := api.Group("/v1")
	v1.Post("/submit", middleware.LimitOctetStreamBodyMiddleware(s.cfg.OctetStreamLimit), s.submitTransactionHandler.SubmitTransaction)
	v1.Get("/getDocumentationForLookupServiceProvider", s.lookupDocumentationHandler.GetDocumentation)

	admin := v1.Group("/admin", middleware.BearerTokenAuthorizationMiddleware(s.cfg.AdminBearerToken))
	admin.Post("/syncAdvertisements", s.syncAdvertisementsHandler.Handle)
}
