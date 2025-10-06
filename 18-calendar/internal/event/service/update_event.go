package service

import (
	"context"
	"errors"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/repo"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/domain"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/dto"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/errutils"
)

func (e *Event) UpdateEvent(ctx context.Context, event dto.UpdateEventRequest, userID int) error {
	const op = "EventService.UpdateEventRequest()"

	domainEvent := domain.Event{
		ID:          event.ID,
		UserID:      userID,
		Date:        event.Date,
		Status:      domain.EventStatus(event.Status),
		Description: event.Description,
	}

	isValidStatus := domain.IsValidStatus(domainEvent.Status)
	if !isValidStatus {
		return errutils.Wrap(op, ErrInvalidStatus)
	}

	err := e.eventRepo.UpdateEvent(ctx, domainEvent)
	if err != nil {
		if errors.Is(err, repo.ErrEventNotFound) {
			return errutils.Wrap(op, ErrEventNotFound)
		}
		return errutils.Wrap(op, err)
	}

	return nil
}
