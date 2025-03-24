// pkg/models/order.go
package models

type Order struct {
	Symbol string
	Side   string
	Price  float64
	Qty    float64
}
