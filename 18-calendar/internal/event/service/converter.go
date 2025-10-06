package service

import (
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/domain"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/dto"
)

func domainToGetEventsResponse(domainEvents []domain.Event) dto.GetEventsResponse {
	events := make([]dto.Event, 0, len(domainEvents))
	for _, e := range domainEvents {
		events = append(events, dto.Event{
			Date:        e.Date,
			Status:      string(e.Status),
			Description: e.Description,
		})
	}

	return dto.GetEventsResponse{
		Events: events,
	}
}
