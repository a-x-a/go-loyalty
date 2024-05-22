package model

import "github.com/pkg/errors"

var (
	ErrInvalidRequestFormat      = errors.New("не верный формат запроса")
	ErrInvalidUsernameOrPassword = errors.New("не верный логин или пароль")
	ErrUsernameAlreadyTaken      = errors.New("логин уже занят")
	ErrNotRegisteredUser         = errors.New("пользователь не зарегистрирован")
)
