package main

import (
	"exreptile/reptile"
	"fmt"
	"time"
)
func main() {
	reptile.ConfigInit()

	go reptile.MarketKLine()
	go reptile.FuturesKLine()

	go reptile.MarketTicker()
	go reptile.FuturesTicker()

	go reptile.MarketTrade()
	go reptile.FuturesTrade()

	go reptile.MarketDepth()
	go reptile.FuturesDepth()

	go reptile.FuturesInfo()

	fmt.Println("******************start******************")
	for {
		time.Sleep(1 * time.Second)
	}
}
