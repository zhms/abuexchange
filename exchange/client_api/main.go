package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
	"xserver/controller"
	"xserver/server"
)

func main() {
    rand.Seed(time.Now().Unix())
	server.Init()
	new(controller.VerifyController).Init()
	new(controller.UserController).Init()
	new(controller.AssetController).Init()

	x := []int{1,5,6,7}
	sort.Slice(x, func(i, j int) bool {
		a := x[i]
		b := x[j]
		return a < b
	})
	fmt.Println(x)

	server.Run()
}

