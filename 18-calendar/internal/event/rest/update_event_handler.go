package rest

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/service"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/dto"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/response"
)

func (h *EventHandler) UpdateEventHandler(c *fiber.Ctx) error {
	const op = "EventHandler.UpdateEventRequest()"

	event := dto.UpdateEventRequest{}
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

	if err := h.event.UpdateEvent(c.Context(), event, userID); err != nil {
		if errors.Is(err, service.ErrInvalidStatus) {
			l.Warn("invalid event status", "status", event.Status)
			return c.Status(fiber.StatusBadRequest).JSON(
				response.Error(msgInvalidStatus))
		}

		if errors.Is(err, service.ErrEventNotFound) {
			l.Warn("event not found", "event_id", event.ID)
			return c.Status(fiber.StatusNotFound).JSON(
				response.Error(msgEventNotFound))
		}

		l.Error("internal error during event updating", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(msgInternalError))
	}

	return c.Status(fiber.StatusOK).JSON(
		response.Success("event updated successfully"))
}
