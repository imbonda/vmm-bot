package models

type RawOrderBook struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}
