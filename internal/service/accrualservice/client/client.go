package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/a-x-a/go-loyalty/internal/service/accrualservice/customerrors"
	"github.com/a-x-a/go-loyalty/internal/service/accrualservice/model"
)

type (
	HTTPClient interface {
		Get(reqURL string) (*http.Response, error)
	}

	AccrualClient struct {
		URL    string
		client HTTPClient
		l      *zap.Logger

		isAvailable atomic.Bool
	}
)

func New(address string, l *zap.Logger) *AccrualClient {
	if !strings.HasPrefix(address, "http") {
		address = fmt.Sprintf("http://%s", address)
	}

	c := AccrualClient{
		URL:    address,
		client: &http.Client{},
		l:      l,
	}

	c.isAvailable.Store(true)

	return &c
}

func (c *AccrualClient) Get(ctx context.Context, number string) (model.AccrualOrder, error) {
	var order model.AccrualOrder

	if !c.isAvailable.Load() {
		return order, customerrors.ErrClientIsNoAvailable
	}

	url := fmt.Sprintf("%s/api/orders/%s", c.URL, number)

	c.l.Info("get order from accrual system", zap.String("URL", url))

	resp, err := c.client.Get(url)
	if err != nil {
		c.l.Debug("failed to get response from accrual system", zap.Error(errors.Wrap(err, "accrualclient.get")))
		return order, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			c.l.Debug("failed to read response from accrual system", zap.Error(errors.Wrap(err, "accrualclient.get")))
			return order, err
		}

		if err := json.Unmarshal(data, &order); err != nil {
			c.l.Debug("failed to unmarshal response from accrual system", zap.Error(errors.Wrap(err, "accrualclient.get")))
			return order, err
		}

		if !order.IsValid() {
			c.l.Debug("invalid accrual order", zap.Any("order", order))
			return order, customerrors.ErrInvalidAccrualOrder
		}

		c.l.Info("get response from accrual system", zap.Any("order", order))

		return order, nil

	case http.StatusNoContent:
		c.l.Info("no content", zap.Int("code", http.StatusNoContent))

		return order, customerrors.ErrNoContent

	case http.StatusTooManyRequests:
		retryHeader := resp.Header.Get("Retry-After")
		retryAfter, err := strconv.Atoi(retryHeader)
		if err != nil {
			return order, customerrors.ErrTooManyRequests
		}

		c.l.Info("too many requests", zap.Int("code", http.StatusNoContent), zap.Int("retry-after", retryAfter))

		go func(wait time.Duration) {
			c.isAvailable.Store(false)
			c.l.Debug("keep client closed", zap.Duration("duration", wait))
			time.Sleep(wait)
			c.isAvailable.Store(true)
			c.l.Debug("open client")
		}(time.Duration(retryAfter) * time.Second)

		return order, customerrors.ErrTooManyRequests
	}

	return order, nil
}
