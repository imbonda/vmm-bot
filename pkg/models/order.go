package models

type OrderAction string

const (
	Buy  OrderAction = "buy"
	Sell OrderAction = "sell"
)

type Order struct {
	Action OrderAction
	Symbol string
	Side   string
	Price  float64
	Qty    float64
}
