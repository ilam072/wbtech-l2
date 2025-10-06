package rest

import (
	"context"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/dto"
	"log/slog"
	"time"
)

const (
	msgInternalError      = "something went wrong, try again later"
	msgInvalidRequestBody = "invalid request body"
	msgInvalidStatus      = "invalid event status"
	msgEventNotFound      = "event not found"
	msgInvalidEventID     = "event id must be a number"
	msgInvalidDateFormat  = "invalid date format, must be YYYY-MM-DD"
)

// api/update_event POST
// api/create_event/ POST
// api/delete_event/:id DELETE
// api/events_for_day?date=2025-08-30 GET
// api/events_for_week?date=2025-08-30 GET
// api/events_for_month?date=2025-08-30 GET

//go:generate mockgen -source=handler.go -destination=../../event/mocks/event_handler_mock.go -package=mocks
type Event interface {
	CreateEvent(ctx context.Context, event dto.CreateEventRequest, userID int) error
	UpdateEvent(ctx context.Context, event dto.UpdateEventRequest, userID int) error
	DeleteEvent(ctx context.Context, eventID int, userID int) error
	GetEventsForDay(ctx context.Context, userID int, date time.Time) (dto.GetEventsResponse, error)
	GetEventsForWeek(ctx context.Context, userID int, start time.Time) (dto.GetEventsResponse, error)
	GetEventsForMonth(ctx context.Context, userID int, start time.Time) (dto.GetEventsResponse, error)
}

type Validator interface {
	Validate(i interface{}) error
	FormatValidationErrors(err error) map[string]string
}

type EventHandler struct {
	log       *slog.Logger
	event     Event
	validator Validator
}

func NewEventHandler(log *slog.Logger, event Event, validator Validator) *EventHandler {
	h := &EventHandler{
		log:       log,
		event:     event,
		validator: validator,
	}

	return h
}
