package models

import (
	"fmt"
	"strconv"
)

type OrderBook struct {
	Symbol    string     `json:"s"`
	Asks      [][]string `json:"a"`
	Bids      [][]string `json:"b"`
	Timestamp int64      `json:"ts"`
	UpdateID  int64      `json:"u"`
	Sequence  int64      `json:"seq"`
	CTime     int64      `json:"cts"`
	bestAsk   *float64
	bestBid   *float64
}

func (b *OrderBook) Ask() (float64, error) {
	if b.bestAsk != nil {
		return *b.bestAsk, nil
	}
	if len(b.Asks) == 0 {
		return 0, nil
	}
	if len(b.Asks[0]) < 2 {
		return 0, fmt.Errorf("invalid order book")
	}
	askPrice, err := strconv.ParseFloat(b.Asks[0][0], 64)
	if err != nil {
		return 0, err
	}
	b.bestAsk = &askPrice
	return askPrice, nil
}

func (b *OrderBook) Bid() (float64, error) {
	if b.bestBid != nil {
		return *b.bestBid, nil
	}
	if len(b.Bids) == 0 {
		return 0, nil
	}
	if len(b.Bids[0]) < 2 {
		return 0, fmt.Errorf("invalid order book")
	}
	bidPrice, err := strconv.ParseFloat(b.Bids[0][0], 64)
	if err != nil {
		return 0, err
	}
	b.bestBid = &bidPrice
	return bidPrice, nil
}

func (b *OrderBook) Spread() (float64, error) {
	askPrice, err := b.Ask()
	if err != nil {
		return 0, err
	}
	bidPrice, err := b.Bid()
	if err != nil {
		return 0, err
	}
	return askPrice - bidPrice, nil
}
