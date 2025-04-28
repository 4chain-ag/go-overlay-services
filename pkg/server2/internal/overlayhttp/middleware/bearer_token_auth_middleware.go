package middleware

import (
	"strings"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/openapi"
	"github.com/gofiber/fiber/v2"
)

// BearerTokenAuthorizationMiddleware is a middleware function that checks if the request
// contains a valid Bearer token in the Authorization header. If the token is invalid or
// missing, it responds with an appropriate error.
func BearerTokenAuthorizationMiddleware(expectedToken string, next fiber.Handler) fiber.Handler {
	const scheme = "Bearer "
	return func(c *fiber.Ctx) error {
		// Retrieve the Authorization header from the request
		auth := c.Get(fiber.HeaderAuthorization)

		// Check if the Authorization header is missing
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(MissingAuthorizationHeaderResponse)
		}

		// Check if the Authorization header does not start with 'Bearer '
		if !strings.HasPrefix(auth, scheme) {
			return c.Status(fiber.StatusUnauthorized).JSON(MissingAuthorizationHeaderBearerTokenValueResponse)
		}

		// Extract the token from the Authorization header
		token := strings.TrimPrefix(auth, scheme)

		// Check if the token does not match the expected token
		if token != expectedToken {
			return c.Status(fiber.StatusForbidden).JSON(InvalidBearerTokenValueResponse)
		}

		// Proceed with the next handler if the token is valid
		return next(c)
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
