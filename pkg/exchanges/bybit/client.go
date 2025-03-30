package bybit

import (
	"context"
	"encoding/json"
	"fmt"

	bybit "github.com/bybit-exchange/bybit.go.api"
	"github.com/go-kit/log"

	"github.com/imbonda/bybit-vmm-bot/pkg/models"
)

type Client struct {
	client *bybit.Client
	logger log.Logger
}

type NewClientInput struct {
	APIKey    string
	APISecret string
	Logger    log.Logger
}

type returnCode int

const (
	successCode returnCode = 0
)

func NewClient(ctx context.Context, input *NewClientInput) (*Client, error) {
	return &Client{
		client: bybit.NewBybitHttpClient(
			input.APISecret,
			input.APISecret,
			bybit.WithBaseURL(bybit.TESTNET),
		),
	}, nil
}

func (api *Client) GetOrderBook(ctx context.Context, symbol string) (*models.OrderBook, error) {
	res, err := api.client.
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
	if err = validateResponse(ctx, res); err != nil {
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

func (api *Client) PlaceOrder(ctx context.Context, order *models.Order) error {
	_, err := api.client.
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

func validateResponse(ctx context.Context, res *bybit.ServerResponse) error {
	if returnCode(res.RetCode) != successCode {
		err := fmt.Errorf("request failed: %v", res.RetMsg)
		return err
	}
	return nil
}
