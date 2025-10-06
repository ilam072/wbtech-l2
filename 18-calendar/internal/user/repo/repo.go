package repo

import "errors"

var (
	ErrUserExists = errors.New("username is already occupied")
)
