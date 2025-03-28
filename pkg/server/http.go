package server

import (
	"fmt"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
)

// HTTPOption defines a functional option for configuring an HTTP server.
// These options allow for flexible setup of middlewares and configurations.
type HTTPOption func(*HTTP)

// WithMiddleware adds net/http-style middleware to the HTTP server.
func WithMiddleware(f func(http.Handler) http.Handler) HTTPOption {
	return func(h *HTTP) {
		h.middlewares = append(h.middlewares, adaptor.HTTPMiddleware(f))
	}
}

// WithFiberMiddleware adds a Fiber-style middleware to the HTTP server.
func WithFiberMiddleware(m fiber.Handler) HTTPOption {
	return func(h *HTTP) {
		h.middlewares = append(h.middlewares, m)
	}
}

// WithCORS allows configuring CORS using a custom or default Fiber config.
func WithCORS(config ...cors.Config) HTTPOption {
	c := cors.ConfigDefault
	if len(config) > 0 {
		c = config[0]
	}
	return WithFiberMiddleware(cors.New(c))
}

// WithConfig sets the HTTP server configuration.
func WithConfig(cfg *config.Config) HTTPOption {
	return func(h *HTTP) {
		h.cfg = cfg
	}
}

// HTTP manages the Fiber server and its configuration.
type HTTP struct {
	middlewares []fiber.Handler
	app         *fiber.App
	cfg         *config.Config
}

// New returns a new HTTP server with the provided options.
func New(opts ...HTTPOption) *HTTP {
	overlayAPI := app.New(NewNoopEngineProvider())
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
	v1.Post("/submit", adaptor.HTTPHandlerFunc(overlayAPI.Commands.SubmitTransactionHandler.Handle))
	v1.Get("/topic-managers", adaptor.HTTPHandlerFunc(overlayAPI.Queries.TopicManagerDocumentationHandler.Handle))

	// Admin:
	admin := v1.Group("/admin")
	admin.Post("/advertisements-sync", adaptor.HTTPHandlerFunc(overlayAPI.Commands.SyncAdvertismentsHandler.Handle))

	return &http
}

// SocketAddr builds the address string for binding.
func (h *HTTP) SocketAddr() string {
	return fmt.Sprintf("%s:%d", h.cfg.Addr, h.cfg.Port)
}

// ListenAndServe starts the Fiber app using the configured socket address.
func (h *HTTP) ListenAndServe() error {
	if err := h.app.Listen(h.SocketAddr()); err != nil {
		return fmt.Errorf("http server: fiber app listen failed: %w", err)
	}
	return nil
}

// App exposes the underlying Fiber app.
func (h *HTTP) App() *fiber.App {
	return h.app
}
