package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/a-x-a/go-loyalty/docs"
	"github.com/a-x-a/go-loyalty/internal/app"
)

//	@title	API «Гофермарт»
//	@version	0.1
//	@description	API сервер накопительной система лояльности «Гофермарт».

//	@host	localhost:8080
//	@BasePath	/api

//	@securityDefinitions.apikey ApiKeyAuth
//	@in header
//	@name	Authorization

func main() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	srv := app.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go srv.Run(ctx)

	<-sigint

	ctxShutdown, cancelShutdown := context.WithTimeout(ctx, time.Second*5)
	defer cancelShutdown()

	srv.Shutdown(ctxShutdown)
}
