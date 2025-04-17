package models

import "fmt"

type RawBookTickers []RawBookTicker

func (t RawBookTickers) LastTicker() (*RawBookTicker, error) {
	if len(t) < 1 {
		return nil, fmt.Errorf("missing tickers")
	}
	return &t[0], nil
}

type RawBookTicker struct {
	Symbol    string `json:"symbol"`
	EventType string `json:"eventType"`
	Time      int    `json:"time"`
	AskPrice  string `json:"askPrice"`
	AskAmount string `json:"askVolume"`
	BidPrice  string `json:"bidPrice"`
	BidAmount string `json:"bidVolume"`
}

type BookTicker struct {
	AskPrice  string
	AskAmount string
	BidPrice  string
	BidAmount string
}

type RawPriceTickers []RawPriceTicker

func (t RawPriceTickers) LastTicker() (*RawPriceTicker, error) {
	if len(t) < 1 {
		return nil, fmt.Errorf("missing tickers")
	}
	return &t[0], nil
}

type RawPriceTicker struct {
	Symbol string     `json:"symbol"`
	Trades []RawTrade `json:"trades"`
}

func (t *RawPriceTicker) LastTrade() (*RawTrade, error) {
	if len(t.Trades) < 1 {
		return nil, fmt.Errorf("missing trades")
	}
	return &t.Trades[0], nil
}

type RawTrade struct {
	Timestamp int    `json:"timestamp"`
	TradeID   string `json:"tradeId"`
	Price     string `json:"price"`
	Amount    string `json:"amount"`
	Type      int    `json:"type"`
	Volume    string `json:"volume"`
}

type PriceTicker struct {
	LastPrice string
}
