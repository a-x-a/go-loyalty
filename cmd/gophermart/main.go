package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/a-x-a/go-loyalty/internal/app"
)

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
