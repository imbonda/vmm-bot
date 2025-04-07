package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/imbonda/vmm-bot/cmd/interfaces"
	"github.com/imbonda/vmm-bot/cmd/service/http/docs"
	"github.com/imbonda/vmm-bot/cmd/service/models"
	"github.com/imbonda/vmm-bot/internal/trader"
)

// @title Trader API
// @version 1.0
// @description This is a sample server to perform symbol trading requests
type TraderBackend struct {
	addr   string
	server *http.Server
	trader interfaces.Trader
	logger log.Logger
}

func NewTraderService(ctx context.Context, input *models.NewTraderServiceInput) (interfaces.TraderService, error) {
	traderClient, err := trader.NewTrader(ctx, &trader.NewTraderInput{
		ExchangeClient:    input.ExchangeClient,
		PriceOracleClient: input.PriceOracleClient,
		Symbol:            input.Trade.Symbol,
		OracleSymbol:      input.Trade.OracleSymbol,
		SpreadMarginMin:   input.Trade.SpreadMarginMin,
		SpreadMarginMax:   input.Trade.SpreadMarginMax,
		TradeAmountMin:    input.Trade.TradeAmountMin,
		TradeAmountMax:    input.Trade.TradeAmountMax,
		PriceDecimals:     input.Trade.PriceDecimals,
		AmountDecimals:    input.Trade.AmountDecimals,
		Logger:            input.Logger,
	})
	if err != nil {
		return nil, err
	}
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, func(config *ginSwagger.Config) {
		config.InstanceName = docs.SwaggerInfoTraderBackend.InstanceName()
	}))

	backend := &TraderBackend{
		addr: input.Executor.ListenAddress,
		server: &http.Server{
			Addr:    input.Executor.ListenAddress,
			Handler: router,
		},
		trader: traderClient,
		logger: input.Logger,
	}

	// Register routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/trade", backend.handleTrade)
	}

	return backend, nil
}

func (b *TraderBackend) Start(ctx context.Context) error {
	go func() {
		level.Info(b.logger).Log("msg", "starting server on port", "address", b.addr)

		err := b.server.ListenAndServe()
		if err != nil {
			level.Error(b.logger).Log("msg", "error starting server", "err", err)
		}
	}()
	return nil
}

func (b *TraderBackend) Shutdown(ctx context.Context) error {
	return b.server.Shutdown(ctx)
}

// @Summary		Trade once for the configure symbol
// @Description	Call the trade once method to execute a trade
// @ID			trade_once
// @Accept		json
// @Produce		json
// @Success		200		{object}	models.TradeOnceOutput
// @Failure		500		{object} 	errorResponse
// @Router			/api/v1/trade [post]
func (b *TraderBackend) handleTrade(c *gin.Context) {
	output, err := b.trader.TradeOnce(c.Request.Context())
	if err != nil {
		level.Error(b.logger).Log("msg", "error executing trade", "err", err)
		c.JSON(http.StatusInternalServerError, errorResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, output)
}
