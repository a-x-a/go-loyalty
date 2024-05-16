package handler

import (
	"io/ioutil"
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

		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		orderNumber := string(body)

		// TODO orderservice.Upload(ctx context.Context, userID int64, number string) error
		// err := h.Service.OrderService.Upload(ctx, userID, orderNumber) error

		return c.JSON(http.StatusOK, echo.Map{"order": orderNumber, "user_id": userID})
		// return responseWithCode(c, http.StatusOK)
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

		orders := userID
		// TODO orderservice.GetAll(ctx context.Context, userID int64) (*Orders, error)
		// orders, err := h.Service.OrderService.GetBalance(ctx, userID) (*Orders, error)

		return c.JSON(http.StatusOK, orders)
	}

	return fn
}
