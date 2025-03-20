package server

import (
	"fmt"
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/queries"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
)

type HTTPOption func(*HTTP)

func WithLoggingMiddleware(f func(http.Handler) http.Handler) HTTPOption {
	return func(h *HTTP) {
		h.middlewares = append(h.middlewares, adaptor.HTTPMiddleware(f))
	}
}

func WithConfig(cfg *Config) HTTPOption {
	return func(h *HTTP) {
		h.cfg = cfg
	}
}

func WithAdminBearerToken(s string) HTTPOption {
	return func(h *HTTP) {
		h.cfg.AdminBearerToken = s
	}
}

type Config struct {
	AdminBearerToken string
	Addr             string
	Port             int
}

func (c *Config) SocketAddr() string { return fmt.Sprintf("%s:%d", c.Addr, c.Port) }

type HTTP struct {
	middlewares []fiber.Handler
	app         *fiber.App
	cfg         *Config
}

func New(opts ...HTTPOption) *HTTP {
	noopEngine := &engine.NoopEngineProvider{}
	overlayAPI := &app.Application{
		Commands: &app.Commands{
			SubmitTransactionHandler: commands.NewSubmitTransactionCommandHandler(noopEngine),
			SyncAdvertismentsHandler: commands.NewSyncAdvertismentsHandler(noopEngine),
		},
		Queries: &app.Queries{
			TopicManagerDocumentationHandler: queries.NewTopicManagerDocumentationHandler(noopEngine),
		},
	}

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
	admin.Use(AdminRoutesAuthorizationMiddleware(http.cfg.AdminBearerToken))
	admin.Post("/advertisements-sync", overlayAPI.Commands.SyncAdvertismentsHandler.Handle)

	return &http
}

func (h *HTTP) ListenAndServe() { log.Fatal(h.app.Listen(h.cfg.SocketAddr())) }
