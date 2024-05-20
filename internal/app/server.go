package app

import (
	"context"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/a-x-a/go-loyalty/internal/config"
	"github.com/a-x-a/go-loyalty/internal/handler"
	"github.com/a-x-a/go-loyalty/internal/service"
	"github.com/a-x-a/go-loyalty/internal/service/userservice"
	"github.com/a-x-a/go-loyalty/internal/storage"
	"github.com/a-x-a/go-loyalty/internal/util"
)

type Server struct {
	e   *echo.Echo
	h   *handler.Handler
	cfg config.ServiceConfig
}

func NewServer() *Server {
	cfg := config.NewServiceConfig()

	dbConn, err := storage.NewConnection(cfg.DatabaseURI, "postgres")
	if err != nil {
		panic("unable to database connection")
	}

	// User service.
	userStorage := storage.NewUserStorage(dbConn)
	userService := userservice.New(userStorage)
	// Order service.
	// Ballance service.
	// Accrual service.
	s := service.New(userService)
	h := handler.New(s)

	return &Server{
		e:   echo.New(),
		h:   h,
		cfg: cfg,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.e.POST("/api/user/register", s.h.RegisterUser())
	s.e.POST("/api/user/login", s.h.Login())

	config := util.NewJWTConfig(s.cfg.Secret)
	r := s.e.Group("/api/user") // Authorized user only
	r.Use(echojwt.WithConfig(config))
	// Orders.
	r.POST("/orders", s.h.UploadOrder())
	r.GET("/orders", s.h.GetAllOrders())
	// Ballance.
	r.GET("/balance", s.h.GetBalance())
	r.POST("/balance/withdraw", s.h.WithdrawBalance())
	r.GET("/withdrawals", s.h.WithdrawalsBalance())

	if err := s.e.Start(s.cfg.RunAddress); err != http.ErrServerClosed {
		s.e.Logger.Info(err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.e.Shutdown(ctx); err != nil {
		s.e.Logger.Info(err)
	}
}
