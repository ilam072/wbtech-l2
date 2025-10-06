package service

import (
	"context"
	"errors"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/repo"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/types/domain"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/types/dto"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/errutils"
	"golang.org/x/crypto/bcrypt"
)

func (u *User) RegisterUser(ctx context.Context, user dto.RegisterUser) error {
	const op = "UserService.RegisterUser()"

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errutils.Wrap(op, err)
	}

	domainUser := domain.User{
		Username:     user.Username,
		PasswordHash: string(passwordHash),
	}

	if err = u.userRepo.CreateUser(ctx, domainUser); err != nil {
		if errors.Is(err, repo.ErrUserExists) {
			return errutils.Wrap(op, ErrUserExists)

		}
		return errutils.Wrap(op, err)
	}

	return nil
}
