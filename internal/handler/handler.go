package handler

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/a-x-a/go-loyalty/internal/model"
)

type (
	APIService interface {
		// UserService
		RegisterUser(ctx context.Context, login, password string) (string, error)
		Login(ctx context.Context, login, password string) (string, error)
		// OrderService
		UploadOrder(ctx context.Context, userID int64, number string) error
		GetAllOrders(ctx context.Context, userID int64) (*model.Orders, error)
		// BallanceService
		GetBalance(ctx context.Context, userID int64) (*model.Balance, error)
		WithdrawBalance(ctx context.Context, userID int64, number string, sum float64) error
		GetWithdrawalsBalance(ctx context.Context, userID int64) (*model.Withdrawals, error)
	}

	Handler struct {
		Service APIService
	}
)

func New(s APIService) *Handler {
	return &Handler{
		Service: s,
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
