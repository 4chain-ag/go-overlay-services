package middleware

import (
	"strings"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/gofiber/fiber/v2"
)

// BearerTokenAuthorizationMiddleware is a middleware function that checks if the request
// contains a valid Bearer token in the Authorization header. If the token is invalid or
// missing, it responds with an appropriate error.
func BearerTokenAuthorizationMiddleware(expectedToken string) fiber.Handler {
	const scheme = "Bearer "
	return func(c *fiber.Ctx) error {
		auth := c.Get(fiber.HeaderAuthorization)
		if auth == "" {
			return NewMissingAuthorizationHeaderError()
		}

		if !strings.HasPrefix(auth, scheme) {
			return NewMissingBearerTokenValueError()
		}

		token := strings.TrimPrefix(auth, scheme)
		if token != expectedToken {
			return NewInvalidBearerTokenValueError()
		}

		return c.Next()
	}
}

// NewMissingAuthorizationHeaderError returns an app.Error indicating that the
// Authorization header is missing from the request.
func NewMissingAuthorizationHeaderError() app.Error {
	const str = "Unauthorized access: Missing Authorization header in the request"
	return app.NewAuthorizationError(str, str)
}

// NewMissingBearerTokenValueError returns an app.Error indicating that the
// Bearer token value is missing from the Authorization header.
func NewMissingBearerTokenValueError() app.Error {
	const str = "Unauthorized access: Missing Authorization header Bearer token value"
	return app.NewAuthorizationError(str, str)
}

// NewInvalidBearerTokenValueError returns an app.Error indicating that the
// Bearer token provided is invalid or not recognized.
func NewInvalidBearerTokenValueError() app.Error {
	const str = "Forbidden access: Invalid Bearer token value"
	return app.NewAccessForbiddenError(str, str)
}
