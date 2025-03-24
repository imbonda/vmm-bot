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
}

func (b *OrderBook) Spread() (float64, error) {
	if len(b.Asks) == 0 || len(b.Bids) == 0 {
		return 0, nil
	}
	if len(b.Asks[0]) < 2 || len(b.Bids[0]) < 2 {
		return 0, fmt.Errorf("invalid order book")
	}
	askPrice, err := strconv.ParseFloat(b.Asks[0][0], 64)
	if err != nil {
		return 0, err
	}
	bidPrice, err := strconv.ParseFloat(b.Bids[0][0], 64)
	if err != nil {
		return 0, err
	}
	return askPrice - bidPrice, nil
}
