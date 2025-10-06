package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/response"
	"time"
)

func (h *EventHandler) GetEventsForWeekHandler(c *fiber.Ctx) error {
	const op = "EventHandler.GetEventsForWeekHandler()"

	date, err := time.Parse(time.DateOnly, c.Query("date"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.Error(msgInvalidDateFormat))
	}

	userID := c.Locals("userID").(int)

	l := h.log.With("user_id", userID, "op", op)

	events, err := h.event.GetEventsForWeek(c.Context(), userID, date)
	if err != nil {
		l.Error("internal error during fetching events for week", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(msgInternalError))
	}

	return c.Status(fiber.StatusOK).JSON(events)
}
