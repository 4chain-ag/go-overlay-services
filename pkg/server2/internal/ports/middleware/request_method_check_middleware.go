package middleware

import (
	"fmt"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
)

// RequestMethodCheckMiddleware creates a middleware that validates incoming requests against
// known route paths and their expected HTTP methods. If the path is unrecognized, it returns
// a 404 Not Found with details. If the method is invalid, it returns a 405 Method Not Allowed
// with a mismatch description. Otherwise, it passes the request to the next handler.
// func RequestMethodCheckMiddleware() fiber.Handler {

// 	return func(c *fiber.Ctx) error {
// 		ctx := c.Context()
// 		actual := string(ctx.Method())
// 		path := string(ctx.Path())
// 		if !routes.exists(path) {
// 			return c.Status(fiber.StatusNotFound).JSON(NewRequestPathNotRecognized()) // to change
// 		}

// 		expected := routes.method(path)
// 		if actual != expected {
// 			return c.Status(fiber.StatusMethodNotAllowed).JSON(NewRequestInvalidMethodResponse(actual, expected))
// 		}
// 		return c.Next()
// 	}
// }

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
func NewRequestPathNotRecognized() openapi.BadRequestResponse {
	return openapi.BadRequestResponse{
		Message: "The submitted request does not correspond to any of supported paths.",
	}
}
