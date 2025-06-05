package server2

import (
	"context"
	"fmt"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/core/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/adapters"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/middleware"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
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

	// ArcApiKey is the API key for ARC service integration.
	ArcApiKey string `mapstructure:"arc_api_key"`

	// ArcCallbackToken is the token for authenticating ARC callback requests.
	ArcCallbackToken string `mapstructure:"arc_callback_token"`

	// OctetStreamLimit defines the maximum allowed bytes read size (in bytes).
	// This limit by default is set to 1GB to protect against excessively large payloads.
	OctetStreamLimit int64 `mapstructure:"octet_stream_limit"`

	// ConnectionReadTimeout defines the maximum duration an active connection is allowed to stay open.
	// Once this threshold is exceeded, the connection will be forcefully closed.
	ConnectionReadTimeout time.Duration `mapstructure:"connection_read_timeout_limit"`
}

// DefaultConfig provides a default configuration with reasonable values for local development.
var DefaultConfig = Config{
	AppName:               "Overlay API v0.0.0",
	Port:                  3000,
	Addr:                  "localhost",
	ServerHeader:          "Overlay API",
	AdminBearerToken:      uuid.NewString(),
	ArcApiKey:             "",
	ArcCallbackToken:      uuid.NewString(),
	OctetStreamLimit:      middleware.ReadBodyLimit1GB,
	ConnectionReadTimeout: 10 * time.Second,
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
		s.provider = provider
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

// WithArcApiKey sets the ARC API key used for ARC service integration.
// It returns a ServerOption that applies this configuration to ServerHTTP.
func WithArcApiKey(apiKey string) ServerOption {
	return func(s *ServerHTTP) {
		s.cfg.ArcApiKey = apiKey
	}
}

// WithArcCallbackToken sets the ARC callback token used for authenticating
// ARC callback requests on the HTTP server.
// It returns a ServerOption that applies this configuration to ServerHTTP.
func WithArcCallbackToken(token string) ServerOption {
	return func(s *ServerHTTP) {
		s.cfg.ArcCallbackToken = token
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

// WithConfig sets the configuration for the HTTP server using the provided Config.
// It initializes a new Fiber application with the specified server settings.
// Returns a ServerOption to apply during server setup.
func WithConfig(cfg Config) ServerOption {
	return func(s *ServerHTTP) {
		s.cfg = cfg
		s.app = newFiberApp(cfg)
	}
}

// ServerHTTP represents the HTTP server instance, including configuration,
// Fiber app instance, middleware stack, and registered request handlers.
type ServerHTTP struct {
	cfg        Config          // cfg holds the server configuration settings.
	app        *fiber.App      // app is the Fiber application instance serving HTTP requests.
	middleware []fiber.Handler // middleware is a list of Fiber middleware functions to be applied globally.

	// Handlers for processing incoming HTTP requests
	registry *ports.HandlerRegistryService
	provider engine.OverlayEngineProvider // provider is stored to create registry after all options applied
}

// SocketAddr builds the address string for binding.
func (s *ServerHTTP) SocketAddr() string {
	return fmt.Sprintf("%s:%d", s.cfg.Addr, s.cfg.Port)
}

// ListenAndServe starts the HTTP server and begins listening on the configured socket address.
// It blocks until the server is stopped or an error occurs.
func (s *ServerHTTP) ListenAndServe(ctx context.Context) error {
	return s.app.Listen(s.SocketAddr())
}

// Shutdown gracefully shuts down the HTTP server using the provided context,
// allowing ongoing requests to complete within the context's deadline.
func (s *ServerHTTP) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

// New creates and configures a new instance of ServerHTTP.
// It initializes the application with default settings and middleware, registers OpenAPI handlers,
// sets up transaction submission and advertisement synchronization handlers using the provided OverlayEngineProvider,
// and applies any optional functional configuration options passed via opts.
func New(opts ...ServerOption) *ServerHTTP {
	srv := &ServerHTTP{
		cfg: DefaultConfig,
		app: newFiberApp(DefaultConfig),
	}

	for _, o := range opts {
		o(srv)
	}

	provider := srv.provider
	if provider == nil {
		provider = adapters.NewNoopEngineProvider()
	}
	srv.registry = ports.NewHandlerRegistryService(provider, ports.ArcIngestHandlerConfig{
		ArcApiKey:        srv.cfg.ArcApiKey,
		ArcCallbackToken: srv.cfg.ArcCallbackToken,
	})

	openapi.RegisterHandlersWithOptions(srv.app, srv.registry, openapi.FiberServerOptions{
		HandlerMiddleware: []fiber.Handler{
			middleware.BearerTokenAuthorizationMiddleware(srv.cfg.AdminBearerToken),
		},
		GlobalMiddleware: middleware.BasicMiddlewareGroup(middleware.BasicMiddlewareGroupConfig{
			EnableStackTrace: true,
			OctetStreamLimit: srv.cfg.OctetStreamLimit,
		}),
	})

	srv.app.Get("/metrics", monitor.New(monitor.Config{Title: "Overlay-services API"}))

	return srv
}

// newFiberApp creates and returns a new instance of a fiber.App with the provided configuration and middleware.
// The app is configured with case-sensitive routing, strict routing, custom server headers, and read timeout settings.
// Additionally, any provided middleware handlers are applied to the app.
func newFiberApp(cfg Config) *fiber.App {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  cfg.ServerHeader,
		AppName:       cfg.AppName,
		ReadTimeout:   cfg.ConnectionReadTimeout,
		ErrorHandler:  ports.ErrorHandler(),
	})

	return app
}
