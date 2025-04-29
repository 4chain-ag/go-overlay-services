package middleware

import (
	"errors"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/ports/openapi"
	"github.com/gofiber/fiber/v2"
)

func ErrorMiddleware(c *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if !errors.As(err, &fiberErr) {
		return c.Status(fiber.StatusInternalServerError).JSON(openapi.Error{Message: "uncategorized error"})
	}

	switch fiberErr.Code {
	case fiber.StatusBadRequest:
		return c.Status(fiber.StatusBadRequest).JSON(openapi.Error{Message: fiberErr.Message})
	case fiber.StatusInternalServerError:
		return c.Status(fiber.StatusInternalServerError).JSON(openapi.Error{Message: fiberErr.Message})
	default:
		return c.Status(fiberErr.Code).JSON(openapi.Error{Message: fiberErr.Message})
	}
}
