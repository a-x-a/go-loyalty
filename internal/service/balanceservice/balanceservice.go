package balanceservice

import (
	"context"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/a-x-a/go-loyalty/internal/config"
	"github.com/a-x-a/go-loyalty/internal/customerrors"
	"github.com/a-x-a/go-loyalty/internal/model"
	"github.com/a-x-a/go-loyalty/internal/storage"
)

type (
	BalanceStorage interface {
		Create(ctx context.Context, uid int64) error
		Get(ctx context.Context, uid int64) (*storage.DTOBalance, error)
		Withdraw(ctx context.Context, uid int64, number string, sum float64) error
		GetWithdrawals(ctx context.Context, uid int64) (*storage.DTOWithdrawals, error)
		Accrual(ctx context.Context, uid int64, sum float64) error
	}

	BalanceService struct {
		storage BalanceStorage
		cfg     config.ServiceConfig
		l       *zap.Logger
	}
)

func New(storage BalanceStorage, cfg config.ServiceConfig, l *zap.Logger) *BalanceService {
	return &BalanceService{storage, cfg, l}
}

func (s *BalanceService) Create(ctx context.Context, uid int64) error {
	s.l.Debug("create user ballance", zap.Int64("uid", uid))

	err := s.storage.Create(ctx, uid)
	if err != nil {
		s.l.Debug("failed to create user ballance", zap.Error(errors.Wrap(err, "storage.create")))
		return err
	}

	return err
}

func (s *BalanceService) Get(ctx context.Context, uid int64) (*model.Balance, error) {
	b, err := s.storage.Get(ctx, uid)
	if err != nil {
		s.l.Debug("failed to get user ballance", zap.Error(errors.Wrap(err, "storage.get")))
		return nil, err
	}

	s.l.Debug("getbalance", zap.Any("balance", b))

	return &model.Balance{
		Current:   b.Current,
		Withdrawn: b.Withdrawn,
	}, nil
}

func (s *BalanceService) Withdraw(ctx context.Context, uid int64, number string, sum float64) error {
	err := s.storage.Withdraw(ctx, uid, number, sum)
	if err != nil {
		s.l.Debug("withdraw error", zap.Error(errors.Wrap(err, "balanceservice.withdraw")))
		return err
	}

	return nil
}

func (s *BalanceService) GetWithdrawals(ctx context.Context, uid int64) (*model.Withdrawals, error) {
	w, err := s.storage.GetWithdrawals(ctx, uid)
	if err != nil {
		s.l.Debug("withdrawals error", zap.Error(errors.Wrap(err, "balanceservice.getwithdrawals")))
		return nil, err
	}

	s.l.Debug("getwithdrawals", zap.Any("withdrawals", w))

	if len(*w) == 0 {
		return nil, customerrors.ErrNoContent
	}

	withdrawals := model.Withdrawals{}
	for _, v := range *w {
		withdrawal := model.Withdrawal{
			Order:       v.Order,
			Sum:         v.Sum,
			ProcessedAt: v.ProcessedAt.Local(),
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	return &withdrawals, nil
}

func (s *BalanceService) Update(ctx context.Context, uid int64, sum float64) error {
	return s.storage.Accrual(ctx, uid, sum)
}
