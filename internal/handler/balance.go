package handler

import (
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

		balance := userID
		// TODO balanceservice.GetBalance(ctx context.Context, userID int64) (*Ballance, error)
		// balance, err := h.Service.BalanceService.GetBalance(ctx, userID) (*Ballance, error)

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

		// TODO balanceservice.Withdraw(ctx context.Context, userID int64, wdr WithdrawRequest) error
		// err := h.Service.BalanceService.Withdraw(ctx, userID, wdr) error

		return c.JSON(http.StatusOK, echo.Map{"data": data, "user_id": userID})
		// return responseWithCode(c, http.StatusOK)
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

		withdrawals := data
		// TODO balanceservice.GetWithdrawals(ctx context.Context, userID int64) (*WithdrawalsResponse, error)
		//withdrawals, err := h.Service.BalanceService.GetWithdrawals(ctx, userID, wdr) (*WithdrawalsResponse, error)

		return c.JSON(http.StatusOK, echo.Map{"withdrawals": withdrawals, "user_id": userID})
		// return c.JSON(http.StatusOK, withdrawals)
	}

	return fn
}
