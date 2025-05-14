package server2

import (
	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/middleware"
	"github.com/gofiber/fiber/v2"
)

// ExampleRoutes demonstrates how to register custom HTTP routes
// TODO: Remove this example and add in unit tests
func ExampleRoutes() {
	server := New(WithConfig(DefaultConfig))

	server.RegisterRoute("GET", "/api/v1/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Hello, World!"})
	})

	server.RegisterRoute("GET", "/api/v1/users/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"userId": c.Params("id")})
	})

	server.RegisterRoute("POST", "/api/v1/data", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"received": true, "bytes": len(c.Body())})
	})

	server.RegisterRoute("GET", "/api/v1/secure",
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"secure": "data"})
		},
		func(c *fiber.Ctx) error {
			if c.Get("X-API-Key") == "" {
				return c.Status(401).JSON(fiber.Map{"error": "API key required"})
			}
			return c.Next()
		},
	)

	// Route with built-in middleware
	server.RegisterRoute("POST", "/api/v1/upload",
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"uploaded": true})
		},
		middleware.LimitOctetStreamBodyMiddleware(1024*1024), // 1MB limit
	)
}
