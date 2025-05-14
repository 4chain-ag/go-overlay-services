package ports

import (
	"fmt"
	"slices"
	"strings"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// BearerTokenAuthorizationHandler returns a fiber.Handler that validates the
// Bearer token present in Authorization header of incoming HTTP requests.
// It also conditionally check if the requests is authorized based on OpenAPI
// security scopes.
func BearerTokenAuthorizationHandler(expectedToken string) fiber.Handler {
	const scheme = "Bearer "
	const userScope = "user"

	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		scopes, ok := ctx.UserValue(openapi.BearerAuthScopes).([]string)
		if !ok {
			return NewAccessScopeAssertionError()
		}
		if len(scopes) == 0 || slices.Contains(scopes, userScope) {
			return nil
		}

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

		return nil
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

// NewAccessScopeAssertionError returns an app.Error indicating that the
// authorization scopes assertion failed, usually due to missing or
// improperly formated OpenAPI scope data in the request context.
func NewAccessScopeAssertionError() app.Error {
	return app.NewUnknownError(
		fmt.Sprintf("Authorization scope assertion failure: expected to get string slice under %s user context key to properly extract the request scope.", openapi.BearerAuthScopes),
		"Unable to process request to the endpoint. Please verify the request content and try again later.")
}
