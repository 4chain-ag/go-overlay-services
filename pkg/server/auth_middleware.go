package server

import "github.com/gofiber/fiber/v2"

func AdminRoutesAuthorizationMiddleware(token string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(MissingAuthorizationHeaderResponse)
		}

		// TODO: Add token matching.
		return c.Next()
	}
}

type AuthorizationMiddlewareResponse struct {
	Message string
}

var MissingAuthorizationHeaderResponse = AuthorizationMiddlewareResponse{
	Message: "Authorization header is missing",
}
