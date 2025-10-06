package service

import (
	"context"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/domain"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/dto"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/errutils"
)

func (e *Event) CreateEvent(ctx context.Context, event dto.CreateEventRequest, userID int) error {
	const op = "EventService.CreateEvent()"

	domainEvent := domain.Event{
		UserID:      userID,
		Date:        event.Date,
		Status:      domain.StatusPlanned,
		Description: event.Description,
	}
	if err := e.eventRepo.CreateEvent(ctx, domainEvent); err != nil {
		return errutils.Wrap(op, err)
	}

	return nil
}
