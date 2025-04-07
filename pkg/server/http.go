package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server/internal/app/jsonutil"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
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

// HTTPOption defines a functional option for configuring an HTTP server.
// These options allow for flexible setup of middlewares and configurations.
type HTTPOption func(*HTTP) error

// WithMiddleware adds custom middleware to the HTTP server.
// The execution order depends on the sequence in which the middlewares are passed
func WithMiddleware(f func(http.Handler) http.Handler) HTTPOption {
	return func(h *HTTP) error {
		h.middleware = append(h.middleware, adaptor.HTTPMiddleware(f))
		return nil
	}
}

// WithConfig sets the configuration for the HTTP server.
func WithConfig(cfg *Config) HTTPOption {
	return func(h *HTTP) error {
		h.cfg = cfg
		h.app = fiber.New(fiber.Config{
			CaseSensitive: true,
			StrictRouting: true,
			ServerHeader:  cfg.ServerHeader,
			AppName:       cfg.AppName,
		})

		return nil
	}
}

// HTTP manages connections to the overlay server instance. It accepts and responds to client sockets,
// using idempotency to improve fault tolerance and mitigate duplicated requests.
// It applies all configured options along with the list of middlewares.
type HTTP struct {
	middleware []fiber.Handler
	app        *fiber.App
	cfg        *Config
}

// New returns an instance of the HTTP server and applies all specified functional options before starting it.
func New(opts ...HTTPOption) (*HTTP, error) {
	overlayAPI, err := app.New(NewNoopEngineProvider())
	if err != nil {
		return nil, fmt.Errorf("failed to create overlay API: %w", err)
	}

	http := &HTTP{
		cfg: &DefaultConfig,
		app: fiber.New(fiber.Config{
			CaseSensitive: true,
			StrictRouting: true,
			ServerHeader:  "Overlay API",
			AppName:       "Overlay API v0.0.0",
		}),
		middleware: []fiber.Handler{
			idempotency.New(),
			cors.New(),
		},
	}

	for _, o := range opts {
		if err := o(http); err != nil {
			return nil, err
		}
	}

	for _, m := range http.middleware {
		http.app.Use(m)
	}

	// Routes...
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

	return http, nil
}

// SocketAddr builds the address string for binding.
func (h *HTTP) SocketAddr() string {
	return fmt.Sprintf("%s:%d", h.cfg.Addr, h.cfg.Port)
}

// ListenAndServe handles HTTP requests from the configured socket address.
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
