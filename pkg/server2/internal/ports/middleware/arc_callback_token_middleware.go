package middleware

import (
	"strings"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

var (
	// MissingArcAuthHeaderResponse is returned when the Authorization header
	// is completely missing from the request.
	MissingArcAuthHeaderResponse = openapi.BadRequestResponse{
		Message: "Authorization header is missing from the request.",
	}

	// MissingArcAuthHeaderValueResponse is returned when the Authorization header
	// is present but doesn't contain a proper Bearer token.
	MissingArcAuthHeaderValueResponse = openapi.BadRequestResponse{
		Message: "Authorization header is present, but the Bearer token is missing.",
	}

	// InvalidArcBearerTokenValueResponse is returned when the provided Bearer token
	// doesn't match the expected token value.
	InvalidArcBearerTokenValueResponse = openapi.BadRequestResponse{
		Message: "The Bearer token provided is invalid.",
	}

	// EndpointNotSupportedResponse is returned when the endpoint is accessed but
	// is not configured in the current service (arcApiKey is empty).
	EndpointNotSupportedResponse = openapi.BadRequestResponse{
		Message: "This endpoint is not supported by the current service configuration.",
	}
)

// unsupportedArcEndpointMiddleware returns a middleware function that responds with a 404 Not Found
// for ARC-related endpoints that are not supported by the current service configuration.
func unsupportedArcEndpointMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(EndpointNotSupportedResponse)
	}
}

// ArcCallbackTokenMiddleware is a middleware that checks the Authorization header for a valid Bearer token.
// It protects the ARC ingest endpoint from unauthorized access.
// It checks for a Bearer token in the Authorization header and compares it to the expected token value.
// The endpoint will return 404 Not Found if arcApiKey is empty, indicating ARC integration is disabled.
func ArcCallbackTokenMiddleware(arcCallbackToken string, arcApiKey string) fiber.Handler {
	const schema = "Bearer "

	if arcApiKey == "" {
		return unsupportedArcEndpointMiddleware()
	}

	return func(c *fiber.Ctx) error {
		if arcCallbackToken == "" {
			return c.Next()
		}

		auth := c.Get(fiber.HeaderAuthorization)
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(MissingArcAuthHeaderResponse)
		}

		if !strings.HasPrefix(auth, schema) || len(auth) <= len(schema) {
			return c.Status(fiber.StatusUnauthorized).JSON(MissingArcAuthHeaderValueResponse)
		}

		token := strings.TrimPrefix(auth, schema)
		if token != arcCallbackToken {
			return c.Status(fiber.StatusForbidden).JSON(InvalidArcBearerTokenValueResponse)
		}

		return c.Next()
	}
}
