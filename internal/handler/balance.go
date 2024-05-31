package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/a-x-a/go-loyalty/internal/customerrors"
)

// Balance godoc
//
//	@Summary	Получение баланса
//	@Description	Получение текущего баланса пользователя
//	@Tags	balance
//	@ID	balance-get
//	@Accept	json
//	@Produce	json
//	@Success	200	{object}	Balance	"Успешная обработка запроса"
//	@Failure	401	"Пользователь не авторизован"
//	@Security	ApiKeyAuth
//	@Router	/user/balance [GET]
func (h *Handler) GetBalance() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		uid, err := getUserID(c)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		balance, err := h.s.GetBalance(ctx, uid)
		if err != nil {
			return responseWithError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, balance)
	}

	return fn
}

type withdrawRequest struct {
	Order string  `json:"order" example:"2377225624" validate:"required"`
	Sum   float64 `json:"sum" example:"751" validate:"required"`
} //	@Name	Withdraw

// Balance-withdraw godoc
//
//	@Summary	Запрос на списание средств
//	@Description	Запрос на списание средств
//	@Tags	balance
//	@ID	balance-withdraw
//	@Accept	json
//	@Produce	json
//	@Param	data	body	Withdraw	true	"`order` — номер заказа, `sum` — сумма баллов к списанию в счёт оплаты"
//	@Success	200	"Успешная обработка запроса"
//	@Failure	401	"Пользователь не авторизован"
//	@Failure	402	"На счету недостаточно средств"
//	@Failure	422	"Неверный номер заказа"
//	@Security	ApiKeyAuth
//	@Router	/user/balance/withdraw [POST]
func (h *Handler) WithdrawBalance() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		uid, err := getUserID(c)
		if err != nil {
			return err
		}

		data := &withdrawRequest{}
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

// Withdrawals godoc
//
//	@Summary	Получение информации о выводе средств
//	@Description	Получение информации о выводе средств
//	@Tags	balance
//	@ID	balance-withdrawals
//	@Accept	json
//	@Produce	json
//	@Success	200	{object}	Withdrawals	"Успешная обработка запроса"
//	@Failure	204	"Нет ни одного списания"
//	@Failure	401	"Пользователь не авторизован"
//	@Security	ApiKeyAuth
//	@Router	/user/withdrawals [GET]
func (h *Handler) WithdrawalsBalance() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		uid, err := getUserID(c)
		if err != nil {
			return err
		}

		// data := &withdrawRequwst{}
		// if err := c.Bind(&data); err != nil {
		// 	return responseWithError(c, http.StatusInternalServerError, err)
		// }

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		withdrawals, err := h.s.GetWithdrawalsBalance(ctx, uid)
		if err != nil {
			switch {
			case errors.Is(err, customerrors.ErrNoContent):
				return responseWithError(c, http.StatusNoContent, err)
			}

			return responseWithError(c, http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, withdrawals)
	}

	return fn
}
