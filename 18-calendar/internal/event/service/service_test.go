package service_test

import (
	"context"
	"errors"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/mocks"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/repo"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/service"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/dto"
)

func TestEventService_CreateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepo(ctrl)
	svc := service.NewEvent(mockRepo)

	ctx := context.Background()
	userID := 42
	req := dto.CreateEventRequest{
		Date:        time.Now(),
		Description: "test event",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().CreateEvent(ctx, gomock.Any()).Return(nil)

		err := svc.CreateEvent(ctx, req, userID)
		assert.NoError(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().CreateEvent(ctx, gomock.Any()).Return(errors.New("db error"))

		err := svc.CreateEvent(ctx, req, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EventService.CreateEvent()")
	})
}

func TestEventService_UpdateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepo(ctrl)
	svc := service.NewEvent(mockRepo)

	ctx := context.Background()
	userID := 42

	req := dto.UpdateEventRequest{
		ID:          1,
		Date:        time.Now(),
		Status:      "planned",
		Description: "update test",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().UpdateEvent(ctx, gomock.Any()).Return(nil)

		err := svc.UpdateEvent(ctx, req, userID)
		assert.NoError(t, err)
	})

	t.Run("invalid status", func(t *testing.T) {
		badReq := req
		badReq.Status = "unknown"

		err := svc.UpdateEvent(ctx, badReq, userID)
		assert.ErrorIs(t, errors.Unwrap(err), service.ErrInvalidStatus)
	})

	t.Run("event not found", func(t *testing.T) {
		mockRepo.EXPECT().UpdateEvent(ctx, gomock.Any()).Return(repo.ErrEventNotFound)

		err := svc.UpdateEvent(ctx, req, userID)
		assert.ErrorIs(t, errors.Unwrap(err), service.ErrEventNotFound)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().UpdateEvent(ctx, gomock.Any()).Return(errors.New("db error"))

		err := svc.UpdateEvent(ctx, req, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EventService.UpdateEventRequest()")
	})
}

func TestEventService_DeleteEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepo(ctrl)
	svc := service.NewEvent(mockRepo)

	ctx := context.Background()
	userID := 42
	eventID := 7

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().DeleteEvent(ctx, eventID, userID).Return(nil)

		err := svc.DeleteEvent(ctx, eventID, userID)
		assert.NoError(t, err)
	})

	t.Run("event not found", func(t *testing.T) {
		mockRepo.EXPECT().DeleteEvent(ctx, eventID, userID).Return(repo.ErrEventNotFound)

		err := svc.DeleteEvent(ctx, eventID, userID)
		assert.ErrorIs(t, errors.Unwrap(err), service.ErrEventNotFound)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().DeleteEvent(ctx, eventID, userID).Return(errors.New("db error"))

		err := svc.DeleteEvent(ctx, eventID, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EventService.DeleteEvent()")
	})
}

func TestEventService_GetEventsForDay(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepo(ctrl)
	svc := service.NewEvent(mockRepo)

	ctx := context.Background()
	userID := 1
	date := time.Now()

	domainEvents := []domain.Event{
		{ID: 1, UserID: userID, Date: date, Status: domain.StatusPlanned, Description: "meeting"},
	}
	expected := dto.GetEventsResponse{
		Events: []dto.Event{
			{Date: date, Status: "planned", Description: "meeting"},
		},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetEventsForDay(ctx, userID, date).Return(domainEvents, nil)

		resp, err := svc.GetEventsForDay(ctx, userID, date)
		assert.NoError(t, err)
		assert.Equal(t, expected, resp)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().GetEventsForDay(ctx, userID, date).Return(nil, errors.New("db error"))

		_, err := svc.GetEventsForDay(ctx, userID, date)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EventService.GetEventsForDay()")
	})
}

func TestEventService_GetEventsForWeek(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepo(ctrl)
	svc := service.NewEvent(mockRepo)

	ctx := context.Background()
	userID := 1
	date := time.Now()

	domainEvents := []domain.Event{
		{ID: 2, UserID: userID, Date: date, Status: domain.StatusDone, Description: "task"},
	}
	expected := dto.GetEventsResponse{
		Events: []dto.Event{
			{Date: date, Status: "done", Description: "task"},
		},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetEventsForWeek(ctx, userID, date).Return(domainEvents, nil)

		resp, err := svc.GetEventsForWeek(ctx, userID, date)
		assert.NoError(t, err)
		assert.Equal(t, expected, resp)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().GetEventsForWeek(ctx, userID, date).Return(nil, errors.New("db error"))

		_, err := svc.GetEventsForWeek(ctx, userID, date)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EventService.GetEventsForWeek()")
	})
}

func TestEventService_GetEventsForMonth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepo(ctrl)
	svc := service.NewEvent(mockRepo)

	ctx := context.Background()
	userID := 1
	date := time.Now()

	domainEvents := []domain.Event{
		{ID: 3, UserID: userID, Date: date, Status: domain.StatusPlanned, Description: "conference"},
	}
	expected := dto.GetEventsResponse{
		Events: []dto.Event{
			{Date: date, Status: "planned", Description: "conference"},
		},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetEventsForMonth(ctx, userID, date).Return(domainEvents, nil)

		resp, err := svc.GetEventsForMonth(ctx, userID, date)
		assert.NoError(t, err)
		assert.Equal(t, expected, resp)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().GetEventsForMonth(ctx, userID, date).Return(nil, errors.New("db error"))

		_, err := svc.GetEventsForMonth(ctx, userID, date)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EventService.GetEventsForMonth()")
	})
}
