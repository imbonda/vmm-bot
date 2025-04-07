package models

type OrderAction string

const (
	Buy  OrderAction = "buy"
	Sell OrderAction = "sell"
)

type Order struct {
	Symbol string
	Price  string
	Qty    string
	Action OrderAction
}
