// pkg/models/order.go
package models

type OrderBook struct {
	Category string `category`
	Symbol   string `symbol`
	List     any    `list`
}
