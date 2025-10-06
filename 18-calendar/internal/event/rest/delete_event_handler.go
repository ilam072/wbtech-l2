package rest

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/service"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/response"
	"strconv"
)

func (h *EventHandler) DeleteEventHandler(c *fiber.Ctx) error {
	const op = "EventHandler.DeleteEvent()"

	eventID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.Error(msgInvalidEventID))
	}

	userID := c.Locals("userID").(int)

	l := h.log.With("user_id", userID, "op", op)

	if err := h.event.DeleteEvent(c.Context(), eventID, userID); err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			l.Warn("event not found", "event_id", eventID)
			return c.Status(fiber.StatusNotFound).JSON(
				response.Error(msgEventNotFound))
		}

		l.Error("internal error during event deletion", "err", err, "event_id", eventID)
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(msgInternalError))
	}

	return c.SendStatus(fiber.StatusNoContent)
}
