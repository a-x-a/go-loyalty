package app

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/a-x-a/go-loyalty/internal/handler"
	"github.com/a-x-a/go-loyalty/internal/service"
)

type Server struct {
	e *echo.Echo
	h *handler.Handler
}

// TODO move to config
const (
	ADDR   = "localhost:9090"
	SECRET = "secret"
)

// TODO move to model or pkg.JWT
type JWTCustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewServer() *Server {
	s := service.New()
	h := handler.New(s)

	return &Server{
		e: echo.New(),
		h: h,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.e.POST("/api/user/register", s.h.RegisterUser())
	s.e.POST("/api/user/login", s.h.Login())

	// Configure middleware with the custom claims type.
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JWTCustomClaims)
		},
		SigningKey:     []byte(SECRET),
		SuccessHandler: handler.SuccessHandler,
	}

	r := s.e.Group("/api/user") // Authorized user only
	r.Use(echojwt.WithConfig(config))
	// Orders.
	r.POST("/orders", s.h.UploadOrder())
	r.GET("/orders", s.h.GetAllOrders())
	// Ballance.
	r.GET("/balance", s.h.GetBalance())
	r.POST("/balance/withdraw", s.h.WithdrawBalance())
	r.GET("/withdrawals", s.h.WithdrawalsBalance())

	if err := s.e.Start(ADDR); err != http.ErrServerClosed {
		s.e.Logger.Info(err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.e.Shutdown(ctx); err != nil {
		s.e.Logger.Info(err)
	}
}
