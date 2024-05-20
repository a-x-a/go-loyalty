package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type registerUser struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Регистрация пользователя.
func (h *Handler) RegisterUser() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		data := &registerUser{}
		if err := c.Bind(&data); err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		token, err := h.Service.RegisterUser(ctx, data.Login, data.Password)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		c.Response().Header().Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

		return responseWithCode(c, http.StatusOK)
	}

	return fn
}

type loginUser struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Аутентификация пользователя.
func (h *Handler) Login() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		data := &loginUser{}
		if err := c.Bind(&data); err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		token, err := h.Service.Login(ctx, data.Login, data.Password)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		c.Response().Header().Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

		return responseWithCode(c, http.StatusOK)
	}

	return fn
}
