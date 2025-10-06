package rest

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/response"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/service"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/types/dto"
	"log/slog"
)

func (h *UserHandler) SignIn(c *fiber.Ctx) error {
	const op = "UserHandler.SignIn()"

	user := dto.LoginUser{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.Error(msgInvalidRequestBody))
	}

	if err := h.validator.Validate(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.Error(h.validator.FormatValidationErrors(err)))
	}

	l := h.log.With("username", user.Username, "op", op)

	token, err := h.user.Login(c.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			l.Warn("invalid login credentials", slog.String("err", err.Error()))
			return c.Status(fiber.StatusUnauthorized).JSON(
				response.Error(msgInvalidCredential))
		}
		l.Error("internal error during login", slog.String("err", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(msgInternalError))
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"status": "OK",
			"token":  token,
		})
}
