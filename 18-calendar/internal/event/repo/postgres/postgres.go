package postgres

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/repo"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/event/types/domain"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/errutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type EventRepo struct {
	pool *pgxpool.Pool
}

func NewEventRepo(db *pgxpool.Pool) *EventRepo {
	return &EventRepo{pool: db}
}

func (r *EventRepo) CreateEvent(ctx context.Context, event domain.Event) error {
	const op = "postgres.CreateEvent()"

	query, args, err := goqu.Insert("events").Rows(event).ToSQL()
	if err != nil {
		return errutils.Wrap(op, err)
	}

	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		return errutils.Wrap(op, err)
	}

	return nil
}

func (r *EventRepo) UpdateEvent(ctx context.Context, event domain.Event) error {
	const op = "postgres.CreateEvent()"

	query, args, err := goqu.Update("events").
		Set(goqu.Record{
			"date":        event.Date,
			"status":      event.Status,
			"description": event.Description,
		}).
		Where(goqu.Ex{"id": event.ID, "user_id": event.UserID}).
		ToSQL()
	if err != nil {
		return errutils.Wrap(op, err)
	}

	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return errutils.Wrap(op, err)
	}
	if res.RowsAffected() == 0 {
		return errutils.Wrap(op, repo.ErrEventNotFound)
	}

	return nil
}

func (r *EventRepo) DeleteEvent(ctx context.Context, eventID int, userID int) error {
	const op = "postgres.DeleteEvent()"

	query, args, err := goqu.Delete("events").Where(goqu.Ex{"id": eventID, "user_id": userID}).ToSQL()
	if err != nil {
		return errutils.Wrap(op, err)
	}

	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return errutils.Wrap(op, err)
	}
	if res.RowsAffected() == 0 {
		return errutils.Wrap(op, repo.ErrEventNotFound)
	}

	return nil
}

func (r *EventRepo) GetEventsForDay(ctx context.Context, userID int, date time.Time) ([]domain.Event, error) {
	const op = "postgres.GetEventsForDay()"

	query, args, err := goqu.From("events").
		Where(goqu.Ex{
			"user_id": userID,
			"date":    date,
		}).ToSQL()
	if err != nil {
		return nil, errutils.Wrap(op, err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errutils.Wrap(op, err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		if err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.Date,
			&event.Status,
			&event.Description,
		); err != nil {
			return nil, errutils.Wrap(op, err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepo) GetEventsForWeek(ctx context.Context, userID int, start time.Time) ([]domain.Event, error) {
	const op = "postgres.GetEventsForWeek()"

	end := start.AddDate(0, 0, 6)

	query, args, err := goqu.From("events").
		Where(
			goqu.C("user_id").Eq(userID),
			goqu.C("date").Between(goqu.Range(start, end)),
		).ToSQL()
	if err != nil {
		return nil, errutils.Wrap(op, err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errutils.Wrap(op, err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		if err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.Date,
			&event.Status,
			&event.Description,
		); err != nil {
			return nil, errutils.Wrap(op, err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepo) GetEventsForMonth(ctx context.Context, userID int, start time.Time) ([]domain.Event, error) {
	const op = "postgres.GetEventsForMonth()"

	end := start.AddDate(0, 1, -1)

	query, args, err := goqu.From("events").
		Where(
			goqu.C("user_id").Eq(userID),
			goqu.C("date").Between(goqu.Range(start, end)),
		).ToSQL()
	if err != nil {
		return nil, errutils.Wrap(op, err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errutils.Wrap(op, err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		if err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.Date,
			&event.Status,
			&event.Description,
		); err != nil {
			return nil, errutils.Wrap(op, err)
		}
		events = append(events, event)
	}

	return events, nil
}
