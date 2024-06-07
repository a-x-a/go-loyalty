package orderservice

import (
	"context"
	"regexp"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/a-x-a/go-loyalty/internal/config"
	"github.com/a-x-a/go-loyalty/internal/customerrors"
	"github.com/a-x-a/go-loyalty/internal/model"
	accrualModel "github.com/a-x-a/go-loyalty/internal/service/accrualservice/model"
	"github.com/a-x-a/go-loyalty/internal/storage"
	"github.com/a-x-a/go-loyalty/pkg/luhn"
)

type (
	OrderStorage interface {
		Add(ctx context.Context, uid int64, number string) error
		GetAll(ctx context.Context, uid int64) (*storage.DTOOrders, error)
		Update(ctx context.Context, order storage.DTOOrder) error
		GetToProcessing(ctx context.Context) (*storage.DTOAccrualOrders, error)
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

	s.l.Debug("get orders", zap.Any("orders", o))

	if len(*o) == 0 {
		return nil, customerrors.ErrNoContent
	}

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
	digitsRegExp := regexp.MustCompile(`^\d+$`)
	if !digitsRegExp.MatchString(number) {
		return customerrors.ErrInvalidRequestFormat
	}

	if luhn.Check(number) {
		return nil
	}

	return customerrors.ErrInvalidOrderNumber
}

func (s *OrderService) GetOrdersToProcessing(ctx context.Context) (*accrualModel.AccrualOrders, error) {
	o, err := s.storage.GetToProcessing(ctx)
	if err != nil {
		s.l.Debug("get orders error", zap.Error(errors.Wrap(err, "orderservice.getorderstoprocessing")))
		return nil, err
	}

	s.l.Debug("get orders", zap.Any("orders", *o))

	orders := accrualModel.AccrualOrders{}
	for _, v := range *o {
		order := accrualModel.AccrualOrder{
			UID:   v.UID,
			Order: v.Number,
		}
		orders = append(orders, order)
	}

	return &orders, nil
}

func (s *OrderService) Update(ctx context.Context, number string, status int, accrual float64) error {
	if err := s.CheckNumber(ctx, number); err != nil {
		s.l.Debug("update order error", zap.Error(errors.Wrap(err, "orderservice.update")))
		return err
	}

	order := storage.DTOOrder{
		Number:  number,
		Status:  status,
		Accrual: accrual,
	}

	return s.storage.Update(ctx, order)
}
