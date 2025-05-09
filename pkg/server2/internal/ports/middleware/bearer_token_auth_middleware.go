package middleware

import (
	"strings"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
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
			return c.Status(fiber.StatusUnauthorized).JSON(MissingAuthorizationHeaderResponse)
		}

		if !strings.HasPrefix(auth, scheme) {
			return c.Status(fiber.StatusUnauthorized).JSON(MissingAuthorizationHeaderBearerTokenValueResponse)
		}

		token := strings.TrimPrefix(auth, scheme)
		if token != expectedToken {
			return c.Status(fiber.StatusForbidden).JSON(InvalidBearerTokenValueResponse)
		}

		return c.Next()
	}
}

// MissingAuthorizationHeaderResponse represents a bad request response when the Authorization header is missing from the request.
var MissingAuthorizationHeaderResponse = openapi.BadRequestResponse{
	Message: "Unauthorized: Missing Authorization header in the request",
}

// MissingAuthorizationHeaderBearerTokenValueResponse represents a bad request response when the Authorization header is present
// but the Bearer token value is missing.
var MissingAuthorizationHeaderBearerTokenValueResponse = openapi.BadRequestResponse{
	Message: "Unauthorized: Missing Authorization header Bearer token value",
}

// InvalidBearerTokenValueResponse represents a bad request response when the Bearer token value provided is invalid.
var InvalidBearerTokenValueResponse = openapi.BadRequestResponse{
	Message: "Forbidden: Invalid Bearer token value",
}
