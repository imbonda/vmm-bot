package api

import (
	"context"
	"encoding/json"

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

func (b *BybitClient) GetOrderBook(ctx context.Context, symbol string) (*models.OrderBook, error) {
	res, err := b.client.
		NewUtaBybitServiceWithParams(
			map[string]interface{}{
				"category": "spot",
				"symbol":   symbol,
			},
		).
		GetOrderBookInfo(ctx)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(res.Result)
	if err != nil {
		return nil, err
	}
	result := &models.OrderBook{}
	err = json.Unmarshal(data, result)
	return result, err
}

func (b *BybitClient) PlaceOrder(ctx context.Context, order *models.Order) error {
	_, err := b.client.
		NewUtaBybitServiceWithParams(
			map[string]interface{}{
				"category":    "linear",
				"symbol":      order.Symbol,
				"side":        order.Action,
				"positionIdx": 0,
				"orderType":   "Limit",
				"qty":         order.Qty,
				"price":       order.Price,
				"timeInForce": "GTC",
			},
		).
		PlaceOrder(context.Background())
	return err
}
