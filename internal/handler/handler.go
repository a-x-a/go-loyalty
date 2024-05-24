package handler

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"

	"github.com/a-x-a/go-loyalty/internal/model"
)

type (
	APIService interface {
		// UserService
		RegisterUser(ctx context.Context, login, password string) (string, error)
		Login(ctx context.Context, login, password string) (string, error)
		// OrderService
		UploadOrder(ctx context.Context, uid int64, number string) error
		GetAllOrders(ctx context.Context, uid int64) (*model.Orders, error)
		// BallanceService
		GetBalance(ctx context.Context, uid int64) (*model.Balance, error)
		WithdrawBalance(ctx context.Context, uid int64, number string, sum float64) error
		GetWithdrawalsBalance(ctx context.Context, uid int64) (*model.Withdrawals, error)
	}

	Handler struct {
		s APIService
	}
)

func New(s APIService) *Handler {
	return &Handler{
		s: s,
	}
}

func responseWithError(c echo.Context, code int, err error) error {
	// resp := fmt.Sprintf("%d: %s", code, err.Error())
	// return c.JSON(code, err)
	return c.NoContent(code)
}

func responseWithCode(c echo.Context, code int) error {
	// c.Response().WriteHeader(code)
	return c.NoContent(code)
}

func getUserId(c echo.Context) (int64, error) {
	uid, ok := c.Get("uid").(int64)
	if !ok {
		return 0, errors.New("user id is not defined")
	}

	return uid, nil
}
