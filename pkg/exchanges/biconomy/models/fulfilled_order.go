package models

type RawFulfilledOrder struct {
	Amount     string  `json:"amount"`
	OrderID    int     `json:"id"`
	Market     string  `json:"market"`
	Price      string  `json:"price"`
	Side       int     `json:"side"`
	CreatedAt  float64 `json:"ctime"`
	ModifiedAt float64 `json:"mtime"`
}
