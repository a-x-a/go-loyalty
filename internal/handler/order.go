package handler

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/a-x-a/go-loyalty/internal/customerrors"
)

// Order-upload godoc
//
//	@Summary	Загрузка номера заказа
//	@Description	Загрузка номера заказа
//	@Tags	orders
//	@ID	orders-upload
//	@Accept	plain
//	@Produce	json
//	@Param	number	body	string	true	"Номер заказа"
//	@Success	200	"Номер заказа уже был загружен"
//	@Success	202	"Новый номер заказа принят в обработку"
//	@Failure	400	"Неверный формат запроса"
//	@Failure	401	"Пользователь не авторизован"
//	@Failure	409	"Номер заказа уже был загружен другим пользователем"
//	@Failure	422	"Неверный номер заказа"
//	@Security	ApiKeyAuth
//	@Router	/user/orders [POST]
func (h *Handler) UploadOrder() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		uid, err := getUserID(c)
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

// Orders-get godoc
//
//	@Summary	Получение списка заказов
//	@Description	Получение списка загруженных номеров заказов
//	@Tags	orders
//	@ID	orders-get
//	@Accept	json
//	@Produce	json
//	@Success	200	{object}	Orders	"Успешная обработка запроса"
//	@Success	204	"Нет данных для ответа"
//	@Failure	401	"Пользователь не авторизован"
//	@Security	ApiKeyAuth
//	@Router	/user/orders [GET]
func (h *Handler) GetAllOrders() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		uid, err := getUserID(c)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		orders, err := h.s.GetAllOrders(ctx, uid)
		if err != nil {
			switch {
			case errors.Is(err, customerrors.ErrNoContent):
				return responseWithCode(c, http.StatusNoContent)
			}

			return responseWithError(c, http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, orders)
	}

	return fn
}
