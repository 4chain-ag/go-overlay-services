package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// BasicMiddlewareGroup returns a slice of standard Fiber middleware handlers
// that provide core HTTP functionality such as request ID assignment,
// idempotency support, CORS handling, panic recovery, structured logging,
// and optional pprof profiling under /api/v1.
//
// Note: Stack traces are enabled in recover middleware for debugging purposes
// and should be disabled in production builds.
func BasicMiddlewareGroup() []fiber.Handler {
	return []fiber.Handler{
		requestid.New(),
		idempotency.New(),
		cors.New(),
		recover.New(recover.Config{EnableStackTrace: true}),
		logger.New(logger.Config{
			Format:     "date=${time} request_id=${locals:requestid} status=${status} method=${method} path=${path} err=${error}​\n",
			TimeFormat: "02-Jan-2006 15:04:05",
		}),
		pprof.New(pprof.Config{Prefix: "/api/v1"}),
	}
}
