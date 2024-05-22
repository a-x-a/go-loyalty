package app

import (
	"context"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/a-x-a/go-loyalty/internal/config"
	"github.com/a-x-a/go-loyalty/internal/handler"
	"github.com/a-x-a/go-loyalty/internal/logger"
	"github.com/a-x-a/go-loyalty/internal/service"
	"github.com/a-x-a/go-loyalty/internal/service/userservice"
	"github.com/a-x-a/go-loyalty/internal/storage"
	"github.com/a-x-a/go-loyalty/internal/util"
)

type Server struct {
	e   *echo.Echo
	h   *handler.Handler
	l   *zap.Logger
	cfg config.ServiceConfig
}

func NewServer() *Server {
	logLevel := "debug"
	log := logger.InitLogger(logLevel)
	defer log.Sync()

	cfg := config.NewServiceConfig()

	dbConn, err := storage.NewConnection(cfg.DatabaseURI, "postgres")
	if err != nil {
		log.Panic("unable to database connection", zap.Error(errors.Wrap(err, "storage.newconnection")))
	}

	if err := migrationRun(cfg.DatabaseURI, log); err != nil {
		log.Panic("unable to init DB", zap.Error(errors.Wrap(err, "migrationrun")))
	}

	// User service.
	userStorage := storage.NewUserStorage(dbConn, log)
	userService := userservice.New(userStorage, cfg, log)
	// Order service.
	// Ballance service.
	// Accrual service.
	s := service.New(userService, log)
	h := handler.New(s)

	return &Server{
		e:   echo.New(),
		h:   h,
		l:   log,
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
		s.l.Warn("failed to start http server", zap.Error(errors.Wrap(err, "echo.start")))
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.e.Shutdown(ctx); err != nil {
		s.l.Warn("server shutdowning error", zap.Error(errors.Wrap(err, "shutdown")))
	}
}
