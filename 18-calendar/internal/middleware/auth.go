package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/response"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/jwt"
	"strings"
)

const bearerPrefix = "Bearer "

func Auth(manager *jwt.Manager) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, bearerPrefix) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(
				response.Error("unauthorized"))

		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)
		userID, err := manager.ParseToken(token)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(
				response.Error("unauthorized"))
		}

		ctx.Locals("userID", userID)

		return ctx.Next()
	}
}
