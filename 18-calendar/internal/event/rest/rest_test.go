package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/mocks"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/rest"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/service"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/dto"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/response"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvent := mocks.NewMockEvent(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	app := fiber.New()
	h := rest.NewEventHandler(slogdiscard.NewDiscardLogger(), mockEvent, mockValidator)
	app.Post("/api/create_event", func(c *fiber.Ctx) error {
		c.Locals("userID", 42) // эмуляция авторизации
		return h.CreateEventHandler(c)
	})

	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/create_event", bytes.NewBuffer([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("validation error", func(t *testing.T) {
		// пустое описание
		body, _ := json.Marshal(dto.CreateEventRequest{
			Date:        time.Now(),
			Description: "",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/create_event", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		mockValidator.EXPECT().
			Validate(gomock.Any()).
			Return(errors.New("validation failed"))
		mockValidator.EXPECT().
			FormatValidationErrors(gomock.Any()).
			Return(map[string]string{"description": "required"})

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("internal error", func(t *testing.T) {
		body, _ := json.Marshal(dto.CreateEventRequest{
			Date:        time.Now(),
			Description: "Team meeting",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/create_event", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		mockValidator.EXPECT().
			Validate(gomock.Any()).
			Return(nil)

		mockEvent.EXPECT().
			CreateEvent(gomock.Any(), gomock.Any(), 42).
			Return(errors.New("db error"))

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("success", func(t *testing.T) {
		body, _ := json.Marshal(dto.CreateEventRequest{
			Date:        time.Now(),
			Description: "Project deadline",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/create_event", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		mockValidator.EXPECT().
			Validate(gomock.Any()).
			Return(nil)

		mockEvent.EXPECT().
			CreateEvent(gomock.Any(), gomock.Any(), 42).
			Return(nil)

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})
}

func TestUpdateEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvent := mocks.NewMockEvent(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	app := fiber.New()
	h := rest.NewEventHandler(slogdiscard.NewDiscardLogger(), mockEvent, mockValidator)

	app.Post("/api/update_event", func(c *fiber.Ctx) error {
		c.Locals("userID", 42)
		return h.UpdateEventHandler(c)
	})

	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/update_event", bytes.NewBuffer([]byte("bad-json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("validation error", func(t *testing.T) {
		body, _ := json.Marshal(dto.UpdateEventRequest{
			ID:          1,
			Date:        time.Now(),
			Status:      "",
			Description: "test event",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/update_event", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation failed"))
		mockValidator.EXPECT().FormatValidationErrors(gomock.Any()).Return(map[string]string{"status": "required"})

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid status", func(t *testing.T) {
		body, _ := json.Marshal(dto.UpdateEventRequest{
			ID:          1,
			Date:        time.Now(),
			Status:      "unknown",
			Description: "test event",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/update_event", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockEvent.EXPECT().UpdateEvent(gomock.Any(), gomock.Any(), 42).Return(service.ErrInvalidStatus)

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("event not found", func(t *testing.T) {
		body, _ := json.Marshal(dto.UpdateEventRequest{
			ID:          999,
			Date:        time.Now(),
			Status:      "planned",
			Description: "test event",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/update_event", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockEvent.EXPECT().UpdateEvent(gomock.Any(), gomock.Any(), 42).Return(service.ErrEventNotFound)

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("internal error", func(t *testing.T) {
		body, _ := json.Marshal(dto.UpdateEventRequest{
			ID:          2,
			Date:        time.Now(),
			Status:      "planned",
			Description: "test event",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/update_event", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockEvent.EXPECT().UpdateEvent(gomock.Any(), gomock.Any(), 42).Return(errors.New("db error"))

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("success", func(t *testing.T) {
		body, _ := json.Marshal(dto.UpdateEventRequest{
			ID:          3,
			Date:        time.Now(),
			Status:      "planned",
			Description: "test event",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/update_event", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockEvent.EXPECT().UpdateEvent(gomock.Any(), gomock.Any(), 42).Return(nil)

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestDeleteEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvent := mocks.NewMockEvent(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	app := fiber.New()
	h := rest.NewEventHandler(slogdiscard.NewDiscardLogger(), mockEvent, mockValidator)

	app.Delete("/api/delete_event/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", 42)
		return h.DeleteEventHandler(c)
	})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/delete_event/not-a-number", nil)

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("event not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/delete_event/123", nil)

		mockEvent.EXPECT().DeleteEvent(gomock.Any(), 123, 42).Return(service.ErrEventNotFound)

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("internal error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/delete_event/124", nil)

		mockEvent.EXPECT().DeleteEvent(gomock.Any(), 124, 42).Return(errors.New("db error"))

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/delete_event/125", nil)

		mockEvent.EXPECT().DeleteEvent(gomock.Any(), 125, 42).Return(nil)

		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func TestGetEventsForDayHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvent := mocks.NewMockEvent(ctrl)

	app := fiber.New()
	h := rest.NewEventHandler(slogdiscard.NewDiscardLogger(), mockEvent, nil)

	events := app.Group("/api")
	events.Get("/events_for_day", func(c *fiber.Ctx) error {
		c.Locals("userID", 42)
		return h.GetEventsForDayHandler(c)
	})

	t.Run("invalid date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/events_for_day?date=not-a-date", nil)
		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var r response.Resp
		_ = json.NewDecoder(resp.Body).Decode(&r)
		assert.Equal(t, "error", r.Status)
		// текст берётся из хендлера: msgInvalidDateFormat
		assert.NotEmpty(t, r.Data)
	})

	t.Run("internal error", func(t *testing.T) {
		q := time.Now().Format(time.DateOnly)
		req := httptest.NewRequest(http.MethodGet, "/api/events_for_day?date="+q, nil)

		mockEvent.EXPECT().
			GetEventsForDay(gomock.Any(), 42, gomock.Any()).
			Return(dto.GetEventsResponse{}, errors.New("db error"))

		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var r response.Resp
		_ = json.NewDecoder(resp.Body).Decode(&r)
		assert.Equal(t, "error", r.Status)
		assert.NotEmpty(t, r.Data)
	})

	t.Run("success", func(t *testing.T) {
		q := time.Now().Format(time.DateOnly)
		req := httptest.NewRequest(http.MethodGet, "/api/events_for_day?date="+q, nil)

		t1 := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)
		t2 := time.Date(2025, 8, 15, 2, 0, 0, 0, time.UTC)
		expected := dto.GetEventsResponse{
			Events: []dto.Event{
				{Date: t1, Status: "planned", Description: "meeting"},
				{Date: t2, Status: "done", Description: "coding"},
			},
		}

		mockEvent.EXPECT().
			GetEventsForDay(gomock.Any(), 42, gomock.Any()).
			Return(expected, nil)

		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var got dto.GetEventsResponse
		_ = json.NewDecoder(resp.Body).Decode(&got)
		assert.Equal(t, expected, got)
	})
}

func TestGetEventsForWeekHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvent := mocks.NewMockEvent(ctrl)

	app := fiber.New()
	h := rest.NewEventHandler(slogdiscard.NewDiscardLogger(), mockEvent, nil)

	events := app.Group("/api")
	events.Get("/events_for_week", func(c *fiber.Ctx) error {
		c.Locals("userID", 42)
		return h.GetEventsForWeekHandler(c)
	})

	t.Run("invalid date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/events_for_week?date=bad", nil)
		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var r response.Resp
		_ = json.NewDecoder(resp.Body).Decode(&r)
		assert.Equal(t, "error", r.Status)
	})

	t.Run("internal error", func(t *testing.T) {
		q := time.Now().Format(time.DateOnly)
		req := httptest.NewRequest(http.MethodGet, "/api/events_for_week?date="+q, nil)

		mockEvent.EXPECT().
			GetEventsForWeek(gomock.Any(), 42, gomock.Any()).
			Return(dto.GetEventsResponse{}, errors.New("db error"))

		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var r response.Resp
		_ = json.NewDecoder(resp.Body).Decode(&r)
		assert.Equal(t, "error", r.Status)
	})

	t.Run("success", func(t *testing.T) {
		q := time.Now().Format(time.DateOnly)
		req := httptest.NewRequest(http.MethodGet, "/api/events_for_week?date="+q, nil)

		t1 := time.Date(2025, 8, 18, 0, 0, 0, 0, time.UTC)
		expected := dto.GetEventsResponse{
			Events: []dto.Event{
				{Date: t1, Status: "planned", Description: "gym"},
			},
		}

		mockEvent.EXPECT().
			GetEventsForWeek(gomock.Any(), 42, gomock.Any()).
			Return(expected, nil)

		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var got dto.GetEventsResponse
		_ = json.NewDecoder(resp.Body).Decode(&got)
		assert.Equal(t, expected, got)
	})
}

func TestGetEventsForMonthHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvent := mocks.NewMockEvent(ctrl)

	app := fiber.New()
	h := rest.NewEventHandler(slogdiscard.NewDiscardLogger(), mockEvent, nil)

	events := app.Group("/api")
	events.Get("/events_for_month", func(c *fiber.Ctx) error {
		c.Locals("userID", 42)
		return h.GetEventsForMonthHandler(c)
	})

	t.Run("invalid date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/events_for_month?date=x", nil)
		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var r response.Resp
		_ = json.NewDecoder(resp.Body).Decode(&r)
		assert.Equal(t, "error", r.Status)
	})

	t.Run("internal error", func(t *testing.T) {
		q := time.Now().Format(time.DateOnly)
		req := httptest.NewRequest(http.MethodGet, "/api/events_for_month?date="+q, nil)

		mockEvent.EXPECT().
			GetEventsForMonth(gomock.Any(), 42, gomock.Any()).
			Return(dto.GetEventsResponse{}, errors.New("db error"))

		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var r response.Resp
		_ = json.NewDecoder(resp.Body).Decode(&r)
		assert.Equal(t, "error", r.Status)
	})

	t.Run("success", func(t *testing.T) {
		q := time.Now().Format(time.DateOnly)
		req := httptest.NewRequest(http.MethodGet, "/api/events_for_month?date="+q, nil)

		t1 := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
		expected := dto.GetEventsResponse{
			Events: []dto.Event{
				{Date: t1, Status: "planned", Description: "conference"},
			},
		}

		mockEvent.EXPECT().
			GetEventsForMonth(gomock.Any(), 42, gomock.Any()).
			Return(expected, nil)

		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var got dto.GetEventsResponse
		_ = json.NewDecoder(resp.Body).Decode(&got)
		assert.Equal(t, expected, got)
	})
}
