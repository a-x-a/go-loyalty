package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type registerUser struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Регистрация пользователя
func (h *Handler) RegisterUser() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		data := &registerUser{}
		if err := c.Bind(&data); err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		token := data
		// TODO userservice.RegisterUser(ctx context.Context, login, password string) (error)
		// err := h.Service.UserService.RegisterUser(ctx, data.Login, data.Password)
		// token, err := h.Service.UserService.Login(cts, data.Login, data.Password)

		c.Response().Header().Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

		return responseWithCode(c, http.StatusOK)
	}

	return fn
}

type loginUser struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Аутентификация пользователя
func (h *Handler) Login() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		data := &loginUser{}
		if err := c.Bind(&data); err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		token := data
		// TODO userservice.Login(ctx context.Context, login, password string) (string, error)
		// token, err := h.Service.UserService.Login(ctx, data.Login, data.Password)

		c.Response().Header().Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

		return responseWithCode(c, http.StatusOK)
	}

	return fn
}
