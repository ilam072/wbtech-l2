package rest

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/response"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/service"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/types/dto"
)

func (h *UserHandler) SignUp(c *fiber.Ctx) error {
	const op = "UserHandler.SignUp()"

	user := dto.RegisterUser{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.Error(msgInvalidRequestBody))
	}

	if err := h.validator.Validate(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.Error(h.validator.FormatValidationErrors(err)))
	}

	l := h.log.With("username", user.Username, "op", op)

	if err := h.user.RegisterUser(c.Context(), user); err != nil {
		if errors.Is(err, service.ErrUserExists) {
			l.Warn("attempt to register existing user")
			return c.Status(fiber.StatusConflict).JSON(
				response.Error(msgUsernameOccupied))
		}

		l.Error("internal error during user registration",
			"err", err,
		)
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(msgInternalError))
	}

	return c.SendStatus(fiber.StatusCreated)
}
