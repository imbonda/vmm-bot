package biconomy

import (
	"context"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-resty/resty/v2"

	"github.com/imbonda/bybit-vmm-bot/pkg/exchanges/biconomy/hooks"
	"github.com/imbonda/bybit-vmm-bot/pkg/models"
	"github.com/imbonda/bybit-vmm-bot/pkg/utils"
)

// API Configuration
const (
	BaseAPIURL = "https://www.biconomy.com"
	APIV1      = "api/v1"
	APIV2      = "api/v2"
)

type returnCode int

const (
	successCode returnCode = 0
)

type tradingSide string

const (
	ASK = "1"
	BID = "2"
)

func resolveSide(action models.OrderAction) tradingSide {
	if action == models.Buy {
		return BID
	}
	return ASK
}

type Client struct {
	v1     *utils.Endpoint
	v2     *utils.Endpoint
	creds  *utils.Credentials
	client *resty.Client
	logger log.Logger
}

type NewClientInput struct {
	APIKey    string
	APISecret string
	Logger    log.Logger
}

func NewClient(ctx context.Context, input *NewClientInput) (*Client, error) {
	v1 := utils.NewEndpoint(APIV1)
	v2 := utils.NewEndpoint(APIV2)
	creds := &utils.Credentials{
		APIKey:    input.APIKey,
		APISecret: input.APISecret,
	}
	client := resty.New().
		SetBaseURL(BaseAPIURL).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("X-SITE-ID", "127")
	// Add credentials to every request.
	client.OnBeforeRequest(hooks.GetSigAuthBeforeRequestHook(client, creds))
	return &Client{
		v1:     v1,
		v2:     v2,
		creds:  creds,
		client: client,
		logger: input.Logger}, nil
}

func (api *Client) GetOrderBook(ctx context.Context, symbol string) (*models.OrderBook, error) {
	var result struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
	}
	resp, err := api.client.R().
		SetResult(&result).
		SetQueryParam("symbol", symbol).
		Get(api.v1.Join("depth"))
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

	payload := map[string]string{
		"market": order.Symbol,
		"side":   string(resolveSide(order.Action)),
		"amount": utils.FormatFloatToString(order.Qty),
		"price":  utils.FormatFloatToString(order.Price),
	}

	resp, err := api.client.R().
		SetFormData(payload).
		SetResult(&result).
		Post(api.v2.Join("private/trade/limit"))

	if err != nil {
		return err
	}

	if resp.IsError() || result.Code != int(successCode) {
		return fmt.Errorf("placeOrder Biconomy API error: %s", result.Message)
	}

	return nil
}
