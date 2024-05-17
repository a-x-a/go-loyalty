package service

import (
	"context"

	"github.com/a-x-a/go-loyalty/internal/model"
)

type (
	Service struct {
		UserService
		OrderService
		BallanceService
	}

	UserService interface {
		RegisterUser(ctx context.Context, login, password string) (string, error)
		Login(ctx context.Context, login, password string) (string, error)
	}

	OrderService interface {
		UploadOrder(ctx context.Context, userID int64, number string) error
		GetAllOrders(ctx context.Context, userID int64) (*model.Orders, error)
	}

	BallanceService interface {
		GetBalance(ctx context.Context, userID int64) (*model.Balance, error)
		WithdrawBalance(ctx context.Context, userID int64, number string, sum float64) error
		GetWithdrawalsBalance(ctx context.Context, userID int64) (*model.Withdrawals, error)
	}
)

func New() *Service {
	return &Service{}
}
