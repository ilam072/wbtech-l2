package dto

import "time"

type CreateEventRequest struct {
	Date        time.Time `json:"date" validate:"required"`
	Description string    `json:"description" validate:"required"`
}

type UpdateEventRequest struct {
	ID          int       `json:"id" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
	Status      string    `json:"status" validate:"required"`
	Description string    `json:"description" validate:"required"`
}

type Event struct {
	Date        time.Time `json:"date"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
}

type GetEventsResponse struct {
	Events []Event
}
