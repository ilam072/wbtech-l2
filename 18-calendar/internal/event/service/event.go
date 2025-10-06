package service

import (
	"context"
	"errors"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/domain"
	"time"
)

//go:generate mockgen -source=event.go -destination=../mocks/service_mocks.go -package=mocks
type EventRepo interface {
	CreateEvent(ctx context.Context, event domain.Event) error
	UpdateEvent(ctx context.Context, event domain.Event) error
	DeleteEvent(ctx context.Context, eventID int, userID int) error

	GetEventsForDay(ctx context.Context, userID int, date time.Time) ([]domain.Event, error)
	GetEventsForWeek(ctx context.Context, userID int, start time.Time) ([]domain.Event, error)
	GetEventsForMonth(ctx context.Context, userID int, start time.Time) ([]domain.Event, error)
}

var (
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidStatus = errors.New("invalid event status")
)

type Event struct {
	eventRepo EventRepo
}

func NewEvent(repo EventRepo) *Event {
	return &Event{eventRepo: repo}
}
