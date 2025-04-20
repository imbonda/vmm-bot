package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/imbonda/vmm-bot/cmd/config"
	"github.com/imbonda/vmm-bot/cmd/service"
)

var (
	cfg    *config.Configuration
	logger log.Logger
)

func init() {
	cfg = &config.Configuration{}
	if err := config.LoadConfig(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	logger = cfg.GetLogger()
}

func main() {
	ctx := context.Background()
	traderService, err := service.GetTraderService(ctx, cfg)
	if err != nil {
		level.Error(logger).Log("msg", "failed to initiate a trader", "err", err)
		return
	}
	level.Info(logger).Log("msg", "trader service starting")
	if err = traderService.Start(ctx); err != nil {
		level.Error(logger).Log("msg", "failed to start the trader", "err", err)
		return
	}

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received.
	select {
	case sig := <-quit:
		level.Info(logger).Log("msg", "received signal", "signal", sig)
	case <-ctx.Done():
		level.Debug(logger).Log("msg", "server canceled. shutdown server ...")
	}

	// gracefully shutdown the server, waiting max 5 seconds for current operations to complete
	ctx1, cancel1 := context.WithTimeout(context.Background(), cfg.Service.GracefulShutdown)
	defer cancel1()
	if err = traderService.Shutdown(ctx1); err != nil {
		level.Error(logger).Log("msg", "server shutdown:", "err", err)
		return
	}

	level.Info(logger).Log("msg", "exiting")
}
