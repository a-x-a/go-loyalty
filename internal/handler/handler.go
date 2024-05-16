package handler

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/a-x-a/go-loyalty/internal/service"
)

type Handler struct {
	Service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func SuccessHandler(c echo.Context) {
	// TODO move to model or pkg.JWT
	type JWTCustomClaims struct {
		UserID int64 `json:"user_id"`
		jwt.RegisteredClaims
	}

	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return
	}

	claims := token.Claims.(*JWTCustomClaims)
	if !ok {
		return
	}

	c.Set("userID", claims.UserID)
}

func responseWithError(c echo.Context, code int, err error) error {
	// resp := fmt.Sprintf("%d: %s", code, err.Error())
	return c.JSON(code, err)
}

func responseWithCode(c echo.Context, code int) error {
	// c.Response().WriteHeader(code)
	return c.NoContent(code)
}

func getUserId(c echo.Context) (int64, error) {
	userID, ok := c.Get("userID").(int64)
	if !ok {
		return 0, errors.New("user id is not defined")
	}

	return userID, nil
}
