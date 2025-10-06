package service

import (
	"context"
	"errors"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/repo"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/errutils"
)

func (e *Event) DeleteEvent(ctx context.Context, eventID int, userID int) error {
	const op = "EventService.DeleteEvent()"

	if err := e.eventRepo.DeleteEvent(ctx, eventID, userID); err != nil {
		if errors.Is(err, repo.ErrEventNotFound) {
			return errutils.Wrap(op, ErrEventNotFound)
		}
		return errutils.Wrap(op, err)
	}

	return nil
}
