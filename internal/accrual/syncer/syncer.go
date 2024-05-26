package syncer

import (
	"context"

	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	accrualModel "github.com/a-x-a/go-loyalty/internal/accrual/model"
)

type (
	AccrualStorage interface{}

	AccrualClient interface {
		Get(ctx context.Context, number string) (*accrualModel.AccrualOrder, error)
	}

	AccrualSyncer struct {
		storage AccrualStorage
		client  AccrualClient
		l       *zap.Logger

		frequency        time.Duration
		accrualRateLimit int
	}
)

var (
	ErrNotFoundOrders = errors.New("not orders to sync")
)

func New(storage AccrualStorage, client AccrualClient, frequency time.Duration, accrualRateLimit int, l *zap.Logger) *AccrualSyncer {
	return &AccrualSyncer{
		storage:          storage,
		client:           client,
		l:                l,
		frequency:        frequency,
		accrualRateLimit: accrualRateLimit,
	}
}

func (s *AccrualSyncer) Start(ctx context.Context) error {
	timer := time.NewTicker(s.frequency)
	defer timer.Stop()

	// i := 0
	for {
		select {
		case <-timer.C:
			// i++
			// s.l.Info("[JOB|%v] Sync order info", i)
			ordersToSyncChan, err := s.getOrdersToSync(ctx)
			if err != nil {
				s.l.Debug("failed to get not processed orders", zap.Error(errors.Wrap(err, "accrualsyncer.getorderstosync")))
				continue
			}

			responceChan := make(chan accrualModel.AccrualOrder, len(ordersToSyncChan))

			var wg sync.WaitGroup
			for w := 1; w <= s.accrualRateLimit; w++ {
				wg.Add(1)
				s.getAccrualOrdersResp(ctx, &wg, ordersToSyncChan, responceChan)
			}

			go func() {
				wg.Wait()
				close(responceChan)
			}()

			for order := range responceChan {
				s.l.Debug("get order from responce channel", zap.Any("order", order))
				err = s.updateOrder(ctx, order)
				if err != nil {
					s.l.Debug("update order", zap.Error(errors.Wrap(err, "accrualsyncer.updateorder")))
					continue
				}
			}

		case <-ctx.Done():
			return nil
		}
	}
}

func (s *AccrualSyncer) getOrdersToSync(ctx context.Context) (chan accrualModel.AccrualOrder, error) {
	// TODO
	tx, err := s.storage.OpenTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open transaction, %w", err)
	}
	defer func() {
		_ = j.repo.Commit(ctx, tx)
	}()

	ordersToSync, err := j.repo.GetNotProcessedOrders(ctx, tx)
	if len(ordersToSync) == 0 {
		return nil, ErrNotFoundOrders
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get not porocessed orders from db: %w", err)
	}
	ordersToSyncCh := make(chan models.Order, len(ordersToSync))

	go func() {
		defer close(ordersToSyncCh)

		for i := 0; i < len(ordersToSync); i++ {
			j.logger.Debug("get order to sync", ordersToSync[i].ID)
			ordersToSyncCh <- ordersToSync[i]
		}
	}()
	return ordersToSyncCh, nil
}

func (s *AccrualSyncer) getAccrualOrdersResp(
	ctx context.Context,
	wg *sync.WaitGroup,
	ordersToSync chan models.Order,
	responceChan chan accrualModel.AccrualOrder) {
	go func() {
		defer wg.Done()

		for orderToSync := range ordersToSync {
			order, err := s.client.Get(ctx, orderToSync.ID)
			if err != nil {
				s.l.Debug("fail to get order", zap.Error(errors.Wrap(err, "accrualsyncer.client.get")))
				continue
			}

			responceChan <- *order
		}
	}()
}

func (s *AccrualSyncer) updateOrder(ctx context.Context, order accrualModel.AccrualOrder) error {
	// TODO
	tx, err := j.repo.OpenTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to open transaction, %w", err)
	}

	order, err := j.repo.GetOrder(ctx, tx, order.OrderID, true) // for_update
	if err != nil {
		_ = j.repo.Rollback(ctx, tx)
		return fmt.Errorf("failed to get order for update, %w", err)
	}

	user, err := j.repo.GetUserByID(ctx, tx, order.UserID, true) // for_update
	if err != nil {
		_ = j.repo.Rollback(ctx, tx)
		return fmt.Errorf("failed to get order for update, %w", err)
	}
	if accrualOrder.Status == accrualModels.REGISTERED.String() {
		accrualOrder.Status = models.NEW.String()
	}
	if err = j.repo.UpdateOrder(ctx, tx, order.ID, accrualOrder.Accrual, accrualOrder.Status); err != nil {
		_ = j.repo.Rollback(ctx, tx)
		return fmt.Errorf("failed to update order, %w", err)
	}
	if accrualOrder.Accrual != 0 {
		if err = j.repo.UpdateUserBalance(ctx, tx, user.ID, user.Balance+accrualOrder.Accrual, user.Withdrawn); err != nil {
			_ = j.repo.Rollback(ctx, tx)
			return fmt.Errorf("failed to update user balance, %w", err)
		}
	}

	if err = j.repo.Commit(ctx, tx); err != nil {
		return fmt.Errorf("failef to commit, %w", err)
	}

	return nil
}
