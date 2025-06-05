package ports

import (
	"strings"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/gofiber/fiber/v2"
)

const scheme = "Bearer "

// HandleArcIngest method decorates the next ArcIngestService with authorization checks.
// ArcIngestHandlerConfig defines the configuration for the ArcIngestHandler.
type ArcIngestHandlerConfig struct {
	ArcApiKey        string
	ArcCallbackToken string
}

// ArcAuthorizationDecorator is a decorator that adds authorization checks to the ArcIngestHandler.
type ArcAuthorizationDecorator struct {
	next Handler
	cfg  ArcIngestHandlerConfig
}

// Handle method decorates the next ArcIngestService with authorization checks.
func (a *ArcAuthorizationDecorator) Handle(c *fiber.Ctx) error {
	auth := c.Get(fiber.HeaderAuthorization)
	if auth == "" {
		return NewArcMissingAuthHeaderError()
	}

	if !strings.HasPrefix(auth, scheme) || len(auth) <= len(scheme) {
		return NewArcMissingBearerTokenError()
	}

	token := strings.TrimPrefix(auth, scheme)
	if token != a.cfg.ArcCallbackToken {
		return NewArcInvalidBearerTokenError()
	}

	return a.next.Handle(c)
}

// NewArcAuthorizationDecorator creates a new ArcAuthorizationDecorator.
func NewArcAuthorizationDecorator(next Handler, cfg ArcIngestHandlerConfig) *ArcAuthorizationDecorator {
	if cfg.ArcApiKey == "" {
		panic("ArcApiKey is required")
	}

	if cfg.ArcCallbackToken == "" {
		panic("ArcCallbackToken is required")
	}

	return &ArcAuthorizationDecorator{
		next: next,
		cfg:  cfg,
	}
}

// NewArcMissingAuthHeaderError returns an app.Error indicating that the
// Authorization header is missing from the ARC callback request.
func NewArcMissingAuthHeaderError() app.Error {
	const str = "Authorization header is missing from the request"
	return app.NewAuthorizationError(str, str)
}

// NewArcMissingBearerTokenError returns an app.Error indicating that the
// Bearer token value is missing from the Authorization header.
func NewArcMissingBearerTokenError() app.Error {
	const str = "Authorization header is present, but the Bearer token is missing"
	return app.NewAuthorizationError(str, str)
}

// NewArcInvalidBearerTokenError returns an app.Error indicating that the
// Bearer token provided is invalid or not recognized.
func NewArcInvalidBearerTokenError() app.Error {
	const str = "The Bearer token provided is invalid"
	return app.NewAccessForbiddenError(str, str)
}
