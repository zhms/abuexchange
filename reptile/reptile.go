package main

import (
	"exreptile/futures"
	"exreptile/market"
	"exreptile/reptile"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	reptile.ConfigInit()

	go futures.FuturesInfo()
	go futures.FuturesKLine()
	go futures.FuturesDepth()
	go futures.FuturesTicker()
	go futures.FuturesTrade()

	go market.MarketKLine()
	go market.MarketDepth()
	go market.MarketTicker()
	go market.MarketTrade()

	fmt.Println("******************start******************")
	for {
		time.Sleep(1 * time.Second)
	}
}
