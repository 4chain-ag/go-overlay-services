package server

import (
	"fmt"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
)

// HTTPOption defines a functional option for configuring an HTTP server.
// These options allow for flexible setup of middlewares and configurations.
type HTTPOption func(*HTTP)

// WithMiddleware adds custom middleware to the HTTP server.
// The execution order depends on the sequence in which the middlewares are passed
func WithMiddleware(f func(http.Handler) http.Handler) HTTPOption {
	return func(h *HTTP) {
		h.middlewares = append(h.middlewares, adaptor.HTTPMiddleware(f))
	}
}

// WithConfig sets the HTTP server configuration based on the given definition.
func WithConfig(cfg *Config) HTTPOption {
	return func(h *HTTP) {
		h.cfg = cfg
	}
}

// Config describes the configuration of the HTTP server instance.
type Config struct {
	Addr string
	Port int
}

// SocketAddr returns the socket address string based on the configured address and port combination.
func (c *Config) SocketAddr() string { return fmt.Sprintf("%s:%d", c.Addr, c.Port) }

// HTTP manages connections to the overlay server instance. It accepts and responds to client sockets,
// using idempotency to improve fault tolerance and mitigate duplicated requests.
// It applies all configured options along with the list of middlewares."
type HTTP struct {
	middlewares []fiber.Handler
	app         *fiber.App
	cfg         *Config
}

// New returns an instance of the HTTP server and applies all specified functional options before starting it.
func New(opts ...HTTPOption) *HTTP {
	overlayAPI := app.New(&engine.NoopEngineProvider{})
	http := HTTP{
		app: fiber.New(fiber.Config{
			CaseSensitive: true,
			StrictRouting: true,
			ServerHeader:  "Overlay API",
			AppName:       "Overlay API v0.0.0",
		}),
		middlewares: []fiber.Handler{idempotency.New()},
	}
	for _, o := range opts {
		o(&http)
	}
	for _, m := range http.middlewares {
		http.app.Use(m)
	}

	// Routes:
	api := http.app.Group("/api")
	v1 := api.Group("/v1")

	// Non-Admin:
	v1.Post("/submit", overlayAPI.Commands.SubmitTransactionHandler.Handle)
	v1.Get("/topic-managers", overlayAPI.Queries.TopicManagerDocumentationHandler.Handle)

	// Admin:
	admin := v1.Group("/admin")
	admin.Post("/advertisements-sync", overlayAPI.Commands.SyncAdvertismentsHandler.Handle)

	return &http
}

// ListenAndServe handles HTTP requests from the configured socket address."
func (h *HTTP) ListenAndServe() error {
	if err := h.app.Listen(h.cfg.SocketAddr()); err != nil {
		return fmt.Errorf("http server: fiber app listen failed: %w", err)
	}
	return nil
}
