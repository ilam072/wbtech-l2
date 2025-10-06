package rest

import (
	"context"
	"github.com/ilam072/wbtech-l2/18-calendar/internal/user/types/dto"
	"log/slog"
)

const (
	msgInternalError      = "something went wrong, try again later"
	msgUsernameOccupied   = "username is already occupied"
	msgInvalidRequestBody = "invalid request body"
	msgInvalidCredential  = "invalid credentials"
)

//go:generate mockgen -source=handler.go -destination=../../user/mocks/user_handler_mock.go -package=mocks
type User interface {
	RegisterUser(context.Context, dto.RegisterUser) error
	Login(context.Context, dto.LoginUser) (string, error)
}

type Validator interface {
	Validate(i interface{}) error
	FormatValidationErrors(err error) map[string]string
}

type UserHandler struct {
	log       *slog.Logger
	user      User
	validator Validator
}

func NewUserHandler(log *slog.Logger, user User, validator Validator) *UserHandler {
	h := &UserHandler{
		log:       log,
		user:      user,
		validator: validator,
	}

	return h
}
