package main

import (
	"math/rand"
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
	server.Run()
}

