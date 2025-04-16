package models

import (
	"fmt"

	"github.com/imbonda/vmm-bot/pkg/utils"
)

type Spread struct {
	Ask float64
	Bid float64
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
		Ask: askPrice,
		Bid: bidPrice,
	}, nil
}

func (s *Spread) Clone() *Spread {
	return &Spread{
		Ask: s.Ask,
		Bid: s.Bid,
	}
}

func (s *Spread) Diff() float64 {
	return s.Ask - s.Bid
}

func (s *Spread) Contains(values ...float64) bool {
	for _, v := range values {
		if v < s.Bid || v > s.Ask {
			return false
		}
	}
	return true
}

func (s *Spread) MarginSpread(lowerPercentage, upperPercentage float64) *Spread {
	ask := s.Bid + s.Diff()*upperPercentage
	bid := s.Bid + s.Diff()*lowerPercentage
	return &Spread{
		Ask: ask,
		Bid: bid,
	}
}
