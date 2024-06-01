package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/a-x-a/go-loyalty/internal/customerrors"
)

type userAccount struct {
	Login    string `json:"login" example:"<login>" validate:"required"`
	Password string `json:"password" example:"<password>" validate:"required"`
} //	@Name	Account

// Register godoc
//
//	@Summary	Регистрация пользователя
//	@Description	Регистрация производится по паре логин/пароль
//	@Tags	user
//	@ID	user-register
//	@Accept	json
//	@Produce	json
//	@Param	data	body	Account	true	"Логин и пароль пользователя"
//	@Success	200	"Пользователь успешно зарегистрирован и аутентифицирован"
//	@Failure	400	"Неверный формат запроса"
//	@Failure	409	"Логин уже занят"
//	@Router	/user/register [post]
func (h *Handler) RegisterUser() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		data := &userAccount{}
		err = FillFromJSON(body, data)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		token, err := h.s.RegisterUser(ctx, data.Login, data.Password)
		if err != nil {
			switch {
			case errors.Is(err, customerrors.ErrInvalidRequestFormat):
				return responseWithError(c, http.StatusBadRequest, err)
			case errors.Is(err, customerrors.ErrInvalidUsernameOrPassword):
				return responseWithError(c, http.StatusBadRequest, err)
			case errors.Is(err, customerrors.ErrUsernameAlreadyTaken):
				return responseWithError(c, http.StatusConflict, err)
			}

			return responseWithError(c, http.StatusInternalServerError, err)
		}

		c.Response().Header().Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

		return responseWithCode(c, http.StatusOK)
	}

	return fn
}

// Login godoc
//
//	@Summary	Аутентификация пользователя
//	@Description	Аутентификация производится по паре логин/пароль
//	@Tags	user
//	@ID	user-login
//	@Accept	json
//	@Produce	json
//	@Param	data	body	Account	true	"Логин и пароль поьзователя"
//	@Success	200	"Пользователь успешно аутентифицирован"
//	@Failure	400	"Неверный формат запроса"
//	@Failure	401	"Неверная пара логин/пароль"
//	@Router	/user/login [post]
func (h *Handler) Login() echo.HandlerFunc {
	var fn = func(c echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		data := &userAccount{}
		err = FillFromJSON(body, data)
		if err != nil {
			return responseWithError(c, http.StatusBadRequest, err)
		}

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		token, err := h.s.Login(ctx, data.Login, data.Password)
		if err != nil {
			switch {
			case errors.Is(err, customerrors.ErrInvalidRequestFormat):
				return responseWithError(c, http.StatusBadRequest, err)
			case errors.Is(err, customerrors.ErrInvalidUsernameOrPassword):
				return responseWithError(c, http.StatusUnauthorized, err)
			}

			return responseWithError(c, http.StatusInternalServerError, err)
		}

		c.Response().Header().Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

		return responseWithCode(c, http.StatusOK)
	}

	return fn
}
