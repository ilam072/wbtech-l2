package service

import (
	"context"
	"errors"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/types/domain"
	"time"
)

//go:generate mockgen -source=user.go -destination=../mocks/service_mocks.go -package=mocks
type UserRepo interface {
	CreateUser(context.Context, domain.User) error
	User(context.Context, string) (domain.User, error)
}

type TokenManager interface {
	NewToken(userID int, ttl time.Duration) (string, error)
	ParseToken(tokenStr string) (int, error)
}

type User struct {
	userRepo UserRepo
	manager  TokenManager

	tokenTTL time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("username is already occupied")
)

func NewUser(repo UserRepo, manager TokenManager, tokenTTL time.Duration) *User {
	return &User{
		userRepo: repo,
		manager:  manager,
		tokenTTL: tokenTTL,
	}
}
