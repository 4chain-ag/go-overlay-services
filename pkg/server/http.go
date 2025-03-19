package server

import (
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/engine"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/commands"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/queries"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Config struct {
	AdminBearerToken string
	Addr             string
	Port             int
}

func (c *Config) SocketAddr() string { return fmt.Sprintf("%s:%d", c.Addr, c.Port) }

type HTTP struct {
	app *fiber.App
	cfg *Config
}

func NewHTTP(cfg *Config) *HTTP {
	noopEngine := &engine.NoopEngineProvider{}
	overlayApp := &app.Application{
		Commands: &app.Commands{
			SubmitTransactionHandler: commands.NewSubmitTransactionCommandHandler(noopEngine),
			SyncAdvertismentsHandler: commands.NewSyncAdvertismentsHandler(noopEngine),
		},
		Queries: &app.Queries{
			TopicManagerDocumentationHandler: queries.NewTopicManagerDocumentationHandler(noopEngine),
		},
	}

	log.SetLevel(log.LevelDebug)
	return &HTTP{app: initFiberApp(overlayApp, cfg.AdminBearerToken), cfg: cfg}
}

func (h *HTTP) ListenAndServe() { log.Fatal(h.app.Listen(h.cfg.SocketAddr())) }

func initFiberApp(overlayApp *app.Application, token string) *fiber.App {
	fiberAPI := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Overlay API",
		AppName:       "Overlay API v0.0.0",
	})

	// Middlewares:
	fiberAPI.Use(idempotency.New())
	fiberAPI.Use(requestid.New())
	fiberAPI.Use(logger.New())

	// Routes:
	api := fiberAPI.Group("/api")
	v1 := api.Group("/v1")

	// Non-Admin:
	v1.Post("/submit", overlayApp.Commands.SubmitTransactionHandler.Handle)
	v1.Get("/topic-managers", overlayApp.Queries.TopicManagerDocumentationHandler.Handle)

	// Admin:
	admin := v1.Group("/admin")
	admin.Use(AdminRoutesAuthorizationMiddleware(token))
	admin.Post("/advertisements-sync", overlayApp.Commands.SyncAdvertismentsHandler.Handle)

	return fiberAPI
}
