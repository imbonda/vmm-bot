package biconomy

import (
	"context"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-resty/resty/v2"

	"github.com/imbonda/bybit-vmm-bot/pkg/models"
)

type credentials struct {
	apiKey    string
	apiSecret string
}

type Client struct {
	credentials
	client *resty.Client
	logger log.Logger
}

type NewClientInput struct {
	APIKey    string
	APISecret string
	Logger    log.Logger
}

// API Configuration
const (
	BaseAPIURL = "https://www.biconomy.com/api/v1/"
)

type returnCode int

const (
	successCode returnCode = 0
)

type tradingSide int

const (
	ASK = 1
	BID = 2
)

func NewClient(ctx context.Context, input *NewClientInput) (*Client, error) {
	credentials := credentials{
		apiKey:    input.APIKey,
		apiSecret: input.APISecret,
	}
	client := resty.New().
		SetBaseURL(BaseAPIURL).
		SetHeader("Content-Type", "application/json")
	// Add credentials to every request.
	client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		request.SetQueryParam("api_key", credentials.apiKey)
		request.SetQueryParam("secret_key", credentials.apiSecret)
		return nil
	})
	return &Client{
		credentials: credentials,
		client:      client,
		logger:      input.Logger}, nil
}

func (api *Client) GetOrderBook(ctx context.Context, symbol string) (*models.OrderBook, error) {
	var result struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
	}
	resp, err := api.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("depth?symbol=%s", symbol))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("orderBook Biconomy API error: %s", resp.Status())
	}
	return &models.OrderBook{
		Symbol: symbol,
		Asks:   result.Asks,
		Bids:   result.Bids,
	}, nil
}

func (api *Client) PlaceOrder(ctx context.Context, order *models.Order) error {
	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Result  struct {
			Amount     string  `json:"amount"`
			OrderID    int     `json:"id"`
			Market     string  `json:"market"`
			Price      string  `json:"price"`
			Side       int     `json:"side"`
			CreatedAt  float64 `json:"ctime"`
			ModifiedAt float64 `json:"mtime"`
		} `json:"result"`
	}

	payload := map[string]interface{}{
		"market": order.Symbol,
		"side":   resolveSide(order.Action),
		"amount": order.Qty,
		"price":  order.Price,
	}

	resp, err := api.client.R().
		SetBody(payload).
		SetResult(&result).
		Post("private/trade/limit")

	if err != nil {
		return err
	}

	if resp.IsError() || result.Code != int(successCode) {
		return fmt.Errorf("placeOrder Biconomy API error: %s", result.Message)
	}

	return nil
}

func resolveSide(action models.OrderAction) tradingSide {
	if action == models.Buy {
		return BID
	}
	return ASK
}
