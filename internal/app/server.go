package app

import (
	"context"
	"net/http"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	accrualclient "github.com/a-x-a/go-loyalty/internal/accrual/client"
	accrualsyncer "github.com/a-x-a/go-loyalty/internal/accrual/syncer"
	"github.com/a-x-a/go-loyalty/internal/config"
	"github.com/a-x-a/go-loyalty/internal/handler"
	"github.com/a-x-a/go-loyalty/internal/logger"
	"github.com/a-x-a/go-loyalty/internal/service"
	"github.com/a-x-a/go-loyalty/internal/service/authservice"
	"github.com/a-x-a/go-loyalty/internal/service/balanceservice"
	"github.com/a-x-a/go-loyalty/internal/service/orderservice"
	"github.com/a-x-a/go-loyalty/internal/storage"
	"github.com/a-x-a/go-loyalty/internal/util"
)

type Server struct {
	e   *echo.Echo
	h   *handler.Handler
	l   *zap.Logger
	s   *accrualsyncer.AccrualSyncer
	cfg config.ServiceConfig
}

func NewServer() *Server {
	logLevel := "debug"
	log := logger.InitLogger(logLevel)
	defer log.Sync()

	cfg := config.NewServiceConfig()
	if len(cfg.AccrualSystemAddress) == 0 {
		log.Panic("not defined accrual system address")
	}

	dbConn, err := storage.NewConnection(cfg.DatabaseURI, "postgres")
	if err != nil {
		log.Panic("unable to database connection", zap.Error(errors.Wrap(err, "storage.newconnection")))
	}

	if err := migrationRun(cfg.DatabaseURI, log); err != nil {
		log.Panic("unable to init DB", zap.Error(errors.Wrap(err, "migrationrun")))
	}

	// User service.
	userStorage := storage.NewUserStorage(dbConn, log)
	userService := authservice.New(userStorage, cfg, log)
	// Order service.
	orderStorage := storage.NewOrderStorage(dbConn, log)
	orderService := orderservice.New(orderStorage, cfg, log)
	// Ballance service.
	balanceStorage := storage.NewBalanceStorage(dbConn, log)
	balanceService := balanceservice.New(balanceStorage, cfg, log)

	accrualClient := accrualclient.New(cfg.AccrualSystemAddress, log)
	accrualSyncer := accrualsyncer.New(orderService, balanceService, nil, accrualClient, time.Second*15, 5, log)

	s := service.New(userService, orderService, balanceService, log)
	h := handler.New(s)

	return &Server{
		e:   echo.New(),
		h:   h,
		l:   log,
		s:   accrualSyncer,
		cfg: cfg,
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		if err := s.s.Start(ctx); err != nil {
			s.l.Panic("failed to start accrual sycronization", zap.Error(errors.Wrap(err, "syncer.start")))
		}
	}()

	s.e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	s.e.POST("/api/user/register", s.h.RegisterUser())
	s.e.POST("/api/user/login", s.h.Login())

	config := util.NewJWTConfig(s.cfg.Secret)
	r := s.e.Group("/api/user")
	r.Use(echojwt.WithConfig(config))

	r.POST("/orders", s.h.UploadOrder())
	r.GET("/orders", s.h.GetAllOrders())

	r.GET("/balance", s.h.GetBalance())
	r.POST("/balance/withdraw", s.h.WithdrawBalance())
	r.GET("/withdrawals", s.h.WithdrawalsBalance())

	s.l.Info("start http server", zap.String("address", s.cfg.RunAddress))

	if err := s.e.Start(s.cfg.RunAddress); err != http.ErrServerClosed {
		s.l.Panic("failed to start http server", zap.Error(errors.Wrap(err, "echo.start")))
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) {
	s.l.Info("start server shutdown...")

	if err := s.e.Shutdown(ctx); err != nil {
		s.l.Warn("server shutdowning error", zap.Error(errors.Wrap(err, "server.shutdown")))
	}

	s.l.Info("successfully server shutdowning")
}
