// cmd/main.go
package main

import (
	"log"

	"github.com/imbonda/bybit-vmm-bot/pkg/trader"
)

func main() {
	t := trader.NewTrader()
	if err := t.Start(); err != nil {
		log.Fatalf("Error starting trader: %v", err)
	}
}
