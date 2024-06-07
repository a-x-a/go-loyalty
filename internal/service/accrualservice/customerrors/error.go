package customerrors

import "errors"

var (
	ErrNotFoundOrders      = errors.New("нет заказов для обработки")
	ErrClientIsNoAvailable = errors.New("client is no available")
	ErrInvalidAccrualOrder = errors.New("invalid accrual order")                      // 500
	ErrNoContent           = errors.New("заказ не зарегистрирован в системе расчёта") // 204
	ErrTooManyRequests     = errors.New("превышено количество запросов к сервису")    // 429
)
