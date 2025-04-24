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
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/overlayhttp"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
type ServerOption func(*Server)

// WithMiddleware adds Fiber middleware to the HTTP server.
func WithMiddleware(f fiber.Handler) ServerOption {
	return func(s *Server) {
		s.middleware = append(s.middleware, f)
	}
}

// WithConfig sets the configuration for the HTTP server.
func WithConfig(cfg *Config) ServerOption {
	return func(s *Server) {
		s.cfg = cfg
		s.app = fiber.New(fiber.Config{
			CaseSensitive: true,
			StrictRouting: true,
			ServerHeader:  cfg.ServerHeader,
			AppName:       cfg.AppName,
		})
	}
}

// WithEngine sets the overlay engine provider for the HTTP server.
func WithEngine(provider engine.OverlayEngineProvider) ServerOption {
	return func(s *Server) {
		s.overlayEngineProvider = provider
	}
}

// Server manages connections to the overlay server instance. It accepts and responds to client sockets,
// using idempotency to improve fault tolerance and mitigate duplicated requests.
// It applies all configured options along with the list of middlewares.
type Server struct {
	cfg                   *Config
	app                   *fiber.App
	middleware            []fiber.Handler
	overlayEngineProvider engine.OverlayEngineProvider
}

// New returns an instance of the HTTP server and applies all specified functional options before starting it.
func New(opts ...ServerOption) *Server {
	srv := &Server{
		cfg:                   &DefaultConfig,
		overlayEngineProvider: NewNoopEngineProvider(),
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
			recover.New(recover.Config{EnableStackTrace: true}),
			logger.New(logger.Config{
				Format:     "${locals:requestid} ${status} - ${method} ${path}​\n",
				TimeFormat: "02-Jan-2006",
			}),
		},
	}

	for _, m := range srv.middleware {
		srv.app.Use(m)
	}
	for _, opt := range opts {
		opt(srv)
	}

	handler := overlayhttp.NewHTTPHandler(srv.overlayEngineProvider)
	openapi.RegisterHandlers(srv.app, handler)

	return srv
}

// SocketAddr builds the address string for binding.
func (s *Server) SocketAddr() string {
	return fmt.Sprintf("%s:%d", s.cfg.Addr, s.cfg.Port)
}

// ListenAndServe starts the HTTP server and listens for termination signals.
// It returns a channel that will be closed once the shutdown is complete.
func (s *Server) ListenAndServe(ctx context.Context) <-chan struct{} {
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
