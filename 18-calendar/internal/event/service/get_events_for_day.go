package service

import (
	"context"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/dto"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/errutils"
	"time"
)

func (e *Event) GetEventsForDay(ctx context.Context, userID int, date time.Time) (dto.GetEventsResponse, error) {
	const op = "EventService.GetEventsForDay()"

	domainEvents, err := e.eventRepo.GetEventsForDay(ctx, userID, date)
	if err != nil {
		return dto.GetEventsResponse{}, errutils.Wrap(op, err)
	}

	return domainToGetEventsResponse(domainEvents), nil
}
