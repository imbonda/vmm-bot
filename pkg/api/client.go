// pkg/api/client.go
package api

import (
	"context"

	bybit "github.com/bybit-exchange/bybit.go.api"
	"github.com/imbonda/bybit-vmm-bot/pkg/models"
)

type BybitClient struct {
	client *bybit.Client
}

func NewBybitClient(apiKey, apiSecret string) *BybitClient {
	return &BybitClient{
		client: bybit.NewBybitHttpClient(
			apiKey,
			apiSecret,
			bybit.WithBaseURL(bybit.TESTNET),
		),
	}
}

func (b *BybitClient) GetOrderBook(symbol string) (models.OrderBook, error) {
	res, err := b.client.
		NewUtaBybitServiceWithParams(
			map[string]interface{}{
				"category": "spot",
				"symbol":   symbol,
				"interval": "1",
			},
		).
		GetMarketKline(context.Background())
	return models.OrderBook{
		Category: (res.Result).(map[string]any)["category"].(string),
		Symbol:   (res.Result).(map[string]any)["symbol"].(string),
		List:     (res.Result).(map[string]any)["list"].(any),
	}, err
}

func (b *BybitClient) PlaceOrder(
	symbol,
	side,
	orderType,
	qty,
	price string,
) (any, error) {
	res, err := b.client.
		NewUtaBybitServiceWithParams(
			map[string]interface{}{
				"category":    "linear",
				"symbol":      "BTCUSDT",
				"side":        "Buy",
				"positionIdx": 0,
				"orderType":   "Limit",
				"qty":         "0.001",
				"price":       "10000",
				"timeInForce": "GTC",
			},
		).
		PlaceOrder(context.Background())
	return res, err
}
