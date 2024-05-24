package customerrors

import "github.com/pkg/errors"

var (
	ErrInvalidRequestFormat = errors.New("не верный формат запроса") // 400

	ErrInvalidUsernameOrPassword = errors.New("не верный логин или пароль")
	ErrUsernameAlreadyTaken      = errors.New("логин уже занят")
	ErrNotRegisteredUser         = errors.New("пользователь не зарегистрирован")

	ErrInsufficientFunds = errors.New("на счету недостаточно средств") // 402
	ErrNotWithdrawals    = errors.New("нет списаний")                  // 204

	ErrInvalidOrderNumber         = errors.New("неверный номер заказа")                              // 422
	ErrOrderUploadedByUser        = errors.New("номер заказа уже был загружен пользователем")        // 200
	ErrOrderInProcess             = errors.New("новый номер заказа принят в обработку")              // 202
	ErrOrderUploadedByAnotherUser = errors.New("номер заказа уже был загружен другим пользователем") // 409
	ErrNotOrders                  = errors.New("нет данных для ответа")                              // 204

)
