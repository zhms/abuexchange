package controller

import "xserver/server"

type TransferController struct {
}

func (c *TransferController) Init() {
	group := server.Http().NewGroup("/transfer")
	{
		group.Post("/list", user_list)
	}
}
