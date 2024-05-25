package orderservice

import (
	"context"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/a-x-a/go-loyalty/internal/config"
	"github.com/a-x-a/go-loyalty/internal/customerrors"
	"github.com/a-x-a/go-loyalty/internal/model"
	"github.com/a-x-a/go-loyalty/internal/storage"
	"github.com/a-x-a/go-loyalty/pkg/luhn"
)

type (
	OrderStorage interface {
		Add(ctx context.Context, uid int64, number string) error
		// Update(ctx context.Context, uid int64) error
		GetAll(ctx context.Context, uid int64) (*storage.DTOOrders, error)
	}

	OrderService struct {
		storage OrderStorage
		cfg     config.ServiceConfig
		l       *zap.Logger
	}
)

func New(storage OrderStorage, cfg config.ServiceConfig, l *zap.Logger) *OrderService {
	return &OrderService{storage, cfg, l}
}

func (s *OrderService) Upload(ctx context.Context, uid int64, number string) error {
	if err := s.CheckNumber(ctx, number); err != nil {
		s.l.Debug("upload order error", zap.Error(errors.Wrap(err, "orderservice.upload")))
		return err
	}

	return s.storage.Add(ctx, uid, number)
}

func (s *OrderService) GetAll(ctx context.Context, uid int64) (*model.Orders, error) {
	o, err := s.storage.GetAll(ctx, uid)
	if err != nil {
		s.l.Debug("get all order error", zap.Error(errors.Wrap(err, "orderservice.getall")))
		return nil, err
	}

	s.l.Debug("getall", zap.Any("orders", o))

	orders := model.Orders{}
	for _, v := range *o {
		order := model.Order{
			Number:     v.Number,
			Status:     model.OrderStatus(v.Status).String(),
			Accrual:    v.Accrual,
			UploadedAt: v.UploadedAt.Local(),
		}
		orders = append(orders, order)
	}
	return &orders, nil
}

func (s *OrderService) CheckNumber(ctx context.Context, number string) error {
	if luhn.Check(number) {
		return nil
	}

	return customerrors.ErrInvalidOrderNumber
}
