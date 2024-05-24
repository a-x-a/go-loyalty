package handler

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/a-x-a/go-loyalty/internal/customerrors"
)

// Загрузка номера заказа.
func (h *Handler) UploadOrder() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		uid, err := getUserId(c)
		if err != nil {
			return err
		}

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		orderNumber := string(body)
		if len(orderNumber) == 0 {
			return responseWithError(c, http.StatusBadRequest, customerrors.ErrInvalidRequestFormat)
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		err = h.s.UploadOrder(ctx, uid, orderNumber)
		if err != nil {
			switch {
			case errors.Is(err, customerrors.ErrOrderUploadedByUser):
				return responseWithCode(c, http.StatusOK)
			case errors.Is(err, customerrors.ErrOrderUploadedByAnotherUser):
				return responseWithError(c, http.StatusConflict, err)
			case errors.Is(err, customerrors.ErrInvalidOrderNumber):
				return responseWithError(c, http.StatusUnprocessableEntity, err)
			}

			return responseWithError(c, http.StatusInternalServerError, err)
		}

		return responseWithCode(c, http.StatusAccepted)
	}

	return fn
}

// Получение списка загруженных номеров заказов
func (h *Handler) GetAllOrders() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		uid, err := getUserId(c)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		orders, err := h.s.GetAllOrders(ctx, uid)
		if err != nil {
			switch {
			case errors.Is(err, customerrors.ErrNotOrders):
				return responseWithCode(c, http.StatusNoContent)
			}

			return responseWithError(c, http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, orders)
	}

	return fn
}
