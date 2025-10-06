package service

import (
	"context"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/types/dto"
	"github.com/ilam072/wbtech-l2/18-calendar/pkg/errutils"
	"golang.org/x/crypto/bcrypt"
)

func (u *User) Login(ctx context.Context, user dto.LoginUser) (string, error) {
	const op = "UserService.Login()"

	domainUser, err := u.userRepo.User(ctx, user.Username)
	if err != nil {
		return "", errutils.Wrap(op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(domainUser.PasswordHash), []byte(user.Password)); err != nil {
		return "", errutils.Wrap(op, ErrInvalidCredentials)
	}

	token, err := u.manager.NewToken(domainUser.ID, u.tokenTTL)
	if err != nil {
		return "", errutils.Wrap(op, err)
	}

	return token, nil
}
