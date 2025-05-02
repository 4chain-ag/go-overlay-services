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
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
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
}

// DefaultConfig provides a default configuration with reasonable values for local development.
var DefaultConfig = Config{
	AppName:          "Overlay API v0.0.0",
	Port:             3000,
	Addr:             "localhost",
	ServerHeader:     "Overlay API",
	AdminBearerToken: uuid.NewString(),
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
func WithEngine(e engine.OverlayEngineProvider) ServerOption {
	return func(s *ServerHTTP) {
		s.submitTransactionHandler = ports.NewSubmitTransactionHandler(e)
		s.advertisementsSyncHandler = ports.NewAdvertisementsSyncHandler(e)
		s.lookupServiceDocumentationHandler = ports.NewLookupServiceDocumentationHandler(e)
		s.topicManagerDocumentationHandler = ports.NewTopicManagerDocumentationHandler(e)
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
	submitTransactionHandler          *ports.SubmitTransactionHandler          // submitTransactionHandler handles transaction submission requests.
	advertisementsSyncHandler         *ports.AdvertisementsSyncHandler         // advertisementsSyncHandler handles advertisement sync requests.
	lookupServiceDocumentationHandler *ports.LookupServiceDocumentationHandler // lookupServiceDocumentationHandler handles lookup service documentation requests.
	topicManagerDocumentationHandler  *ports.TopicManagerDocumentationHandler  // topicManagerDocumentationHandler handles topic manager documentation requests.
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

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

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
	noop := newNoopEngineProvider()
	srv := &ServerHTTP{
		submitTransactionHandler:          ports.NewSubmitTransactionHandler(noop),
		advertisementsSyncHandler:         ports.NewAdvertisementsSyncHandler(noop),
		lookupServiceDocumentationHandler: ports.NewLookupServiceDocumentationHandler(noop),
		topicManagerDocumentationHandler:  ports.NewTopicManagerDocumentationHandler(noop),
		cfg:                               &DefaultConfig,
		app: fiber.New(fiber.Config{
			CaseSensitive: true,
			StrictRouting: true,
			ServerHeader:  "Overlay API",
			AppName:       "Overlay API v0.0.0",
		}),
		middleware: []fiber.Handler{
			requestid.New(),
			idempotency.New(),
			cors.New(),
			recover.New(recover.Config{EnableStackTrace: true}), // TODO: stack trace should be disabled after releasing to the prod.
			logger.New(logger.Config{
				Format:     "date=${time} request_id=${locals:requestid} status=${status} method=${method} path=${path}​\n",
				TimeFormat: "02-Jan-2006 15:04:05",
			}),
			pprof.New(pprof.Config{Prefix: "/api/v1"}),
		},
	}

	for _, m := range srv.middleware {
		srv.app.Use(m)
	}

	for _, opt := range opts {
		opt(srv)
	}

	api := srv.app.Group("/api")
	v1 := api.Group("/v1")
	v1.Post("/submit", srv.submitTransactionHandler.SubmitTransaction)
	v1.Get("/getDocumentationForLookupServiceProvider", srv.lookupServiceDocumentationHandler.GetDocumentation)
	v1.Get("/getDocumentationForTopicManager", srv.topicManagerDocumentationHandler.GetDocumentation)

	admin := v1.Group("/admin", middleware.BearerTokenAuthorizationMiddleware(srv.cfg.AdminBearerToken))
	admin.Post("/syncAdvertisements", srv.advertisementsSyncHandler.Handle)

	return srv
}
