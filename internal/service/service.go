package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/a-x-a/go-loyalty/internal/model"
)

type (
	Service struct {
		UserService
		OrderService
		BallanceService
		l *zap.Logger
	}

	UserService interface {
		Register(ctx context.Context, login, password string) error
		Login(ctx context.Context, login, password string) (string, error)
	}

	OrderService interface {
		UploadOrder(ctx context.Context, uid int64, number string) error
		GetAllOrders(ctx context.Context, uid int64) (*model.Orders, error)
		CheckNumber(ctx context.Context, number string) error
	}

	BallanceService interface {
		Get(ctx context.Context, uid int64) (*model.Balance, error)
		Withdraw(ctx context.Context, uid int64, number string, sum float64) error
		GetWithdrawals(ctx context.Context, uid int64) (*model.Withdrawals, error)
	}
)

func New(userService UserService, balanceService BallanceService, log *zap.Logger) *Service {
	return &Service{
		UserService:     userService,
		BallanceService: balanceService,
		l:               log,
	}
}

// Auth service.
func (s *Service) RegisterUser(ctx context.Context, login, password string) (string, error) {
	err := s.UserService.Register(ctx, login, password)
	if err != nil {
		return "", err
	}

	return s.UserService.Login(ctx, login, password)
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	return s.UserService.Login(ctx, login, password)
}

// Order service.
func (s *Service) UploadOrder(ctx context.Context, uid int64, number string) error {
	return nil
}

func (s *Service) GetAllOrders(ctx context.Context, uid int64) (*model.Orders, error) {
	orders, err := s.OrderService.GetAllOrders(ctx, uid)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Service) CheckOrderNumber(ctx context.Context, number string) error {
	// s.OrderService.CheckNumber(ctx, number)
	return nil
}

// Ballance service.
func (s *Service) GetBalance(ctx context.Context, uid int64) (*model.Balance, error) {
	b, err := s.BallanceService.Get(ctx, uid)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) WithdrawBalance(ctx context.Context, uid int64, number string, sum float64) error {
	if sum == 0 {
		return nil
	}

	// TODO
	// проверить номер заказа orderservice.CheckOrderNumber(ctx context.Context, number string) error
	if err := s.CheckOrderNumber(ctx, number); err != nil {
		return err
	}
	// return customerrors.ErrInvalidOrderNumber
	// выполнить запрос на списание
	err := s.BallanceService.Withdraw(ctx, uid, number, sum)

	return err
}

func (s *Service) GetWithdrawalsBalance(ctx context.Context, uid int64) (*model.Withdrawals, error) {
	withdrawals, err := s.BallanceService.GetWithdrawals(ctx, uid)

	return withdrawals, err
}
