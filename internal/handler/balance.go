package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Получение текущего баланса пользователя
func (h *Handler) GetBalance() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		userID, err := getUserId(c)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		balance, err := h.Service.GetBalance(ctx, userID)
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
		userID, err := getUserId(c)
		if err != nil {
			return err
		}

		data := &withdrawRequwst{}
		if err := c.Bind(&data); err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		err = h.Service.WithdrawBalance(ctx, userID, data.Order, data.Sum)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		return responseWithCode(c, http.StatusOK)
	}

	return fn
}

// Получение информации о выводе средств
func (h *Handler) WithdrawalsBalance() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		userID, err := getUserId(c)
		if err != nil {
			return err
		}

		data := &withdrawRequwst{}
		if err := c.Bind(&data); err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		withdrawals, err := h.Service.GetWithdrawalsBalance(ctx, userID)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, withdrawals)
	}

	return fn
}
