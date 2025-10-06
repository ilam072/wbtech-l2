package service

import (
	"context"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/dto"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/errutils"
	"time"
)

func (e *Event) GetEventsForWeek(ctx context.Context, userID int, start time.Time) (dto.GetEventsResponse, error) {
	const op = "EventService.GetEventsForWeek()"

	domainEvents, err := e.eventRepo.GetEventsForWeek(ctx, userID, start)
	if err != nil {
		return dto.GetEventsResponse{}, errutils.Wrap(op, err)
	}

	return domainToGetEventsResponse(domainEvents), nil
}
