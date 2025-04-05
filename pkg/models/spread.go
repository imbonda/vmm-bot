package models

import (
	"fmt"

	"github.com/imbonda/vmm-bot/pkg/utils"
)

type Spread struct {
	Ask  float64
	Bid  float64
	Diff float64
}

func NewSpread(ask, bid string) (*Spread, error) {
	askPrice, err := utils.ParseFloat(ask)
	if err != nil {
		return nil, fmt.Errorf("failed parse ask price: %s", ask)
	}
	bidPrice, err := utils.ParseFloat(bid)
	if err != nil {
		return nil, fmt.Errorf("failed parse bid price: %s", bid)
	}
	return &Spread{
		Ask:  askPrice,
		Bid:  bidPrice,
		Diff: askPrice - bidPrice,
	}, nil
}
