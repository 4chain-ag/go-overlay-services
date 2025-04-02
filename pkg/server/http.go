package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/4chain-ag/go-overlay-services/pkg/config"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server/app/jsonutil"
	"github.com/4chain-ag/go-overlay-services/pkg/server/mongo"
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
func WithConfig(cfg *config.Config) HTTPOption {
	return func(h *HTTP) {
		h.cfg = cfg
	}
}

// WithMongo sets the MongoDB client for the HTTP server based on the configuration.
func WithMongo() HTTPOption {
	return func(h *HTTP) {
		if h.cfg == nil || h.cfg.Mongo.URI == "" {
			return
		}
		client, err := mongo.Connect(h.cfg)
		if err != nil {
			panic(fmt.Sprintf("MongoDB connect failed: %v", err))
		}
		h.mongo = client
	}
}

// SocketAddr returns the socket address string based on the configured address and port combination.
func (h *HTTP) SocketAddr() string { return fmt.Sprintf("%s:%d", h.cfg.Addr, h.cfg.Port) }

// HTTP manages connections to the overlay server instance. It accepts and responds to client sockets,
// using idempotency to improve fault tolerance and mitigate duplicated requests.
// It applies all configured options along with the list of middlewares."
type HTTP struct {
	middlewares []fiber.Handler
	app         *fiber.App
	cfg         *config.Config
	mongo       *mongo.Client
}

// New returns an instance of the HTTP server and applies all specified functional options before starting it.
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
	v1.Post("/request-foreign-gasp-node", adaptor.HTTPHandlerFunc(overlayAPI.Commands.RequestForeignGASPNodeHandler.Handle))

	// Admin:
	admin := v1.Group("/admin", adaptor.HTTPMiddleware(AdminAuth(http.cfg.AdminBearerToken)))
	admin.Post("/advertisements-sync", adaptor.HTTPHandlerFunc(overlayAPI.Commands.SyncAdvertismentsHandler.Handle))
	admin.Post("/start-gasp-sync", adaptor.HTTPHandlerFunc(overlayAPI.Commands.StartGASPSyncHandler.Handle))

	return &http
}

// ListenAndServe handles HTTP requests from the configured socket address."
func (h *HTTP) ListenAndServe() error {
	if err := h.app.Listen(h.SocketAddr()); err != nil {
		return fmt.Errorf("http server: fiber app listen failed: %w", err)
	}
	return nil
}

// AdminAuth is a middleware that checks the Authorization header for a valid Bearer token.
// protects the HTTP server from unauthorized access.
// It checks for a Bearer token in the Authorization header and compares it to the expected value.
func AdminAuth(expectedToken string) func(http.Handler) http.Handler {
	type AuthorizationFailureResponse struct {
		Status  string `json:"error"`
		Message string `json:"message"`
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				jsonutil.SendHTTPResponse(w, http.StatusUnauthorized, AuthorizationFailureResponse{
					Status:  http.StatusText(http.StatusUnauthorized),
					Message: "Missing Authorization header in the request",
				})
				return
			}

			if !strings.HasPrefix(auth, "Bearer ") {
				jsonutil.SendHTTPResponse(w, http.StatusUnauthorized, AuthorizationFailureResponse{
					Status:  http.StatusText(http.StatusUnauthorized),
					Message: "Missing Authorization header Bearer token value",
				})
				return
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			if token != expectedToken {
				jsonutil.SendHTTPResponse(w, http.StatusForbidden, AuthorizationFailureResponse{
					Status:  http.StatusText(http.StatusForbidden),
					Message: "Invalid Bearer token value",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
