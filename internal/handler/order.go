package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Загрузка номера заказа
func (h *Handler) UploadOrder() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		userID, err := getUserId(c)
		if err != nil {
			return err
		}

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		orderNumber := string(body)

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		err = h.Service.UploadOrder(ctx, userID, orderNumber)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		return responseWithCode(c, http.StatusOK)
	}

	return fn
}

// Получение списка загруженных номеров заказов
func (h *Handler) GetAllOrders() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		userID, err := getUserId(c)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		orders, err := h.Service.GetAllOrders(ctx, userID)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, orders)
	}

	return fn
}
