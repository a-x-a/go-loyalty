package customerrors

import "github.com/pkg/errors"

var (
	ErrInvalidRequestFormat      = errors.New("не верный формат запроса")
	ErrInvalidUsernameOrPassword = errors.New("не верный логин или пароль")
	ErrUsernameAlreadyTaken      = errors.New("логин уже занят")
	ErrNotRegisteredUser         = errors.New("пользователь не зарегистрирован")
	ErrInsufficientFunds         = errors.New("на счету недостаточно средств")
	ErrInvalidOrderNumber        = errors.New("неверный номер заказа")
)
