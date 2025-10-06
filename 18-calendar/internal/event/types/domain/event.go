package domain

import "time"

type EventStatus string

const (
	StatusPlanned  EventStatus = "planned"
	StatusDone     EventStatus = "done"
	StatusCanceled EventStatus = "canceled"
)

type Event struct {
	ID          int         `db:"id"`
	UserID      int         `db:"user_id"`
	Date        time.Time   `db:"date"`
	Status      EventStatus `db:"status"`
	Description string      `db:"description"`
}

func IsValidStatus(s EventStatus) bool {
	switch s {
	case StatusPlanned, StatusCanceled, StatusDone:
		return true
	default:
		return false
	}
}
