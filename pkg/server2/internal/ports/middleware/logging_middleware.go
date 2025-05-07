package middleware

import (
	"errors"
	"log/slog"

	"github.com/4chain-ag/go-overlay-services/pkg/server2/internal/app"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type LoggingMiddlewareConfig struct {
	Component string
	Logger    *slog.Logger
}

func LoggingMiddleware(cfg LoggingMiddlewareConfig) fiber.Handler {
	logger := cfg.Logger.With(slog.String("component", cfg.Component))
	return func(c *fiber.Ctx) error {
		requestID := c.Locals(requestid.ConfigDefault.ContextKey).(string)
		base := logger.With(
			slog.String("path", c.Route().Path),
			slog.String("method", c.Method()),
			slog.String("source_ip", c.IP()),
			slog.String("request_id", requestID),
		)

		err := c.Next()
		if err == nil {
			base.With(slog.String("status", "OK")).Info("log-line")
			return nil
		}

		var target app.Error
		if !errors.As(err, &target) || target.IsZero() {
			return err
		}

		base.With(
			slog.String("err", target.Error()),
			slog.String("service", target.Service()),
			slog.String("type", target.ErrorType().String()),
			slog.Int("status_code", ErrorResponseCodesMapping[target.ErrorType()]),
		).Error("log-line")

		return err
	}
}
