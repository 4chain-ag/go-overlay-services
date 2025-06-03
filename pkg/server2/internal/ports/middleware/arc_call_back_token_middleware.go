package middleware

import (
	"slices"
	"strings"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// ArcCallbackTokenMiddleware returns a fiber.Handler that validates the
// Bearer token present in Authorization header for ARC callback requests.
// It protects the ARC ingest endpoint from unauthorized access.
// It also conditionally checks if the request is authorized based on OpenAPI security scopes.
func ArcCallbackTokenMiddleware(arcCallbackToken string, arcApiKey string) fiber.Handler {
	const scheme = "Bearer "
	const arcCallbackScope = "arcCallBackToken"

	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		scopes, ok := ctx.UserValue(openapi.BearerAuthScopes).([]string)
		if !ok {
			return NewArcBearerAuthScopesAssertionError()
		}
		if len(scopes) == 0 {
			return NewArcEmptyAccessScopesAssertionError()
		}
		if !slices.Contains(scopes, arcCallbackScope) {
			return nil
		}

		if arcApiKey == "" {
			return NewArcEndpointNotSupportedError()
		}

		if arcCallbackToken == "" {
			return nil
		}

		auth := c.Get(fiber.HeaderAuthorization)
		if auth == "" {
			return NewArcMissingAuthHeaderError()
		}

		if !strings.HasPrefix(auth, scheme) || len(auth) <= len(scheme) {
			return NewArcMissingBearerTokenError()
		}

		token := strings.TrimPrefix(auth, scheme)
		if token != arcCallbackToken {
			return NewArcInvalidBearerTokenError()
		}

		return nil
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

// NewArcEndpointNotSupportedError returns an app.Error indicating that the ARC endpoint
// is not supported by the current service configuration.
func NewArcEndpointNotSupportedError() app.Error {
	const str = "This endpoint is not supported by the current service configuration"
	return app.NewIncorrectInputError(str, str)
}

// NewArcBearerAuthScopesAssertionError returns an app.Error indicating that the
// authorization scopes assertion failed for ARC requests.
func NewArcBearerAuthScopesAssertionError() app.Error {
	return app.NewAuthorizationError(
		"Authorization scope assertion failure for ARC endpoint",
		"Unable to process request to the ARC endpoint. Please verify the request content and try again later.")
}

// NewArcEmptyAccessScopesAssertionError returns an app.Error indicating that the
// authorization scope list exists but is empty for ARC requests.
func NewArcEmptyAccessScopesAssertionError() app.Error {
	return app.NewAuthorizationError(
		"Authorization scope assertion failure: empty scopes for ARC endpoint",
		"Unable to process request to the ARC endpoint. Please verify the request content and try again later.")
}
