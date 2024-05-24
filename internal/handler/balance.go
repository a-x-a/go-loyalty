package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/a-x-a/go-loyalty/internal/customerrors"
)

// Получение текущего баланса пользователя
func (h *Handler) GetBalance() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		uid, err := getUserId(c)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		balance, err := h.s.GetBalance(ctx, uid)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, balance)
	}

	return fn
}

type withdrawRequwst struct {
	Order string  `json:"order" validate:"required"`
	Sum   float64 `json:"sum" validate:"required"`
}

// Запрос на списание средств
func (h *Handler) WithdrawBalance() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		uid, err := getUserId(c)
		if err != nil {
			return err
		}

		data := &withdrawRequwst{}
		if err := c.Bind(&data); err != nil {
			return responseWithError(c, http.StatusInternalServerError, err)
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		err = h.s.WithdrawBalance(ctx, uid, data.Order, data.Sum)
		if err != nil {
			switch {
			case errors.Is(err, customerrors.ErrInsufficientFunds):
				return responseWithError(c, http.StatusPaymentRequired, err)
			case errors.Is(err, customerrors.ErrInvalidOrderNumber):
				return responseWithError(c, http.StatusUnprocessableEntity, err)
			}

			return responseWithError(c, http.StatusInternalServerError, err)
		}

		return responseWithCode(c, http.StatusOK)
	}

	return fn
}

// Получение информации о выводе средств
func (h *Handler) WithdrawalsBalance() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		uid, err := getUserId(c)
		if err != nil {
			return err
		}

		data := &withdrawRequwst{}
		if err := c.Bind(&data); err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		withdrawals, err := h.s.GetWithdrawalsBalance(ctx, uid)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, withdrawals)
	}

	return fn
}
