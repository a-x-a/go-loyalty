package app

import (
	"context"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	echoSwagger "github.com/swaggo/echo-swagger" // echo-swagger middleware
	"go.uber.org/zap"

	"github.com/a-x-a/go-loyalty/internal/config"
	"github.com/a-x-a/go-loyalty/internal/handler"
	"github.com/a-x-a/go-loyalty/internal/logger"
	"github.com/a-x-a/go-loyalty/internal/service"
	"github.com/a-x-a/go-loyalty/internal/service/accrualservice"
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
	aw  *accrualservice.AccrualWorker
	cfg config.ServiceConfig
}

func NewServer() *Server {
	logLevel := "debug"
	log := logger.InitLogger(logLevel)
	defer log.Sync()

	cfg := config.NewServiceConfig()
	if len(cfg.AccrualSystemAddress) == 0 {
		log.Warn("not defined accrual system address", zap.String("address", cfg.AccrualSystemAddress))
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

	accrualWorker := accrualservice.New(orderService, balanceService, nil, cfg.AccrualSystemAddress, cfg.AccrualFrequency, cfg.AccrualRateLimit, log)

	s := service.New(userService, orderService, balanceService, log)
	h := handler.New(s)

	return &Server{
		e:   echo.New(),
		h:   h,
		l:   log,
		aw:  accrualWorker,
		cfg: cfg,
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		if err := s.aw.Start(ctx); err != nil {
			s.l.Panic("failed to start accrual worker", zap.Error(errors.Wrap(err, "acrualworker.start")))
		}
	}()

	s.e.GET("/swagger/*", echoSwagger.WrapHandler)

	zipRouter := s.e.Group("/api")
	zipRouter.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	// 	Skipper: func(c echo.Context) bool {
	// 		if strings.Contains(c.Request().URL.Path, "swagger") {
	// 			return true
	// 		}
	// 		return false
	// 	},
	// }))

	zipRouter.POST("/user/register", s.h.RegisterUser())
	zipRouter.POST("/user/login", s.h.Login())

	config := util.NewJWTConfig(s.cfg.Secret)
	securityRouter := zipRouter.Group("/user")
	securityRouter.Use(echojwt.WithConfig(config))

	securityRouter.POST("/orders", s.h.UploadOrder())
	securityRouter.GET("/orders", s.h.GetAllOrders())

	securityRouter.GET("/balance", s.h.GetBalance())
	securityRouter.POST("/balance/withdraw", s.h.WithdrawBalance())
	securityRouter.GET("/withdrawals", s.h.WithdrawalsBalance())

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
