package models

type OrderAction string

const (
	Buy  OrderAction = "buy"
	Sell OrderAction = "sell"
)

type Order struct {
	Symbol string
	Action OrderAction
	Price  float64
	Qty    float64
}
