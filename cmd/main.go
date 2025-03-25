// cmd/main.go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/imbonda/bybit-vmm-bot/cmd/config"
	"github.com/imbonda/bybit-vmm-bot/cmd/internal/trader"
	"github.com/imbonda/bybit-vmm-bot/pkg/exchanges/bybit"
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
	bybitClient, err := bybit.NewClient(ctx, &bybit.NewClientInput{
		APIKey:    cfg.BybitAPIKey,
		APISecret: cfg.BybitAPISecret,
		Logger:    nil,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to create bybit client", "err", err)
		return
	}
	bybitTrader, err := trader.NewTrader(ctx, &trader.NewTraderInput{
		Symbol:               cfg.Symbol,
		ExchangeClient:       bybitClient,
		MinExecutionDuration: cfg.MinExecutionDuration,
		MaxExecutionDuration: cfg.MaxExecutionDuration,
		Logger:               logger,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to initiate a trader", "err", err)
		return
	}
	if err = bybitTrader.Start(ctx); err != nil {
		level.Error(logger).Log("msg", "failed to start the trader", "err", err)
		return
	}

	level.Info(logger).Log("msg", "auto trading started")

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
	ctx1, cancel1 := context.WithTimeout(context.Background(), cfg.GraceFullShutdown)
	defer cancel1()
	if err = bybitTrader.Shutdown(ctx1); err != nil {
		level.Error(logger).Log("msg", "server shutdown:", "err", err)
		return
	}

	level.Info(logger).Log("msg", "cve matcher server exiting")
}
