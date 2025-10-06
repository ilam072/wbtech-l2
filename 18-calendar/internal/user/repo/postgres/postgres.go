package postgres

import (
	"context"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/repo"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/types/domain"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/errutils"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user domain.User) error {
	const op = "repo.postgres.CreateUser()"

	query, args, err := goqu.Insert("users").Rows(user).ToSQL()
	if err != nil {
		return errutils.Wrap(op, err)
	}

	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return errutils.Wrap(op, repo.ErrUserExists)
			}
		}
		return errutils.Wrap(op, err)
	}

	return nil
}

func (r *UserRepo) User(ctx context.Context, username string) (domain.User, error) {
	const op = "repo.postgres.User()"

	query, args, err := goqu.From("users").
		Select("id", "username", "password_hash").
		Where(goqu.Ex{"username": username}).ToSQL()

	if err != nil {
		return domain.User{}, errutils.Wrap(op, err)
	}

	row := r.pool.QueryRow(ctx, query, args...)
	if err != nil {
		return domain.User{}, errutils.Wrap(op, err)
	}

	var user domain.User
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
	); err != nil {
		return domain.User{}, errutils.Wrap(op, err)
	}

	return user, nil
}
