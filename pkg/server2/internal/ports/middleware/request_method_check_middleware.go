package middleware

import (
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// routeMethods maps request paths to their expected HTTP methods.
type routeMethods map[string]string

// exists checks whether a route path is registered.
func (r routeMethods) exists(s string) bool {
	_, ok := r[s]
	return ok
}

// method returns the expected HTTP method for the given path.
func (r routeMethods) method(s string) string {
	return r[s]
}

// details returns a copy of the route-method mappings as a map[string]any,
// useful for error details in response payloads.
func (r routeMethods) details() map[string]any {
	details := make(map[string]any)
	for path, method := range r {
		details[path] = method
	}
	return details
}

// newRouteMethods returns the list of all valid routes and their expected HTTP methods.
func newRouteMethods() routeMethods {
	return routeMethods{
		"/api/v1/submit":                   fiber.MethodPost,
		"/api/v1/admin/syncAdvertisements": fiber.MethodPost,
	}
}

// RequestMethodCheckMiddleware creates a middleware that validates incoming requests against
// known route paths and their expected HTTP methods. If the path is unrecognized, it returns
// a 404 Not Found with details. If the method is invalid, it returns a 405 Method Not Allowed
// with a mismatch description. Otherwise, it passes the request to the next handler.
func RequestMethodCheckMiddleware() fiber.Handler {
	routes := newRouteMethods()

	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		actual := string(ctx.Method())
		path := string(ctx.Path())
		if !routes.exists(path) {
			return c.Status(fiber.StatusNotFound).JSON(NewRequestPathNotRecognized(routes.details()))
		}

		expected := routes.method(path)
		if actual != expected {
			return c.Status(fiber.StatusMethodNotAllowed).JSON(NewRequestInvalidMethodResponse(actual, expected))
		}
		return c.Next()
	}
}

// NewRequestInvalidMethodResponse creates a bad request response when the submitted request uses an invalid HTTP method.
// It takes the actual HTTP method used and the expected HTTP method, and returns an openapi.BadRequestResponse
// with a message indicating the mismatch.
func NewRequestInvalidMethodResponse(actual string, expected string) openapi.BadRequestResponse {
	return openapi.BadRequestResponse{
		Message: fmt.Sprintf("The submitted request does not correspond to expected method type: %s, HTTP request method used: %s", expected, actual),
	}
}

// NewRequestPathNotRecognized creates a bad request response when the request path is not recognized.
// It includes the list of supported paths and their corresponding methods in the response details.
func NewRequestPathNotRecognized(details map[string]any) openapi.BadRequestResponse {
	return openapi.BadRequestResponse{
		Details: &details,
		Message: "The submitted request does not correspond to any of supported paths.",
	}
}
