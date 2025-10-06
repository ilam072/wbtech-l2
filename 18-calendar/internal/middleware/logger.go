package middleware

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"time"
)

func Logger(l *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		l.Info(
			"request completed",
			"method", c.Method(),
			"url", c.OriginalURL(),
			"status", c.Response().StatusCode(),
			"duration", time.Since(start).String(),
		)

		return err
	}
}
