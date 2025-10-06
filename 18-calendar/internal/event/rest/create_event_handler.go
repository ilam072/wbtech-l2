package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/dto"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/response"
)

func (h *EventHandler) CreateEventHandler(c *fiber.Ctx) error {
	const op = "EventHandler.CreateEventRequest()"

	event := dto.CreateEventRequest{}
	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.Error(msgInvalidRequestBody))
	}

	if err := h.validator.Validate(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.Error(h.validator.FormatValidationErrors(err)))
	}

	userID := c.Locals("userID").(int)

	l := h.log.With("user_id", userID, "op", op)

	if err := h.event.CreateEvent(c.Context(), event, userID); err != nil {
		l.Error("internal error during event creation", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(msgInternalError))
	}

	return c.SendStatus(fiber.StatusCreated)
}
