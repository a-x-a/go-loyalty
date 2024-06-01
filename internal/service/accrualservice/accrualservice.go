package accrualservice

import (
	"context"

	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/a-x-a/go-loyalty/internal/service/accrualservice/client"
	"github.com/a-x-a/go-loyalty/internal/service/accrualservice/customerrors"
	accrualErr "github.com/a-x-a/go-loyalty/internal/service/accrualservice/customerrors"
	accrualModel "github.com/a-x-a/go-loyalty/internal/service/accrualservice/model"
)

type (
	AccrualStorage interface{}

	AccrualClient interface {
		Get(ctx context.Context, number string) (accrualModel.AccrualOrder, error)
	}

	Services struct {
		Order   OrderService
		Balance BalanceService
	}

	OrderService interface {
		GetOrdersToProcessing(ctx context.Context) (*accrualModel.AccrualOrders, error)
		Update(ctx context.Context, number string, status int, accrual float64) error
	}

	BalanceService interface {
		Update(ctx context.Context, uid int64, accrual float64) error
	}

	AccrualWorker struct {
		storage AccrualStorage
		client  AccrualClient
		l       *zap.Logger

		frequency        time.Duration
		accrualRateLimit int

		Services Services
	}
)

func New(orderService OrderService, balanceService BalanceService,
	storage AccrualStorage, accrualSystemAddress string,
	frequency time.Duration, accrualRateLimit int, l *zap.Logger) *AccrualWorker {
	accrualClient := client.New(accrualSystemAddress, l)

	return &AccrualWorker{
		storage:          storage,
		client:           accrualClient,
		l:                l,
		frequency:        frequency,
		accrualRateLimit: accrualRateLimit,

		Services: Services{
			Order:   orderService,
			Balance: balanceService,
		},
	}
}

func (s *AccrualWorker) Start(ctx context.Context) error {
	timer := time.NewTicker(s.frequency)
	defer timer.Stop()

	var i int
	for {
		select {
		case <-timer.C:
			i++
			s.l.Info("processing orders", zap.Int("worker", i))
			ordersToProcessingChan, err := s.getOrdersToProceessing(ctx)
			if err != nil {
				s.l.Debug("failed to get not processed orders", zap.Error(errors.Wrap(err, "accrualworker.getorderstosync")))
				continue
			}

			responceChan := make(chan accrualModel.AccrualOrder, len(ordersToProcessingChan))

			wg := sync.WaitGroup{}
			wg.Add(s.accrualRateLimit)
			for w := 1; w <= s.accrualRateLimit; w++ {
				s.getAccrualOrdersResp(ctx, &wg, ordersToProcessingChan, responceChan)
			}

			go func() {
				wg.Wait()
				close(responceChan)
			}()

			for order := range responceChan {
				s.l.Debug("get order from response channel", zap.Any("order", order))

				err = s.updateOrder(ctx, order)
				if err != nil {
					s.l.Debug("update order", zap.Error(errors.Wrap(err, "accrualworker.updateorder")))
					continue
				}

				if order.Accrual != 0 {
					err = s.updateBalance(ctx, order)
					if err != nil {
						s.l.Debug("update balance", zap.Error(errors.Wrap(err, "accrualworker.updatebalance")))
					}
				}
			}

		case <-ctx.Done():
			return nil
		}
	}
}

func (s *AccrualWorker) getOrdersToProceessing(ctx context.Context) (chan accrualModel.AccrualOrder, error) {
	ordersToProceessing, err := s.Services.Order.GetOrdersToProcessing(ctx)
	if err != nil {
		s.l.Debug("failed to get not porocessing orders",
			zap.Error(errors.Wrap(err, "accrualservice.getorderstoproceed")))

		return nil, err
	}

	if len(*ordersToProceessing) == 0 {
		return nil, accrualErr.ErrNotFoundOrders
	}

	ordersToProcessingChan := make(chan accrualModel.AccrualOrder, len(*ordersToProceessing))

	go func() {
		defer close(ordersToProcessingChan)

		for _, v := range *ordersToProceessing {
			s.l.Debug("order to proceed", zap.String("order", v.Order))
			ordersToProcessingChan <- v
		}
	}()

	return ordersToProcessingChan, nil
}

func (s *AccrualWorker) getAccrualOrdersResp(ctx context.Context, wg *sync.WaitGroup,
	ordersToProcessingChan chan accrualModel.AccrualOrder, responceChan chan accrualModel.AccrualOrder) {
	go func() {
		defer wg.Done()

		for orderToProcessing := range ordersToProcessingChan {
			order, err := s.client.Get(ctx, orderToProcessing.Order)
			if err != nil {
				if !errors.Is(err, customerrors.ErrNoContent) {
					s.l.Debug("fail to get order", zap.Error(errors.Wrap(err, "accrualworker.client.get")))
					continue
				}
				order.Status = accrualModel.PROCESSING.String()
			}

			order.UID = orderToProcessing.UID

			responceChan <- order
		}
	}()
}

func (s *AccrualWorker) updateOrder(ctx context.Context, order accrualModel.AccrualOrder) error {
	status := order.GetStatusIndex()
	if status < 1 {
		return accrualErr.ErrInvalidAccrualOrder
	}

	return s.Services.Order.Update(ctx, order.Order, status.Index(), order.Accrual)
}

func (s *AccrualWorker) updateBalance(ctx context.Context, order accrualModel.AccrualOrder) error {
	return s.Services.Balance.Update(ctx, order.UID, order.Accrual)
}
