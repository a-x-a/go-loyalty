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
		BalanceService
		l *zap.Logger
	}

	UserService interface {
		Register(ctx context.Context, login, password string) (int64, error)
		Login(ctx context.Context, login, password string) (string, error)
	}

	OrderService interface {
		Upload(ctx context.Context, uid int64, number string) error
		GetAll(ctx context.Context, uid int64) (*model.Orders, error)
		CheckNumber(ctx context.Context, number string) error
	}

	BalanceService interface {
		Create(ctx context.Context, uid int64) error
		Get(ctx context.Context, uid int64) (*model.Balance, error)
		Withdraw(ctx context.Context, uid int64, number string, sum float64) error
		GetWithdrawals(ctx context.Context, uid int64) (*model.Withdrawals, error)
	}
)

func New(userService UserService, orderService OrderService, balanceService BalanceService, log *zap.Logger) *Service {
	return &Service{
		UserService:    userService,
		OrderService:   orderService,
		BalanceService: balanceService,
		l:              log,
	}
}

// Auth service.
func (s *Service) RegisterUser(ctx context.Context, login, password string) (string, error) {
	uid, err := s.UserService.Register(ctx, login, password)
	if err != nil {
		return "", err
	}

	if uid > 0 {
		if err := s.BalanceService.Create(ctx, uid); err != nil {
			return "", err
		}
	}

	return s.UserService.Login(ctx, login, password)
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	return s.UserService.Login(ctx, login, password)
}

// Order service.
func (s *Service) UploadOrder(ctx context.Context, uid int64, number string) error {
	return s.OrderService.Upload(ctx, uid, number)
}

func (s *Service) GetAllOrders(ctx context.Context, uid int64) (*model.Orders, error) {
	orders, err := s.OrderService.GetAll(ctx, uid)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Service) CheckOrderNumber(ctx context.Context, number string) error {
	return s.OrderService.CheckNumber(ctx, number)
}

// Ballance service.
func (s *Service) GetBalance(ctx context.Context, uid int64) (*model.Balance, error) {
	b, err := s.BalanceService.Get(ctx, uid)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) WithdrawBalance(ctx context.Context, uid int64, number string, sum float64) error {
	if sum == 0 {
		return nil
	}

	// проверить номер заказа.
	if err := s.CheckOrderNumber(ctx, number); err != nil {
		return err
	}

	// выполнить запрос на списание.
	err := s.BalanceService.Withdraw(ctx, uid, number, sum)

	return err
}

func (s *Service) GetWithdrawalsBalance(ctx context.Context, uid int64) (*model.Withdrawals, error) {
	withdrawals, err := s.BalanceService.GetWithdrawals(ctx, uid)

	return withdrawals, err
}
