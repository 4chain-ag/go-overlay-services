package middleware

import (
	"errors"
	"fmt"
	"time"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

// ErrorResponseMiddlewareConfig defines the configuration for mapping specific application error types
// to their corresponding OpenAPI-compliant error responses.
//
// This configuration allows the application to declaratively associate known error types
// with structured and consistent HTTP responses.
type ErrorResponseMiddlewareConfig struct {
	ErrorTypeIncorrectInput  openapi.Error // Response for invalid or malformed client input (HTTP 400)
	ErrorTypeProviderFailure openapi.Error // Response for downstream service or dependency failures (HTTP 500)
	RequestTimeoutResponse   time.Duration // Timeout duration used to construct a 408 Request Timeout response
}

// ErrorResponseCodesMapping maps custom application error types to corresponding HTTP status codes.
// This mapping is used to convert internal error representations into appropriate HTTP responses.
//
// Mapping:
//   - ErrorTypeAuthorization    → 401 Unauthorized
//   - ErrorTypeIncorrectInput   → 400 Bad Request
//   - ErrorTypeOperationTimeout → 408 Request Timeout
//   - ErrorTypeProviderFailure  → 500 Internal Server Error
var ErrorResponseCodesMapping = map[app.ErrorType]int{
	app.ErrorTypeAuthorization:    fiber.StatusUnauthorized,
	app.ErrorTypeIncorrectInput:   fiber.StatusBadRequest,
	app.ErrorTypeOperationTimeout: fiber.StatusRequestTimeout,
	app.ErrorTypeProviderFailure:  fiber.StatusInternalServerError,
}

// ErrorResponseMiddleware returns a Fiber middleware that intercepts application-specific errors
// and maps them to appropriate OpenAPI error responses based on their error type.
// If the error matches a known application error type and is non-zero,
// the middleware responds with the corresponding structured error and appropriate HTTP status code.

// Supported error types:
//   - ErrorTypeIncorrectInput:      returns HTTP 400 Bad Request
//   - ErrorTypeOperationTimeout:    returns HTTP 408 Request Timeout
//   - ErrorTypeProviderFailure:     returns HTTP 500 Internal Server Error
//   - ErrorTypeUnknown (or others): returns a generic 500 Internal Server Error (`UnhandledErrorTypeResponse`)
//
// If the error is not recognized as an `app.Error`, a fallback response is returned.
func ErrorResponseMiddleware(cfg ErrorResponseMiddlewareConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			var target app.Error
			if !errors.As(err, &target) || target.IsZero() {
				return c.Status(fiber.StatusInternalServerError).JSON(UnhandledErrorTypeResponse)
			}

			code := ErrorResponseCodesMapping[target.ErrorType()]
			switch target.ErrorType() {
			case app.ErrorTypeIncorrectInput:
				return c.Status(code).JSON(cfg.ErrorTypeIncorrectInput)

			case app.ErrorTypeOperationTimeout:
				return c.Status(code).JSON(NewRequestTimeoutResponse(cfg.RequestTimeoutResponse))

			case app.ErrorTypeProviderFailure:
				return c.Status(code).JSON(cfg.ErrorTypeProviderFailure)

			case app.ErrorTypeUnknown:
				return c.Status(code).JSON(UnhandledErrorTypeResponse)
			}
		}

		return nil
	}
}

// UnhandledErrorTypeResponse is the default response returned when an error occurs
// that does not match any known or handled ErrorType.
// It represents a generic internal server error to avoid exposing internal details to the client.
var UnhandledErrorTypeResponse = openapi.InternalServerErrorResponse{
	Message: "An internal error occurred during processing the request. Please try again later or contact the support team.",
}

// NewRequestTimeoutResponse creates a timeout response when a request exceeds the allowed time limit.
// It takes the timeout duration and returns an openapi.RequestTimeoutResponse with a message indicating that the request timed out.
func NewRequestTimeoutResponse(limit time.Duration) openapi.RequestTimeoutResponse {
	return openapi.BadRequestResponse{
		Message: fmt.Sprintf("The submitted request exceeded the timeout limit of %d seconds.", int(limit.Seconds())),
	}
}
